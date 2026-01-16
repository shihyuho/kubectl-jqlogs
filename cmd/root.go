package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
	Short: "A kubectl plugin to format json logs",
	Long: `A kubectl plugin that behaves like kubectl logs but formats JSON logs using jq.
It behaves like a filter: Non-JSON lines are printed as-is, while JSON lines are formatted.
All standard kubectl logs flags (e.g., -f, --tail, -p) are supported and passed through to kubectl.`,
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
			fmt.Printf("kubectl-jqlogs version %s\n", Version)
			os.Exit(0)
		}

		// Run kubectl logs
		// We expect 'kubectl' to be in the PATH.
		runCmd := exec.Command("kubectl", append([]string{"logs"}, kubectlArgs...)...)

		// Inherit env?
		runCmd.Env = os.Environ()
		runCmd.Stderr = os.Stderr

		stdoutPipe, err := runCmd.StdoutPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating stdout pipe: %v\n", err)
			os.Exit(1)
		}

		if err := runCmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting kubectl logs: %v\n", err)
			os.Exit(1)
		}

		// Process stream
		if err := jqlogs.ProcessStream(stdoutPipe, os.Stdout, jqQuery, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing logs: %v\n", err)
		}

		if err := runCmd.Wait(); err != nil {
			// If kubectl logs failed (e.g. pod not found), it exits with non-zero.
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}
			os.Exit(1)
		}
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
	rootCmd.Flags().BoolP("color-output", "C", false, "colorize JSON")
	rootCmd.Flags().BoolP("monochrome-output", "M", false, "monochrome (don't colorize JSON)")
	rootCmd.Flags().Bool("yaml-output", false, "output as YAML")
	rootCmd.Flags().BoolP("sort-keys", "S", false, "sort keys of objects on output")
	rootCmd.Flags().BoolP("ascii-output", "a", false, "output ASCII with escaped characters")
	rootCmd.Flags().Bool("unbuffered", false, "flush output stream after each JSON object")
	rootCmd.Flags().Bool("seq", false, "use the RS/LF for input/output separators")
}
