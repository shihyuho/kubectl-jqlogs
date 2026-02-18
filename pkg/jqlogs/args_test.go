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
			name:            "With Yaml Flag (short -y)",
			args:            []string{"-y", "pod"},
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
			name:            "Mixed Flags",
			args:            []string{"-r", "--indent", "4", "pod", "--", ".msg"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     ".msg",
			wantOpts: JqFlagOptions{
				Raw:    true,
				Indent: 4,
			},
			wantHelp:    false,
			wantVersion: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKubectlArgs, gotJqQuery, gotOpts, gotHelp, gotVersion := ParseArgs(tt.args)

			assertStringSliceEqual(t, "kubectlArgs", gotKubectlArgs, tt.wantKubectlArgs)

			if gotJqQuery != tt.wantJqQuery {
				t.Errorf("ParseArgs() jqQuery = %v, want %v", gotJqQuery, tt.wantJqQuery)
			}
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

// assertStringSliceEqual is a test helper that compares two string slices,
// treating nil and empty slices as equal.
func assertStringSliceEqual(t *testing.T, field string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("ParseArgs() %s = %v, want %v", field, got, want)
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("ParseArgs() %s = %v, want %v", field, got, want)
			return
		}
	}
}
