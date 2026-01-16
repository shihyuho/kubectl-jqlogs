package jqlogs

import "strings"

// JqFlagOptions holds flags specific to the jq processor
type JqFlagOptions struct {
	RawInput   bool // -R / --raw-input
	Compact    bool // -c / --compact-output
	Color      bool // -C / --color-output
	Monochrome bool // -M / --monochrome-output
	Yaml       bool // --yaml-output
	SortKeys   bool // -S / --sort-keys
	Unbuffered bool // --unbuffered
	Seq        bool // --seq
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
		case "-R", "--raw-input":
			opts.RawInput = true
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
		case "--unbuffered":
			opts.Unbuffered = true
			continue
		case "--seq":
			opts.Seq = true
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
