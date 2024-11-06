package main

import (
	"testing"
)

func TestPrepareWarkingFile(t *testing.T) {
	llName, asmName, executableName := prepareWorkingFile("out")
	if llName != "out.ll" {
		t.Errorf("have = %s, want = %s", "out.ll", llName)
	}
	if asmName != "out.s" {
		t.Errorf("have = %s, want = %s", "out.s", asmName)
	}
	if executableName != "out" {
		t.Errorf("have = %s, want = %s", "out", executableName)
	}
}
