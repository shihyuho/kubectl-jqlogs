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
			name:            "With Supported but Ignored Flags",
			args:            []string{"-S", "--unbuffered", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantOpts:        JqFlagOptions{SortKeys: true, Unbuffered: true},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKubectlArgs, gotJqQuery, gotOpts, gotHelp, gotVersion := ParseArgs(tt.args)
			if !reflect.DeepEqual(gotKubectlArgs, tt.wantKubectlArgs) {
				t.Errorf("ParseArgs() kubectlArgs = %v, want %v", gotKubectlArgs, tt.wantKubectlArgs)
			}
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
