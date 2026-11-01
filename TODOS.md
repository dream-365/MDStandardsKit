# TODOS

## Quality Scorecard

### Calibration set and threshold machinery

**What:** Build `gold.v1.json`, `-calibrate`, and the two threshold flags.

**Why:** Needed before Stage 1/2 can gate on the score; v1 is report-only.

**Context:** Deferred from the v1 scorecard (D1). The design's calibration section specifies 5 gold specs, core weighted kappa >= 0.6, and the 25th percentile of human-approved specs. Start by defining the gold schema and the kappa inputs, including `na` and missing ratings.

**Effort:** M
**Priority:** P2
**Depends on:** v1 rubric, reviewer, audit, and report shipped and run on real specs.

### Release tarball for the doc-writer skill

**What:** Build `doc-writer-skill.tar.gz` in the release workflow.

**Why:** The design originally promised it; the release workflow currently uploads only binaries.

**Context:** Dropped from v1 (D5). Manual copy is the v1 install path. Add a release step that packages `.claude/skills/doc-writer/`, `tools/qualitycheck/`, and the rubric.

**Effort:** S
**Priority:** P3
**Depends on:** v1 skill and tool stable.
