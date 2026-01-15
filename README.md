# kubectl-jqlogs

**Make your JSON logs readable again.**

`kubectl-jqlogs` works exactly like `kubectl logs`, but with a built-in `jq` engine that automatically prettifies JSON output. No pipes, no extra tools—just clean, queryable logs out of the box.

![License](https://img.shields.io/github/license/shihyuho/kubectl-jqlogs)
![Release](https://img.shields.io/github/v/release/shihyuho/kubectl-jqlogs)

## Features

- **Drop-in Replacement**: Supports standard `kubectl logs` flags (e.g., `-f`, `--tail`, `-n`).
- **Auto-Detection**: Automatically detects if a log line is valid JSON. Non-JSON lines are printed as-is.
- **Built-in JQ**: Uses [gojq](https://github.com/itchyny/gojq), a pure Go implementation of jq. No external dependencies required.
- **Custom Queries**: Filter and transform your logs on the fly using standard jq syntax.
- **Raw Output**: Support `-r` flag to output raw strings, perfect for readable stack traces.

## Installation

### Using Krew (Recommended)

Once the plugin is released, you can install it via [Krew](https://krew.sigs.k8s.io/):

```bash
kubectl krew install jqlogs
```

Or install from the Release manifest directly:

```bash
kubectl krew install --manifest=kubectl-jqlogs.yaml
```

### Using Go

```bash
go install github.com/shihyuho/kubectl-jqlogs@latest
```

### Manual Installation

Download the binary for your platform from the [Releases](https://github.com/shihyuho/kubectl-jqlogs/releases) page and place it in your `$PATH`.

## Usage

The syntax is identical to `kubectl logs`, with an optional jq query at the end.

### Basic Usage

View logs with automatic JSON formatting:

```bash
kubectl jqlogs -n my-namespace my-pod
```

### With jq Query

To use a jq query, separate it with `--` and provide the query.

**Filter by level:**

```bash
kubectl jqlogs -n my-namespace my-pod -- .level
```

### Extended JQ Flags

`kubectl-jqlogs` supports standard standard `gojq` flags:

- `-r`, `--raw-output`: Output raw strings, not JSON texts.
- `-c`, `--compact-output`: Compact output (no pretty print).
- `-C`, `--color-output`: Colorize JSON output.
- `--yaml-output`: Output as YAML.
- `--arg name value`: Set a variable `$name` to the string `value`.
- `--argjson name value`: Set a variable `$name` to the JSON `value`.

#### Examples

**Compact Output:**
```bash
kubectl jqlogs -c -n my-ns my-pod
# {"level":"info","msg":"hello"}
```

**YAML Output:**
```bash
kubectl jqlogs --yaml-output -n my-ns my-pod
# level: info
# msg: hello
```

**Using Variables:**
```bash
kubectl jqlogs -n my-ns my-pod --arg env prod -- 'select(.environment == $env)'
```

**Simple Field Selection:**

You can select multiple fields simply by listing them separated by spaces. The plugin will automatically format them.

```bash
kubectl jqlogs -n my-namespace my-pod -- .level .message
# Output: "info An error occurred"

# For fields starting with '@', you can use them directly:
kubectl jqlogs -n my-namespace my-pod -- .@timestamp .message
# Output: "2026-01-15... An error occurred"
```

**Raw Output (Readable Stack Traces):**

Use `-r` to output raw strings without quotes, which renders newlines (`\n`) correctly.

```bash
kubectl jqlogs -r -n my-namespace my-pod -- .message
# Output:
# Error: ...
#   at com.example...
```

**Select messages directly:**

```bash
kubectl jqlogs -n my-namespace my-pod -- 'select(.level=="error") | .msg'
```

### Streaming Logs

Follow logs with `-f`:

```bash
kubectl jqlogs -f -n my-namespace my-pod
```

## How it Works

`kubectl-jqlogs` acts as a wrapper around `kubectl logs`. It executes the native command, captures the output stream, and processes each line:

1. If the line is valid JSON, it applies the specified jq query (default is `.`) and pretty-prints the result.
2. If the line is not JSON, it is printed verbatim.
3. If the query fails for a line, an error is reported but the stream continues.

## License

MIT
