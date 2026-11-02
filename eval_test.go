package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestEvalPolicyHonorsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := evalPolicy(ctx, policyDir(t), Corpus{}); err == nil {
		t.Fatal("evalPolicy returned nil error for a canceled context")
	}
}

func TestEvalTimeoutExitsTwo(t *testing.T) {
	old := policyEvalTimeout
	policyEvalTimeout = -time.Second
	t.Cleanup(func() { policyEvalTimeout = old })

	root := t.TempDir()
	writeFile(t, root, "docs/specs/001.md", doc("Approved", "docs/architecture/003.md"))
	writeFile(t, root, "docs/architecture/003.md", doc("Published"))

	code, _, errs := runCLI(t, root, "--format", "json", "--phase", "0")
	if code != 2 {
		t.Fatalf("exit = %d, want 2; stderr=%s", code, errs)
	}
	if !strings.Contains(errs, "policy.load_error") {
		t.Fatalf("stderr = %q, want policy.load_error", errs)
	}
}
