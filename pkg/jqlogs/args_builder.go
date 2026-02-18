package jqlogs

import (
	"fmt"
	"strconv"
)

// jqProgramName is the placeholder for os.Args[0] when invoking gojq/cli.
// gojq/cli uses os.Args[0] as the program name in usage messages.
const jqProgramName = "jq"

// BuildJqArgs constructs the arguments for the underlying jq execution
func BuildJqArgs(jqQuery string, opts JqFlagOptions) []string {
	// Strategy: jq -R -r 'try (fromjson | <query>) catch .'
	// -R: Raw Input (read lines as strings)
	// -r: Raw Output (print fallback strings without quotes)

	// We always use -R and -r by default in the wrapper strategy
	// Pre-allocate capacity: 3 base flags + up to 5 optional flags + 1 query
	args := make([]string, 0, 9)
	args = append(args, jqProgramName, "-R", "-r")

	if opts.Compact {
		args = append(args, "-c")
	}
	if opts.Color {
		args = append(args, "-C")
	}
	if opts.Monochrome {
		args = append(args, "-M")
	}
	if opts.Yaml {
		args = append(args, "--yaml-output")
	}
	if opts.Tab {
		args = append(args, "--tab")
	}
	if opts.Indent > 0 {
		args = append(args, "--indent", strconv.Itoa(opts.Indent))
	}

	// Prepare Query
	if jqQuery == "" {
		jqQuery = "."
	}
	// Apply SmartQuery transformation
	jqQuery = SmartQuery(jqQuery)

	// Wrap Query for Hybrid Mode
	// Note: try/catch in jq passes the *error message* to the catch block, not the original input.
	// So we must bind the input to a variable first: . as $line | try (fromjson | ...) catch $line
	//
	// Handling Raw Output (-r):
	// We globally enable -r to ensure the 'catch $line' part prints raw strings (no quotes) for non-JSON logs.
	// However, for the 'fromjson' part (valid JSON logs), we want to respect the user's choice:
	// - If User specified -r: We don't need to do anything, global -r handles it.
	// - If User did NOT specify -r: Global -r would strip quotes from JSON strings, which is wrong.
	//   So we pipe the result to `if type=="string" then tojson else . end` to re-add quotes for strings.
	jqLogic := jqQuery
	if !opts.Raw {
		jqLogic = fmt.Sprintf("(%s) | if type==\"string\" then tojson else . end", jqQuery)
	}

	wrappedQuery := fmt.Sprintf(". as $line | try (fromjson | %s) catch $line", jqLogic)
	args = append(args, wrappedQuery)

	return args
}
