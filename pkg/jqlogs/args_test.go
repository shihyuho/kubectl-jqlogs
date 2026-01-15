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
		wantRaw         bool
		wantHelp        bool
		wantVersion     bool
	}{
		{
			name:            "Basic Usage",
			args:            []string{"-n", "ns", "pod"},
			wantKubectlArgs: []string{"-n", "ns", "pod"},
			wantJqQuery:     "",
			wantRaw:         false,
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Query",
			args:            []string{"pod", "--", ".level"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     ".level",
			wantRaw:         false,
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Raw Flag",
			args:            []string{"-r", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantRaw:         true,
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Raw Flag Long",
			args:            []string{"--raw-output", "pod"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantRaw:         true,
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "With Help Flag",
			args:            []string{"--help"},
			wantKubectlArgs: []string{},
			wantJqQuery:     "",
			wantRaw:         false,
			wantHelp:        true,
			wantVersion:     false,
		},
		{
			name:            "With Short Help Flag",
			args:            []string{"-h"},
			wantKubectlArgs: []string{},
			wantJqQuery:     "",
			wantRaw:         false,
			wantHelp:        true,
			wantVersion:     false,
		},
		{
			name:            "With Version Flag",
			args:            []string{"--version"},
			wantKubectlArgs: []string{},
			wantJqQuery:     "",
			wantRaw:         false,
			wantHelp:        false,
			wantVersion:     true,
		},
		{
			name:            "With Short Version Flag",
			args:            []string{"-v"},
			wantKubectlArgs: []string{},
			wantJqQuery:     "",
			wantRaw:         false,
			wantHelp:        false,
			wantVersion:     true,
		},
		{
			name:            "Help Mixed",
			args:            []string{"pod", "-h"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "",
			wantRaw:         false,
			wantHelp:        true,
			wantVersion:     false,
		},
		{
			name:            "With Raw Flag and Query",
			args:            []string{"-r", "pod", "--", ".msg"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     ".msg",
			wantRaw:         true,
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "Raw Flag After Dash (Should Ignore)",
			args:            []string{"pod", "--", "-r"},
			wantKubectlArgs: []string{"pod"},
			wantJqQuery:     "-r",
			wantRaw:         false,
			wantHelp:        false,
			wantVersion:     false,
		},
		{
			name:            "Complex Args",
			args:            []string{"-f", "-n", "ns", "-r", "pod", "--", ".a", ".b"},
			wantKubectlArgs: []string{"-f", "-n", "ns", "pod"},
			wantJqQuery:     ".a .b",
			wantRaw:         true,
			wantHelp:        false,
			wantVersion:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKubectlArgs, gotJqQuery, gotRaw, gotHelp, gotVersion := ParseArgs(tt.args)
			if !reflect.DeepEqual(gotKubectlArgs, tt.wantKubectlArgs) {
				// Handle nil vs empty slice comparison
				if len(gotKubectlArgs) == 0 && len(tt.wantKubectlArgs) == 0 {
					// ok
				} else {
					t.Errorf("ParseArgs() kubectlArgs = %v, want %v", gotKubectlArgs, tt.wantKubectlArgs)
				}
			}
			if gotJqQuery != tt.wantJqQuery {
				t.Errorf("ParseArgs() jqQuery = %v, want %v", gotJqQuery, tt.wantJqQuery)
			}
			if gotRaw != tt.wantRaw {
				t.Errorf("ParseArgs() raw = %v, want %v", gotRaw, tt.wantRaw)
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
