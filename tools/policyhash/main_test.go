package main

import (
	"os"
	"path/filepath"
	"testing"
)

// The golden value locks the wire format: two spaces between the file hash and
// the slash-separated, sorted relative path, one line per file, then one
// sha256 over the concatenation.
func TestHashDir(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	mustWrite(t, "policy/b.rego", "package b\n")
	mustWrite(t, "policy/data/standard.json", "{}\n")
	mustWrite(t, "policy/a.rego", "package a\n")

	got, err := hashDir("policy")
	if err != nil {
		t.Fatal(err)
	}
	const want = "8797ffb8448a8ac6e1225aa2b6bf751c9add144a0d80eb030af0907a26c3a104"
	if got != want {
		t.Fatalf("hashDir = %s, want %s", got, want)
	}
}

func TestHashDirEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.Mkdir("policy", 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := hashDir("policy"); err == nil {
		t.Fatal("hashDir on an empty directory succeeded")
	}
}

func mustWrite(t *testing.T, name, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
