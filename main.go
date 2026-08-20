package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outDir := flag.String("o", ".", "directory to write tangled files into")
	dryRun := flag.Bool("n", false, "print tangled files to stdout instead of writing")
	verbose := flag.Bool("v", false, "print extraction matrix and written paths")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `linum — literate programming with line-number labels

Read a Markdown narrative whose code is tagged as [file:N], sort fragments
by N within each file, and emit the assembled sources.

  `+"`code`"+` [app.py:10]              inline snippet
  `+"```python [app.py:50]"+`           labelled fence

Usage:
  linum [flags] <document.md>

Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	srcPath := flag.Arg(0)
	data, err := os.ReadFile(srcPath)
	if err != nil {
		fatal(err)
	}

	snips, err := parseMarkdown(string(data))
	if err != nil {
		fatal(err)
	}
	if len(snips) == 0 {
		fatal(fmt.Errorf("%s: no labelled snippets found", srcPath))
	}

	groups, err := groupAndSort(snips)
	if err != nil {
		fatal(err)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "Target File\tLabel\tType\tPreview\n")
		for _, g := range groups {
			for _, s := range g.Snippets {
				kind := "Block"
				if s.Inline {
					kind = "Inline"
				}
				preview := oneLine(s.Code, 48)
				fmt.Fprintf(os.Stderr, "%s\t%d\t%s\t%s\n", s.File, s.Label, kind, preview)
			}
		}
	}

	for _, g := range groups {
		body := assemble(g)
		dest := filepath.Join(*outDir, filepath.FromSlash(g.File))
		if *dryRun {
			fmt.Printf("===== %s =====\n%s", g.File, body)
			if !hasTrailingNL(body) {
				fmt.Println()
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			fatal(err)
		}
		if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
			fatal(err)
		}
		if *verbose {
			fmt.Fprintf(os.Stderr, "wrote %s (%d snippets)\n", dest, len(g.Snippets))
		}
	}
}

func oneLine(s string, max int) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func hasTrailingNL(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\n'
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "linum: %v\n", err)
	os.Exit(1)
}
