package runner

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestParseFileList(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	files := []string{"main.tf", "provider.tf", "variables.tf"}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(tempDir, file), []byte(""), 0o644); err != nil {
			t.Fatalf("unexpected error writing %s: %v", file, err)
		}
	}

	tests := []struct {
		name        string
		fileList    []string
		ignoreRegex string
		want        []string
	}{
		{
			name:     "no ignore regex returns all explicit files",
			fileList: []string{filepath.Join(tempDir, "main.tf"), filepath.Join(tempDir, "provider.tf")},
			want:     []string{filepath.Join(tempDir, "main.tf"), filepath.Join(tempDir, "provider.tf")},
		},
		{
			name:        "ignore regex excludes matching explicit files",
			fileList:    []string{filepath.Join(tempDir, "main.tf"), filepath.Join(tempDir, "provider.tf"), filepath.Join(tempDir, "variables.tf")},
			ignoreRegex: `.*/provider\.tf`,
			want:        []string{filepath.Join(tempDir, "main.tf"), filepath.Join(tempDir, "variables.tf")},
		},
		{
			name:        "ignore regex excludes files found by walking a directory",
			fileList:    []string{tempDir},
			ignoreRegex: `.*/provider\.tf`,
			want:        []string{filepath.Join(tempDir, "main.tf"), filepath.Join(tempDir, "variables.tf")},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseFileList(tt.fileList, tt.ignoreRegex)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			sort.Strings(got)
			sort.Strings(tt.want)

			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}
