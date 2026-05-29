# Planner Protocol

Planner-only rules. Read `PROTOCOL.md`, this file, `GUIDANCE.md`, matching
lessons, and the artifacts Director provides.

## Planning

Write `.agent-relay/runs/<KEY>-PLAN.md` from `templates/plan.md`.

The `PLAN` starts with a short `Director Brief` containing goal, scope, success
criteria, risks, required checks, and a minimal `Executor Prompt`. Director reads
this brief by default before delegating to Executor, so keep it self-contained.

Define success criteria, validation, and risks in the `PLAN`. Do not expand
scope beyond Director's request. Return ambiguity to Director.

## Review

Review the `RUN` against the `PLAN` and write
`.agent-relay/runs/<KEY>-REVIEW-<NN>.md` from `templates/review.md`.

A review is evidence for the user's decision, not approval. Mark findings as:

- `blocker`: must be fixed before approval can be requested.
- `nit`: non-blocking issue or cleanup.

Include three to five concrete manual checks for the user.

## Report To Director

Report artifact path, outcome, validation status, blockers, nits, risks, manual
checks when reviewing, and a concise suggested summary. Do not close work or use
Director-owned tools.
