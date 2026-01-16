package jqlogs

import (
	"bytes"
	"strings"
	"testing"
)

func TestProcessStream(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		query       string
		wantOutput  string
		wantContain string // for flexible matching (e.g. error messages)
	}{
		{
			name:       "Simple JSON",
			input:      `{"key": "value"}`,
			query:      ".",
			wantOutput: "{\n  \"key\": \"value\"\n}\n",
		},
		{
			name:       "Non-JSON",
			input:      "Plain text log",
			query:      ".",
			wantOutput: "Plain text log\n",
		},
		{
			name:       "Mixed Content",
			input:      "Start\n{\"foo\": \"bar\"}\nEnd",
			query:      ".",
			wantOutput: "Start\n{\n  \"foo\": \"bar\"\n}\nEnd\n",
		},
		{
			name:       "JQ Query",
			input:      `{"level": "info", "msg": "hello"}`,
			query:      ".msg",
			wantOutput: "\"hello\"\n",
		},
		{
			name:       "Simple Field Selection",
			input:      `{"level": "info", "app": "myapp"}`,
			query:      ".level .app",
			wantOutput: "\"info myapp\"\n",
		},
		{
			name:       "Valid JQ Priority",
			input:      `{"level": "info", "app": "myapp"}`,
			query:      ".level",
			wantOutput: "\"info\"\n",
		},
		{
			name:       "Complex JQ Priority",
			input:      `{"items": [{"name": "a"}, {"name": "b"}]}`,
			query:      ".items[] | .name",
			wantOutput: "\"a\"\n\"b\"\n",
		},
		{
			name:       "At Prefix Field",
			input:      `{"@timestamp": "time"}`,
			query:      ".@timestamp",
			wantOutput: "\"time\"\n",
		},
		{
			name:        "Invalid JQ Query",
			input:       `{"a": 1}`,
			query:       ".[",
			wantContain: "invalid jq query", // This returns error from NewProcessor/Run but actually ProcessStream returns error immediately for invalid query.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			var output bytes.Buffer

			// Default empty options for these tests
			err := ProcessStream(input, &output, tt.query, JqFlagOptions{})

			if tt.wantContain != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantContain)
				} else if !strings.Contains(err.Error(), tt.wantContain) {
					t.Errorf("expected error containing %q, got %q", tt.wantContain, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := output.String(); got != tt.wantOutput {
				t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, tt.wantOutput)
			}
		})
	}
}

func TestProcessStream_Options(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		query      string
		opts       JqFlagOptions
		wantOutput string
	}{
		{
			name:       "Raw String with Newline",
			input:      `{"msg": "line1\nline2"}`,
			query:      ".msg",
			opts:       JqFlagOptions{Raw: true},
			wantOutput: "line1\nline2\n",
		},
		{
			name:       "Non-Raw String (JSON Quoted)",
			input:      `{"msg": "line1\nline2"}`,
			query:      ".msg",
			opts:       JqFlagOptions{Raw: false},
			wantOutput: "\"line1\\nline2\"\n",
		},
		{
			name:       "Raw Object (Fallback to JSON)",
			input:      `{"a": 1}`,
			query:      ".",
			opts:       JqFlagOptions{Raw: true},
			wantOutput: "{\n  \"a\": 1\n}\n",
		},

		{
			name:       "Yaml Output",
			input:      `{"a": 1, "b": "text"}`,
			query:      ".",
			opts:       JqFlagOptions{Yaml: true},
			wantOutput: "a: 1\nb: text\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(tt.input)
			var output bytes.Buffer

			err := ProcessStream(input, &output, tt.query, tt.opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := output.String(); got != tt.wantOutput {
				t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, tt.wantOutput)
			}
		})
	}
}
