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

	ordered, err := dependencyResolver(makefile.Stages, makefile.Stages[0].Target)
	if err != nil {
		return err
	}

	defaultTarget := makefile.Stages[0].Target
	runTargets := []string{}

	for _, t := range ordered { //run the first target and its dependencies only
		runTargets = append(runTargets, t)
		if t == defaultTarget {
			break
		}
	}

	for _, target := range runTargets {
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
