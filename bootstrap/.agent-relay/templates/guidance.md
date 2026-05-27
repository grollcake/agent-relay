# Guidance

Use this file for durable project guidance that should survive across agent
sessions.

Do not use this file for progress tracking, current task status, temporary
plans, or next steps. Use `relay.log` and `.agent-relay/runs/` artifacts for
that.

## Stable Project Context

- <long-lived project purpose or domain context>
- <important architectural or product context that is unlikely to change often>

## User Instructions

- <standing user preference or instruction>

## Constraints

- <hard technical, product, compatibility, or operational constraint>

## Do Not

- <thing future agents must not do>

## Security And Privacy

- <data, credential, privacy, or operational safety rule>

## Conventions

- When LeadAI delegates to PlanAI or ExecAI with an agent/subagent tool, run the
  delegated agent in the background if the tool supports it. In Claude Code,
  set `run_in_background: true` for every PlanAI/ExecAI delegation, including
  review delegation and follow-up SendMessage calls, so LeadAI can respond to
  the user while delegated work continues.
- Treat `REVIEW` as evidence for the user's decision, never as approval. Only
  explicit user approval allows LeadAI to write and append `CLOSE`; user
  `FEEDBACK` after any review is a normal pipeline step.
- PlanAI and ExecAI notify LeadAI when artifacts are complete and provide a
  suggested event summary; only LeadAI writes `relay.log`.
- Keep delegation prompts brief: pass explicit constraints and source/artifact
  paths, not copied context or derived criteria. PlanAI defines criteria and
  risks in `PLAN`; ExecAI and the reviewer follow `PLAN`.
