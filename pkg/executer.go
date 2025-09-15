package makego

import (
	"fmt"
	"os"
	"os/exec"
)

func getStageByTarget(stages []Stage, target string) *Stage {
	for i, s := range stages {
		if s.Target == target {
			return &stages[i]
		}
	}
	return nil
}

func ExecuteMakefile(path string, targets []string) error {
	makefile, err := parseMakefile(path)
	if err != nil {
		return err
	}

	defaultTarget := makefile.Stages[0].Target

	if len(targets) == 0 {
		targets = []string{defaultTarget}
	}

	for _, target := range targets {
		ordered, err := dependencyResolver(makefile.Stages, target)
		if err != nil {
			return err
		}
		for _, t := range ordered {
			stage := getStageByTarget(makefile.Stages, t)
			if stage == nil {
				continue
			}
			for _, cmd := range stage.Commands {
				fmt.Print(cmd + "\n")
				cmd := exec.Command("sh", "-c", cmd)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					return err
				}
			}
		}

	}

	return nil

}
