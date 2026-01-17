package jqlogs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// JqFlagOptions holds flags specific to the jq processor
type JqFlagOptions struct {
	Raw        bool     // -r / --raw-output
	Compact    bool     // -c / --compact-output
	Color      bool     // -C / --color-output
	Monochrome bool     // -M / --monochrome-output
	Yaml       bool     // --yaml-output
	Tab        bool     // --tab
	Indent     int      // --indent n
	Args       []string // --arg name value
	JsonArgs   []string // --argjson name value
}

// ParseArgs parses the command line arguments
func ParseArgs(args []string) (kubectlArgs []string, jqQuery string, opts JqFlagOptions, help bool, version bool) {
	// Manually scan for flags
	filteredArgs := []string{}
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
		case "--yaml-output":
			opts.Yaml = true
			continue
		case "--tab":
			opts.Tab = true
			continue
		case "--indent":
			if i+1 < len(args) {
				val, err := strconv.Atoi(args[i+1])
				if err == nil {
					opts.Indent = val
					i++ // Consume value
					continue
				} else {
					fmt.Fprintf(os.Stderr, "Warning: invalid argument for --indent: %v\n", args[i+1])
				}
			} else {
				fmt.Fprintf(os.Stderr, "Warning: missing argument for --indent\n")
			}
			// If missing value or invalid, do NOT continue, allow it to be appended to filteredArgs?
			// Actually if it's invalid intended for jqlogs, we might want to consume it to avoid kubectl error,
			// or let it pass to kubectl. Standard behavior: if it looks like our flag, consume it.
			// But if we warned, maybe we should still consume it?
			// If we fail to parse, currently it falls through to 'append(filteredArgs, arg)' which adds '--indent'
			// and then next iteration adds the invalid value. This passes '--indent value' to kubectl.
			// Kubectl will likely fail or complain. This seems acceptable as "we didn't handle it, so maybe kubectl handles it".
			// Review decision: Just warn is sufficient, let fallback happen.
		case "--arg":
			if i+2 < len(args) {
				name := args[i+1]
				val := args[i+2]
				opts.Args = append(opts.Args, name, val)
				i += 2 // Consume name and value
				continue
			}
		case "--argjson":
			if i+2 < len(args) {
				name := args[i+1]
				val := args[i+2]
				opts.JsonArgs = append(opts.JsonArgs, name, val)
				i += 2 // Consume name and value
				continue
			}

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
