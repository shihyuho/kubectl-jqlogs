package jqlogs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/itchyny/gojq"
	"sigs.k8s.io/yaml"
)

// ProcessStream reads lines from r, tries to parse them as JSON,
// logs them using the provided jq query, and writes output to w.
// If a line is not JSON, it is written to w as is.
func ProcessStream(r io.Reader, w io.Writer, queryString string, opts JqFlagOptions) error {
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
			// Not JSON, join custom logic?
			// JQ usually only processes JSON inputs or string if -R (raw input)
			// But here "Non-JSON lines are printed as-is".
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
				fmt.Fprintf(w, "jq error: %v\n", err)
				continue
			}

			// Format output based on options

			// Raw Output (-r)
			if opts.Raw {
				if str, ok := v.(string); ok {
					fmt.Fprintln(w, str)
					continue
				}
			}

			// YAML Output
			if opts.Yaml {
				output, err := yaml.Marshal(v)
				if err != nil {
					fmt.Fprintf(w, "yaml marshal error: %v\n", err)
					continue
				}
				// yaml.Marshal usually adds newline
				fmt.Fprint(w, string(output))
				continue
			}

			// JSON Output
			var output []byte
			if opts.Compact {
				output, err = json.Marshal(v)
			} else {
				output, err = json.MarshalIndent(v, "", "  ")
			}

			if err != nil {
				fmt.Fprintf(w, "marshal error: %v\n", err)
				continue
			}

			// Color Output (-C or auto-detect if we supported it, but now manual -C / -M)
			// If -M (monochrome) is true, disable color (default is no color anyway unless -C).
			// If -C (color) is true, colorize.
			if opts.Color && !opts.Monochrome {
				printColoredJson(w, output)
			} else {
				fmt.Fprintln(w, string(output))
			}
		}
	}
	return scanner.Err()
}

func printColoredJson(w io.Writer, data []byte) {
	colorizeJson(w, string(data))
}

func colorizeJson(w io.Writer, s string) {
	keyColor := color.New(color.FgBlue, color.Bold).SprintFunc()
	stringColor := color.New(color.FgGreen).SprintFunc()
	numberColor := color.New(color.FgCyan).SprintFunc()
	boolColor := color.New(color.FgYellow).SprintFunc()

	// Simple state machine
	var (
		i      int
		length = len(s)
	)

	for i < length {
		char := s[i]

		switch {
		case char == '"':
			// Start of string
			start := i
			i++
			// Find end of string (ignoring escaped quotes)
			for i < length {
				if s[i] == '"' && s[i-1] != '\\' {
					break
				}
				i++
			}
			if i < length {
				i++ // consume closing quote
			}

			// Check if next non-whitespace is colon
			isKey := false
			j := i
			for j < length {
				if s[j] != ' ' && s[j] != '\t' && s[j] != '\n' && s[j] != '\r' {
					if s[j] == ':' {
						isKey = true
					}
					break
				}
				j++
			}

			strVal := s[start:i]
			if isKey {
				fmt.Fprint(w, keyColor(strVal))
			} else {
				fmt.Fprint(w, stringColor(strVal))
			}

		case (char >= '0' && char <= '9') || char == '-':
			// Number
			start := i
			i++
			for i < length {
				c := s[i]
				if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
					i++
				} else {
					break
				}
			}
			fmt.Fprint(w, numberColor(s[start:i]))

		case char == 't' || char == 'f' || char == 'n':
			// boolean or null tries
			start := i
			i++
			for i < length {
				c := s[i]
				if c >= 'a' && c <= 'z' {
					i++
				} else {
					break
				}
			}
			val := s[start:i]
			if val == "true" || val == "false" || val == "null" {
				fmt.Fprint(w, boolColor(val))
			} else {
				fmt.Fprint(w, val)
			}

		default:
			fmt.Fprint(w, string(char))
			i++
		}
	}
	fmt.Fprintln(w) // ensure newline at end
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
