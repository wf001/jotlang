package types

import "github.com/llir/llvm/ir/types"

// Get LLVM type from corresponding custom type
func (node *Node) GetLLVMType() types.Type {
	var retType types.Type
	if node.IsType(TY_INT32) {
		retType = types.I32
	} else if node.IsType(TY_STR) {
		retType = types.I8Ptr
	} else if node.IsType(TY_NIL) {
		retType = types.I32
	}
	return retType
}
