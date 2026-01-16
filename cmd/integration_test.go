package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/itchyny/gojq/cli"
	"github.com/shihyuho/kubectl-jqlogs/pkg/jqlogs"
)

// captureOutput captures stdout while running the function f
func captureOutput(f func() int) (string, int) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), exitCode
}

func TestJqIntegration(t *testing.T) {
	// Setup Sample Input
	// Note: We are simulating the "cli.Run" part, so we need to provide input via Stdin
	inputLogs := `{"level":"info","msg":"hello"}
{"level":"error","msg":"fail"}
Plain Text Line
{"level":"info","nested":{"foo":"bar"}}`

	tests := []struct {
		name          string
		jqQuery       string
		opts          jqlogs.JqFlagOptions
		wantExitCode  int
		wantOutput    string
		containOutput []string
	}{
		{
			name:         "Default (Hybrid, Pretty Print)",
			jqQuery:      ".",
			opts:         jqlogs.JqFlagOptions{}, // Defaults
			wantExitCode: 0,
			// Expect pretty printed JSON and raw internal failure fallback
			containOutput: []string{
				`"level": "info"`,
				`"msg": "hello"`,
				`Plain Text Line`, // Fallback
			},
		},
		{
			name:         "Compact Output (-c)",
			jqQuery:      ".",
			opts:         jqlogs.JqFlagOptions{Compact: true},
			wantExitCode: 0,
			wantOutput: `{"level":"info","msg":"hello"}
{"level":"error","msg":"fail"}
Plain Text Line
{"level":"info","nested":{"foo":"bar"}}
`,
		},
		{
			name:         "Smart Query",
			jqQuery:      ".level .msg",
			opts:         jqlogs.JqFlagOptions{},
			wantExitCode: 0,
			// Output: "info hello" \n "error fail"
			wantOutput: `info hello
error fail
Plain Text Line
info null
`,
		},
		{
			name:         "Raw Output (-r)",
			jqQuery:      ".msg",
			opts:         jqlogs.JqFlagOptions{Raw: true},
			wantExitCode: 0,
			// Output strings without quotes
			// Non-strings (null) are printed as null
			wantOutput: `hello
fail
Plain Text Line
null
`,
		},
		{
			name:         "YAML Output",
			jqQuery:      ".",
			opts:         jqlogs.JqFlagOptions{Yaml: true},
			wantExitCode: 0,
			containOutput: []string{
				"level: info",
				"msg: hello",
				"nested:",
				"  foo: bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock Args
			args := jqlogs.BuildJqArgs(tt.jqQuery, tt.opts)
			oldArgs := os.Args
			os.Args = args
			defer func() { os.Args = oldArgs }()

			// Mock Stdin
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			w.Write([]byte(inputLogs))
			w.Close()
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()

			// Run
			output, exitCode := captureOutput(func() int {
				return cli.Run()
			})

			// Assert
			if exitCode != tt.wantExitCode {
				t.Errorf("Exit Code = %d, want %d", exitCode, tt.wantExitCode)
			}

			if tt.wantOutput != "" && output != tt.wantOutput {
				t.Errorf("Output =\n%q\nwant\n%q", output, tt.wantOutput)
			}

			for _, sub := range tt.containOutput {
				if !strings.Contains(output, sub) {
					t.Errorf("Output missing substring %q. Got:\n%s", sub, output)
				}
			}
		})
	}
}
