package codegen

import (
	"os"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
	"github.com/wf001/modo/util"
)

type assembler struct {
	program *mTypes.Program
}

type context struct {
	mod      *ir.Module
	function *ir.Func
	block    *ir.Block
	prog     *mTypes.Program
	scope    *mTypes.Node
}

func newI32(s string) *constant.Int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newI32: %s", err)
	}
	return constant.NewInt(types.I32, i)
}

func newArray(length uint64) *types.ArrayType {
	return types.NewArray(length, types.I8)
}

func newStr(block *ir.Block, n *mTypes.Node) *ir.InstGetElementPtr {
	strType := newArray(n.Len)
	strPtr := block.NewAlloca(strType)

	strGEP := block.NewGetElementPtr(strType, strPtr)

	block.NewStore(constant.NewCharArray([]byte(n.Val)), strPtr)

	return strGEP
}

func isNumericIR(arg value.Value) bool {
	return util.EqualType(arg, (*constant.Int)(nil)) ||
		util.EqualType(arg, (*ir.InstAdd)(nil)) ||
		util.EqualType(arg, (*ir.InstLoad)(nil)) ||
		// TODO: %d -> %s
		util.EqualType(arg, (*ir.InstICmp)(nil))
}

func isStringIR(arg value.Value) bool {
	return util.EqualType(arg, (*ir.InstGetElementPtr)(nil))
}

func (ctx *context) gen(node *mTypes.Node) value.Value {

	if node.IsKind(mTypes.ND_DECLARE) {
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_VAR_DECLARE) {
		if node.Val == "main" {
			// means declaring main function regarded as entrypoint

			fnc := ctx.mod.NewFunc(
				"main",
				types.I32,
			)
			llBlock := fnc.NewBlock("")

			ctx.function = fnc
			ctx.block = llBlock
			res := ctx.gen(node.Child)
			llBlock.NewCall(res)
			llBlock.NewRet(newI32("0"))

		} else {
			// means declaring global variable or function named except main
			retType := types.I32 // TODO: to be changable
			funcName := node.GetFuncName()

			fnc := ctx.mod.NewFunc(
				funcName,
				retType,
			)
			llBlock := fnc.NewBlock("")

			ctx.function = fnc
			res := ctx.gen(node.Child)
			node.FuncPtr = fnc
			llBlock.NewRet(res)

		}
	} else if node.IsKind(mTypes.ND_VAR_REFERENCE) {
		// PERFORMANCE: too redundant
		// find in local variable which is declared with let
		for s := ctx.scope; s != nil; s = s.Next {
			if s.Val == node.Val {
				if s.Child.IsType(mTypes.TY_INT32) {
					node.Type = mTypes.TY_INT32
					return ctx.block.NewLoad(types.I32, s.VarPtr)

				} else if s.Child.IsType(mTypes.TY_STR) {
					node.Type = mTypes.TY_STR
					return s.VarPtr

				} else {
					log.Panic("unresolved variable: have %+v, %+s", s)
				}
			}
		}

		// find in global variable which is declared with def
		for declare := ctx.prog.Declares; declare != nil; declare = declare.Next {
			if declare.Child.Val == node.Val {
				return ctx.block.NewCall(declare.Child.FuncPtr)
			}
		}

		log.Panic("unresolved symbol: '%s'", node.Val)

	} else if node.IsKind(mTypes.ND_LAMBDA) {
		funcFn := ctx.mod.NewFunc(
			node.GetFuncName(),
			types.I32,
		)
		llBlock := funcFn.NewBlock(node.GetBlockName("fn.entry"))

		ctx.function = funcFn
		ctx.block = llBlock
		res := ctx.gen(node.Child)
		if res != nil {
			llBlock.NewRet(res)
		}
		return funcFn

	} else if node.IsKind(mTypes.ND_BIND) {
		for varDeclare := node.Bind; varDeclare != nil; varDeclare = varDeclare.Next {
			v := ctx.gen(varDeclare.Child)
			if ctx.scope.Child == nil {
				ctx.scope = varDeclare
			} else {
				ctx.scope.Next = varDeclare
			}

			if isNumericIR(v) {
				varDeclare.VarPtr = ctx.block.NewAlloca(types.I32)
				ctx.block.NewStore(v, varDeclare.VarPtr)

			} else if isStringIR(v) {
				varDeclare.VarPtr = v
			}

		}
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_EXPR) {
		var res value.Value
		for child := node.Child; child != nil; child = child.Next {
			res = ctx.gen(child)
		}
		return res

	} else if node.IsKind(mTypes.ND_LIBCALL) {
		// means calling standard library
		arg := ctx.gen(node.Child)
		node.Child.IRValue = arg

		for n := node.Child.Next; n != nil; n = n.Next {
			arg := ctx.gen(n)
			n.IRValue = arg
		}

		libFunc := lib.LibInsts[node.Val]
		return libFunc(ctx.block, ctx.prog.BuiltinLibs, node.Child)

	} else if node.IsKindNary() {
		// nary takes more than 2 arguments
		child := node.Child
		fst := ctx.gen(child)

		child = child.Next
		snd := ctx.gen(child)

		nary := operatorInsts[node.Kind]
		res := nary(ctx.block, fst, snd)

		for child = child.Next; child != nil; child = child.Next {
			fst = res
			snd = ctx.gen(child)
			res = nary(ctx.block, fst, snd)
		}
		return res

	} else if node.IsKindBinary() {
		// binary takes exactly 2 arguments
		child := node.Child
		fst := ctx.gen(child)

		child = child.Next
		snd := ctx.gen(child)

		binary := operatorInsts[node.Kind]
		res := binary(ctx.block, fst, snd)
		return res

	} else if node.IsKind(mTypes.ND_IF) {
		condBlock := ctx.function.NewBlock(node.GetBlockName("if.cond"))
		ctx.block.NewBr(condBlock)

		thenBlock := ctx.function.NewBlock(node.GetBlockName("if.then"))
		elseBlock := ctx.function.NewBlock(node.GetBlockName("if.else"))
		exitBlock := ctx.function.NewBlock(node.GetBlockName("if.exit"))

		ctx.block = condBlock
		cond := ctx.gen(node.Cond)

		exitBlock.NewRet(newI32("0"))

		thenBlock.NewBr(exitBlock)
		elseBlock.NewBr(exitBlock)

		ctx.block = thenBlock
		ctx.gen(node.Then)

		ctx.block = elseBlock
		ctx.gen(node.Else)

		condBlock.NewCondBr(cond, thenBlock, elseBlock)

	} else if node.IsKind(mTypes.ND_SCALAR) {
		if node.IsType(mTypes.TY_INT32) {
			return newI32(node.Val)

		} else if node.IsType(mTypes.TY_STR) {
			return newStr(ctx.block, node)

		} else {
			log.Panic("unresolved Scalar: have %+v", node)
		}

	} else {
		log.Panic("unresolved Nodekind: have %+v", node)
	}
	return nil
}

func constructModule(prog *mTypes.Program) *ir.Module {
	module := ir.NewModule()
	prog.BuiltinLibs = &mTypes.BuiltinLibProp{}
	lib.DeclareBuiltin(module, prog.BuiltinLibs)

	for declare := prog.Declares; declare != nil; declare = declare.Next {
		c := &context{
			mod:   module,
			scope: &mTypes.Node{},
			prog:  prog,
		}
		c.gen(declare)
	}

	return module
}

func Construct(program *mTypes.Program) *assembler {
	return &assembler{
		program: program,
	}
}

func (a assembler) GenFrontend(llName string, asmName string) {
	log.DebugMessage("ir module constructing")
	module := constructModule(a.program)
	log.DebugMessage("ir module constructed")
	log.Debug("[IR]\n%s\n", module.String())

	err := os.WriteFile(llName, []byte(module.String()), 0600)
	if err != nil {
		log.Panic("fail to write ll: %+v", map[string]interface{}{"err": err, "llName": llName})
	}
	log.Debug("written ll: %s", llName)

}
