# TODOS

## Quality Scorecard

### Calibration set and threshold machinery

**What:** Build `gold.v1.json`, `-calibrate`, and the two threshold flags.

**Why:** Needed before Stage 1/2 can gate on the score; v1 is report-only.

**Context:** Deferred from the v1 scorecard (D1). The scaffold exists at `.claude/skills/doc-writer/spec/quality/calibration/gold.v1.json`: SPEC-001 is gold with human adjudication (human-vs-first core kappa 1.000; human-vs-second core kappa 0.651), and two low-score anchors remain candidates. Four more human-scored specs and the blocking-threshold decision remain.

**Effort:** M
**Priority:** P2
**Depends on:** Ready to continue: v1 shipped, real reviews and one human adjudication exist.

## Completed

### Release tarball for the doc-writer skill

**What:** Build `doc-writer-skill.tar.gz` in the release workflow.

**Why:** The design originally promised it; the release workflow previously uploaded only binaries.

**Context:** The release job now packages `.claude/skills/doc-writer/`, `tools/qualitycheck/`, `policy/`, `go.mod`, and `go.sum` on every `v*` tag.

**Completed:** 2026-11-02 (no version source)
