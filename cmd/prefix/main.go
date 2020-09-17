package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"moul.io/prefix"
)

func main() {
	if err := run(os.Args); err != nil {
		if err != flag.ErrHelp {
			log.Fatalf("error: %v", err)
		}
		os.Exit(1)
	}
}

var (
	flags  = flag.NewFlagSet("prefix", flag.ContinueOnError)
	format = flags.String("format", "{{DEFAULT}} ", "format string")
)

func run(args []string) error {
	flags.Usage = usage

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// parse args
	var input io.Reader
	{
		remainingArgs := flags.Args()
		switch {
		case len(remainingArgs) == 0:
			input = os.Stdin
		case len(remainingArgs) == 1 && remainingArgs[0] == "-":
			input = os.Stdin
		case len(remainingArgs) == 1:
			f, err := os.Open(remainingArgs[0])
			if err != nil {
				return err
			}
			defer f.Close()
			input = f
		case len(remainingArgs) > 2:
			return fmt.Errorf("usage: prefix FILE")
		}
	}

	// configure prefixer
	var prefixer prefix.LinePrefixer
	{
		prefixer = prefix.New(*format)
	}

	// stream and prefix input
	{
		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(prefixer.PrefixLine(line))
		}
	}
	return nil
}

func usage() {
	// usage
	{
		fmt.Fprintln(os.Stderr, `USAGE`)
		fmt.Fprintln(os.Stderr, `  prefix [flags] file`)
		fmt.Fprintln(os.Stderr)
	}

	// flags
	{
		fmt.Fprintln(os.Stderr, `FLAGS`)
		flags.PrintDefaults()
		fmt.Fprintln(os.Stderr)
	}

	// syntax
	{
		fmt.Fprintln(os.Stderr, `SYNTAX`)
		keys := []string{}
		for key := range prefix.AvailablePatterns {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintf(os.Stderr, "  %-35s %s\n", key, prefix.AvailablePatterns[key])
		}
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "  the following helpers are also available:")
		fmt.Fprintln(os.Stderr, "  - from the text/template library    https://golang.org/pkg/text/template/")
		fmt.Fprintln(os.Stderr, "  - from the sprig project            https://github.com/masterminds/sprig#usage")
		fmt.Fprintln(os.Stderr)
	}

	// presets
	{
		fmt.Fprintln(os.Stderr, `PRESETS`)
		keys := []string{}
		for key := range prefix.AvailablePresets {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintf(os.Stderr, "  %-20s %s\n", key, prefix.AvailablePresets[key])
		}
		fmt.Fprintln(os.Stderr)
	}

	// examples
	{
		fmt.Fprintln(os.Stderr, `EXAMPLES`)
		fmt.Fprintln(os.Stderr, `  prefix apache.log`)
		fmt.Fprintln(os.Stderr, `  prefix -format=">>>" apache.log`)
		fmt.Fprintln(os.Stderr, `  tail -f apache.log | prefix -`)
		fmt.Fprintln(os.Stderr, `  my-cool-program 2>&1 | prefix -format="#{{.LineNumber5}} " -`)
	}
}
