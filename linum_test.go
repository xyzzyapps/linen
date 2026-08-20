package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const explainer = `# Building a CSV Stream Summarizer

We start by defining the executable entry point at label ` + "`100`" + `.

` + "```python [app.py:100]\n" + `if __name__ == "__main__":
    main()
` + "```" + `

Before execution, we need system imports at label ` + "`10`" + `. We can add inline snippets like ` + "`import sys` [app.py:10]" + ` directly in prose, followed by module-specific imports:

` + "```python [app.py:12]\n" + `import csv
` + "```" + `

Next, we write our core calculation helper at label ` + "`30`" + `.

` + "```python [app.py:30]\n" + `def compute_avg(data):
    return sum(data) / len(data) if data else 0.0
` + "```" + `

Now we define ` + "`main()`" + ` at label ` + "`50`" + `, which consumes ` + "`sys.stdin`" + ` and calls ` + "`compute_avg()`" + `.

` + "```python [app.py:50]\n" + `def main():
    reader = csv.reader(sys.stdin)
    vals = [float(row[0]) for row in reader if row]
    print(f"Average: {compute_avg(vals):.2f}")
` + "```" + `

We also define default configuration settings in a separate target file (` + "`settings.json`" + `):

` + "```json [settings.json:10]\n" + `{
  "precision": 2,
  "stream_mode": true
}
` + "```" + `
`

func TestParseExplainer(t *testing.T) {
	snips, err := parseMarkdown(explainer)
	if err != nil {
		t.Fatal(err)
	}
	if len(snips) != 6 {
		t.Fatalf("got %d snippets, want 6", len(snips))
	}
	var inline *Snippet
	for i := range snips {
		if snips[i].Inline {
			inline = &snips[i]
		}
	}
	if inline == nil || inline.Code != "import sys" || inline.Label != 10 {
		t.Fatalf("inline snippet = %+v", inline)
	}
}

func TestAssembleExplainer(t *testing.T) {
	snips, err := parseMarkdown(explainer)
	if err != nil {
		t.Fatal(err)
	}
	groups, err := groupAndSort(snips)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, g := range groups {
		got[g.File] = assemble(g)
	}

	wantPy := `import sys
import csv

def compute_avg(data):
    return sum(data) / len(data) if data else 0.0

def main():
    reader = csv.reader(sys.stdin)
    vals = [float(row[0]) for row in reader if row]
    print(f"Average: {compute_avg(vals):.2f}")

if __name__ == "__main__":
    main()
`
	wantJSON := `{
  "precision": 2,
  "stream_mode": true
}
`
	if got["app.py"] != wantPy {
		t.Errorf("app.py mismatch\nGOT:\n%s\nWANT:\n%s", got["app.py"], wantPy)
	}
	if got["settings.json"] != wantJSON {
		t.Errorf("settings.json mismatch\nGOT:\n%s\nWANT:\n%s", got["settings.json"], wantJSON)
	}
}

func TestDuplicateLabel(t *testing.T) {
	src := "```go [a.go:1]\nfoo\n```\n```go [a.go:1]\nbar\n```\n"
	snips, err := parseMarkdown(src)
	if err != nil {
		t.Fatal(err)
	}
	_, err = groupAndSort(snips)
	if err == nil {
		t.Fatal("expected duplicate label error")
	}
}

func TestWriteFiles(t *testing.T) {
	dir := t.TempDir()
	md := filepath.Join(dir, "explainer.md")
	if err := os.WriteFile(md, []byte(explainer), 0o644); err != nil {
		t.Fatal(err)
	}
	snips, err := parseMarkdown(explainer)
	if err != nil {
		t.Fatal(err)
	}
	groups, err := groupAndSort(snips)
	if err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	for _, g := range groups {
		dest := filepath.Join(out, g.File)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dest, []byte(assemble(g)), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	py, err := os.ReadFile(filepath.Join(out, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(py), "import sys") {
		t.Fatal("app.py missing import sys")
	}
}
