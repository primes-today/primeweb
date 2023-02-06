package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fardog/primebot"
)

type fileInterface struct {
	path string
}

func (f fileInterface) Fetch(ctx context.Context) (*primebot.Status, error) {
	raw, err := os.ReadFile(f.path)
	if err != nil {
		return &primebot.Status{}, err
	}

	status, err := strconv.ParseUint(string(raw), 10, 64)
	if err != nil {
		return &primebot.Status{}, err
	}

	return &primebot.Status{
		LastStatus: status,
		Posted:     time.Now(),
	}, nil
}

func main() {
	flag.Usage = func() {
		_, exe := filepath.Split(os.Args[0])
		fmt.Fprint(os.Stderr, "A file output for primebot.")
		fmt.Fprintf(os.Stderr, "Usage:\n\n  %s <in_path> <out_path>\n\nOptions:\n\n", exe)
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	f := fileInterface{path: flag.Arg(0)}
	status, err := f.Fetch(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read status from file: %s, %s", flag.Arg(0), err.Error())
		os.Exit(1)
	}

	fs, err := os.OpenFile(flag.Arg(1), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open output file: %s, %s", flag.Arg(1), err.Error())
		os.Exit(1)
	}
	defer fs.Close()

	g := primebot.NewProbablyPrimeGenerator(status.LastStatus + 1)
	next, err := g.Generate(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate next prime: %s", err.Error())
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d", next)
	fmt.Fprintf(fs, "%d", next)
}
