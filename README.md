# Burnfmt

Burnfmt is a code formatter for the Burn programming language, similar to `gofmt` for Go.

## Features

- Formats Burn code (.bn files) with consistent formatting
- Can be used as a command-line tool or library
- Supports writing back to source files or standard output
- Applies opinionated, consistent style rules automatically

## Installation

```sh
# Build from source
cd burnfmt
go build -o burnfmt

# Install globally (optional)
go install github.com/burnlang/burnfmt@latest
```

## Usage

```sh
# Format a file and print to stdout
burnfmt file.bn

# Format a file and write back to source
burnfmt -w file.bn

# Format multiple files
burnfmt -w file1.bn file2.bn

# Format from stdin
cat file.bn | burnfmt
```

## Formatting Rules

Burnfmt applies the following consistent formatting rules:

- 2-space indentation
- Proper spacing around operators and punctuation
- Consistent brace placement
- Logical grouping of declarations
- Appropriate line breaks for readability

## Integration

Add burnfmt to your workflow:

- **Git pre-commit hook**: Ensure all committed code follows standard formatting
- **CI/CD pipeline**: Verify formatting as part of automated tests
- **Editor integration**: Configure with VS Code, Vim, or other editors

## License

[MIT License](LICENSE)