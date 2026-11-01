package main

import (
	"errors"
	"path"
	"strings"
)

var (
	errEmptyPath    = errors.New("empty path")
	errAbsolutePath = errors.New("absolute path")
	errEscapesRoot  = errors.New("path escapes repository root")
	errURL          = errors.New("url is not a document path")
)

// NormalizeRepoPath applies the bounded normalization agreed in D6 and D12:
// convert backslashes, collapse duplicate and trailing slashes, resolve "..",
// reject absolute paths and URLs, and reject anything that resolves outside
// the repository root. Case is preserved on every platform so macOS and Linux
// produce identical findings.
func NormalizeRepoPath(p string) (string, error) {
	if p == "" {
		return "", errEmptyPath
	}
	if strings.Contains(p, "://") {
		return "", errURL
	}

	q := strings.ReplaceAll(p, `\`, "/")
	clean := path.Clean(q)

	switch {
	case clean == "." || clean == "":
		return "", errEmptyPath
	case strings.HasPrefix(clean, "/"):
		return "", errAbsolutePath
	case len(clean) >= 2 && clean[1] == ':':
		return "", errAbsolutePath // Windows drive letter
	case clean == ".." || strings.HasPrefix(clean, "../"):
		return "", errEscapesRoot
	}
	return clean, nil
}
