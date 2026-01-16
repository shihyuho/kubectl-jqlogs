package jqlogs

import "fmt"

// BuildJqArgs constructs the arguments for the underlying jq execution
func BuildJqArgs(jqQuery string, opts JqFlagOptions) []string {
	// Strategy: jq -R -r 'try (fromjson | <query>) catch .'
	// -R: Raw Input (read lines as strings)
	// -r: Raw Output (print fallback strings without quotes)

	// We always use -R and -r by default in the wrapper strategy
	args := []string{"jq", "-R", "-r"}

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
	if opts.SortKeys {
		args = append(args, "-S")
	}
	if opts.Unbuffered {
		args = append(args, "--unbuffered")
	}
	if opts.Seq {
		args = append(args, "--seq")
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
	wrappedQuery := fmt.Sprintf(". as $line | try (fromjson | (%s)) catch $line", jqQuery)
	args = append(args, wrappedQuery)

	return args
}
