package jqlogs

import "strings"

// JqFlagOptions holds flags specific to the jq processor
type JqFlagOptions struct {
	Raw        bool     // -r / --raw-output
	Compact    bool     // -c / --compact-output
	Color      bool     // -C / --color-output
	Monochrome bool     // -M / --monochrome-output
	Yaml       bool     // --yaml-output
	SortKeys   bool     // -S / --sort-keys
	Ascii      bool     // -a / --ascii-output
	Unbuffered bool     // --unbuffered
	Seq        bool     // --seq
	Args       []string // Alternating key, value for --arg
	ArgJson    []string // Alternating key, value for --argjson
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
		case "-S", "--sort-keys":
			opts.SortKeys = true
			continue
		case "-a", "--ascii-output":
			opts.Ascii = true
			continue
		case "--unbuffered":
			opts.Unbuffered = true
			continue
		case "--seq":
			opts.Seq = true
			continue
		case "--arg":
			if i+2 < len(args) {
				opts.Args = append(opts.Args, args[i+1], args[i+2])
				i += 2
			}
			continue
		case "--argjson":
			if i+2 < len(args) {
				opts.ArgJson = append(opts.ArgJson, args[i+1], args[i+2])
				i += 2
			}
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
