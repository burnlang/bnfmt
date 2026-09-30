<p align="center">
    <img src="https://raw.githubusercontent.com/burnlang/burn/master/assets/logo.svg" alt="Burn logo" width="128">
</p>

# Burnfmt (Go, deprecated)

> [!WARNING]
> This Go build of `burnfmt` is deprecated and no longer maintained. Use the `burnfmt` that ships with Burn.

The Burn toolchain installs `burnfmt`, a formatter written in Burn itself, next to `burn`, `burni` and `burnc`:

```sh
curl -fsSL https://raw.githubusercontent.com/burnlang/burn/master/install.sh | sh
```

It always matches the installed Burn version, so there is nothing extra to install or keep in sync.

## Moving over

1. Remove the Go build: `rm "$(go env GOPATH)/bin/burnfmt"`
2. Install Burn with the command above, which puts `burnfmt` in `~/.burn/bin`
3. Keep using the same commands:

```sh
burnfmt file.bn              # print the formatted file
burnfmt -w file.bn           # format in place
burnfmt --check src/*.bn     # list files that are not formatted, exit code 1 if any
cat file.bn | burnfmt        # format standard input
```

It accepts the same flags as the Go build (`-w`, `-l`/`--check`, `-h`, `-v`), so scripts and editor setups keep
working. `burn fmt` does the same without a separate command.

See the [burnfmt documentation](https://github.com/burnlang/burn/blob/master/docs/tooling/burnfmt.mdx) for details.

## Existing installs

The Go build still works: it runs `burn fmt` under the hood and prints a deprecation notice to standard error.
It gets no new features or fixes.

## License

[MIT License](LICENSE)
