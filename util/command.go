package io

import (
	"bytes"
	"os/exec"
)

func RunCommand(name string, arg ...string) (string, error, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.String(), err, stderr.String()
}
