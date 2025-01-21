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
	"github.com/wf001/modo/pkg/codegen/lib/core"
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
		newI32("0"),
		newI32("0"),
	)
	ctx.block.NewStore(strGEP, strPtr)
	str := ctx.block.NewLoad(types.I8Ptr, strPtr)
	ctx.prog.GlobalStr = append(ctx.prog.GlobalStr, str)
	return str
}

func (ctx *context) genVarDeclare(node *mTypes.Node) value.Value {
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
		retType := node.GetLLVMType()
		funcName := node.GetFuncName()

		var arg []value.Value
		var argp []*ir.Param

		// define arguments type of function
		for a := node.Child.Args; a != nil; a = a.Next {
			t := a.GetLLVMType()

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
		child := ctx.gen(node.Child)
		node.FuncPtr = fnc

		if node.Child.IsKind(mTypes.ND_LAMBDA) {
			lambda := llBlock.NewCall(child, arg...)

			if lambda.Type().Equal(types.Void) {
				llBlock.NewRet(nil)
			} else {
				llBlock.NewRet(lambda)
			}
		} else {
			llBlock.NewRet(child)
		}

	}
	return nil
}

func (ctx *context) genVarReference(node *mTypes.Node) value.Value {
	// PERFORMANCE: too redundant
	// TODO: prohibit same name identifier between global var, binded variable and function argument

	// find in variable which is passed as function argument
	for arg := ctx.argument; arg != nil; arg = arg.Next {
		if arg.Val == node.Val {
			var param value.Value
			for i := 0; i < len(ctx.function.Params); i = i + 1 {
				if ctx.function.Params[i].LocalIdent.LocalName == node.Val {
					param = ctx.function.Params[i]
				}
			}
			node.Type = arg.Type
			return param

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

			} else {
				log.Panic("unresolved NodeType: have %+v", node)
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

	return nil
}

func (ctx *context) genLambda(node *mTypes.Node) value.Value {
	isParentMain := ctx.function.GlobalName == "main"
	unnamedFuncName := node.GetUnnamedFuncName()
	fnEntryBlockName := "fn.entry"

	if isParentMain {
		funcFn := ctx.mod.NewFunc(
			unnamedFuncName,
			types.Void,
			ctx.function.Params...,
		)
		llBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName))

		ctx.function = funcFn
		ctx.block = llBlock

		ctx.gen(node.Child)
		ctx.block.NewRet(nil)

		return funcFn

	} else {
		funcFn := ctx.mod.NewFunc(
			unnamedFuncName,
			ctx.function.Sig.RetType,
			ctx.function.Params...,
		)
		llBlock := funcFn.NewBlock(node.GetBlockName(fnEntryBlockName))

		ctx.function = funcFn
		ctx.block = llBlock

		res := ctx.gen(node.Child)

		if ctx.block.Term == nil {
			ctx.block.NewRet(res)
		}
		return funcFn
	}
}

func (ctx *context) genBranch(
	node *mTypes.Node,
	isVoid bool,
	name string,
	condRet value.Value,
) *ir.Block {

	newBlock := ctx.function.NewBlock(name)
	ctx.block = newBlock
	res := ctx.gen(node)
	if res != nil && !isVoid {
		ctx.block.NewStore(res, condRet)
	}

	return newBlock
}

func (ctx *context) genCondition(node *mTypes.Node) value.Value {
	condBlock := ctx.function.NewBlock(node.GetBlockName("if.cond"))
	ctx.block.NewBr(condBlock)

	exitBlock := node.NextBlock
	isVoid := ctx.function.Sig.RetType.Equal(types.Void)

	// cond
	ctx.block = condBlock
	// NOTE: is it the type truely?
	if !isVoid {
		node.CondRet = ctx.block.NewAlloca(ctx.function.Sig.RetType)
	}
	cond := ctx.gen(node.Cond)

	// exit
	// NOTE: is it the type truely?
	if ctx.function.Sig.RetType.Equal(types.Void) {
		exitBlock.NewRet(nil)
	} else {
		exitBlock.NewRet(exitBlock.NewLoad(ctx.function.Sig.RetType, node.CondRet))
	}

	thenBlock := ctx.genBranch(node.Then, isVoid, node.GetBlockName("if.then"), node.CondRet)
	ctx.block.NewBr(exitBlock)

	elseBlock := ctx.genBranch(node.Else, isVoid, node.GetBlockName("if.else"), node.CondRet)
	ctx.block.NewBr(exitBlock)

	condBlock.NewCondBr(cond, thenBlock, elseBlock)

	return condBlock

}

func (ctx *context) gen(node *mTypes.Node) value.Value {
	log.Debug(log.GREEN(fmt.Sprintf("%+v \"%+v\"", node.Kind, node.Val)))
	if node.IsKind(mTypes.ND_DECLARE) {
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_VAR_DECLARE) {
		ctx.genVarDeclare(node)

	} else if node.IsKind(mTypes.ND_VAR_REFERENCE) {
		return ctx.genVarReference(node)

	} else if node.IsKind(mTypes.ND_LAMBDA) {
		return ctx.genLambda(node)

	} else if node.IsKind(mTypes.ND_BIND) {
		// add node.Bind to last element of ctx.scope
		if ctx.scope == nil {
			ctx.scope = node.Bind
		} else {
			lastScope := ctx.scope.GetLastNode()
			lastScope.Next = node.Bind
		}

		for bind := node.Bind; bind != nil; bind = bind.Next {
			child := ctx.gen(bind.Child)

			if bind.IsType(mTypes.TY_INT32) {
				bind.VarPtr = ctx.block.NewAlloca(types.I32)
				ctx.block.NewStore(child, bind.VarPtr)

			} else if bind.IsType(mTypes.TY_STR) {
				bind.VarPtr = child

			} else {
				log.Panic("unresolved NodeType: have %+v", node)
			}

		}
		return ctx.gen(node.Child)

	} else if node.IsKind(mTypes.ND_EXPR) {
		var res value.Value

		for child := node.Child; child != nil; child = child.Next {
			child.NextBlock = ctx.function.NewBlock(child.GetBlockName("expr"))
			ctx.gen(child)

			ctx.block.NewBr(child.NextBlock)
			ctx.block = child.NextBlock
		}
		return res

	} else if node.IsKind(mTypes.ND_IF) {
		return ctx.genCondition(node)

	} else if node.IsKind(mTypes.ND_LIBCALL) {
		// means calling standard library
		arg := ctx.gen(node.Child)
		node.Child.IRValue = arg

		for n := node.Child.Next; n != nil; n = n.Next {
			arg := ctx.gen(n)
			n.IRValue = arg
		}

		if node.Child.IRValue.Type().Equal(types.Void) {
			n := &mTypes.Node{
				Val: "nil",
			}
			node.Child.IRValue = newStr(ctx, n)
		}

		libFunc := core.LibInsts[node.Val]
		return libFunc(ctx.block, ctx.prog.BuiltinLibs, node.Child)

	} else if node.IsKind(mTypes.ND_FUNCCALL) {
		var arg []value.Value
		// generate ir of their arguments
		for node := node.Child; node != nil; node = node.Next {
			arg = append(arg, ctx.gen(node))
		}

		for i := 0; i < len(ctx.mod.Funcs); i = i + 1 {
			if ctx.mod.Funcs[i].GlobalName == node.GetFuncName() {
				return ctx.block.NewCall(ctx.mod.Funcs[i], arg...)
			}

		}
		log.Panic("unresolved function name: have %+v", node)

	} else if node.IsKind(mTypes.ND_SCALAR) {
		if node.IsType(mTypes.TY_INT32) {
			return newI32(node.Val)

		} else if node.IsType(mTypes.TY_STR) {
			return newStr(ctx, node)

		} else if node.IsType(mTypes.TY_NIL) {
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
			mod:  module,
			prog: prog,
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
