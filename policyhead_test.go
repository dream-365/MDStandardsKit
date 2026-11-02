package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func copyPolicy(t *testing.T, root string) {
	t.Helper()
	if out, err := exec.Command("cp", "-R", "policy", filepath.Join(root, "policy")).CombinedOutput(); err != nil {
		t.Fatalf("copy policy: %v\n%s", err, out)
	}
}

func TestMaterializePolicyFromHEADReadsCommittedContent(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	gitRun(t, root, "init")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-m", "init")

	standardPath := filepath.Join(root, "policy", "data", "standard.json")
	committed, err := os.ReadFile(standardPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(standardPath, append(committed, []byte("\n")...), 0o644); err != nil {
		t.Fatal(err)
	}

	dir, dirty, cleanup, err := materializePolicyFromHEAD(root, filepath.Join(root, "policy"), "")
	if err != nil {
		t.Fatalf("materializePolicyFromHEAD: %v", err)
	}
	defer cleanup()
	if !dirty {
		t.Fatal("dirty = false, want true")
	}
	got, err := os.ReadFile(filepath.Join(dir, "data", "standard.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, committed) {
		t.Fatalf("materialized standard.json differs from HEAD:\ngot  %q\nwant %q", got, committed)
	}
}

func TestPolicyFromHeadDowngradesOnDirtyPolicy(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	gitRun(t, root, "init")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-m", "init")

	standardPath := filepath.Join(root, "policy", "data", "standard.json")
	f, err := os.OpenFile(standardPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	code, out, errs := runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"--phase", "1",
		"--format", "json",
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	if !strings.Contains(errs, "report-only") {
		t.Fatalf("stderr = %q, want report-only warning", errs)
	}
	r := decodeReport(t, out)
	if len(r.Findings) != 1 {
		t.Fatalf("findings = %+v, want the violation still reported", r.Findings)
	}
}

func TestPolicyFromHeadDowngradesOnBaselineOutsideRoot(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/missing.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	gitRun(t, root, "init")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-m", "init")

	outside := filepath.Join(t.TempDir(), "baseline.json")
	if err := os.WriteFile(outside, []byte("[]"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, out, errs := runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"--phase", "1",
		"--baseline", outside,
		"--format", "json",
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s out=%s", code, errs, out)
	}
	if !strings.Contains(errs, "report-only") {
		t.Fatalf("stderr = %q, want report-only warning", errs)
	}
}

func TestDirtyPolicyWithExtractErrorIsReportOnly(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	writeFile(t, root, "docs/specs/bad.md", "# no frontmatter\n")
	gitRun(t, root, "init")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-m", "init")

	standardPath := filepath.Join(root, "policy", "data", "standard.json")
	f, err := os.OpenFile(standardPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	code, out, errs := runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"--phase", "1",
		"--format", "json",
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stderr=%s", code, errs)
	}
	r := decodeReport(t, out)
	if len(r.ExtractErrors) != 1 {
		t.Fatalf("extract errors = %+v, want 1", r.ExtractErrors)
	}
}

func TestDirtyPolicyWithZeroEdgesStillFailsClosed(t *testing.T) {
	root := t.TempDir()
	copyPolicy(t, root)
	writeFile(t, root, "docs/specs/001.md", doc("Approved"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))
	gitRun(t, root, "init")
	gitRun(t, root, "add", ".")
	gitRun(t, root, "commit", "-m", "init")

	standardPath := filepath.Join(root, "policy", "data", "standard.json")
	f, err := os.OpenFile(standardPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("\n"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	code, _, errs := runCLI(t, root,
		"--policy", filepath.Join(root, "policy"),
		"--policy-from-head",
		"--phase", "1",
		"--format", "json",
	)
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "zero dependency edges") {
		t.Fatalf("stderr = %q, want zero-edge coverage error", errs)
	}
}

func TestPolicyFromHeadOutsideRootFailsClosed(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	code, _, errs := runCLI(t, root,
		"--policy", outside,
		"--policy-from-head",
		"--phase", "0",
	)
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "must be inside --root") {
		t.Fatalf("stderr = %q, want inside-root error", errs)
	}
}
