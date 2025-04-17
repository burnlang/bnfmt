package cmd

import (
    "fmt"
    "io"
    "os"
    "strings"
    
    "github.com/burnlang/burnfmt/pkg/format"
)

func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
    nonOptions, options := parseArgs(args)
    
    if options["help"] {
        printUsage(stdout)
        return 0
    }
    
    if options["version"] {
        fmt.Fprintf(stdout, "Burnfmt v%s\n", getVersion())
        return 0
    }
    
    
    if len(nonOptions) == 0 {
        return formatStdin(stdin, stdout, stderr, options["write"])
    }
    
    
    for _, path := range nonOptions {
        if err := formatFile(path, stdout, stderr, options["write"]); err != nil {
            fmt.Fprintf(stderr, "Error formatting %s: %v\n", path, err)
            return 1
        }
    }
    
    return 0
}

func formatStdin(stdin io.Reader, stdout, stderr io.Writer, write bool) int {
    
    data, err := io.ReadAll(stdin)
    if err != nil {
        fmt.Fprintf(stderr, "Error reading stdin: %v\n", err)
        return 1
    }
    
    
    formatted, err := format.Format(string(data))
    if err != nil {
        fmt.Fprintf(stderr, "Error formatting: %v\n", err)
        return 1
    }
    
    
    _, err = fmt.Fprint(stdout, formatted)
    if err != nil {
        fmt.Fprintf(stderr, "Error writing to stdout: %v\n", err)
        return 1
    }
    
    return 0
}

func formatFile(path string, stdout, stderr io.Writer, write bool) error {
    
    if !strings.HasSuffix(path, ".bn") {
        fmt.Fprintf(stderr, "Warning: File %s does not have the .bn extension\n", path)
    }
    
    
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("error reading file: %v", err)
    }
    
    
    formatted, err := format.Format(string(data))
    if err != nil {
        return fmt.Errorf("error formatting: %v", err)
    }
    
    
    if write {
        err = os.WriteFile(path, []byte(formatted), 0644)
        if err != nil {
            return fmt.Errorf("error writing to file: %v", err)
        }
        fmt.Fprintf(stdout, "Formatted %s\n", path)
    } else {
        
        _, err = fmt.Fprint(stdout, formatted)
        if err != nil {
            return fmt.Errorf("error writing to stdout: %v", err)
        }
    }
    
    return nil
}


func parseArgs(args []string) ([]string, map[string]bool) {
    nonOptions := []string{}
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
            case "-l", "--list":
                options["list"] = true
            case "-s", "--simplify":
                options["simplify"] = true
            }
        } else {
            nonOptions = append(nonOptions, arg)
        }
    }
    
    return nonOptions, options
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
    fmt.Fprintln(w, "  -l, --list     List files that would be formatted differently")
    fmt.Fprintln(w, "  -s, --simplify Simplify code when formatting")
    fmt.Fprintln(w, "")
    fmt.Fprintln(w, "Examples:")
    fmt.Fprintln(w, "  burnfmt file.bn          Show formatted output for file.bn")
    fmt.Fprintln(w, "  burnfmt -w file.bn       Format file.bn and overwrite")
    fmt.Fprintln(w, "  cat file.bn | burnfmt    Format code from standard input")
}

func getVersion() string {
    return "0.1.0"
}