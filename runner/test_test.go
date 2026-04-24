package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileList_IgnoreAppliedToExplicitFiles(t *testing.T) {
	// Create a temporary directory with test files.
	dir := t.TempDir()

	keep := filepath.Join(dir, "keep.yaml")
	ignore := filepath.Join(dir, "provider.tf")

	if err := os.WriteFile(keep, []byte("key: value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ignore, []byte(`provider "aws" {}\n`), 0o600); err != nil {
		t.Fatal(err)
	}

	// Pass both files explicitly (not as a directory) with an ignore regex
	// that matches provider.tf. Prior to the fix, the ignore regex was only
	// evaluated when conftest walked a directory, so provider.tf would have
	// been included even though it matched --ignore.
	files, err := parseFileList([]string{keep, ignore}, `.*/provider\.tf`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, f := range files {
		if f == ignore {
			t.Errorf("expected %q to be excluded by the ignore regex, but it was included", ignore)
		}
	}

	found := false
	for _, f := range files {
		if f == keep {
			found = true
		}
	}
	if !found {
		t.Errorf("expected %q to be included, but it was not", keep)
	}
}

func TestParseFileList_IgnoreAppliedToDirectoryFiles(t *testing.T) {
	// Existing behaviour: ignore regex works when conftest walks a directory.
	dir := t.TempDir()

	keep := filepath.Join(dir, "keep.yaml")
	ignore := filepath.Join(dir, "provider.tf")

	if err := os.WriteFile(keep, []byte("key: value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ignore, []byte(`provider "aws" {}\n`), 0o600); err != nil {
		t.Fatal(err)
	}

	files, err := parseFileList([]string{dir}, `.*/provider\.tf`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, f := range files {
		if f == ignore {
			t.Errorf("expected %q to be excluded by the ignore regex, but it was included", ignore)
		}
	}
}

func TestParseFileList_InvalidIgnoreRegex(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "keep.yaml")
	if err := os.WriteFile(f, []byte("key: value\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := parseFileList([]string{f}, `[invalid`)
	if err == nil {
		t.Error("expected an error for an invalid ignore regex, got nil")
	}
}
