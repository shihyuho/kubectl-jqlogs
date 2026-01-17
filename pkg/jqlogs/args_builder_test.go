package jqlogs

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildJqArgs(t *testing.T) {
	tests := []struct {
		name     string
		jqQuery  string
		opts     JqFlagOptions
		wantArgs []string
	}{
		{
			name:    "Default",
			jqQuery: "",
			opts:    JqFlagOptions{},
			wantArgs: []string{
				"jq", "-R", "-r",
				`. as $line | try (fromjson | (.) | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "Raw output flag",
			jqQuery: ".msg",
			opts:    JqFlagOptions{Raw: true},
			wantArgs: []string{
				"jq", "-R", "-r",
				". as $line | try (fromjson | .msg) catch $line",
			},
		},
		{
			name:    "Without Raw output flag (default)",
			jqQuery: ".msg",
			opts:    JqFlagOptions{Raw: false},
			wantArgs: []string{
				"jq", "-R", "-r",
				". as $line | try (fromjson | (.msg) | if type==\"string\" then tojson else . end) catch $line",
			},
		},
		{
			name:    "Compact",
			jqQuery: ".",
			opts:    JqFlagOptions{Compact: true},
			wantArgs: []string{
				"jq", "-R", "-r", "-c",
				`. as $line | try (fromjson | (.) | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "Color",
			jqQuery: ".level",
			opts:    JqFlagOptions{Color: true},
			wantArgs: []string{
				"jq", "-R", "-r", "-C",
				`. as $line | try (fromjson | (.level) | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "Yaml Output",
			jqQuery: ".msg",
			opts:    JqFlagOptions{Yaml: true},
			wantArgs: []string{
				"jq", "-R", "-r", "--yaml-output",
				`. as $line | try (fromjson | (.msg) | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "Smart Query",
			jqQuery: ".level .msg",
			opts:    JqFlagOptions{},
			wantArgs: []string{
				"jq", "-R", "-r",
				`. as $line | try (fromjson | ("\(.level) \(.msg)") | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "New Flags: Tab and Indent",
			jqQuery: ".",
			opts:    JqFlagOptions{Tab: true, Indent: 4},
			wantArgs: []string{
				"jq", "-R", "-r", "--tab", "--indent", "4",
				`. as $line | try (fromjson | (.) | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "New Flags: Args",
			jqQuery: ".",
			opts:    JqFlagOptions{Args: []string{"env", "prod", "user", "admin"}},
			wantArgs: []string{
				"jq", "-R", "-r", "--arg", "env", "prod", "--arg", "user", "admin",
				`. as $line | try (fromjson | (.) | if type=="string" then tojson else . end) catch $line`,
			},
		},
		{
			name:    "New Flags: JsonArgs",
			jqQuery: ".",
			opts:    JqFlagOptions{JsonArgs: []string{"config", `{"debug":true}`}},
			wantArgs: []string{
				"jq", "-R", "-r", "--argjson", "config", `{"debug":true}`,
				`. as $line | try (fromjson | (.) | if type=="string" then tojson else . end) catch $line`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildJqArgs(tt.jqQuery, tt.opts)
			if !reflect.DeepEqual(got, tt.wantArgs) {
				t.Errorf("BuildJqArgs() = \n%v\nwant \n%v", strings.Join(got, "\n"), strings.Join(tt.wantArgs, "\n"))
			}
		})
	}
}
