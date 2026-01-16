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
				`. as $line | try (fromjson | (.)) catch $line`,
			},
		},
		{
			name:    "Compact",
			jqQuery: ".",
			opts:    JqFlagOptions{Compact: true},
			wantArgs: []string{
				"jq", "-R", "-r", "-c",
				`. as $line | try (fromjson | (.)) catch $line`,
			},
		},
		{
			name:    "Color",
			jqQuery: ".level",
			opts:    JqFlagOptions{Color: true},
			wantArgs: []string{
				"jq", "-R", "-r", "-C",
				`. as $line | try (fromjson | (.level)) catch $line`,
			},
		},
		{
			name:    "Yaml Output",
			jqQuery: ".msg",
			opts:    JqFlagOptions{Yaml: true},
			wantArgs: []string{
				"jq", "-R", "-r", "--yaml-output",
				`. as $line | try (fromjson | (.msg)) catch $line`,
			},
		},
		{
			name:    "Smart Query",
			jqQuery: ".level .msg",
			opts:    JqFlagOptions{},
			wantArgs: []string{
				"jq", "-R", "-r",
				`. as $line | try (fromjson | ("\(.level) \(.msg)")) catch $line`,
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
