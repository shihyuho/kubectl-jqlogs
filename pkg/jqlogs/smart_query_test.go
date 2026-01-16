package jqlogs

import (
	"testing"
)

func TestSmartQuery(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Empty",
			input: "",
			want:  "",
		},
		{
			name:  "Dot",
			input: ".",
			want:  ".",
		},
		{
			name:  "Single Field",
			input: ".level",
			want:  ".level",
		},
		{
			name:  "Complex JQ",
			input: ".items[] | .name",
			want:  ".items[] | .name",
		},
		{
			name:  "Simple Field Selection",
			input: ".level .app",
			want:  "\"\\(.level) \\(.app)\"",
		},
		{
			name:  "Three Fields",
			input: ".a .b .c",
			want:  "\"\\(.a) \\(.b) \\(.c)\"",
		},
		{
			name:  "At Prefix Single",
			input: ".@timestamp",
			want:  ".\"@timestamp\"",
		},
		{
			name:  "At Prefix Mixed",
			input: ".level .@timestamp",
			want:  "\"\\(.level) \\(.\"@timestamp\")\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SmartQuery(tt.input); got != tt.want {
				t.Errorf("SmartQuery(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
