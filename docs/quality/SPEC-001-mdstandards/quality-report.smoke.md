# Spec Quality Report

- Target: `docs/specs/SPEC-001-mdstandards/spec.md`
- Rubric: 1.0.0 (`acb806c85c7a357f1cdfeaea80a8dc951333aefe5b630fe8576e0e0ce7f25947`)
- Target SHA-256: `bb8d4081ca8b4738feeeea3a8c22353dabf47fdeb466803ef4593b8b4d0f1d56`
- Reviewer session: `reviewer-smoke`

**Headline score (core): 0.0**  |  Diagnostic score (all): 0.0

> The auditor verifies evidence presence, section placement, and score arithmetic. It does not judge whether evidence is sufficient or true.

> Normalization covers full-width ASCII, NBSP, and ideographic space; it is not full Unicode NFKC.

| Dimension | Headline | Diagnostic | Weight |
|---|---:|---:|---:|
| 背景落地 | 0.0 | 0.0 | 2 |
| 痛点证据 | 0.0 | 0.0 | 3 |
| 目标可验收 | 0.0 | 0.0 | 3 |
| 非目标边界 | 0.0 | 0.0 | 2 |
| 架构自洽 | 0.0 | 0.0 | 2 |
| 详细设计完整 | 0.0 | 0.0 | 2 |
| NFR 可度量 | 0.0 | 0.0 | 2 |
| 上线可回滚 | 0.0 | 0.0 | 2 |
| 替代方案诚实 | 0.0 | 0.0 | 2 |
| 跨章自洽 | 0.0 | 0.0 | 2 |
| 证据卫生 | 0.0 | 0.0 | 3 |

## 背景落地 (background)

- **bg.current_state** [core, level 0] — 1.1 是否写清了模块职责、核心流程与当前技术栈？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **bg.failure_chain** [core, level 0] — 1.2 是否有从触发到后果的分步失败链？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **bg.trigger** [core, level 0] — 1.3 是否解释了为什么现在做？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **bg.non_goals** [core, level 0] — 1.4 是否写了不做的后果与本次非目标？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 痛点证据 (pain_evidence)

- **bg.pain_evidence** [core, level 0] — 1.2 是否有至少一条可核实的证据（事故、度量或引述）？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 目标可验收 (goal_verifiability)

- **goal.verifiable** [core, level 0] — 2.1 每个 Goal 是否都有验收方式？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **goal.quantified** [core, level 0] — 2.1 是否至少有一个 Goal 有数字/单位/明确阈值？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **goal.count** [extended, level 0] — 2.1 Goals 是否控制在 3–6 个？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 非目标边界 (non_goal_boundary)

- **nongoal.labeled** [core, level 0] — 2.2 每条 Non-Goal 是否标了本次不做或 Deferred to Phase 2？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **nongoal.not_dumping_ground** [extended, level 0] — 2.2 非目标是否能追溯到痛点或边界？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 架构自洽 (architecture)

- **arch.solution_shape** [core, level 0] — 第 3 章是否有拓扑或等价的模块/数据流描述？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **arch.topology** [extended, level 0] — 3.1 拓扑图是否让每个模块有一句职责？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **arch.change_set** [extended, level 0] — 第 3 章是否写清新增/改造/保持不变？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **arch.rationale** [extended, level 0] — 关键选型是否有理由，详细对比指向第 7 章？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 详细设计完整 (detailed_design)

- **design.contract** [core, level 0] — 第 4 章是否至少有一份数据模型或接口契约？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **design.data_model** [extended, level 0] — 4.3 数据模型是否有字段类型、必填与校验？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **design.interface** [extended, level 0] — 4.4 接口是否有请求/响应/状态码或方法签名？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **design.sequence** [extended, level 0] — 4.5 复杂流程是否有时序图或等价分步说明？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **design.errors** [extended, level 0] — 4.6 每个错误码是否有触发场景与处理策略？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## NFR 可度量 (nfr)

- **nfr.measurable** [core, level 0] — 第 5 章每条 NFR 是否有数字、单位或验收方式？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **nfr.coverage** [core, level 0] — 第 5 章是否覆盖性能、可靠性、可观测性、安全，或说明不适用？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **nfr.honesty** [core, level 0] — 未实测的指标是否标了 [estimated] 或 未实现？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 上线可回滚 (rollout)

- **rollout.release** [core, level 0] — 6.1 是否有发布策略与灰度节奏？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **rollout.migration** [core, level 0] — 6.2 是否有迁移/兼容方案，或说明无需迁移？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **rollout.smoke** [core, level 0] — 6.3 是否有冒烟清单与观察指标？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **rollout.rollback_trigger** [core, level 0] — 6.4 回滚触发条件是否是可观测数字？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **rollout.rollback_steps** [core, level 0] — 6.4 回滚步骤是否可执行？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 替代方案诚实 (alternatives)

- **alts.status_quo** [core, level 0] — 第 7 章是否有 Status Quo 及其成本/风险？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **alts.rejected** [core, level 0] — 第 7 章是否至少有一个被否决方案，且淘汰原因对应当前约束？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **alts.matrix** [core, level 0] — 第 7 章是否有与约束对应的权衡矩阵？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **alts.consistency** [core, level 0] — 第 7 章结论是否与 Goals、架构选型一致？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 跨章自洽 (cross_consistency)

- **x.goals_trace** [core, level 0] — 后续章节是否逐条回应 Goals？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **x.nfr_arch** [core, level 0] — NFR 与架构、选型是否一致？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **x.rollback_release** [core, level 0] — 回滚与发布策略是否一致？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## 证据卫生 (evidence_hygiene)

- **x.no_placeholders** [core, level 0] — 是否没有把 无 / ... / 待补充 / TBD 当作实质回答？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **x.no_code** [core, level 0] — 是否没有堆砌实现代码、测试日志或偶发报错？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **x.numbers_sourced** [core, level 0] — 关键数字是否有出处，或标了 [estimated]？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run
- **x.no_fabricated_refs** [core, level 0] — 是否没有捏造 SPEC 引用、depends_on 目标或 M<N>？
  - rationale: smoke run; not a content verdict
  - missing: not reviewed in this smoke run

## Audit failures

None.
