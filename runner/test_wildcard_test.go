package runner

import (
	"reflect"
	"sort"
	"testing"
)

func TestExpandWildcardNamespaces(t *testing.T) {
	available := []string{
		"main",
		"k8s.simple.deployment",
		"k8s.simple.hpa",
		"k8s.combined.deployment",
		"k8s.combined.hpa",
		"terraform.simple",
		"terraform.combined",
	}

	tests := []struct {
		name      string
		patterns  []string
		available []string
		want      []string
	}{
		{
			name:      "no wildcards returns patterns as-is",
			patterns:  []string{"main", "terraform.simple"},
			available: available,
			want:      []string{"main", "terraform.simple"},
		},
		{
			name:      "star matches all",
			patterns:  []string{"*"},
			available: available,
			want:      available,
		},
		{
			name:      "prefix wildcard",
			patterns:  []string{"k8s.simple.*"},
			available: available,
			want:      []string{"k8s.simple.deployment", "k8s.simple.hpa"},
		},
		{
			name:      "suffix wildcard",
			patterns:  []string{"*.combined"},
			available: available,
			want:      []string{"terraform.combined"},
		},
		{
			name:      "mixed wildcard and literal",
			patterns:  []string{"main", "k8s.combined.*"},
			available: available,
			want:      []string{"main", "k8s.combined.deployment", "k8s.combined.hpa"},
		},
		{
			name:      "wildcard with no matches returns nothing for that pattern",
			patterns:  []string{"nonexistent.*"},
			available: available,
			want:      nil,
		},
		{
			name:      "duplicate patterns deduplicated",
			patterns:  []string{"main", "main"},
			available: available,
			want:      []string{"main"},
		},
		{
			name:      "wildcard deduplicates across patterns",
			patterns:  []string{"k8s.*.*", "k8s.simple.*"},
			available: available,
			want: []string{
				"k8s.simple.deployment",
				"k8s.simple.hpa",
				"k8s.combined.deployment",
				"k8s.combined.hpa",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandWildcardNamespaces(tt.patterns, tt.available)
			// Sort both slices for stable comparison when order doesn't matter.
			// The last sub-test asserts order implicitly via sorted comparison.
			gotSorted := make([]string, len(got))
			copy(gotSorted, got)
			sort.Strings(gotSorted)

			wantSorted := make([]string, len(tt.want))
			copy(wantSorted, tt.want)
			sort.Strings(wantSorted)

			if !reflect.DeepEqual(gotSorted, wantSorted) {
				t.Errorf("expandWildcardNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
