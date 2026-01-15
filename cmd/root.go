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

  # With complex jq query (select and pipe)
  kubectl jqlogs -n my-ns my-pod -- 'select(.level=="error") | .message'`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse arguments using helper
		kubectlArgs, jqQuery, rawOutput, help, version := jqlogs.ParseArgs(args)

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
		// Construct command: kubectl logs [args...]
		// Only if the user didn't provide 'logs' subcommand?
		// The plugin is called as `kubectl jqlogs`. KubeCtl calls the binary `kubectl-jqlogs`.
		// Checks if the first arg is 'logs'. If the user types `kubectl jqlogs logs ...`, we might not want to duplicate.
		// But usually plugins replace a whole verb or add a new one.
		// Here `jqlogs` is the new verb.
		// So we essentially want to run `kubectl logs ...`.

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
		if err := jqlogs.ProcessStream(stdoutPipe, os.Stdout, jqQuery, rawOutput); err != nil {
			// If processing error, maybe just print it?
			// But ProcessStream already prints non-JSON logs.
			// The error return from ProcessStream is usually scanner error.
			fmt.Fprintf(os.Stderr, "Error processing logs: %v\n", err)
		}

		if err := runCmd.Wait(); err != nil {
			// If kubectl logs failed (e.g. pod not found), it exits with non-zero.
			// We should exit with same code if possible.
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
	// Initialize flags here if necessary, but since we are wrapping kubectl logs,
	// we might want DisableFlagParsing: true to forward everything.
	// However, we might need to parse the jq query.
	// Strategy:
	// If the last arg starts with '.', it might be a jq query.
	// But kubectl logs also has many flags.
	// For now, let's keep it simple and refine in the logic implementation phase.
	rootCmd.Flags().BoolVarP(&rawOutput, "raw-output", "r", false, "output raw strings, not JSON texts")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "print the version")
}
