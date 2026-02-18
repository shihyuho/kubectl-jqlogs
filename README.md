# kubectl-jqlogs

[English](README.md) | [繁體中文](README-zh_TW.md)

**Make JSON logs readable again.**

`kubectl-jqlogs` works exactly like `kubectl logs`, but with a built-in `jq` engine that automatically prettifies JSON output. No pipes, no extra tools—just clean, queryable logs out of the box.

[![License](https://img.shields.io/github/license/shihyuho/kubectl-jqlogs)](LICENSE)
[![Release](https://img.shields.io/github/v/release/shihyuho/kubectl-jqlogs)](https://github.com/shihyuho/kubectl-jqlogs/releases/latest)



## ✨ Features

- **Hybrid Log Processing**: Seamlessly handles mixed content. JSON logs are formatted, while standard text logs are printed as-is. No more "parse error" crashing your pipe.
- **Smart Query Syntax**: Simplified syntax for field selection. Use `.level .msg` instead of complicated string interpolation.
- **Power of [gojq](https://github.com/itchyny/gojq)**: 
  - **Extensions**: Supports `-y`/`--yaml-output` to convert JSON logs to YAML.
  - **Precision**: Handles large numbers and heavy calculations with arbitrary precision (using `math/big`).
  - **Standard Compliance**: Fully implements pure jq language without external C dependencies.
- **Native Experience**: Supports all standard `kubectl logs` flags and auto-completion.
- **Raw Output**: Built-in support for Raw Output (`-r`) and Colorized output.

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

### Supported Flags

`kubectl-jqlogs` supports standard `gojq` flags:

- `-r`, `--raw-output`: Output raw strings, not JSON texts.
- `-c`, `--compact-output`: Compact instead of pretty-printed output.
- `-C`, `--color-output`: Colorize JSON.
- `-M`, `--monochrome-output`: Monochrome (don't colorize JSON).
- `-y`, `--yaml-output`: Output as YAML.
- `--tab`: Use tabs for indentation.
- `--indent n`: Use n spaces for indentation (0-7, default: 2).

#### Examples

**YAML Output:**
```bash
kubectl jqlogs -y -n my-ns my-pod
# level: info
# msg: hello
```

**Simple Field Selection:**

This is a **custom capability** added by `kubectl-jqlogs` to simplify log inspection. Instead of writing verbose `jq` structures (like `{msg: .msg}` or `"\(.msg)"`), you can simply list the fields you want (e.g., `.level .msg`), and the plugin handles the formatting.

> **Note**: This is purely an enhancement. **All standard `jq` syntax is fully supported**, so you can still use complex filters, `select()`, `map()`, and pipes whenever needed.

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

## Shell Alias

To save time, usage of a shell alias is recommended. Replace `kubectl logs` with a shorter command like `klo`:

### Bash / Zsh

```bash
# Add to your .bashrc or .zshrc
alias klo='kubectl jqlogs'
```

Now you can simply run:

```bash
klo -n my-ns my-pod
```

## How it Works

`kubectl-jqlogs` acts as a wrapper around `kubectl logs`. It executes the native command, captures the output stream, and processes each line:

1. If the line is valid JSON, it applies the specified jq query (default is `.`) and pretty-prints the result.
2. If the line is not JSON, it is printed verbatim.
3. If the jq query fails for a JSON line (e.g., the field doesn't exist), the original line is printed as-is.

## License

[MIT](LICENSE)
