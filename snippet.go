package main

import (
	"fmt"
	"sort"
	"strings"
)

// Snippet is one labelled fragment destined for a target file.
type Snippet struct {
	File   string
	Label  int
	Code   string
	Inline bool
	Line   int // 1-based source line in the markdown document
}

// FileSnippets groups snippets that belong to one output file.
type FileSnippets struct {
	File     string
	Snippets []Snippet
}

func groupAndSort(snips []Snippet) ([]FileSnippets, error) {
	byFile := map[string][]Snippet{}
	order := []string{}
	seen := map[string]bool{}

	for _, s := range snips {
		if !seen[s.File] {
			seen[s.File] = true
			order = append(order, s.File)
		}
		byFile[s.File] = append(byFile[s.File], s)
	}

	out := make([]FileSnippets, 0, len(order))
	for _, file := range order {
		group := byFile[file]
		sort.SliceStable(group, func(i, j int) bool {
			return group[i].Label < group[j].Label
		})
		labels := map[int]Snippet{}
		for _, s := range group {
			if prev, ok := labels[s.Label]; ok {
				return nil, fmt.Errorf("%s: duplicate label %d (markdown lines %d and %d)",
					file, s.Label, prev.Line, s.Line)
			}
			labels[s.Label] = s
		}
		out = append(out, FileSnippets{File: file, Snippets: group})
	}
	return out, nil
}

// assemble joins snippets for one file in label order.
// A blank line is inserted when successive labels jump by 10 or more,
// matching the usual literate gap between "nearby" and "distant" slots.
func assemble(group FileSnippets) string {
	var b strings.Builder
	var prevLabel *int
	for i, s := range group.Snippets {
		code := strings.TrimRight(s.Code, "\r\n")
		if prevLabel != nil && s.Label-*prevLabel >= 10 {
			if !strings.HasSuffix(b.String(), "\n\n") {
				b.WriteByte('\n')
			}
		}
		b.WriteString(code)
		if i < len(group.Snippets)-1 || !strings.HasSuffix(code, "\n") {
			b.WriteByte('\n')
		}
		l := s.Label
		prevLabel = &l
	}
	text := b.String()
	if !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	return text
}
