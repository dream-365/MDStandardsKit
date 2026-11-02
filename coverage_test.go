package main

import (
	"strings"
	"testing"
)

func TestZeroParsedFailsClosedInPhase1(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "docs/specs/bad.md", "# no frontmatter\n")

	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "1")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "no governed document parsed") {
		t.Fatalf("stderr = %q, want zero-parsed coverage error", errs)
	}
}

func TestZeroGovernedDocumentsFailsClosedPhase0(t *testing.T) {
	root := t.TempDir()

	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "matched no files") {
		t.Fatalf("stderr = %q, want zero-governed coverage error", errs)
	}
}
