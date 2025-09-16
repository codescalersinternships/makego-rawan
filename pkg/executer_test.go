package makego

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeMakefile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "Makefile")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp Makefile: %v", err)
	}

	return path
}

func TestExecuteMakefileIntegration(t *testing.T) {
	t.Run("default target with dependencies only", func(t *testing.T) {
		makefileContent := `
hello: hello1 hello2
	cat hello1.txt
	cat hello2.txt

hello1:
	echo "Hello 1, Make!" > hello1.txt

hello2:
	echo "Hello 2, Make!" > hello2.txt

clean:
	rm -f hello1.txt hello2.txt
`
		path := writeMakefile(t, makefileContent)

		err := ExecuteMakefile(path, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		data, err := os.ReadFile("hello1.txt")
		if err != nil {
			t.Fatalf("expected hello1.txt to exist %v", err)
		}
		if !strings.Contains(string(data), "Hello 1, Make!") {
			t.Errorf("unexpected hello1.txt content, got %q", string(data))
		}

		data, err = os.ReadFile("hello2.txt")
		if err != nil {
			t.Fatalf("expected hello2.txt to exist %v", err)
		}
		if !strings.Contains(string(data), "Hello 2, Make!") {
			t.Errorf("unexpected hello2.txt content, got %q", string(data))
		}
	})

	t.Run("run target clean", func(t *testing.T) {
		makefileContent := `
clean:
	rm -f dummy.txt
`
		path := writeMakefile(t, makefileContent)

		if err := os.WriteFile("dummy.txt", []byte{}, 0644); err != nil {
			t.Fatalf("failed to write dummy file: %v", err)
		}

		err := ExecuteMakefile(path, []string{"clean"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = os.Stat("dummy.txt")
		if !os.IsNotExist(err) {
			t.Errorf("clean target didn't remove dummy file")
		}
	})

	t.Run("fails on invalid command", func(t *testing.T) {
		makefileContent := `
invalid:
	invalidcommand
`
		path := writeMakefile(t, makefileContent)

		err := ExecuteMakefile(path, []string{"invalid"})
		if err == nil {
			t.Fatalf("expected error for invalid command, got nil")
		}
	})
}
