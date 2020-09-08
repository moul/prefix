package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

func run(args []string) error {
	flags := flag.NewFlagSet("prefix", flag.ExitOnError)
	format := flags.String("format", prefix.DefaultFormat, "format string")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// parse args
	var input io.Reader
	{
		remainingArgs := flags.Args()
		fmt.Println(remainingArgs)
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
