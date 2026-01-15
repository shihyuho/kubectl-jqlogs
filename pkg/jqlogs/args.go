package jqlogs

import "strings"

// ParseArgs parses the command line arguments to extract:
// 1. kubectlArgs: arguments to be passed to kubectl logs
// 2. jqQuery: the jq query string (after --)
// 3. raw: whether raw output is enabled (-r / --raw-output)
// 4. help: whether help is requested (-h / --help)
// 5. version: whether version is requested (-v / --version)
func ParseArgs(args []string) (kubectlArgs []string, jqQuery string, raw bool, help bool, version bool) {
	// Manually scan for flags
	var filteredArgs []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			// Stop parsing flags, append the rest
			filteredArgs = append(filteredArgs, args[i:]...)
			break
		}
		if arg == "-r" || arg == "--raw-output" {
			raw = true
			continue
		}
		if arg == "-h" || arg == "--help" {
			help = true
			continue
		}
		if arg == "-v" || arg == "--version" {
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

	return kubectlArgs, jqQuery, raw, help, version
}
