package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	files, options := parseArgs(args)

	if options["help"] {
		printUsage(stdout)
		return 0
	}

	if options["version"] {
		fmt.Fprintf(stdout, "Burnfmt v%s\n", getVersion())
		return 0
	}

	burn, err := findBurn()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	fmtArgs := []string{"fmt"}
	switch {
	case options["list"]:
		fmtArgs = append(fmtArgs, "--check")
	case options["write"]:
		fmtArgs = append(fmtArgs, "-w")
	}

	if len(files) == 0 {
		return run(burn, fmtArgs, stdin, stdout, stderr)
	}

	code := 0
	for _, path := range files {
		if !strings.HasSuffix(path, ".bn") {
			fmt.Fprintf(stderr, "Warning: File %s does not have the .bn extension\n", path)
		}
		if c := run(burn, append(append([]string{}, fmtArgs...), path), stdin, stdout, stderr); c != 0 {
			code = c
			continue
		}
		if options["write"] {
			fmt.Fprintf(stdout, "Formatted %s\n", path)
		}
	}
	return code
}

func run(burn string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cmd := exec.Command(burn, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(stderr, "Error running %s: %v\n", burn, err)
		return 1
	}
	return 0
}

func findBurn() (string, error) {
	name := "burn"
	if os.PathSeparator == '\\' {
		name = "burn.exe"
	}
	if p := os.Getenv("BURN_PATH"); p != "" {
		return p, nil
	}
	if exe, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exe), name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	if p, err := exec.LookPath(name); err == nil {
		return p, nil
	}
	return "", errors.New("burnfmt: the `burn` executable was not found; install Burn or set BURN_PATH")
}

func parseArgs(args []string) ([]string, map[string]bool) {
	files := []string{}
	options := make(map[string]bool)

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-h", "--help":
				options["help"] = true
			case "-v", "--version":
				options["version"] = true
			case "-w", "--write":
				options["write"] = true
			case "-l", "--list", "--check":
				options["list"] = true
			}
		} else {
			files = append(files, arg)
		}
	}

	return files, options
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Burnfmt - A code formatter for the Burn programming language")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  burnfmt [options] [file...]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "  -h, --help     Show this help message")
	fmt.Fprintln(w, "  -v, --version  Show version information")
	fmt.Fprintln(w, "  -w, --write    Write result back to source file(s)")
	fmt.Fprintln(w, "  -l, --list     List files that are not formatted (exit code 1 if any)")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  burnfmt file.bn          Show formatted output for file.bn")
	fmt.Fprintln(w, "  burnfmt -w file.bn       Format file.bn and overwrite")
	fmt.Fprintln(w, "  cat file.bn | burnfmt    Format code from standard input")
}

func getVersion() string {
	return "2.0.0"
}
