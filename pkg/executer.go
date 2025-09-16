package makego

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
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
		ordered, err := dependencyResolver(makefile.Stages, target) // get the ordered list of targets to execute
		if err != nil {
			return err
		}

		done := make(map[string]bool)

		for i := 0; i < len(ordered); { 
			var wg sync.WaitGroup // wait group to wait for all goroutines to finish
			ch := make(chan error, len(ordered)) // channel to collect errors from goroutines
			level := []string{}
			for j := i; j < len(ordered); j++ { 
				stage := getStageByTarget(makefile.Stages, ordered[j])
				if stage == nil {
					continue
				}
				ready := true
				for _, dep := range stage.Dependencies { // check if all dependencies of the current target are done
					if !done[dep] {
						ready = false
						break
					}
				}

				if ready {
					level = append(level, ordered[j]) 
				}

			}
			for _, t := range level {
				wg.Add(1) // increment wait group counter (means wait for one more goroutine)
				go func(target string) {
					defer wg.Done()
					stage := getStageByTarget(makefile.Stages, target)
					err := runStage(stage)
					if err != nil {
						ch <- err
						return
					}
					done[target] = true

				}(t)
			}
			wg.Wait()
			close(ch)
			if len(ch) > 0 {
				return <-ch
			}
			i += len(level)
		}
	}
	return nil
}

func runStage(stage *Stage) error {
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
	return nil
}
