package jqlogs

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestRunner_Run_Success(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runner := &Runner{
		Stdout: &stdout,
		Stderr: &stderr,
		// Mock ExecKubectl
		ExecKubectl: func(args []string, stdout io.Writer, stderr io.Writer) error {
			stdout.Write([]byte("mock log line 1\n"))
			return nil
		},
		// Mock ExecJq
		ExecJq: func(args []string, stdin io.Reader) int {
			data, _ := io.ReadAll(stdin)
			// In our mock, jq just capitalizes the input to prove it ran
			stdout.Write([]byte(strings.ToUpper(string(data))))
			return 0
		},
	}

	kubectlArgs := []string{"-n", "default", "pod"}
	jqQuery := "."
	opts := JqFlagOptions{}

	// Add a WaitGroup or Channel to ensure the goroutine writing to stdout finishes
	// Since ExecJq is synchronous and returns right away in our mock,
	// the test hits the assertions BEFORE the Runner's scanner goroutine finishes copying to stdout.
	done := make(chan bool)
	go func() {
		exitCode := runner.Run(kubectlArgs, jqQuery, opts)
		if exitCode != 0 {
			t.Errorf("expected exit code 0, got %d", exitCode)
		}
		done <- true
	}()
	<-done

	// We still need a tiny delay because Runner.Run returns ExecJq immediately,
	// but the scanner is running asynchronously. A better way in tests is to wait for the pipe to close.
	// Since our runner does not expose the waitgroup, we'll just sleep briefly for the test.
	time.Sleep(50 * time.Millisecond)

	outStr := stdout.String()
	// Because of filtering, plain text is NO LONGER sent to jq, so it won't be uppercased.
	if !strings.Contains(outStr, "mock log line 1") {
		t.Errorf("expected output to contain 'mock log line 1', got %q", outStr)
	}
}

func TestRunner_Run_Filtering(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// We need an io.Pipe to simulate jq reading stream, but writing to memory to avoid deadlocks
	// in the mock ExecJq returning immediately
	var mockJqOut bytes.Buffer

	runner := &Runner{
		Stdout: &stdout,
		Stderr: &stderr,
		ExecKubectl: func(args []string, out io.Writer, err io.Writer) error {
			out.Write([]byte("plain text log 1\n"))
			out.Write([]byte("{\"level\":\"info\",\"msg\":\"json log 1\"}\n"))
			out.Write([]byte("  [1, 2, 3]\n")) // starts with space then array
			out.Write([]byte("plain text log 2\n"))
			return nil
		},
		ExecJq: func(args []string, stdin io.Reader) int {
			// Read from passed stdin (which should only contain JSON lines now)
			// We copy it to stdout with a prefix to indicate JQ processed it

			// To avoid blocking, we read everything first
			data, _ := io.ReadAll(stdin)

			// Now we write out the simulated JQ output
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			for _, l := range lines {
				if l != "" {
					mockJqOut.Write([]byte("JQ_PROCESSED: " + l + "\n"))
				}
			}

			// Simulate jq output streaming back to stdout
			io.Copy(&stdout, &mockJqOut)
			return 0
		},
	}

	exitCode := runner.Run([]string{}, ".", JqFlagOptions{})
	if exitCode != 0 {
		t.Errorf("expected 0, got %d", exitCode)
	}

	// Wait a tiny bit for the stdout writes from the scanner's plain text to finish flushing
	time.Sleep(50 * time.Millisecond)

	outStr := stdout.String()

	// What we EXPECT is that the Runner sends JSON to JQ and plain text to stdout.
	// Since our Mock JQ prefixes everything IT receives with JQ_PROCESSED,
	// if the Runner is broken and sends EVERYTHING to JQ, the plain text will ALSO
	// have JQ_PROCESSED.

	// Plain text should NOT have JQ_PROCESSED prefix
	if strings.Contains(outStr, "JQ_PROCESSED: plain text log 1") {
		t.Errorf("Plain text log was sent to JQ. Output: %s", outStr)
	}
	if !strings.Contains(outStr, "plain text log 1") { // newline check removed to be safe on cross-platform trims
		t.Errorf("Missing untouched plain text log. Output: %s", outStr)
	}

	// JSON lines SHOULD have JQ_PROCESSED prefix
	if !strings.Contains(outStr, "JQ_PROCESSED: {\"level\":\"info\",\"msg\":\"json log 1\"}") {
		t.Errorf("Missing JQ processed json. Output: %s", outStr)
	}
	// For the array, we expect the leading spaces to be preserved because gojq will receive it as is
	if !strings.Contains(outStr, "JQ_PROCESSED:   [1, 2, 3]") {
		t.Errorf("Missing JQ processed array. Output: %s", outStr)
	}
}
