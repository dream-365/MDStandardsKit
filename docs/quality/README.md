# Quality reports

These are **smoke reports**: the pipeline ran end-to-end with a zero-level findings skeleton. They prove the tool, not the content quality of the target specs. A real independent reviewer pass is the next step; until then, treat the headline score as 0.0 by construction, not as a verdict.

- `SPEC-001-mdstandards/` — `docs/specs/SPEC-001-mdstandards/spec.md`
- `doc-writing-skill/` — `docs/designs/doc-writing-skill.md`
- `example-spec/` — `testdata/example-repo/docs/specs/example-spec.md`

The real review writes `quality-findings.json` / `quality-report.md` (without `.smoke`) using `reviewer.md` and `rubric.v1.yaml`.
