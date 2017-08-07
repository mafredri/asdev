package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/go-snakecase"
)

const (
	acURL = "https://www.asustor.com/en/app_central"
)

func main() {
	flag.Parse()
	pkg := flag.Arg(0)
	name := flag.Arg(1)

	if pkg == "" {
		fmt.Fprintf(os.Stderr, "error: package name (first argument) must be provided\n")
		os.Exit(1)
	}

	if name == "" {
		fmt.Fprintf(os.Stderr, "error: filename (second argument) must be provided\n")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	b, err := run(ctx)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(name, b, 0644)
	if err != nil {
		panic(err)
	}
}

var (
	catRe = regexp.MustCompile(`<option value="[^"]+\?type=([0-9]+)" ?>([^<]+)</option>`)
)

type category struct {
	id   string
	desc string
}

func (c category) snake() string {
	return snakecase.Snakecase(c.desc)
}

type generator struct {
	buf bytes.Buffer
}

func (g *generator) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(&g.buf, format, a...)
}

func (g *generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Printf("error: could not format: %v", err)
		return g.buf.Bytes()
	}
	return src
}

func run(ctx context.Context) ([]byte, error) {
	resp, err := http.Get(acURL)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cats []category
	for _, m := range catRe.FindAllStringSubmatch(string(b), -1) {
		cats = append(cats, category{
			id:   m[1],
			desc: strings.Trim(m[2], " "),
		})
	}

	var g generator
	g.Printf(`// Code generated by catgen. DO NOT EDIT.

package main

type category string

var (
	categories = []string{`)
	for _, c := range cats {
		g.Printf(`
		%q,`, c.snake())
	}
	g.Printf(`
	}
)

func (c category) ID() string {
	switch c {`)
	for _, c := range cats {
		g.Printf(`
	case %q:
		return %q`, c.snake(), c.id)
	}
	g.Printf(`
	default:
		panic("category " + string(c) + " does not exist")
	}
}

func (c category) Name() string {
	switch c {`)
	for _, c := range cats {
		g.Printf(`
	case %q:
		return %q`, c.snake(), c.desc)
	}
	g.Printf(`
	default:
		panic("category " + string(c) + " does not exist")
	}
}
`)

	return g.format(), nil
}