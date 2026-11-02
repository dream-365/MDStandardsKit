# Spike: tasks 的边界 — document_types 还是独立 taskscheck

Date: 2026-11-02
Mode: paper spike (read-only; no production code changed)
Status: RECOMMENDATION

## Question

Can `tasks.md` validation reuse the `document_types` seam that the idea type just added, instead of shipping a third binary?

## Evidence

1. **The seam is repo-scoped, not document-scoped.**
   `main.go:45` defines `repoConfig{Type string}`; `loadRepoConfig` (`main.go:261`) reads one `.mdstandards.json` and returns one type. `run` (`main.go:120-151`) selects exactly one `DocumentType` profile and one set of globs. `Corpus.DocumentType` (`extract.go:60`) is a corpus-level field.

2. **Tasks must coexist with spec in the same repo.**
   MDStandardsKit contains `docs/specs/**/spec.md` and would add `docs/specs/**/tasks.md` in the same tree. A repo-level `type: tasks` would make mdstandards govern only task files and stop governing spec files. The idea repo can be single-type because it contains only ideas; this repo cannot.

3. **Reusing the seam requires a multi-type redesign, not a free extension.**
   To carry both, `.mdstandards.json` would need a `types` list or a glob-to-type map; extraction would run per type; the corpus and report would need a merged shape; and the spec `-dump-corpus` bytes would have to stay stable. That is a change to mdstandards' config contract while the idea line is still landing.

4. **The graph work stays in Go under either boundary.**
   Task table parsing, cycle detection, transitive reachability, write-set conflicts and critical path are not Rego-shaped. The `document_types` path would still add a Go extractor for tasks and either Go-computed graph facts or new Rego over a new corpus shape.

## Options

| Option | Reuse | New mdstandards risk | New release surface |
|---|---|---|---|
| A separate `taskscheck` binary | go.mod, dispatcher and release pattern; duplicates frontmatter parsing, path normalization, report and exit scaffolding | none | one binary, four assets, shim entries |
| B `tasks` document type | discovery, extraction shell, report, exit codes | high: multi-type config, corpus merge, spec byte-stability | none |
| C multi-type config first, then tasks type | most | highest: redesign of `repoConfig` and extraction | none |

## Recommendation

**Option A — keep the separate `taskscheck` binary for v1.**

The `document_types` seam looked like a free reuse path, but it is repo-scoped and tasks must coexist with spec in the same repository. Option B is not an extension; it is a multi-type redesign of mdstandards while the idea line is still landing. Option A protects the stable spec contract and costs one binary plus duplicated parsing, which is bounded and testable.

Revisit Option C only when a third document type appears and the config is already being generalized.

## Not changed

No production code was modified for this spike. `make check` remains green. The first real `tasks.md` at `docs/specs/SPEC-002-idea-type/tasks.md` is the format test, not part of this boundary decision.
