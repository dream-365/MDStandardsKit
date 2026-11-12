package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func typedDoc(typ, status string, deps ...string) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("type: " + typ + "\n")
	b.WriteString("status: " + status + "\n")
	b.WriteString("owner: example\n")
	if len(deps) > 0 {
		b.WriteString("depends_on:\n")
		for _, d := range deps {
			b.WriteString("  - " + d + "\n")
		}
	}
	b.WriteString("---\n\n# Body\n")
	return b.String()
}

func multiRepo(t *testing.T, config string, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	writeFile(t, root, ".mdstandards.json", config)
	for rel, content := range files {
		writeFile(t, root, rel, content)
	}
	return root
}

func testdataMultiRepo(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("testdata/mixed-repo")
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func TestMixedRepoPasses(t *testing.T) {
	code, out, errs := runCLI(t, testdataMultiRepo(t), "-phase", "1", "-format", "json")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	rep := decodeReport(t, out)
	if len(rep.Findings) != 0 || len(rep.ExtractErrors) != 0 {
		t.Fatalf("findings=%+v extract_errors=%+v", rep.Findings, rep.ExtractErrors)
	}
	if rep.Stats.Matched != 6 || rep.Stats.Parsed != 6 || rep.Stats.Edges != 1 {
		t.Fatalf("stats = %+v, want matched=6 parsed=6 edges=1", rep.Stats)
	}
}

func TestLegacyArchitectureTypeRejected(t *testing.T) {
	root := multiRepo(t, `{"type":"architecture"}`, nil)
	code, _, errs := runCLI(t, root, "-phase", "0", "-format", "json")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
}

func TestMultiDiscoveryBaseline(t *testing.T) {
	root := multiRepo(t, `{"types":["spec"]}`, map[string]string{
		"docs/specs/001/spec.md":   typedDoc("spec", "Approved"),
		"docs/architecture/001.md": typedDoc("architecture", "Approved"),
	})
	code, out, _ := runCLI(t, root, "-dump-corpus")
	if code != 0 {
		t.Fatalf("dump exit = %d", code)
	}
	var c Corpus
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range c.ExtractErrors {
		if e.Path == "docs/architecture/001.md" && strings.Contains(e.Reason, `"architecture" is not enabled`) {
			found = true
		}
	}
	if !found {
		t.Fatalf("extract_errors = %+v, want unenabled architecture discovered via top-level globs", c.ExtractErrors)
	}
}

func TestMultiErrorReasons(t *testing.T) {
	profiles := map[string]DocumentType{
		"spec":  {RuleSet: "spec"},
		"index": {RuleSet: "none", Globs: []string{"docs/specs/README.md"}},
	}
	route := multiRoute(profiles, []string{"docs/specs/**/*.md"})
	cases := []struct {
		name string
		fm   map[string]any
		want string
	}{
		{name: "missing", fm: map[string]any{}, want: "missing frontmatter type"},
		{name: "non-string", fm: map[string]any{"type": []any{"spec"}}, want: "non-empty string"},
		{name: "empty", fm: map[string]any{"type": "  "}, want: "non-empty string"},
		{name: "not enabled", fm: map[string]any{"type": "note"}, want: "not enabled"},
		{name: "path mismatch", fm: map[string]any{"type": "index"}, want: "does not allow path"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := route("docs/specs/001.md", tc.fm)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestLegacyRollback(t *testing.T) {
	files := map[string]string{
		"docs/specs/001/spec.md":   typedDoc("spec", "Approved", "docs/architecture/003.md"),
		"docs/architecture/003.md": typedDoc("architecture", "Approved"),
	}
	root := multiRepo(t, `{"types":["spec","architecture"]}`, files)
	if code, _, errs := runCLI(t, root, "-phase", "1", "-format", "json"); code != 0 {
		t.Fatalf("multi exit = %d; stderr=%s", code, errs)
	}

	writeFile(t, root, ".mdstandards.json", `{"type":"spec"}`)
	if code, _, errs := runCLI(t, root, "-phase", "1", "-format", "json"); code != 0 {
		t.Fatalf("legacy type rewrite exit = %d; stderr=%s", code, errs)
	}

	if err := os.Remove(filepath.Join(root, ".mdstandards.json")); err != nil {
		t.Fatal(err)
	}
	if code, _, errs := runCLI(t, root, "-phase", "1", "-format", "json"); code != 0 {
		t.Fatalf("config removal exit = %d; stderr=%s", code, errs)
	}
}

func TestMultiDumpCorpus(t *testing.T) {
	code, out, errs := runCLI(t, testdataMultiRepo(t), "-dump-corpus")
	if code != 0 {
		t.Fatalf("exit = %d; stderr=%s", code, errs)
	}
	var c Corpus
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatalf("invalid corpus JSON: %v\n%s", err, out)
	}
	if c.Mode != "multi" {
		t.Fatalf("mode = %q, want multi", c.Mode)
	}
	if c.DocumentType != "" {
		t.Fatalf("document_type = %q, want empty in multi mode", c.DocumentType)
	}
	types := map[string]string{}
	for _, d := range c.Documents {
		types[d.Path] = d.Type
	}
	want := map[string]string{
		"docs/specs/001-example/spec.md":            "spec",
		"docs/architecture/example-architecture.md": "architecture",
		"docs/specs/001-example/changelog.md":       "changelog",
		"docs/specs/001-example/tasks.md":           "tasks",
		"docs/specs/README.md":                      "index",
		"2026/2026-11-01/example-idea.md":           "idea",
	}
	for path, typ := range want {
		if types[path] != typ {
			t.Fatalf("type[%s] = %q, want %q", path, types[path], typ)
		}
	}
}

func TestLegacyDumpHasNoMode(t *testing.T) {
	code, out, errs := runCLI(t, "testdata/example-repo", "-dump-corpus")
	if code != 0 {
		t.Fatalf("exit = %d; stderr=%s", code, errs)
	}
	if strings.Contains(out, `"mode"`) {
		t.Fatalf("legacy corpus unexpectedly contains mode:\n%s", out)
	}
}

func TestMultiRoutingErrors(t *testing.T) {
	allTypes := `{"types":["spec","architecture","changelog","tasks","index","idea"]}`
	cases := []struct {
		name   string
		config string
		files  map[string]string
		reason string
	}{
		{
			name:   "missing type",
			config: allTypes,
			files: map[string]string{
				"docs/specs/001.md": "---\nstatus: Approved\n---\n\n# Body\n",
			},
			reason: "missing frontmatter type",
		},
		{
			name:   "unknown type",
			config: allTypes,
			files: map[string]string{
				"docs/specs/001.md": typedDoc("note", "Approved"),
			},
			reason: `document type "note" is not enabled`,
		},
		{
			name:   "non-string type",
			config: allTypes,
			files: map[string]string{
				"docs/specs/001.md": "---\ntype: [spec]\n---\n\n# Body\n",
			},
			reason: "must be a non-empty string",
		},
		{
			name:   "path mismatch",
			config: `{"types":["spec","index"]}`,
			files: map[string]string{
				"docs/specs/001.md": typedDoc("index", "Approved"),
			},
			reason: `does not allow path`,
		},
		{
			name:   "disabled type",
			config: `{"types":["spec"]}`,
			files: map[string]string{
				"docs/architecture/001.md": typedDoc("architecture", "Approved"),
			},
			reason: `"architecture" is not enabled`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := multiRepo(t, tc.config, tc.files)
			code, _, errs := runCLI(t, root, "-phase", "1", "-format", "json")
			if code != 2 || !strings.Contains(errs, "no governed document parsed") {
				t.Fatalf("phase 1 exit = %d stderr=%q, want 2 with zero-parsed coverage error", code, errs)
			}

			code, out, _ := runCLI(t, root, "-dump-corpus")
			if code != 0 {
				t.Fatalf("dump exit = %d", code)
			}
			var c Corpus
			if err := json.Unmarshal([]byte(out), &c); err != nil {
				t.Fatalf("invalid corpus JSON: %v\n%s", err, out)
			}
			if len(c.ExtractErrors) != 1 || !strings.Contains(c.ExtractErrors[0].Reason, tc.reason) {
				t.Fatalf("extract_errors = %+v, want reason containing %q", c.ExtractErrors, tc.reason)
			}

			code, out, _ = runCLI(t, root, "-phase", "0", "-format", "json")
			if code != 0 {
				t.Fatalf("phase 0 exit = %d, want 0; out=%s", code, out)
			}
			rep := decodeReport(t, out)
			if len(rep.ExtractErrors) != 1 || !strings.Contains(rep.ExtractErrors[0].Reason, tc.reason) {
				t.Fatalf("phase 0 extract_errors = %+v, want reason containing %q", rep.ExtractErrors, tc.reason)
			}
		})
	}
}

func TestMultiCanonicalType(t *testing.T) {
	root := multiRepo(t, `{"types":["spec"]}`, map[string]string{
		"docs/specs/001.md": "---\ntype: \" spec \"\nstatus: Approved\nowner: example\n---\n\n# Body\n",
	})
	code, out, errs := runCLI(t, root, "-dump-corpus")
	if code != 0 {
		t.Fatalf("exit = %d; stderr=%s", code, errs)
	}
	var c Corpus
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	if len(c.Documents) != 1 || c.Documents[0].Type != "spec" {
		t.Fatalf("documents = %+v, want canonical type spec", c.Documents)
	}
}

func TestMultiRuleScopeSpecAndIdea(t *testing.T) {
	root := multiRepo(t, `{"types":["spec","architecture","changelog","idea"]}`, map[string]string{
		"docs/specs/001/spec.md":      typedDoc("spec", "Approved", "docs/architecture/missing.md"),
		"docs/specs/001/changelog.md": typedDoc("changelog", "Implemented", "docs/architecture/missing.md"),
		"2026/2026-11-01/idea.md":     ideaContent("AABBCCDD", "", "2026-11-01", "萌芽"),
	})
	// The idea fixture has no type; add it.
	writeFile(t, root, "2026/2026-11-01/idea.md", "---\ntype: idea\n"+strings.TrimPrefix(ideaContent("AABBCCDD", "", "2026-11-01", "萌芽"), "---\n"))

	code, out, errs := runCLI(t, root, "-phase", "1", "-format", "json")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stderr=%s out=%s", code, errs, out)
	}
	rep := decodeReport(t, out)
	var specMissing, changelogMissing, ideaTitle bool
	for _, f := range rep.Findings {
		switch {
		case f.Rule == "refs.missing_target" && f.Document == "docs/specs/001/spec.md":
			specMissing = true
		case f.Rule == "refs.missing_target" && f.Document == "docs/specs/001/changelog.md":
			changelogMissing = true
		case f.Rule == "idea.empty_title" && f.Document == "2026/2026-11-01/idea.md":
			ideaTitle = true
		}
	}
	if !specMissing || changelogMissing || !ideaTitle {
		t.Fatalf("findings = %+v; specMissing=%v changelogMissing=%v ideaTitle=%v", rep.Findings, specMissing, changelogMissing, ideaTitle)
	}
}

func TestMultiCoverageNoneEdgesDoNotSatisfyZeroEdge(t *testing.T) {
	root := multiRepo(t, `{"types":["spec","changelog"]}`, map[string]string{
		"docs/specs/001/spec.md":      typedDoc("spec", "Approved"),
		"docs/specs/001/changelog.md": typedDoc("changelog", "Implemented", "docs/architecture/missing.md"),
	})
	code, _, errs := runCLI(t, root, "-phase", "1", "-format", "json")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "zero dependency edges") {
		t.Fatalf("stderr = %q, want zero dependency edges", errs)
	}
}

func TestMultiCoverageGraphTypeWithZeroDocs(t *testing.T) {
	root := multiRepo(t, `{"types":["spec","idea"]}`, map[string]string{
		"2026/2026-11-01/idea.md": ideaContent("AABBCCDD", "测试想法", "2026-11-01", "萌芽"),
	})
	writeFile(t, root, "2026/2026-11-01/idea.md", "---\ntype: idea\n"+strings.TrimPrefix(ideaContent("AABBCCDD", "测试想法", "2026-11-01", "萌芽"), "---\n"))
	code, out, errs := runCLI(t, root, "-phase", "1", "-format", "json")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	if !strings.Contains(errs, `enabled document type "spec" matched no parsed documents`) {
		t.Fatalf("stderr = %q, want zero-doc warning for spec", errs)
	}
}

func TestMultiCoverageZeroDocWarning(t *testing.T) {
	root := multiRepo(t, `{"types":["spec","architecture","idea"]}`, map[string]string{
		"2026/2026-11-01/idea.md": ideaContent("AABBCCDD", "测试想法", "2026-11-01", "萌芽"),
	})
	writeFile(t, root, "2026/2026-11-01/idea.md", "---\ntype: idea\n"+strings.TrimPrefix(ideaContent("AABBCCDD", "测试想法", "2026-11-01", "萌芽"), "---\n"))
	code, _, errs := runCLI(t, root, "-phase", "1", "-format", "json")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	for _, want := range []string{
		`enabled document type "architecture" matched no parsed documents`,
		`enabled document type "spec" matched no parsed documents`,
	} {
		if !strings.Contains(errs, want) {
			t.Fatalf("stderr = %q, want %q", errs, want)
		}
	}
}

func TestMultiConfigErrors(t *testing.T) {
	cases := []struct {
		name   string
		config string
	}{
		{name: "both keys", config: `{"type":"spec","types":["spec"]}`},
		{name: "empty types", config: `{"types":[]}`},
		{name: "duplicate types", config: `{"types":["spec","spec"]}`},
		{name: "unknown profile", config: `{"types":["note"]}`},
		{name: "empty object", config: `{}`},
		{name: "legacy architecture", config: `{"type":"architecture"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := multiRepo(t, tc.config, nil)
			code, _, errs := runCLI(t, root, "-phase", "0", "-format", "json")
			if code != 2 {
				t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
			}
		})
	}
}

func writePolicyVersion(t *testing.T, root, version string) {
	t.Helper()
	path := filepath.Join(root, "policy", "data", "standard.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out := strings.Replace(string(b), `"standard_version": "2"`, `"standard_version": "`+version+`"`, 1)
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMultiRequiresVersionTwo(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writePolicyVersion(t, root, "1")
	writeFile(t, root, ".mdstandards.json", `{"types":["changelog"]}`)

	code, _, errs := runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"-phase", "1", "-format", "json")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "requires standard_version 2") {
		t.Fatalf("stderr = %q", errs)
	}
}

func TestStandardVersionMatrix(t *testing.T) {
	t.Run("unsupported version", func(t *testing.T) {
		root := t.TempDir()
		copyPolicy(t, root)
		writePolicyVersion(t, root, "3")
		code, _, errs := runCLI(t, root,
			"--policy", filepath.Join(root, "policy"),
			"-phase", "1", "-format", "json")
		if code != 2 || !strings.Contains(errs, "unsupported standard_version") {
			t.Fatalf("exit=%d stderr=%q", code, errs)
		}
	})

	for _, version := range []string{"1", "2"} {
		t.Run("legacy spec accepts v"+version, func(t *testing.T) {
			root := t.TempDir()
			copyPolicy(t, root)
			writePolicyVersion(t, root, version)
			writeFile(t, root, "docs/specs/001.md", typedDoc("spec", "Approved", "docs/architecture/003.md"))
			writeFile(t, root, "docs/architecture/003.md", typedDoc("architecture", "Published"))
			code, out, errs := runCLI(t, root,
				"--policy", filepath.Join(root, "policy"),
				"-phase", "1", "-format", "json")
			if code != 0 {
				t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
			}
		})
	}
}

func runPolicyFromHead(t *testing.T, root string) (int, string, string) {
	t.Helper()
	return runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"-phase", "1", "-format", "json")
}

func TestPolicyFromHeadConfig(t *testing.T) {
	baseFiles := map[string]string{
		"docs/specs/001/spec.md":   typedDoc("spec", "Approved", "docs/architecture/003.md"),
		"docs/architecture/003.md": typedDoc("architecture", "Approved"),
		".mdstandards.json":        `{"types":["spec","architecture"]}`,
	}
	newRepo := func(t *testing.T, files map[string]string, commitConfig bool) string {
		t.Helper()
		root := t.TempDir()
		copyPolicy(t, root)
		for rel, content := range files {
			writeFile(t, root, rel, content)
		}
		gitRun(t, root, "init")
		gitRun(t, root, "add", ".")
		gitRun(t, root, "commit", "-m", "init")
		return root
	}

	t.Run("clean config is not dirty", func(t *testing.T) {
		root := newRepo(t, baseFiles, true)
		code, out, errs := runPolicyFromHead(t, root)
		if code != 0 {
			t.Fatalf("exit = %d; stderr=%s out=%s", code, errs, out)
		}
		if strings.Contains(errs, "differs from HEAD") {
			t.Fatalf("unexpected dirty warning: %s", errs)
		}
	})

	t.Run("modified config downgrades", func(t *testing.T) {
		root := newRepo(t, map[string]string{
			"docs/specs/001/spec.md":   typedDoc("spec", "Approved", "docs/architecture/missing.md"),
			"docs/architecture/003.md": typedDoc("architecture", "Approved"),
			".mdstandards.json":        `{"types":["spec","architecture"]}`,
		}, true)
		writeFile(t, root, ".mdstandards.json", `{"types":["spec","architecture","changelog"]}`)
		code, out, errs := runPolicyFromHead(t, root)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (dirty downgrade); stderr=%s out=%s", code, errs, out)
		}
		if !strings.Contains(errs, "differs from HEAD") {
			t.Fatalf("stderr = %q, want dirty warning", errs)
		}
	})

	t.Run("untracked config downgrades", func(t *testing.T) {
		root := newRepo(t, map[string]string{
			"docs/specs/001/spec.md":   typedDoc("spec", "Approved", "docs/architecture/missing.md"),
			"docs/architecture/003.md": typedDoc("architecture", "Approved"),
		}, false)
		writeFile(t, root, ".mdstandards.json", `{"types":["spec","architecture"]}`)
		code, out, errs := runPolicyFromHead(t, root)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (dirty downgrade); stderr=%s out=%s", code, errs, out)
		}
		if !strings.Contains(errs, "differs from HEAD") {
			t.Fatalf("stderr = %q, want dirty warning", errs)
		}
	})

	t.Run("ignored config downgrades", func(t *testing.T) {
		root := newRepo(t, map[string]string{
			"docs/specs/001/spec.md":   typedDoc("spec", "Approved", "docs/architecture/missing.md"),
			"docs/architecture/003.md": typedDoc("architecture", "Approved"),
			".gitignore":               ".mdstandards.json\n",
		}, false)
		writeFile(t, root, ".mdstandards.json", `{"types":["spec","architecture"]}`)
		code, out, errs := runPolicyFromHead(t, root)
		if code != 0 {
			t.Fatalf("exit = %d, want 0 (dirty downgrade); stderr=%s out=%s", code, errs, out)
		}
		if !strings.Contains(errs, "differs from HEAD") {
			t.Fatalf("stderr = %q, want dirty warning", errs)
		}
	})

	t.Run("absent config is clean legacy", func(t *testing.T) {
		root := newRepo(t, map[string]string{
			"docs/specs/001/spec.md":   typedDoc("spec", "Approved", "docs/architecture/003.md"),
			"docs/architecture/003.md": typedDoc("architecture", "Approved"),
		}, false)
		code, out, errs := runPolicyFromHead(t, root)
		if code != 0 {
			t.Fatalf("exit = %d; stderr=%s out=%s", code, errs, out)
		}
		if strings.Contains(errs, "differs from HEAD") {
			t.Fatalf("unexpected dirty warning: %s", errs)
		}
	})
}
