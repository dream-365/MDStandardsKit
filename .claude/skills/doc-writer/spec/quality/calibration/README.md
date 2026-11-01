# Calibration

`gold.v1.json` is the candidate gold set. Each entry has:

- `target`: the reviewed document.
- `reference_findings`: the first independent reviewer pass.
- `candidate_findings`: an optional second independent pass for inter-rater agreement.
- `human_levels`: `item.id -> 0/1/2`, filled by a human adjudicator. Empty until adjudicated.
- `status`: `candidate` or `gold`.

Current state: 3 candidate entries. SPEC-001 has two independent passes; the other two are low-score anchors. Human adjudication and two more specs are needed before any threshold can be set.

## Run

```bash
qualitycheck -calibrate -rubric <rubric.v1.yaml> \
  -reference <first-findings.json> \
  -findings <second-findings.json> \
  -out calibration-report.md
```

The report gives exact agreement, quadratic weighted kappa, a 3×3 confusion matrix, per-dimension agreement, and the item-level disagreements. It measures inter-rater agreement, not human truth.

## Current result

SPEC-001 is now `gold`. Human adjudication on 2026-11-02 set the three disagreements to the strict reading: `bg.current_state=1`, `arch.topology=1`, `design.sequence=1`; the agreed levels were adopted for the other items.

- Inter-rater (first vs second): 38 compared, 92.1% exact, overall quadratic weighted kappa 0.623; core 29 compared, 96.6% exact, weighted kappa 0.651.
- Human vs first: 94.7% exact, weighted kappa 0.771; core 100% exact, weighted kappa 1.000.
- Human vs second: 97.4% exact, weighted kappa 0.894; core 96.6% exact, weighted kappa 0.651.

The other two entries remain `candidate`. Four more human-scored specs are needed to reach the five-spec promotion gate. Per-dimension kappa can read low at high exact agreement when a dimension has little level variance; read the exact agreement alongside it.

## Promotion criteria

From the design:

- At least 5 human-scored specs.
- Core items weighted kappa >= 0.6.
- Every disagreement adjudicated.
- Threshold = the 25th percentile of human-approved specs, never a fixed 80/100.

Do not enable a blocking threshold before these criteria hold.
