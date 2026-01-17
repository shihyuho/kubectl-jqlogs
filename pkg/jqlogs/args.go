package jqlogs

import (
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
				}
			}
			// If missing value or invalid, treat as normal arg or ignore (jq would error)
			// Here we just keep it in kubectl args if parsing fails, or better, fail?
			// For robustness acting as wrapper, let's just ignore opt setting if invalid
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
