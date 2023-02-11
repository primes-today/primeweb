package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/fardog/primebot"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func (f fileInterface) Post(ctx context.Context, status uint64) error {
	fs, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fs.Close()

	fmt.Fprintf(fs, "%d", status)

	return nil
}

type templateWriter struct {
	template       string
	outputTemplate string
}

func (t templateWriter) Post(ctx context.Context, status uint64) error {
	tpl, err := template.ParseFiles(t.template)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(fmt.Sprintf(t.outputTemplate, status), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	p := message.NewPrinter(language.English)

	type templateData struct {
		Date      string
		Integer   uint64
		Formatted string
	}

	tpl.Execute(out, templateData{
		time.Now().UTC().Format(time.RFC3339),
		status,
		p.Sprintf("%d", status),
	})

	return nil
}

type multiPoster struct {
	posters []primebot.Poster
}

func (m multiPoster) Post(ctx context.Context, status uint64) error {
	for _, p := range m.posters {
		err := p.Post(ctx, status)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	f := fileInterface{path: "_current"}
	status, err := f.Fetch(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read existing status from file: %s", err)
		os.Exit(1)
	}

	p := multiPoster{posters: []primebot.Poster{
		fileInterface{path: "_current"},
		templateWriter{template: "archetypes/prime.md", outputTemplate: "content/prime/%d.md"},
	}}

	g := primebot.NewProbablyPrimeGenerator(status.LastStatus + 1)
	next, err := g.Generate(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate next prime: %s", err)
		os.Exit(1)
	}

	err = p.Post(context.Background(), next)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write post template: %s", err)
	}
}
