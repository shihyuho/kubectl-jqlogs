package jqlogs

import (
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

// SmartQuery tries to detect if the query is a simple list of fields like ".level .app_name"
// and transforms it into a string interpolation query like "\(.level) \(.app_name)".
// It also provides syntactic sugar for @-prefixed fields (e.g. .@timestamp -> ."@timestamp").
// It prioritizes simple column selection if the query looks like a list of simple fields.
func SmartQuery(q string) string {
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
			builder.WriteString(`\(`)
			builder.WriteString(part)
			builder.WriteString(`)`)
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
