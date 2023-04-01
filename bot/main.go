package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"text/template"
	"time"

	"github.com/primes-today/primebot"
)

type fileInterface struct {
	path string
}

func (f fileInterface) Fetch(ctx context.Context) (*primebot.Status, error) {
	raw, err := os.ReadFile(f.path)
	if err != nil {
		return &primebot.Status{}, err
	}

	status, success := (&big.Int{}).SetString(string(raw), 10)
	if !success {
		return &primebot.Status{}, fmt.Errorf("failed to convert string to big.Int: %s", raw)
	}

	return &primebot.Status{
		LastStatus: status,
		Posted:     time.Now(),
	}, nil
}

func (f fileInterface) Post(ctx context.Context, status *big.Int) error {
	fs, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fs.Close()

	fs.WriteString(status.Text(10))

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

func (t *templateWriter) Post(ctx context.Context, status *big.Int) error {
	tpl, err := template.ParseFiles(t.template)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(fmt.Sprintf(t.outputTemplate, status), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	type templateData struct {
		Date      string
		Integer   string
		Formatted string
	}

	if t.now.IsZero() {
		t.SetNow(time.Now())
	}

	tpl.Execute(out, templateData{
		t.now.UTC().Format(time.RFC3339),
		status.Text(10),
		format(status),
	})

	return nil
}

type multiPoster struct {
	posters []primebot.Poster
}

func (m multiPoster) Post(ctx context.Context, status *big.Int) error {
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

	files, err := os.ReadDir("content/primes")
	if err != nil {
		panic(err)
	}

	expected := len(files)

	checked := 0
	g := primebot.NewCompositeGenerator(big.NewInt(0))
	for {
		status, err := g.Generate(context.Background())
		if err != nil {
			panic(err)
		}

		if checked%1000 == 0 {
			fmt.Fprintf(os.Stderr, "checking status %d, %s\n", checked, status)
		}

		if status.Cmp(max.LastStatus) > 0 {
			if expected != checked {
				panic(fmt.Sprintf("expected %d files, but checked %d", expected, checked))
			}

			fmt.Fprintf(os.Stdout, "reached end status %d, checked %d files\n", max.LastStatus, checked)
			break
		}

		_, err = os.ReadFile(fmt.Sprintf("content/primes/%d.md", status))
		if err != nil {
			panic(err)
		}

		checked = checked + 1
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
	max, success := (&big.Int{}).SetString(flag.Arg(3), 10)
	if !success {
		panic(fmt.Errorf("unable to parse big.Int from arg: %s", flag.Arg(3)))
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

	g := primebot.NewCompositeGenerator(status.LastStatus)
	span := end_d.Sub(start_d)
	var statuses []*big.Int

	for {
		status, err := g.Generate(context.Background())
		if err != nil {
			panic(err)
		}

		if status.Cmp(max) >= 0 {
			fmt.Fprintf(os.Stderr, "status %d larger than max %d, breaking", status, max)
			break
		}

		statuses = append(statuses, status)
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
		status, success := (&big.Int{}).SetString(record.Tweet.FullText, 10)
		if !success {
			panic(fmt.Errorf("unable to parse big.Int from tweet record: %s", record.Tweet.FullText))
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

	g := primebot.NewCompositeGenerator(status.LastStatus)
	next, err := g.Generate(context.Background())
	if err != nil {
		panic(err)
	}

	err = p.Post(context.Background(), next)
	if err != nil {
		panic(err)
	}

	fmt.Printf("posted latest prime %s", next)
}

func reverse(s string) (result string) {
	for _, r := range s {
		result = string(r) + result
	}
	return
}

func format(i *big.Int) string {
	var f string
	s := reverse(i.Text(10))
	for i, r := range s {
		if i != 0 && i%3 == 0 {
			f = f + "," + string(r)
			continue
		}
		f = f + string(r)
	}

	return reverse(f)
}
