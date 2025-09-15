package makego

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getStageByTarget(stages []Stage, target string) *Stage {
	for i, s := range stages {
		if s.Target == target {
			return &stages[i]
		}
	}
	return nil
}

func ExecuteMakefile(path string) error {
	makefile, err := parseMakefile(path)
	if err != nil {
		return err
	}

	ordered, err := dependencyResolver(makefile.Stages)
	if err != nil {
		return err
	}

	for _, target := range ordered {
		stage := getStageByTarget(makefile.Stages, target)
		if stage == nil {
			continue
		}
		for _, cmd := range stage.Commands {
			printCmd := true
			if strings.HasPrefix(strings.TrimSpace(cmd), "echo") &&
				(strings.Contains(cmd, ">") || strings.Contains(cmd, ">>")) {
				printCmd = false
			}

			if printCmd {
				fmt.Println(cmd)
			}
			cmd := exec.Command("sh", "-c", cmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil

}
