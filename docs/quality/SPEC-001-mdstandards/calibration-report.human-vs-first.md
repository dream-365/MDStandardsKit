# Calibration Report

- Reference: `docs/quality/SPEC-001-mdstandards/quality-findings.human.json`
- Candidate: `docs/quality/SPEC-001-mdstandards/quality-findings.json`
- Compared: 38
- Exact agreement: 94.7%
- Weighted kappa (quadratic): 0.771
- Core items: compared 29, exact 100.0%, weighted kappa 1.000

| ref\cand | 0 | 1 | 2 |
|---|---:|---:|---:|
| 0 | 0 | 0 | 0 |
| 1 | 0 | 4 | 2 |
| 2 | 0 | 0 | 32 |

## Dimensions

| Dimension | Compared | Exact | Agreement | Weighted kappa |
|---|---:|---:|---:|---:|
| 背景落地 | 4 | 4 | 100.0% | 1.000 |
| 痛点证据 | 1 | 1 | 100.0% | 1.000 |
| 目标可验收 | 3 | 3 | 100.0% | 1.000 |
| 非目标边界 | 2 | 2 | 100.0% | 1.000 |
| 架构自洽 | 4 | 3 | 75.0% | 0.500 |
| 详细设计完整 | 5 | 4 | 80.0% | 0.000 |
| NFR 可度量 | 3 | 3 | 100.0% | 1.000 |
| 上线可回滚 | 5 | 5 | 100.0% | 1.000 |
| 替代方案诚实 | 4 | 4 | 100.0% | 1.000 |
| 跨章自洽 | 3 | 3 | 100.0% | 1.000 |
| 证据卫生 | 4 | 4 | 100.0% | 1.000 |

## Disagreements

- `arch.topology`: reference 1, candidate 2
- `design.sequence`: reference 1, candidate 2

> Inter-rater agreement, not human truth. A high kappa means the instrument is reproducible, not that its levels are correct.
