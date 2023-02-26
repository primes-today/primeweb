package main

import (
	"context"
	"encoding/json"
	"flag"
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
	now            time.Time
}

func (t *templateWriter) SetNow(now time.Time) {
	t.now = now
}

func (t *templateWriter) Post(ctx context.Context, status uint64) error {
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

	if t.now.IsZero() {
		t.SetNow(time.Now())
	}

	tpl.Execute(out, templateData{
		t.now.UTC().Format(time.RFC3339),
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

func validate() {
	f := fileInterface{path: "_current"}
	max, err := f.Fetch(context.Background())
	if err != nil {
		panic(err)
	}

	checked := 0
	g := primebot.NewProbablyPrimeGenerator(2)
	for {
		status, err := g.Generate(context.Background())
		if err != nil {
			panic(err)
		}

		if status > max.LastStatus {
			fmt.Fprintf(os.Stdout, "reached end status %d, checked %d files\n", max.LastStatus, checked)
			break
		}

		_, err = os.ReadFile(fmt.Sprintf("content/primes/%d.md", status))
		if err != nil {
			panic(err)
		}

		checked = checked + 1
		g.SetStart(status + 1)
	}
}

func backfill() {
	start_d, err := time.Parse(time.RFC3339, flag.Arg(1))
	if err != nil {
		panic(err)
	}
	end_d, err := time.Parse(time.RFC3339, flag.Arg(2))
	if err != nil {
		panic(err)
	}
	max, err := strconv.ParseUint(flag.Arg(3), 10, 64)
	if err != nil {
		panic(err)
	}

	f := fileInterface{path: "_current"}
	status, err := f.Fetch(context.Background())
	if err != nil {
		panic(err)
	}

	tpl := templateWriter{template: "archetypes/prime.md", outputTemplate: "content/primes/%d.md"}
	p := multiPoster{posters: []primebot.Poster{
		fileInterface{path: "_current"},
		&tpl,
	}}

	g := primebot.NewProbablyPrimeGenerator(status.LastStatus + 1)
	span := end_d.Sub(start_d)
	var statuses []uint64

	for {
		status, err := g.Generate(context.Background())
		if err != nil {
			panic(err)
		}

		if status >= max {
			fmt.Fprintf(os.Stderr, "status %d larger than max %d, breaking", status, max)
			break
		}

		statuses = append(statuses, status)
		g.SetStart(status + 1)
	}

	sep := int(span.Seconds()) / len(statuses)

	fmt.Fprint(os.Stderr, statuses)
	for i, status := range statuses {
		tpl.SetNow(start_d.Add(time.Second * time.Duration(sep*(i+1))))
		p.Post(context.Background(), status)
	}
}

func importFile() {
	type Tweet struct {
		CreatedAt string `json:"created_at"`
		FullText  string `json:"full_text"`
	}

	type Record struct {
		Tweet Tweet `json:"tweet"`
	}

	raw, err := os.ReadFile(flag.Arg(1))
	if err != nil {
		panic(err)
	}

	var records []Record
	err = json.Unmarshal(raw, &records)
	if err != nil {
		panic(err)
	}

	tpl := templateWriter{template: "archetypes/prime.md", outputTemplate: "content/primes/%d.md"}
	p := multiPoster{posters: []primebot.Poster{
		fileInterface{path: "_current"},
		&tpl,
	}}

	for _, record := range records {
		status, err := strconv.ParseUint(record.Tweet.FullText, 10, 64)
		if err != nil {
			panic(err)
		}
		t, err := time.Parse(time.RubyDate, record.Tweet.CreatedAt)
		if err != nil {
			panic(err)
		}
		tpl.SetNow(t)
		p.Post(context.Background(), status)
	}
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "backfill" {
		backfill()
		return
	} else if flag.Arg(0) == "import_file" {
		importFile()
		return
	} else if flag.Arg(0) == "validate" {
		validate()
		return
	}

	f := fileInterface{path: "_current"}
	status, err := f.Fetch(context.Background())
	if err != nil {
		panic(err)
	}

	p := multiPoster{posters: []primebot.Poster{
		fileInterface{path: "_current"},
		&templateWriter{template: "archetypes/prime.md", outputTemplate: "content/primes/%d.md"},
	}}

	g := primebot.NewProbablyPrimeGenerator(status.LastStatus + 1)
	next, err := g.Generate(context.Background())
	if err != nil {
		panic(err)
	}

	err = p.Post(context.Background(), next)
	if err != nil {
		panic(err)
	}
}
