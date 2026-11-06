// Command policyhash prints the policy bundle hash used by --version.
//
// The hash is defined by the distribution design: every regular file under the
// policy directory contributes "<file sha256>  <repo-relative path>\n" (two
// spaces), files are ordered by their slash-separated relative path, and the
// sha256 of that byte stream is the policy hash. CI and the shim build fallback
// both call this command so the value cannot drift between build paths.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	dir := "policy"
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "usage: policyhash [dir]")
		os.Exit(2)
	}
	if len(os.Args) == 2 {
		dir = os.Args[1]
	}
	hash, err := hashDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "policyhash:", err)
		os.Exit(1)
	}
	fmt.Println(hash)
}

// hashDir computes the policy hash for files under rel, resolved from the
// current working directory. Paths are reported relative to the working
// directory, so the result does not depend on where the tree is checked out.
func hashDir(rel string) (string, error) {
	base, err := filepath.Abs(rel)
	if err != nil {
		return "", err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var names []string
	err = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		name, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}
		names = append(names, filepath.ToSlash(name))
		return nil
	})
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no regular files under %s", rel)
	}
	sort.Strings(names)

	h := sha256.New()
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(cwd, filepath.FromSlash(name)))
		if err != nil {
			return "", err
		}
		sum := sha256.Sum256(b)
		fmt.Fprintf(h, "%s  %s\n", hex.EncodeToString(sum[:]), name)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
