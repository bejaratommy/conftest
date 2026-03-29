package runner

import "os"

// defaultDataDir is the directory automatically loaded as data when no --data flag is provided.
const defaultDataDir = "data"

// resolveDataPaths returns the given paths unchanged when they are non-empty.
// When paths is empty, it checks whether a "data" directory exists in the current
// working directory and, if so, returns it as the sole element.  This mirrors the
// behaviour of the --policy flag, which defaults to "policy".
func resolveDataPaths(paths []string) []string {
	if len(paths) > 0 {
		return paths
	}
	if info, err := os.Stat(defaultDataDir); err == nil && info.IsDir() {
		return []string{defaultDataDir}
	}
	return paths
}
