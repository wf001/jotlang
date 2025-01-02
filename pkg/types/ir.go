package types

import "github.com/llir/llvm/ir/types"

// Get LLVM type from corresponding custom type
func (node *Node) GetLLVMType() types.Type {
	var retType types.Type

	if node.IsType(TY_INT32) {
		return types.I32
	} else if node.IsType(TY_STR) {
		return types.I8Ptr
	} else if node.IsType(TY_NIL) {
		return types.I32
	}
	return retType
}
