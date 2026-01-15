package jqlogs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/itchyny/gojq"
)

// ProcessStream reads lines from r, tries to parse them as JSON,
// logs them using the provided jq query, and writes output to w.
// If a line is not JSON, it is written to w as is.
func ProcessStream(r io.Reader, w io.Writer, queryString string, raw bool) error {
	if queryString == "" {
		queryString = "."
	}

	queryString = smartQuery(queryString)

	query, err := gojq.Parse(queryString)
	if err != nil {
		return fmt.Errorf("invalid jq query: %w", err)
	}
	code, err := gojq.Compile(query)
	if err != nil {
		return fmt.Errorf("failed to compile jq query: %w", err)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var v interface{}
		if err := json.Unmarshal(line, &v); err != nil {
			// Not JSON, just print the line
			fmt.Fprintln(w, string(line))
			continue
		}

		iter := code.Run(v)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				// If processing fails, fallback to printing original line with error warning?
				// Or print the error? jq behavior is to print error.
				fmt.Fprintf(w, "jq error: %v\n", err)
				continue
			}

			// Format output
			if raw {
				if str, ok := v.(string); ok {
					fmt.Fprintln(w, str)
					continue
				}
			}

			// If v is a string and we want raw output, handle it?
			// For now, always MarshalIndent for readability as per goal.
			output, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				fmt.Fprintf(w, "marshal error: %v\n", err)
				continue
			}
			fmt.Fprintln(w, string(output))
		}
	}
	return scanner.Err()
}

// smartQuery tries to detect if the query is a simple list of fields like ".level .app_name"
// and transforms it into a string interpolation query like "\(.level) \(.app_name)".
// It prioritizes the original query if it is valid jq.
// smartQuery tries to detect if the query is a simple list of fields like ".level .app_name"
// and transforms it into a string interpolation query like "\(.level) \(.app_name)".
// It also provides syntactic sugar for @-prefixed fields (e.g. .@timestamp -> ."@timestamp").
// It prioritizes simple column selection if the query looks like a list of simple fields.
func smartQuery(q string) string {
	parts := strings.Fields(q)
	if len(parts) == 0 {
		return q
	}

	// Pre-process: Fix .@ syntax
	anyFixed := false
	for i, part := range parts {
		if strings.HasPrefix(part, ".@") {
			// Auto-fix: .@field -> ."@field"
			parts[i] = "." + fmt.Sprintf("%q", part[1:])
			anyFixed = true
		}
	}

	// If single part, return it (whether fixed or not)
	// This allows ".@timestamp" -> ".\"@timestamp\"" (valid JQ)
	// And ".level" -> ".level" (valid JQ)
	if len(parts) == 1 {
		if anyFixed {
			return parts[0]
		}
		// If not fixed, prefer original 'q' to preserve whitespace/formatting of complex queries,
		// though for single part fields it doesn't matter much.
		// But allow falling through to "Try parse as is".
		if _, err := gojq.Parse(q); err == nil {
			return q
		}
		// If original invalid (e.g. .@timestamp but somehow missed?), use part?
		// No, just return q.
		return q
	}

	// Heuristic for Simple Mode (multiple parts)
	isSimple := true
	for _, part := range parts {
		if !strings.HasPrefix(part, ".") {
			isSimple = false
			break
		}
		// Check for characters that imply complex JQ logic
		// Note: " is allowed now because we might have introduced it in pre-process
		if strings.ContainsAny(part, `|[](){},`) {
			isSimple = false
			break
		}
	}

	if isSimple {
		var builder strings.Builder
		builder.WriteString(`"`)
		for i, part := range parts {
			if i > 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(`\(` + part + `)`)
		}
		builder.WriteString(`"`)
		transformed := builder.String()

		// verify validity
		if _, err := gojq.Parse(transformed); err == nil {
			return transformed
		}
	}

	// Fallback to original
	return q
}
