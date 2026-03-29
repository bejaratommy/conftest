package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDataPaths(t *testing.T) {
	t.Run("returns provided paths unchanged when non-empty", func(t *testing.T) {
		paths := []string{"/some/path", "/other/path"}
		got := resolveDataPaths(paths)
		if len(got) != 2 || got[0] != paths[0] || got[1] != paths[1] {
			t.Errorf("expected %v, got %v", paths, got)
		}
	})

	t.Run("returns empty slice when no paths and no data dir exists", func(t *testing.T) {
		// Change to a temp dir that has no "data" subdirectory.
		tmp := t.TempDir()
		orig, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chdir(orig) })
		if err := os.Chdir(tmp); err != nil {
			t.Fatal(err)
		}

		got := resolveDataPaths(nil)
		if len(got) != 0 {
			t.Errorf("expected empty slice, got %v", got)
		}
	})

	t.Run("returns default data dir when it exists and no paths provided", func(t *testing.T) {
		tmp := t.TempDir()
		orig, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chdir(orig) })
		if err := os.Chdir(tmp); err != nil {
			t.Fatal(err)
		}
		// Create the "data" directory.
		if err := os.Mkdir(filepath.Join(tmp, defaultDataDir), 0755); err != nil {
			t.Fatal(err)
		}

		got := resolveDataPaths([]string{})
		if len(got) != 1 || got[0] != defaultDataDir {
			t.Errorf("expected [%q], got %v", defaultDataDir, got)
		}
	})

	t.Run("does not auto-load when data path is a file, not a dir", func(t *testing.T) {
		tmp := t.TempDir()
		orig, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chdir(orig) })
		if err := os.Chdir(tmp); err != nil {
			t.Fatal(err)
		}
		// Create a file named "data" (not a directory).
		if err := os.WriteFile(filepath.Join(tmp, defaultDataDir), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		got := resolveDataPaths([]string{})
		if len(got) != 0 {
			t.Errorf("expected empty slice when 'data' exists as a file, got %v", got)
		}
	})
}
