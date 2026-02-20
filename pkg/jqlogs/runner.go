package jqlogs

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/itchyny/gojq/cli"
)

// Runner manages the execution pipeline
type Runner struct {
	Stdout      io.Writer
	Stderr      io.Writer
	ExecKubectl func(args []string, stdout io.Writer, stderr io.Writer) error
	ExecJq      func(args []string, stdin io.Reader) int
}

// NewDefaultRunner creates a runner with real dependencies
func NewDefaultRunner() *Runner {
	return &Runner{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		ExecKubectl: func(args []string, stdout io.Writer, stderr io.Writer) error {
			cmd := exec.Command("kubectl", append([]string{"logs"}, args...)...)
			cmd.Env = os.Environ()
			cmd.Stdout = stdout
			cmd.Stderr = stderr
			if err := cmd.Start(); err != nil {
				return err
			}
			return cmd.Wait()
		},
		ExecJq: func(args []string, stdin io.Reader) int {
			// Save originals
			oldArgs := os.Args
			oldStdin := os.Stdin
			defer func() {
				os.Args = oldArgs
				os.Stdin = oldStdin
			}()

			os.Args = args

			// We must pass an *os.File to gojq because it checks if it's a terminal.
			// When piping in real life, it is an *os.File.
			if f, ok := stdin.(*os.File); ok {
				os.Stdin = f
			}
			return cli.Run()
		},
	}
}

// Run executes the kubectl -> stream filter -> jq logs pipeline. Returns exit code.
func (r *Runner) Run(kubectlArgs []string, jqQuery string, opts JqFlagOptions) int {
	// Pipe between kubectl and our Scanner
	kPr, kPw, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(r.Stderr, "Error creating pipe: %v\n", err)
		return 1
	}

	// Pipe between our Scanner and JQ
	jqPr, jqPw, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(r.Stderr, "Error creating pipe: %v\n", err)
		return 1
	}

	// 1. Start kubectl asynchronously
	go func() {
		defer kPw.Close()
		err := r.ExecKubectl(kubectlArgs, kPw, r.Stderr)
		if err != nil {
			// Note: If kubectl fails (e.g. pod not found), standard error is already written to r.Stderr.
			// gojq will read EOF and exit normally.
		}
	}()

	// 2. Start Scanner (Filter) asynchronously
	go func() {
		defer jqPw.Close() // closing this tells JQ we are done sending JSON
		scanner := bufio.NewScanner(kPr)

		// To handle very long lines, allow up to 1MB line buffer
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				fmt.Fprintln(r.Stdout)
				continue
			}

			// Pre-filter logic: skip leading whitespaces to find first char
			isJSON := false
			for _, b := range line {
				if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
					continue
				}
				if b == '{' || b == '[' {
					isJSON = true
				}
				break
			}

			if isJSON {
				// Send to JQ pipe
				jqPw.Write(line)
				jqPw.Write([]byte{'\n'})
			} else {
				// Bypass JQ, print directly to Stdout.
				r.Stdout.Write(line)
				r.Stdout.Write([]byte{'\n'})
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(r.Stderr, "Error reading log stream: %v\n", err)
		}
	}()

	// 3. Run JQ synchronously
	jqArgs := BuildJqArgs(jqQuery, opts)
	return r.ExecJq(jqArgs, jqPr)
}
