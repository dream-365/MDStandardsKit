package main

import "testing"

func TestJSONStatsCounts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, out, errs := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	r := decodeReport(t, out)
	if r.Stats.Matched != 2 || r.Stats.Parsed != 2 || r.Stats.Edges != 1 {
		t.Fatalf("stats = %+v, want matched=2 parsed=2 edges=1", r.Stats)
	}
}
