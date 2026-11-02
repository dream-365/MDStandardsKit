package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDependsOnDedupedPreservesFirstOrder(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved",
		"docs/architecture/b.md",
		"./docs/architecture/a.md",
		"docs/architecture/b.md/",
		"docs/architecture/../architecture/a.md",
	))
	writeFile(t, root, "docs/architecture/a.md", doc("Published"))
	writeFile(t, root, "docs/architecture/b.md", doc("Published"))

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "0", "--dump-corpus")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	var c Corpus
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatalf("invalid corpus JSON: %v\n%s", err, out)
	}
	if c.Stats.Edges != 2 {
		t.Fatalf("edges = %d, want 2 after dedupe", c.Stats.Edges)
	}
	for _, d := range c.Documents {
		if d.Path != "docs/specs/001.md" {
			continue
		}
		want := []string{"docs/architecture/b.md", "docs/architecture/a.md"}
		if len(d.DependsOn) != len(want) {
			t.Fatalf("depends_on = %v, want %v", d.DependsOn, want)
		}
		for i := range want {
			if d.DependsOn[i] != want[i] {
				t.Fatalf("depends_on = %v, want %v", d.DependsOn, want)
			}
		}
		return
	}
	t.Fatal("docs/specs/001.md not found in corpus")
}

func TestFileSizeLimitProducesExtractError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	writeFile(t, root, "docs/specs/big.md", "---\ntype: spec\nstatus: Approved\n---\n"+strings.Repeat("x", maxDocBytes+1))

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0 in phase 0", code)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 1 {
		t.Fatalf("extract errors = %+v, want 1", r.ExtractErrors)
	}
	if !strings.Contains(r.ExtractErrors[0].Reason, "limit is 1048576") {
		t.Fatalf("reason = %q, want size limit", r.ExtractErrors[0].Reason)
	}
}

func TestYAMLAliasGuardProducesExtractError(t *testing.T) {
	cases := map[string]string{
		"stars": strings.Repeat("*", 101),
		"merge": strings.Repeat("<<", 11),
	}
	for name, payload := range cases {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
			writeFile(t, root, "docs/architecture/003.md", doc("Published"))
			writeFile(t, root, "docs/specs/alias.md",
				"---\ntype: spec\nstatus: Approved\nvalue: \""+payload+"\"\n---\n# Body\n")

			code, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
			if code != 0 {
				t.Fatalf("exit = %d, want 0 in phase 0", code)
			}
			r := decodeReport(t, out)
			if len(r.ExtractErrors) != 1 {
				t.Fatalf("extract errors = %+v, want 1", r.ExtractErrors)
			}
			if !strings.Contains(r.ExtractErrors[0].Reason, "possible YAML alias expansion") {
				t.Fatalf("reason = %q, want alias guard", r.ExtractErrors[0].Reason)
			}
		})
	}
}

func TestInvalidYAMLProducesExtractError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	writeFile(t, root, "docs/specs/bad.md", "---\ntype: spec\nstatus: [unclosed\n---\n# Body\n")

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0 in phase 0", code)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 1 {
		t.Fatalf("extract errors = %+v, want 1", r.ExtractErrors)
	}
	if !strings.Contains(r.ExtractErrors[0].Reason, "invalid YAML") {
		t.Fatalf("reason = %q, want invalid YAML", r.ExtractErrors[0].Reason)
	}
}

func TestNonMappingFrontmatterProducesExtractError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	writeFile(t, root, "docs/specs/scalar.md", "---\nnull\n---\n# Body\n")

	code, out, _ := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0 in phase 0", code)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 1 {
		t.Fatalf("extract errors = %+v, want 1", r.ExtractErrors)
	}
	if !strings.Contains(r.ExtractErrors[0].Reason, "not a mapping") {
		t.Fatalf("reason = %q, want non-mapping error", r.ExtractErrors[0].Reason)
	}
}
