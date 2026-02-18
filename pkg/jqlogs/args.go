package jqlogs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// JqFlagOptions holds flags specific to the jq processor
type JqFlagOptions struct {
	Raw        bool // -r / --raw-output
	Compact    bool // -c / --compact-output
	Color      bool // -C / --color-output
	Monochrome bool // -M / --monochrome-output
	Yaml       bool // --yaml-output
	Tab        bool // --tab
	Indent     int  // --indent n
}

// ParseArgs parses the command line arguments
func ParseArgs(args []string) (kubectlArgs []string, jqQuery string, opts JqFlagOptions, help bool, version bool) {
	// Manually scan for flags, separating jqlogs-specific flags from kubectl flags.
	// Note: We perform two passes over args:
	//   Pass 1: Strip jqlogs flags and stop at "--" (appending remainder as-is)
	//   Pass 2: Find "--" in filteredArgs to split kubectlArgs from jqQuery
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			// Stop parsing flags, append the rest
			filteredArgs = append(filteredArgs, args[i:]...)
			break
		}

		// Flag parsing
		switch arg {
		case "-r", "--raw-output":
			opts.Raw = true
			continue
		case "-c", "--compact-output":
			opts.Compact = true
			continue
		case "-C", "--color-output":
			opts.Color = true
			continue
		case "-M", "--monochrome-output":
			opts.Monochrome = true
			continue
		case "-y", "--yaml-output":
			opts.Yaml = true
			continue
		case "--tab":
			opts.Tab = true
			continue
		case "--indent":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --indent requires an argument\n")
				os.Exit(1)
			}
			val, err := strconv.Atoi(args[i+1])
			if err != nil || val < 0 || val > 7 {
				// gojq supports indent values 0-7
				fmt.Fprintf(os.Stderr, "Error: --indent requires an integer between 0 and 7, got: %q\n", args[i+1])
				os.Exit(1)
			}
			opts.Indent = val
			i++ // Consume value
			continue

		case "-h", "--help":
			help = true
			continue
		case "-v", "--version":
			version = true
			continue
		}

		filteredArgs = append(filteredArgs, arg)
	}

	// Find -- separator
	dashIndex := -1
	for i, arg := range filteredArgs {
		if arg == "--" {
			dashIndex = i
			break
		}
	}

	if dashIndex != -1 {
		kubectlArgs = filteredArgs[:dashIndex]
		if dashIndex+1 < len(filteredArgs) {
			jqQuery = strings.Join(filteredArgs[dashIndex+1:], " ")
		}
	} else {
		kubectlArgs = filteredArgs
	}

	return kubectlArgs, jqQuery, opts, help, version
}
