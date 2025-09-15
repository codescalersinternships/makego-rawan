package makego

import (
	"reflect"
	"testing"
)

func TestParseMakefile(t *testing.T) {
	testcases := []struct {
		name         string
		path         string
		expectError  error
		expectedFile *Makefile
	}{
		{
			name:         "Invalid Makefile format, more than one colon in target line",
			path:         "../testdata/makefile_invalid1",
			expectError:  ErrInvalidMakefile,
			expectedFile: nil,
		},
		{
			name:         "Invalid Makefile format, command without target",
			path:         "../testdata/makefile_invalid2",
			expectError:  ErrInvalidMakefile,
			expectedFile: nil,
		},
		{
			name:        "Valid Makefile",
			path:        "../testdata/makefile",
			expectError: nil,
			expectedFile: &Makefile{
				Stages: []Stage{
					{Target: "hello", Dependencies: []string{"hello1.txt", "hello2.txt"}, Commands: []string{"cat hello1.txt", "cat hello2.txt"}},
					{Target: "hello1.txt", Dependencies: []string{}, Commands: []string{`echo "Hello 1, Make!" > hello1.txt`}},
					{Target: "hello2.txt", Dependencies: []string{}, Commands: []string{`echo "Hello 2, Make!" > hello2.txt`}},
					{Target: "clean", Dependencies: []string{}, Commands: []string{"rm -f hello1.txt", "rm -f hello2.txt"}},
				},
			}},
	}
	for _, tc := range testcases {
		t.Run(tc.path, func(t *testing.T) {
			makefile, err := parseMakefile(tc.path)
			if err != tc.expectError {
				t.Errorf("expected makefile: %v, got: %v", tc.expectError, err)
			}

			if makefile == nil && tc.expectedFile == nil {
				return
			}

			if !reflect.DeepEqual(makefile, tc.expectedFile) {
				t.Fatalf("expected makefile: %+v, got: %+v", tc.expectedFile, makefile)
			}
		})
	}
}
