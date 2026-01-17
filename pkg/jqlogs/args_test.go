package jqlogs

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantKubectlArgs []string
		wantJqQuery     string
		wantOpts        JqFlagOptions
		wantHelp        bool
		wantVersion     bool
	}{
		{
			name:            "Basic Usage",
			args:            []string{"-n", "ns", "pod"},
			wantKubectlArgs: []string{"-n", "ns", "pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Query",
			args:            []string{"pod", "--", ".level"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     ".level",
			wantOpts:        JqFlagOptions{},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Raw Output Flag",
			args:            []string{"-r", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Raw: true},
			wantHelp:        false,
			wantVersion:     false,
		},

		{
			name:            "With Compact Flag",
			args:            []string{"-c", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Compact: true},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Yaml Flag",
			args:            []string{"--yaml-output", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Yaml: true},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Color Flag",
			args:            []string{"-C", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Color: true},
			wantHelp:        false,
			wantVersion:     false,
		},

		{
			name:            "With Help Flag",
			args:            []string{"--help"},
			wantKubectlArgs: []string{},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{},
			wantHelp:        true,
			wantVersion:     false,
		},
		{
			name:            "With Version Flag",
			args:            []string{"--version"},
			wantKubectlArgs: []string{},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{},
			wantHelp:        false,
			wantVersion:     true,
		},
		{
			name:            "Help Mixed",
			args:            []string{"pod", "-h"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{},
			wantHelp:        true,
			wantVersion:     false,
		},
		{
			name:            "With Raw Output Flag and Query",
			args:            []string{"-r", "pod", "--", ".msg"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     ".msg",
			wantOpts:        JqFlagOptions{Raw: true},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "Raw Output Flag After Dash (Should Ignore)",
			args:            []string{"pod", "--", "-r"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "-r",
			wantOpts:        JqFlagOptions{},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Tab Flag",
			args:            []string{"--tab", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Tab: true},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Indent Flag",
			args:            []string{"--indent", "4", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Indent: 4},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Indent Flag Missing Value",
			args:            []string{"--indent"},
			wantKubectlArgs: []string{"--indent"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Arg Flag",
			args:            []string{"--arg", "name", "value", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{Args: []string{"name", "value"}},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With ArgJson Flag",
			args:            []string{"--argjson", "idx", "123", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{JsonArgs: []string{"idx", "123"}},
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "Mixed Flags",
			args:            []string{"-r", "--indent", "4", "--arg", "foo", "bar", "pod", "--", ".msg"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     ".msg",
			wantOpts: JqFlagOptions{
				Raw:    true,
				Indent: 4,
				Args:   []string{"foo", "bar"},
			},
			wantHelp:    false,
			wantVersion: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKubectlArgs, gotJqQuery, gotOpts, gotHelp, gotVersion := ParseArgs(tt.args)
			if (len(gotKubectlArgs) == 0 && len(tt.wantKubectlArgs) == 0) || reflect.DeepEqual(gotKubectlArgs, tt.wantKubectlArgs) {
				// ok, handle empty slice vs nil slice implicitly logic if needed, but DeepEqual handles nil vs empty non-nil strictly
				// Let's rely on DeepEqual but normalize empty slices if needed. Go's append creates non-nil.
				// For the purpose of this test assuming ParseArgs returns empty slice not nil for empty.
			} else {
				t.Errorf("ParseArgs() kubectlArgs = %v, want %v", gotKubectlArgs, tt.wantKubectlArgs)
			}

			// Re-check kubectl args strictly using logic, but here let's focus on logic
			if !equalStringSlices(gotKubectlArgs, tt.wantKubectlArgs) {
				t.Errorf("ParseArgs() kubectlArgs = %v, want %v", gotKubectlArgs, tt.wantKubectlArgs)
			}

			if gotJqQuery != tt.wantJqQuery {
				t.Errorf("ParseArgs() jqQuery = %v, want %v", gotJqQuery, tt.wantJqQuery)
			}
			// Use special helper for opts because slice order in Args/JsonArgs matters
			if !reflect.DeepEqual(gotOpts, tt.wantOpts) {
				t.Errorf("ParseArgs() opts = %+v, want %+v", gotOpts, tt.wantOpts)
			}
			if gotHelp != tt.wantHelp {
				t.Errorf("ParseArgs() help = %v, want %v", gotHelp, tt.wantHelp)
			}
			if gotVersion != tt.wantVersion {
				t.Errorf("ParseArgs() version = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

// Helper to handle nil vs empty slice diffs if any
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
