# kubectl-jqlogs

**Make your JSON logs readable again.**

`kubectl-jqlogs` works exactly like `kubectl logs`, but with a built-in `jq` engine that automatically prettifies JSON output. No pipes, no extra tools—just clean, queryable logs out of the box.

![License](https://img.shields.io/github/license/shihyuho/kubectl-jqlogs)
![Release](https://img.shields.io/github/v/release/shihyuho/kubectl-jqlogs)



## ✨ Features

- **Hybrid Log Processing**: Seamlessly handles mixed content. JSON logs are formatted, while standard text logs are printed as-is. No more "parse error" crashing your pipe.
- **Smart Query Syntax**: Simplified syntax for field selection. Use `.level .msg` instead of complicated string interpolation.
- **Power of [gojq](https://github.com/itchyny/gojq)**: 
  - **Extensions**: Supports `--yaml-output` to convert JSON logs to YAML.
  - **Precision**: Handles large numbers and heavy calculations with arbitrary precision (using `math/big`).
  - **Standard Compliance**: Fully implements pure jq language without external C dependencies.
- **Native Experience**: Supports all standard `kubectl logs` flags and auto-completion.
- **Raw Input**: Built-in support for Raw Input (`-R`) and Colorized output.

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

- `-R`, `--raw-input`: Read each line as string instead of JSON.
- `-C`, `--color-output`: Colorize JSON output.
- `--yaml-output`: Output as YAML.

#### Examples

**YAML Output:**
```bash
kubectl jqlogs --yaml-output -n my-ns my-pod
# level: info
# msg: hello
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
kubectl jqlogs -n my-namespace my-pod -- .message
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
