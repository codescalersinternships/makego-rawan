package makego

import (
	"errors"
	"os"
	"strings"
)

var ErrCircularDependency = errors.New("circular dependency detected")
var ErrInvalidMakefile = errors.New("invalid makefile format")

type Stage struct {
	Target       string
	Dependencies []string
	Commands     []string
}

type Makefile struct {
	Stages []Stage
}

func readMakefile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func parseMakefile(path string) (*Makefile, error) {
	content, err := readMakefile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(content, "\n")
	var stages []Stage
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, ":") { // target line
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return nil, ErrInvalidMakefile
			}

			target := strings.TrimSpace(parts[0])
			dependencies_str := strings.TrimSpace(parts[1])
			dependencies := strings.Fields(dependencies_str)
			stages = append(stages, Stage{Target: target, Dependencies: dependencies})

		} else if strings.HasPrefix(line, "\t") { // command line
			if len(stages) == 0 {
				return nil, ErrInvalidMakefile
			}
			command := strings.TrimSpace(line)
			currStage := &stages[len(stages)-1]
			currStage.Commands = append(currStage.Commands, command)
		} else {
			return nil, ErrInvalidMakefile
		}
	}
	return &Makefile{Stages: stages}, nil
}
