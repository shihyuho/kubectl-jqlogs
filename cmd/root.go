package cmd

import (
	"fmt"
	"os"

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

		runner := jqlogs.NewDefaultRunner()
		exitCode := runner.Run(kubectlArgs, jqQuery, opts)
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
	// DisableFlagParsing is set on rootCmd, so cobra never parses these flags.
	// They are registered here solely to populate the --help output with accurate
	// flag descriptions. Actual flag parsing is done manually in ParseArgs.
	rootCmd.Flags().BoolVarP(&rawOutput, "raw-output", "r", false, "output raw strings, not JSON texts")
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "print the version")

	// Register other supported flags for help visibility
	rootCmd.Flags().BoolP("compact-output", "c", false, "compact instead of pretty-printed output")
	rootCmd.Flags().BoolP("color-output", "C", false, "colorize JSON")
	rootCmd.Flags().BoolP("monochrome-output", "M", false, "monochrome (don't colorize JSON)")
	rootCmd.Flags().BoolP("yaml-output", "y", false, "output as YAML")
	rootCmd.Flags().Bool("tab", false, "use tabs for indentation")
	rootCmd.Flags().Int("indent", 2, "use n spaces for indentation (0-7)")
}
