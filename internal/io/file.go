package io

import (
	"fmt"
	"os/exec"

	"github.com/wf001/modo/internal/log"
)

func PrepareWorkingFile(artifactFilePrefix string, currentTime int64) (string, string, string) {
	if artifactFilePrefix == "" {
		generated := "generated"
		artifactDir := fmt.Sprintf("%s/%d", generated, currentTime)
		out, err := exec.Command("mkdir", "-p", artifactDir).CombinedOutput()
		if err != nil {
			log.Panic("fail to make directory: %+v", map[string]interface{}{"err": err, "out": out, "artifactDir": artifactDir})
		}
		log.Debug(log.YELLOW("make dir: %s"), artifactDir)

		artifactFilePrefix = fmt.Sprintf("%s/out", artifactDir)
	}
	log.Debug(log.YELLOW("artifactFilePrefix = %s"), artifactFilePrefix)
	log.Info("persist all of build artifact in %s", artifactFilePrefix)

	llName := fmt.Sprintf("%s.ll", artifactFilePrefix)
	asmName := fmt.Sprintf("%s.s", artifactFilePrefix)
	executableName := fmt.Sprintf("%s", artifactFilePrefix)

	return llName, asmName, executableName
}
