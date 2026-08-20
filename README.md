# linum

Literate programming with line-number labels.

The name is from *line number* and Latin *linum* (flax — the fiber you weave). The Markdown narrative is the weft; numeric labels are the warp. You write code in story order. **linum** extracts every labelled fragment, groups by target file, sorts by label, and emits real source files.

## Syntax

**Labelled fence** — a fenced block tagged with `[file:N]`:

````markdown
```python [app.py:50]
def main():
    ...
```
````

**Inline snippet** — a code span followed by the same tag:

```markdown
We can add `import sys` [app.py:10] directly in prose.
```

Unlabelled fences and ordinary backticks are left alone.

## Assembly

1. Extract every `[file:N]` fence and inline.
2. Group fragments by `file`.
3. Sort each group by `N` ascending.
4. Join the fragments. If successive labels jump by **10 or more**, insert a blank line.

That is why an entry point written at the top of the document as `[app.py:100]` lands at the bottom of `app.py`, while nearby imports at `10` and `12` stay adjacent.

Duplicate labels in the same file are an error. An unclosed labelled fence is an error.

## Install

```text
go install github.com/manic/linum@latest
```

Or from this repo:

```text
go build -o linum .
```

## Usage

```text
linum [flags] <document.md>
```

| Flag | Meaning |
|------|---------|
| `-o dir` | Write tangled files under `dir` (default `.`) |
| `-n` | Dry-run: print files to stdout instead of writing |
| `-v` | Print the extraction matrix (file, label, type, preview) |

Examples:

```text
linum examples/explainer.md
linum -o out examples/explainer.md
linum -n -v examples/explainer.md
```

## Example

`examples/explainer.md` tells the story of a CSV stream summarizer out of execution order: the `__main__` guard first, then imports (including an inline `import sys`), then helpers, then `main()`, then a second target `settings.json`.

```text
go run . -n examples/explainer.md
```

emits:

**`app.py`**

```python
import sys
import csv

def compute_avg(data):
    return sum(data) / len(data) if data else 0.0

def main():
    reader = csv.reader(sys.stdin)
    vals = [float(row[0]) for row in reader if row]
    print(f"Average: {compute_avg(vals):.2f}")

if __name__ == "__main__":
    main()
```

**`settings.json`**

```json
{
  "precision": 2,
  "stream_mode": true
}
```

## Tests

```text
go test ./...
```
