# Agent Relay Protocol

Minimum rules for installed projects. Background, bootstrap, and template
details belong in the repository README and guide.

## Roles

- `Leader`: user communication, classification, scope and risk decisions,
  delegation, result interpretation, and final report. Not a passive relay.
- `Planner`: writes `PLAN`, reviews `RUN`, and writes `DONE` when accepted.
- `Runner`: implements `PLAN`, validates work, and writes `RUN`. Returns
  ambiguity to the Leader without expanding scope.

Planner and Runner communicate only through the Leader.

## Read Before Work

When joining or resuming, read `AGENTS.md`, this file, `GUIDANCE.md`,
`LESSON-LEARNED.md`, existing `lesson-learned/` records, the last 50 lines of
`relay.log`, and latest open-round artifacts if any.

Before starting recordable work, Leader, Planner, and Runner must reread:

1. `.agent-relay/GUIDANCE.md`
2. `.agent-relay/LESSON-LEARNED.md`
3. existing records under `.agent-relay/lesson-learned/`

Clearly excluded requests may be answered without this check. If the request
requires file changes, investigation, design judgment, or project-specific
guidance, complete this check before classifying it as `Trivial` or `Standard`.
Within one continuous session, do not reread `relay.log` before every message.

## Work Classes

- Excluded from records: simple Q&A, short explanation, or brainstorming.
- `Trivial`: minor localized edit, Agent Relay bootstrap, or Agent Relay update
  sync. Leader records `REQUEST -> RUN_DONE`; no completion approval is required.
  Bootstrap and update are not excluded as meta work.
- `Standard`: multi-file work, design judgment, or work needing verification.
  Use Planner -> Runner -> Planner review, preferably in the background. Leader
  stays responsible for user replies during delegation; do not poll or sleep to
  wait for completion.

## Event Timeline

`relay.log` is append-only. Keep summaries short and link artifacts with `path`.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

Use KST `YYYY-MM-DDTHH:MM:SS`, four random lowercase letters for `task-id`, and
only `REQUEST`, `FEEDBACK`, `PLAN`, `RUN_ST`, `RUN_ED`, `REVIEW`, `DONE`, `RUN_DONE`. Leader
creates one `task-id` at `REQUEST` and reuses it through `DONE` for the same
Standard work, including any `FEEDBACK` before `DONE` approval. Use a new
`task-id` for each new `REQUEST`. Leader direct flow is
`REQUEST -> RUN_DONE`; Standard flow is
`REQUEST` -> `PLAN` -> `RUN_ST` -> `RUN_ED` -> `REVIEW` -> `DONE`.
`FEEDBACK` is recorded when the user reports feedback or defects instead of
approving `DONE`. Keep the same `task-id` and artifact key. After `FEEDBACK`,
Leader asks whether to add to the current run or start a new run. For obvious
defects, add to the current run without asking. **Add to current run**: retry
`RUN_ST` -> `RUN_ED` -> `REVIEW-<NN>` within the last `RUN-<NN>` scope and
existing `PLAN`; updating `RUN-<NN>.md` is allowed before `DONE` approval.
**New run**: proceed with `RUN-<NN+1>`. Pad `event` and `role`
to fixed width of 8 characters with trailing spaces. Preserve older event lines
even if their format differs.

`RUN_ST` marks run start. Leader appends it when delegating to Runner. No `path`
required; include `RUN-<NN>` in `summary`. `RUN_ED` marks run completion.
Runner appends it with required `path` after writing `RUN-<NN>.md`. Each round
`<NN>` is one `RUN_ST` paired with the next `RUN_ED`. On `blocker`, append
another `RUN_ST`/`RUN_ED` pair under the same `task-id`.

## Round Artifacts

Store all round artifacts in `.agent-relay/runs/` using one stable
`<YYYYMMDD>-<HHMM>-<SLUG>` key:

- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-PLAN.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-RUN-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-REVIEW-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-DONE.md`

`<NN>` starts at `01`. Never overwrite an older round, except when adding to
the current run after `FEEDBACK` before `DONE` approval; then updating
`RUN-<NN>.md` is allowed. Runner never writes
`DONE`; Planner writes it only when the matching review has no `blocker`.
Each `RUN` records changed files, change summary, validation, and unresolved
risks. `<YYYYMMDD>` and `<HHMM>` come from KST date and 24-hour minute at
`REQUEST` (no separators in `<HHMM>`). Example: `20260526-1430-diary-write`.
Leader chooses `<SLUG>` as lowercase kebab-case. `task-id` identifies log
events for one Standard work; `<YYYYMMDD>-<HHMM>-<SLUG>` identifies run
artifacts. Use one `task-id` and one artifact key together for the same
Standard work. Use the matching template in
`.agent-relay/templates/` for every round artifact.

Artifact creation and `relay.log` event append are separate required actions. A
stage is complete only after both its artifact is written and its matching event
is appended. The artifact author appends the matching event: Planner appends
`PLAN` and `REVIEW`, Runner appends `RUN_ED`, and Leader appends `REQUEST`,
`FEEDBACK`, `RUN_ST`, `RUN_DONE`, and the final `DONE` after user approval. Leader verifies the
previous event exists before delegating the next stage.

## Standard Pipeline

1. Leader classifies the request.
2. Planner writes `PLAN`.
3. Leader appends `RUN_ST` and delegates to Runner.
4. Runner implements, validates, writes `RUN-01`, and appends `RUN_ED`.
5. Planner reviews and writes `REVIEW-01`.
6. If there is no `blocker`, Planner writes `DONE`.
7. Leader reports result, nits, risks, and `DONE` path to the user.
8. After explicit user approval, Leader appends the `DONE` event.
9. If the user gives feedback or reports defects instead of approval, Leader
   appends `FEEDBACK`. For obvious defects, add to the current run; otherwise
   ask the user to choose current run or new run. Then resume from `RUN_ST`.
10. After `DONE` approval, Leader may propose `.agent-relay/GUIDANCE.md` updates or
   `.agent-relay/lesson-learned/` additions. Add only items the user accepts.
11. If there is a `blocker`, Leader appends `RUN_ST` again and repeat `RUN-<NN>` and
    matching `REVIEW-<NN>` without user approval until `REVIEW-03`.
12. If blockers remain after `REVIEW-03`, Leader asks the user to choose retry,
    plan revision, limited acceptance, or stop.

For Standard work, user involvement is required only for final `DONE` approval,
feedback or defects before `DONE` approval (`FEEDBACK`), current-run vs new-run
choice after `FEEDBACK`, remaining blockers after `REVIEW-03`, or when Leader
determines a user decision is needed. Leader continues `FEEDBACK` reruns and
intermediate `RUN`/`REVIEW` rounds without further user approval until the next
`DONE` approval request.

`Trivial` work closes with `RUN_DONE` without completion approval. If it reveals
durable guidance or reusable lessons, still add only user-accepted updates.

## Context Refresh

Leader asks the user whether to replace a Planner or Runner instance if:

- five or more follow-up messages have accumulated in one instance;
- the task topic clearly changes;
- responses slow down or the instance confuses earlier context.

## Prompt And Report Contract

Planner and Runner prompts must include:

- goal, scope, success criteria, validation, and exact artifact path;
- input artifact paths;
- prohibition on out-of-scope work;
- instruction to return ambiguity to the Leader instead of guessing.

Reports should include artifact path, outcome, validation status, blockers or
risks, nits when applicable, and any user decision required. Runner reports
must also list items returned to Leader as out-of-scope. Leader keeps only artifact
paths and minimum decision data unless ambiguity requires more.

## Guidance, Lessons, And Security

- `GUIDANCE.md`: durable instructions, constraints, preferences, conventions,
  security rules, and prohibitions only.
- `lesson-learned/`: reusable mistakes, solutions, and validation knowledge from
  completed work only. Use `templates/lesson-learned.md` and save records as
  `.agent-relay/lesson-learned/<YYYYMMDD>-<slug>.md`.
- Task progress stays in `relay.log` and round artifacts.
- Guidance and lesson updates require user acceptance.
- Never store secrets, credentials, customer data, personal information,
  sensitive internal information, or production secrets under `.agent-relay/`.

## Git And Updates

Commit `.agent-relay/` to Git. Do not add it to `.gitignore`. When updating,
preserve `GUIDANCE.md`, `LESSON-LEARNED.md`, `lesson-learned/`, `relay.log`,
`runs/`, and non-Agent-Relay instructions in tool instruction files. After a
successful update, Leader appends `REQUEST -> RUN_DONE` to `relay.log` with the
before/after `VERSION` in `summary`. Bootstrap and update are recording targets,
not excluded meta work.
