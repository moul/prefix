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
	var input io.Reader

	// parse args
	{
		switch {
		case len(args) == 1:
			input = os.Stdin
		case len(args) == 2 && args[1] == "-":
			input = os.Stdin
		case len(args) == 2:
			f, err := os.Open(args[1])
			if err != nil {
				return err
			}
			defer f.Close()
			input = f
		case len(args) > 3:
			return fmt.Errorf("usage: prefix FILE")
		}
	}

	// configure prefixer
	var prefixer prefix.LinePrefixer
	{
		prefixer = prefix.New(prefix.DefaultFormat)
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
