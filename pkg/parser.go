package makego

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var ErrCircularDependency = errors.New("circular dependency detected")

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

func buildGraph(stages []Stage) map[string][]string {
	graph := make(map[string][]string)
	for _, s := range stages {
		graph[s.Target] = s.Dependencies
	}
	return graph
}

func dependencyResolver(stages []Stage) ([]string, error) {
	graph := buildGraph(stages)
	visited := make(map[string]bool)
	stack := make([]string, 0)
	orderedTargets:= make([]string,0)

	for target := range graph {
		if !visited[target] {
			err := dfs(target, graph, visited, stack, &orderedTargets)
			if err != nil {
				return nil, err
			}
		}
	}
	return orderedTargets, nil
}

func dfs(node string, graph map[string][]string, visited map[string]bool, stack []string, orderedTargets *[]string) error {
	if visited[node] {
		return nil
	}
	if slices.Contains(stack, node) {
		return ErrCircularDependency
	}

	stack = append(stack, node)

	for _, dep := range graph[node] {
		err := dfs(dep, graph, visited, stack, orderedTargets)
		if err != nil {
			return err
		}
	}
	visited[node] = true
	*orderedTargets = append(*orderedTargets, node)
	return nil
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
			dependencies := strings.Fields(dependencies_str)
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

func ExecuteMakefile(path string) error {
	makefile, err := ParseMakefile(path)
	if err != nil {
		return err
	}

	ordered, err := dependencyResolver(makefile.Stages)
	if err != nil {
		return err
	}
	fmt.Println("Execution order:", ordered)

	return nil

}
