package makego

import "slices"

func buildGraph(stages []Stage) map[string][]string {
	graph := make(map[string][]string)
	for _, s := range stages {
		graph[s.Target] = s.Dependencies
	}
	return graph
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

func dependencyResolver(stages []Stage) ([]string, error) {
	graph := buildGraph(stages)
	visited := make(map[string]bool)
	stack := make([]string, 0)
	orderedTargets := make([]string, 0)

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
