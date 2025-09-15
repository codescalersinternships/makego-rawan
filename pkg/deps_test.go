package makego

import (
	"reflect"
	"testing"
)

func TestDependencyResolver(t *testing.T) {
	testcases := []struct {
		name        string
		stages      []Stage
		expectError error
		expected    []string
	}{
		{
			name: "Simple linear dependencies",
			stages: []Stage{
				{Target: "A", Dependencies: []string{"B"}, Commands: []string{"echo A"}},
				{Target: "B", Dependencies: []string{"C"}, Commands: []string{"echo B"}},
				{Target: "C", Dependencies: []string{}, Commands: []string{"echo C"}},
			},
			expectError: nil,
			expected:    []string{"C", "B", "A"},
		},
		{
			name: "Circular dependency",
			stages: []Stage{
				{Target: "A", Dependencies: []string{"B"}, Commands: []string{"echo A"}},
				{Target: "B", Dependencies: []string{"A"}, Commands: []string{"echo B"}},
			},
			expectError: ErrCircularDependency,
			expected:    nil,
		},
		{
			name: "Indirect Circular dependency",
			stages: []Stage{
				{Target: "A", Dependencies: []string{"B"}, Commands: []string{"echo A"}},
				{Target: "B", Dependencies: []string{"C"}, Commands: []string{"echo C"}},
				{Target: "C", Dependencies: []string{"D"}, Commands: []string{"echo D"}},
				{Target: "D", Dependencies: []string{"A"}, Commands: []string{"echo A"}},
			},
			expectError: ErrCircularDependency,
			expected:    nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ordered, err := dependencyResolver(tc.stages)
			if err != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}

			if !reflect.DeepEqual(ordered, tc.expected) {
				t.Fatalf("expected order: %v, got: %v", tc.expected, ordered)
			}
		})
	}
}
