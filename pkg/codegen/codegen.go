package codegen

import (
	"fmt"
	"os"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/wf001/modo/pkg/codegen/lib"
	"github.com/wf001/modo/pkg/log"
	mTypes "github.com/wf001/modo/pkg/types"
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
	argument *mTypes.Node
}

func newI32(s string) *constant.Int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		log.Panic("fail to newI32: %s", err)
	}
	return constant.NewInt(types.I32, i)
}

func newStr(ctx *context, n *mTypes.Node) *ir.InstLoad {
	strConst := constant.NewCharArrayFromString(n.Val)
	globalStr := ctx.mod.NewGlobalDef(fmt.Sprintf(".str.%d", len(ctx.mod.Globals)), strConst)
	globalStr.Linkage = enum.LinkagePrivate
	globalStr.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalStr.Immutable = true
	globalStr.Align = 1

	strPtr := ctx.block.NewAlloca(types.I8Ptr)
	strGEP := ctx.block.NewGetElementPtr(
		types.NewArray(strConst.Typ.Len, types.I8),
		globalStr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	ctx.block.NewStore(strGEP, strPtr)
	str := ctx.block.NewLoad(types.I8Ptr, strPtr)
	ctx.prog.GlobalStr = append(ctx.prog.GlobalStr, str)
	return str
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

			// define function return type
			var retType types.Type
			if node.IsType(mTypes.TY_INT32) {
				retType = types.I32
			} else if node.IsType(mTypes.TY_STR) {
				retType = types.I8Ptr
			} else if node.IsType(mTypes.TY_NIL) {
				retType = types.I32
			}
			funcName := fmt.Sprintf("fn.%s", node.Val)

			var arg []value.Value
			var argp []*ir.Param

			// define arguments type of function
			for a := node.Child.Args; a != nil; a = a.Next {
				var t types.Type
				if a.IsType(mTypes.TY_INT32) {
					t = types.I32
				} else if a.IsType(mTypes.TY_STR) {
					t = types.I8Ptr
				} else if a.IsType(mTypes.TY_NIL) {
					t = types.I32
				}

				arg = append(arg, ir.NewParam(a.Val, t))
				argp = append(argp, ir.NewParam(a.Val, t))
			}

			fnc := ctx.mod.NewFunc(
				funcName,
				retType,
				argp...,
			)
			llBlock := fnc.NewBlock("")

			ctx.function = fnc
			ctx.argument = node.Child.Args
			ctx.block = llBlock
			res := ctx.gen(node.Child)
			if node.Child.IsKind(mTypes.ND_LAMBDA) {
				r := llBlock.NewCall(res, arg...)
				node.FuncPtr = fnc
				llBlock.NewRet(r)
			} else {
				node.FuncPtr = fnc
				llBlock.NewRet(res)
			}

		}
	} else if node.IsKind(mTypes.ND_VAR_REFERENCE) {
		// PERFORMANCE: too redundant
		// TODO: prohibit same name identify between global var, binded variable and function argument

		// find in variable which is passed as function argument
		for arg := ctx.argument; arg != nil; arg = arg.Next {
			if arg.Val == node.Val {
				var param value.Value
				for i := 0; i < len(ctx.function.Params); i = i + 1 {
					if ctx.function.Params[i].LocalIdent.LocalName == node.Val {
						param = ctx.function.Params[i]
					}
				}
				if arg.IsType(mTypes.TY_INT32) {
					node.Type = mTypes.TY_INT32
					return param

				} else if arg.IsType(mTypes.TY_STR) {
					node.Type = mTypes.TY_STR
					return param
				}
			}
		}

		// find in local variable which is declared with let
		for scope := ctx.scope; scope != nil; scope = scope.Next {
			if scope.Val == node.Val {
				if scope.Child.IsType(mTypes.TY_INT32) {
					node.Type = mTypes.TY_INT32
					return ctx.block.NewLoad(types.I32, scope.VarPtr)

				} else if scope.Child.IsType(mTypes.TY_STR) {
					node.Type = mTypes.TY_STR
					return scope.VarPtr
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
			ctx.function.Sig.RetType,
			ctx.function.Params...,
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
		for bind := node.Bind; bind != nil; bind = bind.Next {
			v := ctx.gen(bind.Child)
			if ctx.scope.Child == nil {
				ctx.scope = bind
			} else {
				ctx.scope.Next = bind
			}

			if bind.Type == mTypes.TY_INT32 {
				bind.VarPtr = ctx.block.NewAlloca(types.I32)
				ctx.block.NewStore(v, bind.VarPtr)

			} else if bind.Type == mTypes.TY_STR {
				bind.VarPtr = v
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

	} else if node.IsKind(mTypes.ND_FUNCCALL) {
		var arg []value.Value
		for node := node.Child; node != nil; node = node.Next {
			arg = append(arg, ctx.gen(node))
		}

		for i := 0; i < len(ctx.mod.Funcs); i = i + 1 {
			if ctx.mod.Funcs[i].GlobalName == fmt.Sprintf("fn.%s", node.Val) {
				return ctx.block.NewCall(ctx.mod.Funcs[i], arg...)
			}

		}
		log.Panic("unresolved function name: have %+v", node)

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

		// cond
		ctx.block = condBlock
		// NOTE: is it the type truely?
		node.CondRet = ctx.block.NewAlloca(ctx.function.Sig.RetType)
		cond := ctx.gen(node.Cond)

		// exit
		// NOTE: is it the type truely?
		exitBlock.NewRet(exitBlock.NewLoad(ctx.function.Sig.RetType, node.CondRet))

		ctx.block = thenBlock
		ctx.block.NewBr(exitBlock)
		res := ctx.gen(node.Then)
		if res != nil {
			ctx.block.NewStore(res, node.CondRet)
		}

		ctx.block = elseBlock
		ctx.block.NewBr(exitBlock)
		res = ctx.gen(node.Else)
		if res != nil {
			ctx.block.NewStore(res, node.CondRet)
		}

		condBlock.NewCondBr(cond, thenBlock, elseBlock)

	} else if node.IsKind(mTypes.ND_SCALAR) {
		if node.IsType(mTypes.TY_INT32) {
			return newI32(node.Val)

		} else if node.IsType(mTypes.TY_STR) {
			return newStr(ctx, node)

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

func (a assembler) GenIntermediates(llName string, asmName string) {
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
