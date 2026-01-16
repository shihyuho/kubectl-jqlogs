package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/itchyny/gojq/cli"
	"github.com/shihyuho/kubectl-jqlogs/pkg/jqlogs"
	"github.com/spf13/cobra"
)

// data structure to hold flags or configuration
var (
	rawOutput   bool
	versionFlag bool
	// Version is injected by build flags
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-jqlogs",
	Short: "Readable, colorful JSON logs via jq",
	Long: `A wrapper for 'kubectl logs' with a built-in jq engine (gojq).
It features Hybrid Log Processing (handling standard and JSON logs seamlessly),
Smart Query syntax, and extends jq with YAML output and arbitrary precision math.`,
	Example: `  # Supports all standard kubectl logs flags (e.g. follow, tail)
  kubectl jqlogs -f --tail=50 -n my-ns my-pod

  # Basic usage (auto-format JSON)
  kubectl jqlogs -n my-ns my-pod

  # With simple field selection
  kubectl jqlogs -n my-ns my-pod -- .level .message

  # With Raw Output (readable stack traces)
  kubectl jqlogs -r -n my-ns my-pod -- .message
  
  # Output as YAML
  kubectl jqlogs --yaml-output -n my-ns my-pod

  # With complex jq query (select and pipe)
  kubectl jqlogs -n my-ns my-pod -- 'select(.level=="error") | .message'`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse arguments using helper
		kubectlArgs, jqQuery, opts, help, version := jqlogs.ParseArgs(args)

		if help {
			cmd.Help()
			os.Exit(0)
		}

		if version {
			fmt.Println(Version)
			os.Exit(0)
		}

		// Prepare Pipe for kubectl logs -> gojq
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating pipe: %v\n", err)
			os.Exit(1)
		}

		// Run kubectl logs
		runCmd := exec.Command("kubectl", append([]string{"logs"}, kubectlArgs...)...)
		runCmd.Env = os.Environ()
		runCmd.Stderr = os.Stderr
		runCmd.Stdout = w // Write to pipe

		if err := runCmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting kubectl logs: %v\n", err)
			os.Exit(1)
		}

		// Close write end after command finishes ensures cli.Run receives EOF
		go func() {
			runCmd.Wait()
			w.Close()
		}()

		// Construct gojq arguments
		jqArgs := jqlogs.BuildJqArgs(jqQuery, opts)

		// Delegate to gojq/cli
		// Override os.Args and os.Stdin
		os.Args = jqArgs
		os.Stdin = r // Read from pipe

		exitCode := cli.Run()
		os.Exit(exitCode)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize flags here primarily for Help documentation.
	// Actual parsing is done manually in ParseArgs to support DisableFlagParsing.
	rootCmd.Flags().BoolVarP(&rawOutput, "raw-output", "r", false, "output raw strings, not JSON texts")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "print the version")

	// Register other supported flags for help visibility
	rootCmd.Flags().BoolP("compact-output", "c", false, "compact instead of pretty-printed output")
	rootCmd.Flags().BoolP("color-output", "C", false, "colorize JSON")
	rootCmd.Flags().BoolP("monochrome-output", "M", false, "monochrome (don't colorize JSON)")
	rootCmd.Flags().Bool("yaml-output", false, "output as YAML")
	rootCmd.Flags().BoolP("sort-keys", "S", false, "sort keys of objects on output")
	rootCmd.Flags().Bool("unbuffered", false, "flush output stream after each JSON object")
	rootCmd.Flags().Bool("seq", false, "use the RS/LF for input/output separators")
}
