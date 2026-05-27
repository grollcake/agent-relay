# Agent Relay Protocol

Minimum rules for installed projects. Background, bootstrap, and template
details belong in the repository README and guide.

## Roles

- `Director`: user communication, classification, scope and risk decisions,
  delegation, result interpretation, final report, and every `relay.log` append.
  Not a passive relay. Director delegates to Planner and Executor in the background
  and immediately returns control to the user with a short status. Director remains
  available for new user requests while delegated work runs.
- `Planner`: writes `PLAN` and reviews `RUN`. A `REVIEW` is evidence for the
  user's decision, not approval; Planner never approves completion or writes `CLOSE`.
- `Executor`: implements `PLAN`, validates work, and writes `RUN`. Returns
  ambiguity to Director without expanding scope.

Planner and Executor communicate only through Director.

## Read Before Work

When joining or resuming, read the active instruction file (`AGENTS.md`,
`CLAUDE.md`, or both when present), this file, `GUIDANCE.md`,
the `LESSON-LEARNED.md` index, lesson records selected from that index for the
current resumed scope, the last 50 lines of `relay.log`, and latest open-round
artifacts if any. The
`<agent-relay-rules>...</agent-relay-rules>` block in `AGENTS.md` and
`CLAUDE.md` must remain identical so either file can stand alone.

At the start of each recordable phase, its responsible role must reread:

1. `.agent-relay/GUIDANCE.md`
2. the `.agent-relay/LESSON-LEARNED.md` index
3. only records under `.agent-relay/lesson-learned/` whose `Applies When` or
   `Trigger / Symptom` matches that phase's current scope

Clearly excluded requests may be answered without this check. If the request
requires file changes, investigation, design judgment, or project-specific
guidance, Director completes this check before classifying it as `Direct` or
`Standard`. Before their delegated phases, Planner and Executor each repeat the
index selection for their own potentially expanded scope rather than relying
only on lessons selected by an earlier role.
Within one continuous session, do not reread `relay.log` before every message.

## Work Classes

- Excluded from records: simple Q&A, short explanation, or brainstorming.
- `Direct`: minor localized edit, Agent Relay bootstrap, or Agent Relay update
  sync. Director records `REQUEST -> RUN_DONE`; no completion approval is required.
  Bootstrap and update are not excluded as meta work.
- `Standard`: multi-file work, design judgment, or work needing verification.
  After classification, follow the session Git branch strategy before appending
  `REQUEST`. If a dedicated task branch is used, keep the full workflow, records,
  artifacts, and implementation on that branch through approval, then commit the
  approved state and automatically merge it into the recorded base branch. If no
  task branch is used, keep the workflow on the current branch. Use
  Planner -> Executor -> Planner review, preferably in the background. Director stays
  responsible for user replies during delegation. After each delegation, Director
  must respond to the user instead of waiting for completion; do not block on
  polling, sleep, or delegated work.

At session start, Director asks whether to use a Git branch strategy for this
Agent Relay session: always use branches, do not use branches, or ask per task.

## Event Timeline

`relay.log` is append-only. Keep summaries short and link artifacts with `path`.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

Use local system time in `YYYY-MM-DDTHH:MM:SS` format, four random lowercase letters for `task-id`, and
only `REQUEST`, `PLANNED`, `EXECUTED`, `REVIEW`, `FEEDBACK`, `CLOSE`, `RUN_DONE`.
Director creates one `task-id` at `REQUEST` and reuses it through `CLOSE` for the same
Standard work, including any `FEEDBACK` before `CLOSE` approval. Use a new
`task-id` for each new `REQUEST`. Director direct flow is
`REQUEST -> RUN_DONE`; Standard flow is
`REQUEST` -> `PLANNED` -> `EXECUTED` -> `REVIEW` -> `CLOSE`.
For Standard work with a dedicated task branch, Director switches to that branch
before appending `REQUEST`; no task event is written to the base branch before
the approved task branch is merged. If the session branch strategy does not use
task branches, Director appends `REQUEST` on the current branch.
`FEEDBACK` is recorded when the user reports feedback or defects instead of
approving `CLOSE`. Keep the same `task-id` and artifact key. After `FEEDBACK`,
Director asks whether to add to the current run or start a new run. For obvious
defects, add to the current run without asking. **Add to current run**: retry
`EXECUTED` -> `REVIEW-<NN>` within the last `RUN-<NN>` scope and
existing `PLAN`; updating `RUN-<NN>.md` is allowed before `CLOSE` approval.
**New run**: proceed with `RUN-<NN+1>`. Pad `event` to 8 characters and
`role` to 8 characters with trailing spaces. Preserve older event lines
even if their format differs.

Director appends every `relay.log` event through `director-tool append` during
routine workflow. Use these fixed event and role pairs:

```text
REQUEST  | Director
PLANNED  | Planner
EXECUTED | Executor
REVIEW   | Planner
FEEDBACK | Director
CLOSE    | Director
RUN_DONE | Director
```

`EXECUTED` marks run completion. Director appends it with required `path` after
Executor writes `RUN-<NN>.md`. Each round `<NN>` is one `EXECUTED` followed by the
matching `REVIEW`. On `blocker`, append another `EXECUTED` under the same
`task-id` after Executor completes the next run.

## Round Artifacts

Store all round artifacts in `.agent-relay/runs/` using one stable
`<YYYYMMDD>-<HHMM>-<SLUG>` key:

- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-PLAN.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-RUN-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-REVIEW-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-CLOSE.md`

`<NN>` starts at `01`. Never overwrite an older round, except when adding to
the current run after `FEEDBACK` before `CLOSE` approval; then updating
`RUN-<NN>.md` is allowed. Executor never writes
`CLOSE`; Director writes it only after explicit user approval.
Executor may save `RUN-<NN>.md` as a checkpoint before long validation. A
checkpoint RUN with TODO fields is not complete; Director must not append
`EXECUTED` until TODOs are resolved.
Each `RUN` records changed files, change summary, validation, and unresolved
risks. `<YYYYMMDD>` and `<HHMM>` come from the local system date and 24-hour minute at
`REQUEST` (no separators in `<HHMM>`). Example: `20260526-1430-diary-write`.
Director chooses `<SLUG>` as lowercase kebab-case. `task-id` identifies log
events for one Standard work; `<YYYYMMDD>-<HHMM>-<SLUG>` identifies run
artifacts. Use one `task-id` and one artifact key together for the same
Standard work. Use the matching template in
`.agent-relay/templates/` for every round artifact.

Artifact creation and `relay.log` event append are separate required actions. A
stage is complete only after its artifact is written and Director appends and
verifies its matching event. Planner and Executor notify Director when their
artifacts are complete and provide a suggested event summary; they never append
`relay.log` or claim that an event has been appended. `.agent-relay/scripts/director-tool`
is Director-owned. Director uses it for routine `new-round`, `feedback`,
`append`, `gate`, `status`, and `subagent-prompt`; delegates must not use it to
mutate relay state. Read logs manually only when the tool is missing or fails.

Director must run `.agent-relay/scripts/director-tool gate ...` before delegating the
next stage. If the required prior event is missing, stop delegation and append or
repair the Director-owned log action.

Routine script commands:

| Stage | Command |
| --- | --- |
| Start Standard work | `.agent-relay/scripts/director-tool new-round <slug> --summary <text>` |
| Build delegation prompt | `.agent-relay/scripts/director-tool subagent-prompt <plan|exec|review> --task-id <id> --key <key> [--run-number <NN>]` |
| Append event | `.agent-relay/scripts/director-tool append <EVENT> --task-id <id> --role <role> --summary <text> [--path <path>]` |
| Check next gate | `.agent-relay/scripts/director-tool gate <before-execute|before-review|before-approval> --task-id <id>` |
| Record feedback | `.agent-relay/scripts/director-tool feedback --task-id <id> --summary <text>` |
| Inspect open task | `.agent-relay/scripts/director-tool status [--task-id <id>]` |
| Lint relay state | `.agent-relay/scripts/relay-lint` |
| Merge Agent Relay block | `.agent-relay/scripts/merge-agent-block <target-file> <source-file>` |
| Update Agent Relay | `.agent-relay/scripts/update-agent-relay --upstream <repo> [--apply]` |

Required gates:

| Next stage | Required prior event |
| --- | --- |
| Delegate Executor | `PLANNED` |
| Delegate Planner review | `EXECUTED` |
| Request user approval | `REVIEW` |
| Append final `CLOSE` event | explicit user approval |

## Standard Pipeline

1. Director classifies the request without appending a Standard task event yet.
2. Director applies the session Git branch strategy. If using a dedicated task
   branch, record the current base branch, create the task branch, and switch to
   it. Otherwise continue on the current branch.
3. Director appends `REQUEST`.
4. Planner writes `PLAN` with a top `Director Brief`; Director appends `PLANNED`.
5. Director reads only the `Director Brief` by default, verifies it is complete,
   and delegates to Executor using its `Executor Prompt`.
6. Executor implements, validates, and writes `RUN-01`; Director appends `EXECUTED`.
7. Director verifies `EXECUTED`, then Planner reviews and writes `REVIEW-01`;
   Director appends `REVIEW`.
8. Director verifies `REVIEW`. If there is no `blocker`, the work is ready for a
   user decision; this is not approval.
9. Director reports result, validation, actionable nits or risks if present,
   three to five concrete manual check cases for the user, and the `REVIEW`
   path to the user, then requests approval.
10. After explicit user approval, Director writes `CLOSE`, appends the `CLOSE`
    event, and commits
    the approved task state when appropriate. If a dedicated task branch was used,
    automatically merge it into its recorded base branch. If committing or merging
    cannot be completed cleanly, do not force it; report the blocker.
11. If the user gives feedback or reports defects instead of approval, Director
   appends `FEEDBACK`. For obvious defects, add to the current run; otherwise
   ask the user to choose current run or new run. Then resume from Executor work.
12. After `CLOSE` approval, Director may propose `.agent-relay/GUIDANCE.md` updates or
   `.agent-relay/lesson-learned/` additions. Add only items the user accepts.
13. If there is a `blocker`, Director delegates to Executor again, then repeats `RUN-<NN>` and
    matching `REVIEW-<NN>` without user approval until `REVIEW-03`.
14. If blockers remain after `REVIEW-03`, Director asks the user to choose retry,
    plan revision, limited acceptance, or stop.

For Standard work, user involvement is required only for final `CLOSE` approval,
feedback or defects before `CLOSE` approval (`FEEDBACK`), current-run vs new-run
choice after `FEEDBACK`, remaining blockers after `REVIEW-03`, or when Director
determines a user decision is needed. Director continues `FEEDBACK` reruns and
intermediate `RUN`/`REVIEW` rounds without further user approval until the next
`CLOSE` approval request. Successful approval includes automatic merge into the
recorded base branch only when a dedicated task branch was used; no separate
merge confirmation is requested in that case.

`REVIEW` is evidence for the user, not approval. Only the user may approve
`CLOSE`. `FEEDBACK` after any `REVIEW` round is a normal pipeline step, not an
exception or reversal of a completed task.

`Direct` work closes with `RUN_DONE` without completion approval. If it reveals
durable guidance or reusable lessons, still add only user-accepted updates.

## Context Refresh

Director asks the user whether to replace a Planner or Executor instance if:

- five or more follow-up messages have accumulated in one instance;
- the task topic clearly changes;
- responses slow down or the instance confuses earlier context.

## Prompt And Report Contract

Director builds delegation prompts from `director-tool subagent-prompt` plus
only the explicit task requirements. Planner derives criteria, validation, and
risks in `PLAN`; Executor and review use `PLAN` without Director restating it.
Delegates do not expand scope, guess through ambiguity, run relay-state-changing
`director-tool` commands, or write `relay.log`; they notify Director with
completion and a suggested summary.

Reports should include artifact path, outcome, validation status, blockers or
risks, nits when applicable, three to five suggested manual check cases for the
user when reporting a `REVIEW`, requested event name and suggested summary for
Director to append, and any user decision required. Executor reports must
also list items returned to Director as out-of-scope. Director keeps only artifact
paths and minimum decision data unless ambiguity requires more. Planner must put
a short `Director Brief` at the top of each `PLAN` with goal, scope, success
criteria, risks, required checks, and a minimal `Executor Prompt`. Director reads
only that brief by default before delegating to Executor. Director may read the full
`PLAN` only when the brief is missing, incomplete, inconsistent, high-risk, or a
user decision requires detailed inspection.

## User-Facing Reports

Default to short user-facing reports. For `Direct` work, report the outcome,
key changed scope, and validation in one to three sentences. Do not list every
created or preserved file, narrate protocol steps, or add empty risk and next-step
sections unless the user asks or action is required.

For `Standard` work, an approval request should expose only the outcome,
validation status, actionable nits or risks, three to five user manual checks,
and the `REVIEW` path. Detailed changes and evidence remain in artifacts unless
the user requests them.

## Guidance, Lessons, And Security

- `GUIDANCE.md`: durable instructions, constraints, preferences, conventions,
  security rules, and prohibitions only.
- `lesson-learned/`: reusable mistakes, solutions, and validation knowledge from
  completed work only. Use `templates/lesson-learned.md` and save records as
  `.agent-relay/lesson-learned/<YYYYMMDD>-<trigger-or-symptom>.md`. Each record
  includes `Applies When` and `Trigger / Symptom`. Add each accepted record to
  `LESSON-LEARNED.md`, which is the searchable index of actual lesson records.
- Task progress stays in `relay.log` and round artifacts.
- Guidance and lesson updates require user acceptance.
- Never store secrets, credentials, customer data, personal information,
  sensitive internal information, or production secrets under `.agent-relay/`.

## Git And Updates

Commit `.agent-relay/` to Git. Do not add it to `.gitignore`. When updating,
read `.agent-relay/HOW-TO-UPDATE.md` first. Bootstrap and update are recording
targets, not excluded meta work.
