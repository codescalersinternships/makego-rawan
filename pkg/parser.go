package makego

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Stage struct {
	Target       string
	Dependencies []string
	Commands     []string
}

type Makefile struct {
	Stages []Stage
}

func readMakefile(path string) (string, error) {
	filename := filepath.Base(path)
	if strings.ToLower(filename) != "makefile" {
		return "", fmt.Errorf("the specified file is not a Makefile")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ParseMakefile(path string) (*Makefile, error) {
	content, err := readMakefile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(content, "\n")
	var stages []Stage
	for i, line := range lines {
		// line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, ":") { // target line
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return nil, fmt.Errorf("syntax error on line %d: %s", i+1, line)
			}

			target := strings.TrimSpace(parts[0])
			dependencies_str := strings.TrimSpace(parts[1])
			dependencies := strings.Split(dependencies_str, " ")
			stages = append(stages, Stage{Target: target, Dependencies: dependencies})

		} else if strings.HasPrefix(line, "\t") { // command line
			if len(stages) == 0 {
				return nil, fmt.Errorf("command without target on line %d: %s", i+1, line)
			}
			command := strings.TrimSpace(line)
			currStage := &stages[len(stages)-1]
			currStage.Commands = append(currStage.Commands, command)
		} else {
			return nil, fmt.Errorf("unrecognized line format on line %d: %s", i+1, line)
		}
	}
	return &Makefile{Stages: stages}, nil
}
