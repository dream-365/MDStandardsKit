package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"

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

	// Idea-only facts. They stay empty for spec documents so the
	// -dump-corpus output for the existing corpus is unchanged.
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Date        string    `json:"date,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Sections    []Section `json:"sections,omitempty"`
	Missing     []string  `json:"missing,omitempty"`
	DateInvalid bool      `json:"date_invalid,omitempty"`
	TagsInvalid bool      `json:"tags_invalid,omitempty"`
}

type Section struct {
	Title       string `json:"title"`
	Level       int    `json:"level"`
	BodyChars   int    `json:"body_chars"`
	Empty       bool   `json:"empty"`
	Placeholder bool   `json:"placeholder,omitempty"`
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
	DocumentType  string         `json:"document_type,omitempty"`
	Mode          string         `json:"mode,omitempty"`
}

type docRoute func(path string, fm map[string]any) (DocumentType, string, error)

type routingError struct{ msg string }

func (e routingError) Error() string { return e.msg }

// Extract walks every governed glob and turns matching files into a corpus.
// Every governed file must produce either a document or an extract error; a
// file is never silently dropped (design doc:131).
func Extract(root string, governed []string, docType string, profile DocumentType) (Corpus, error) {
	return extract(root, governed, func(_ string, _ map[string]any) (DocumentType, string, error) {
		return profile, docType, nil
	}, "", nil)
}

// ExtractMulti routes each file by its frontmatter type. The top-level
// governed globs are only the discovery baseline; the route decides which
// profile applies and rejects files whose path does not match that type.
func ExtractMulti(root string, governed []string, profiles map[string]DocumentType, specGlobs []string) (Corpus, error) {
	return extract(root, governed, multiRoute(profiles, specGlobs), multiModeName, profiles)
}

func multiRoute(profiles map[string]DocumentType, specGlobs []string) docRoute {
	return func(path string, fm map[string]any) (DocumentType, string, error) {
		raw, ok := fm["type"]
		if !ok || raw == nil {
			return DocumentType{}, "", routingError{"missing frontmatter type in multi-type mode"}
		}
		s, ok := raw.(string)
		if !ok {
			return DocumentType{}, "", routingError{"frontmatter type must be a non-empty string"}
		}
		t := strings.TrimSpace(s)
		if t == "" {
			return DocumentType{}, "", routingError{"frontmatter type must be a non-empty string"}
		}
		profile, ok := profiles[t]
		if !ok {
			return DocumentType{}, "", routingError{fmt.Sprintf("document type %q is not enabled in .mdstandards.json", t)}
		}
		globs := profile.Globs
		if t == "spec" {
			globs = specGlobs
		}
		for _, g := range globs {
			matched, err := doublestar.Match(g, path)
			if err != nil {
				return DocumentType{}, "", fmt.Errorf("glob %q: %w", g, err)
			}
			if matched {
				return profile, t, nil
			}
		}
		return DocumentType{}, "", routingError{fmt.Sprintf("document type %q does not allow path %s", t, path)}
	}
}

func extract(root string, governed []string, route docRoute, mode string, profiles map[string]DocumentType) (Corpus, error) {
	c := Corpus{Documents: []Document{}, ExtractErrors: []ExtractError{}, Mode: mode}
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

			doc, effectiveType, err := readDocument(fsys, np, route)
			if err != nil {
				c.ExtractErrors = append(c.ExtractErrors, ExtractError{Path: np, Reason: err.Error()})
				continue
			}
			if mode == multiModeName {
				doc.Type = effectiveType
				if profiles[effectiveType].Graph {
					c.Stats.Edges += len(doc.DependsOn)
				}
			} else {
				c.Stats.Edges += len(doc.DependsOn)
			}
			c.Documents = append(c.Documents, doc)
			c.Stats.Parsed++
		}
	}

	sort.Slice(c.Documents, func(i, j int) bool { return c.Documents[i].Path < c.Documents[j].Path })
	sort.Slice(c.ExtractErrors, func(i, j int) bool { return c.ExtractErrors[i].Path < c.ExtractErrors[j].Path })
	return c, nil
}

func readDocument(fsys fs.FS, path string, route docRoute) (Document, string, error) {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		return Document{}, "", err
	}
	if info.Size() > maxDocBytes {
		return Document{}, "", fmt.Errorf("file is %d bytes, limit is %d", info.Size(), maxDocBytes)
	}
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return Document{}, "", err
	}
	fm, err := parseFrontmatter(b)
	if err != nil {
		return Document{}, "", err
	}
	profile, effectiveType, err := route(path, fm)
	if err != nil {
		return Document{}, "", err
	}
	doc, err := documentFromMap(path, fm, bodyAfterFrontmatter(b), effectiveType, profile)
	return doc, effectiveType, err
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

func documentFromMap(path string, m map[string]any, body, effectiveType string, profile DocumentType) (Document, error) {
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

	if raw, ok := m["depends_on"]; ok && raw != nil {
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
	}

	ruleSet := profile.RuleSet
	if ruleSet == "" {
		ruleSet = effectiveType
	}
	if ruleSet == "idea" {
		dateField := profile.DateField
		if dateField == "" {
			dateField = "date"
		}
		if v, ok := m["id"].(string); ok {
			d.ID = strings.TrimSpace(v)
		}
		if v, ok := m["title"].(string); ok {
			d.Title = strings.TrimSpace(v)
		}
		if raw, ok := m[dateField]; ok && raw != nil {
			var valid bool
			d.Date, valid = parseDate(raw)
			if !valid {
				d.DateInvalid = true
			}
		}
		if raw, ok := m["tags"]; ok && raw != nil {
			if tags, ok := parseTags(raw); ok {
				d.Tags = tags
			} else {
				d.TagsInvalid = true
			}
		}
		// title is checked by the dedicated empty-title rule; the rest come
		// from the type profile so the extractor stays type-agnostic.
		for _, field := range profile.RequiredFields {
			if field == "title" {
				continue
			}
			check := field
			if field == "date" {
				check = dateField
			}
			if raw, ok := m[check]; !ok || raw == nil {
				d.Missing = appendUnique(d.Missing, field)
			}
		}
		if d.ID == "" {
			d.Missing = appendUnique(d.Missing, "id")
		}
		if d.Status == "" {
			d.Missing = appendUnique(d.Missing, "status")
		}
		if body != "" {
			d.Sections = parseSections(body, profile.PlaceholderMarkers)
		}
	}
	return d, nil
}

func parseDate(v any) (string, bool) {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if _, err := time.Parse("2006-01-02", s); err != nil {
			return s, false
		}
		return s, true
	case time.Time:
		// YAML resolves an unquoted date-only scalar to time.Time. A real
		// datetime carries a non-zero clock and is not a date-only value.
		if t.Hour() != 0 || t.Minute() != 0 || t.Second() != 0 || t.Nanosecond() != 0 {
			return t.Format("2006-01-02"), false
		}
		return t.Format("2006-01-02"), true
	default:
		return "", false
	}
}

func parseTags(v any) ([]string, bool) {
	items, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		s, ok := item.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

func appendUnique(xs []string, s string) []string {
	for _, x := range xs {
		if x == s {
			return xs
		}
	}
	return append(xs, s)
}

func bodyAfterFrontmatter(b []byte) string {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return ""
	}
	rest := s[4:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return ""
	}
	after := rest[idx+4:]
	return strings.TrimPrefix(after, "\n")
}

func parseSections(body string, markers []string) []Section {
	lines := strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n")
	var sections []Section
	var cur *Section
	var buf []string
	inFence := false
	var fence string

	flush := func() {
		if cur == nil {
			return
		}
		text := strings.TrimSpace(strings.Join(buf, "\n"))
		cur.BodyChars = len([]rune(text))
		cur.Empty = text == ""
		for _, marker := range markers {
			if marker != "" && strings.Contains(text, marker) {
				cur.Placeholder = true
				break
			}
		}
		sections = append(sections, *cur)
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if inFence {
			if strings.HasPrefix(trimmed, fence) {
				inFence = false
				fence = ""
			}
			buf = append(buf, line)
			continue
		}
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = true
			if strings.HasPrefix(trimmed, "```") {
				fence = "```"
			} else {
				fence = "~~~"
			}
			buf = append(buf, line)
			continue
		}
		if level, title, ok := parseATXHeading(line); ok {
			flush()
			cur = &Section{Title: title, Level: level}
			buf = nil
			continue
		}
		if cur != nil {
			buf = append(buf, line)
		}
	}
	flush()
	return sections
}

func parseATXHeading(line string) (int, string, bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trimmed, "#") {
		return 0, "", false
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == '#' {
		n++
	}
	if n > 6 || n == len(trimmed) {
		return 0, "", false
	}
	if trimmed[n] != ' ' && trimmed[n] != '\t' {
		return 0, "", false
	}
	title := strings.TrimSpace(trimmed[n:])
	title = strings.TrimSpace(strings.TrimRight(title, "#"))
	if title == "" {
		return 0, "", false
	}
	return n, title, true
}
