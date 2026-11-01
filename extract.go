package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

const maxDocBytes = 1 << 20 // 1 MiB

type Document struct {
	Path      string   `json:"path"`
	Type      string   `json:"type,omitempty"`
	Status    string   `json:"status,omitempty"`
	Owner     string   `json:"owner,omitempty"`
	DependsOn []string `json:"depends_on"`
}

type ExtractError struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type Stats struct {
	Matched int `json:"documents_matched"`
	Parsed  int `json:"documents_parsed"`
	Edges   int `json:"edges_evaluated"`
}

type Corpus struct {
	Documents     []Document     `json:"documents"`
	ExtractErrors []ExtractError `json:"extract_errors"`
	Stats         Stats          `json:"stats"`
}

// Extract walks every governed glob and turns matching files into a corpus.
// Every governed file must produce either a document or an extract error; a
// file is never silently dropped (design doc:131).
func Extract(root string, governed []string) (Corpus, error) {
	c := Corpus{Documents: []Document{}, ExtractErrors: []ExtractError{}}
	fsys := os.DirFS(root)
	seen := make(map[string]bool)

	for _, g := range governed {
		matches, err := doublestar.Glob(fsys, g, doublestar.WithFilesOnly())
		if err != nil {
			return c, fmt.Errorf("glob %q: %w", g, err)
		}
		sort.Strings(matches)
		for _, m := range matches {
			np, err := NormalizeRepoPath(m)
			if err != nil {
				c.ExtractErrors = append(c.ExtractErrors, ExtractError{Path: m, Reason: err.Error()})
				c.Stats.Matched++
				continue
			}
			if seen[np] {
				continue
			}
			seen[np] = true
			c.Stats.Matched++

			doc, err := readDocument(fsys, np)
			if err != nil {
				c.ExtractErrors = append(c.ExtractErrors, ExtractError{Path: np, Reason: err.Error()})
				continue
			}
			c.Documents = append(c.Documents, doc)
			c.Stats.Parsed++
			c.Stats.Edges += len(doc.DependsOn)
		}
	}

	sort.Slice(c.Documents, func(i, j int) bool { return c.Documents[i].Path < c.Documents[j].Path })
	sort.Slice(c.ExtractErrors, func(i, j int) bool { return c.ExtractErrors[i].Path < c.ExtractErrors[j].Path })
	return c, nil
}

func readDocument(fsys fs.FS, path string) (Document, error) {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		return Document{}, err
	}
	if info.Size() > maxDocBytes {
		return Document{}, fmt.Errorf("file is %d bytes, limit is %d", info.Size(), maxDocBytes)
	}
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return Document{}, err
	}
	fm, err := parseFrontmatter(b)
	if err != nil {
		return Document{}, err
	}
	return documentFromMap(path, fm)
}

func parseFrontmatter(b []byte) (map[string]any, error) {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return nil, errors.New("missing frontmatter")
	}
	rest := s[4:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return nil, errors.New("frontmatter has no closing delimiter")
	}
	body := rest[:idx]

	// Cheap guard against YAML alias expansion before handing bytes to the parser.
	if strings.Count(body, "*") > 100 || strings.Count(body, "<<") > 10 {
		return nil, errors.New("frontmatter rejected: possible YAML alias expansion")
	}

	var m map[string]any
	if err := yaml.Unmarshal([]byte(body), &m); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	if m == nil {
		return nil, errors.New("frontmatter is not a mapping")
	}
	return m, nil
}

func documentFromMap(path string, m map[string]any) (Document, error) {
	d := Document{Path: path, DependsOn: []string{}}
	if v, ok := m["type"].(string); ok {
		d.Type = v
	}
	if v, ok := m["status"].(string); ok {
		d.Status = v
	}
	if v, ok := m["owner"].(string); ok {
		d.Owner = v
	}

	raw, ok := m["depends_on"]
	if !ok || raw == nil {
		return d, nil
	}
	items, ok := raw.([]any)
	if !ok {
		return Document{}, errors.New("depends_on must be a list of paths")
	}

	seen := make(map[string]bool, len(items))
	for _, item := range items {
		s, ok := item.(string)
		if !ok {
			return Document{}, errors.New("depends_on entries must be strings")
		}
		n, err := NormalizeRepoPath(s)
		if err != nil {
			return Document{}, fmt.Errorf("depends_on %q: %w", s, err)
		}
		if seen[n] {
			continue
		}
		seen[n] = true
		d.DependsOn = append(d.DependsOn, n)
	}
	return d, nil
}
