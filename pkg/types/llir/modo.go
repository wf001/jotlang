package llirTypes

import "github.com/llir/llvm/ir"

type LLIRAlloca struct {
	Alloca *ir.InstAlloca
  Next *LLIRAlloca
}

// TODO: implement AllocaInst interface
