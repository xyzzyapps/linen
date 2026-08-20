package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	fenceOpenRe = regexp.MustCompile(`^` + "```" + `([^\n\[]*)\[([^\]\s:]+):(\d+)\]\s*$`)
	fenceBareRe = regexp.MustCompile(`^` + "```" + `\s*$`)
	inlineRe    = regexp.MustCompile("`([^`\n]+)`\\s*\\[([^\\]\\s:]+):(\\d+)\\]")
)

// parseMarkdown extracts labelled fenced blocks and inline snippets.
func parseMarkdown(src string) ([]Snippet, error) {
	lines := strings.Split(src, "\n")
	var snips []Snippet
	inFence := false
	var fenceFile string
	var fenceLabel int
	var fenceStart int
	var fenceBody []string

	for i, line := range lines {
		lineNo := i + 1
		trimmed := strings.TrimRight(line, "\r")

		if inFence {
			if fenceBareRe.MatchString(trimmed) || strings.HasPrefix(trimmed, "```") {
				code := strings.Join(fenceBody, "\n")
				snips = append(snips, Snippet{
					File:   fenceFile,
					Label:  fenceLabel,
					Code:   strings.TrimRight(code, "\n"),
					Inline: false,
					Line:   fenceStart,
				})
				inFence = false
				fenceBody = nil
				continue
			}
			fenceBody = append(fenceBody, trimmed)
			continue
		}

		if m := fenceOpenRe.FindStringSubmatch(trimmed); m != nil {
			label, err := strconv.Atoi(m[3])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid label %q", lineNo, m[3])
			}
			inFence = true
			fenceFile = m[2]
			fenceLabel = label
			fenceStart = lineNo
			fenceBody = nil
			continue
		}

		for _, m := range inlineRe.FindAllStringSubmatch(trimmed, -1) {
			label, err := strconv.Atoi(m[3])
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid label %q", lineNo, m[3])
			}
			snips = append(snips, Snippet{
				File:   m[2],
				Label:  label,
				Code:   strings.TrimSpace(m[1]),
				Inline: true,
				Line:   lineNo,
			})
		}
	}

	if inFence {
		return nil, fmt.Errorf("line %d: unclosed labelled fence for %s:%d", fenceStart, fenceFile, fenceLabel)
	}
	return snips, nil
}
