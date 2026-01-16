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

  # With Raw Input (readable stack traces by default)
  kubectl jqlogs -n my-ns my-pod -- .message
  
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
		// Strategy: jq -R -r 'try (fromjson | <query>) catch .'
		// -R: Raw Input (read lines as strings)
		// -r: Raw Output (print fallback strings without quotes)
		jqArgs := []string{"jq", "-R", "-r"}

		// Add flags
		if opts.Compact {
			jqArgs = append(jqArgs, "-c")
		}
		if opts.Color {
			jqArgs = append(jqArgs, "-C")
		}
		if opts.Monochrome {
			jqArgs = append(jqArgs, "-M")
		}
		if opts.Yaml {
			jqArgs = append(jqArgs, "--yaml-output")
		}
		if opts.SortKeys {
			jqArgs = append(jqArgs, "-S")
		}
		if opts.Unbuffered {
			jqArgs = append(jqArgs, "--unbuffered")
		}
		if opts.Seq {
			jqArgs = append(jqArgs, "--seq")
		}

		// Prepare Query
		if jqQuery == "" {
			jqQuery = "."
		}
		// Apply SmartQuery transformation
		jqQuery = jqlogs.SmartQuery(jqQuery)

		// Wrap Query based on RawInput flag
		var wrappedQuery string
		if opts.RawInput {
			// -R specified: Treat input purely as strings (disable hybrid auto-parse)
			wrappedQuery = jqQuery
		} else {
			// Default: Hybrid Mode
			// try (fromjson | query) catch .
			wrappedQuery = fmt.Sprintf("try (fromjson | (%s)) catch .", jqQuery)
		}

		jqArgs = append(jqArgs, wrappedQuery)

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
	rootCmd.Flags().BoolVarP(&rawOutput, "raw-input", "R", false, "read each line as string instead of JSON")
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
