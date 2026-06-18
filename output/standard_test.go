package output

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestStandard(t *testing.T) {
	tests := []struct {
		name        string
		input       CheckResults
		expected    []string
		showSkipped bool
	}{
		{
			name: "records failures, warnings and skipped",
			input: CheckResults{
				{
					FileName:  "foo.yaml",
					Namespace: "namespace",
					Warnings:  []Result{{Message: "first warning"}},
					Failures:  []Result{{Message: "first failure"}},
				},
			},
			expected: []string{
				"WARN - foo.yaml - namespace - first warning",
				"FAIL - foo.yaml - namespace - first failure",
				"",
				"2 tests, 0 passed, 1 warning, 1 failure, 0 exceptions",
				"",
			},
		},
		{
			name: "skips filenames for stdin",
			input: CheckResults{
				{
					FileName:  "-",
					Namespace: "namespace",
					Warnings:  []Result{{Message: "first warning"}},
					Failures:  []Result{{Message: "first failure"}},
				},
			},
			expected: []string{
				"WARN - - namespace - first warning",
				"FAIL - - namespace - first failure",
				"",
				"2 tests, 0 passed, 1 warning, 1 failure, 0 exceptions",
				"",
			},
		},
		{
			name: "records failures, warnings and skipped",
			input: CheckResults{
				{
					FileName:  "foo.yaml",
					Namespace: "namespace",
					Warnings:  []Result{{Message: "first warning"}},
					Failures:  []Result{{Message: "first failure"}},
					Skipped:   []Result{{Message: "first skipped"}},
				},
			},
			showSkipped: true,
			expected: []string{
				"WARN - foo.yaml - namespace - first warning",
				"FAIL - foo.yaml - namespace - first failure",
				"",
				"3 tests, 0 passed, 1 warning, 1 failure, 0 exceptions, 1 skipped",
				"",
			},
		},
		{
			name: "skips filenames for stdin",
			input: CheckResults{
				{
					FileName:  "-",
					Namespace: "namespace",
					Warnings:  []Result{{Message: "first warning"}},
					Failures:  []Result{{Message: "first failure"}},
				},
			},
			showSkipped: true,
			expected: []string{
				"WARN - - namespace - first warning",
				"FAIL - - namespace - first failure",
				"",
				"2 tests, 0 passed, 1 warning, 1 failure, 0 exceptions, 0 skipped",
				"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strings.Join(tt.expected, "\n")

			buf := new(bytes.Buffer)
			standard := Standard{Writer: buf, NoColor: true, ShowSkipped: tt.showSkipped}
			if err := standard.Output(tt.input); err != nil {
				t.Fatal("output standard:", err)
			}

			actual := buf.String()

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("Unexpected output. expected %v actual %v", expected, actual)
			}
		})
	}
}

func TestStandardTraceQuiet(t *testing.T) {
	results := CheckResults{
		{
			FileName:  "test.yaml",
			Namespace: "main",
			Queries: []QueryResult{
				{
					Query:  "data.main.passing",
					Traces: []string{"TRACE passing query"},
				},
				{
					Query:   "data.main.deny",
					Results: []Result{{Message: "boom"}},
					Traces:  []string{"TRACE failing query"},
				},
			},
		},
	}

	tests := []struct {
		name      string
		quiet     bool
		wantTrace []string
		skipTrace []string
	}{
		{
			name:      "without quiet traces all queries",
			quiet:     false,
			wantTrace: []string{"TRACE passing query", "TRACE failing query"},
		},
		{
			name:      "with quiet traces only failing queries",
			quiet:     true,
			wantTrace: []string{"TRACE failing query"},
			skipTrace: []string{"TRACE passing query"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			standard := Standard{Writer: buf, NoColor: true, Tracing: true, Quiet: tt.quiet}
			if err := standard.outputTraceOnly(results); err != nil {
				t.Fatal("output trace:", err)
			}

			actual := buf.String()
			for _, want := range tt.wantTrace {
				if !strings.Contains(actual, want) {
					t.Errorf("expected trace to contain %q, got %q", want, actual)
				}
			}
			for _, skip := range tt.skipTrace {
				if strings.Contains(actual, skip) {
					t.Errorf("expected trace to omit %q, got %q", skip, actual)
				}
			}
		})
	}
}
