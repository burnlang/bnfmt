# Burnfmt

Burnfmt is a code formatter for the Burn programming language, similar to `gofmt` for Go.

## Features

- Formats Burn code (.bn files) with consistent formatting
- Can be used as a command-line tool or library
- Supports writing back to source files or standard output
- Applies opinionated, consistent style rules automatically

## How it works

Burn 2 has a formatter built into the compiler (`burn fmt`). Burnfmt is a small Go command that keeps the
familiar `burnfmt` interface and runs `burn fmt` under the hood, so both always produce identical output.
Burnfmt looks for the `burn` executable in `$BURN_PATH`, next to its own binary, and on `$PATH`.

## Installation

```sh
go install github.com/burnlang/burnfmt@latest
```

or build from source:

```sh
git clone https://github.com/burnlang/bnfmt.git
cd bnfmt
go build -o burnfmt
```

Burn 2 must be installed as well.

## Usage

```sh
burnfmt file.bn          # print formatted code
burnfmt -w file.bn       # format in place
burnfmt -l *.bn          # list files that are not formatted, exit code 1 if any
cat file.bn | burnfmt    # format standard input
```

## Formatting Rules

- 4-space indentation based on block depth
- One space around binary operators, after commas and colons, and before `{`
- No space inside parentheses and brackets, or after unary `-` and `!`
- A line that opens a block and continues is split after the `{`
- A blank line between a closing `}` and the next declaration
- Consecutive blank lines are collapsed into one
- Comments and string contents are left untouched

## Integration

Add burnfmt to your workflow:

- **Git pre-commit hook**: Ensure all committed code follows standard formatting
- **CI/CD pipeline**: Verify formatting as part of automated tests
- **Editor integration**: Configure with VS Code, Vim, or other editors

## License

[MIT License](LICENSE)