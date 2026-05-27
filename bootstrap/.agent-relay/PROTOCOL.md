# Agent Relay Protocol

Minimum rules for installed projects. Background, bootstrap, and template
details belong in the repository README and guide.

## Roles

- `LeadAI`: user communication, classification, scope and risk decisions,
  delegation, result interpretation, final report, and every `relay.log` append.
  Not a passive relay. LeadAI delegates to PlanAI and ExecAI in the background
  and remains immediately available for user requests while delegated work runs.
- `PlanAI`: writes `PLAN`, reviews `RUN`, and writes `CLOSE` artifacts when accepted.
- `ExecAI`: implements `PLAN`, validates work, and writes `RUN`. Returns
  ambiguity to LeadAI without expanding scope.

PlanAI and ExecAI communicate only through LeadAI.

## Read Before Work

When joining or resuming, read the active instruction file (`AGENTS.md`,
`CLAUDE.md`, or both when present), this file, `GUIDANCE.md`,
`LESSON-LEARNED.md`, existing `lesson-learned/` records, the last 50 lines of
`relay.log`, and latest open-round artifacts if any. The
`<agent-relay-rules>...</agent-relay-rules>` block in `AGENTS.md` and
`CLAUDE.md` must remain identical so either file can stand alone.

Before starting recordable work, LeadAI, PlanAI, and ExecAI must reread:

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
  sync. LeadAI records `REQUEST -> RUN_DONE`; no completion approval is required.
  Bootstrap and update are not excluded as meta work.
- `Standard`: multi-file work, design judgment, or work needing verification.
  After classification, follow the session Git branch strategy before appending
  `REQUEST`. If a dedicated task branch is used, keep the full workflow, records,
  artifacts, and implementation on that branch through approval, then commit the
  approved state and automatically merge it into the recorded base branch. If no
  task branch is used, keep the workflow on the current branch. Use
  PlanAI -> ExecAI -> PlanAI review, preferably in the background. LeadAI stays
  responsible for user replies during delegation; do not block on polling or
  sleep while waiting for delegated work.

At session start, LeadAI asks whether to use a Git branch strategy for this
Agent Relay session: always use branches, do not use branches, or ask per task.

## Event Timeline

`relay.log` is append-only. Keep summaries short and link artifacts with `path`.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

Use KST `YYYY-MM-DDTHH:MM:SS`, four random lowercase letters for `task-id`, and
only `REQUEST`, `PLAN`, `EXECUTE`, `REVIEW`, `FEEDBACK`, `CLOSE`, `RUN_DONE`.
LeadAI creates one `task-id` at `REQUEST` and reuses it through `CLOSE` for the same
Standard work, including any `FEEDBACK` before `CLOSE` approval. Use a new
`task-id` for each new `REQUEST`. LeadAI direct flow is
`REQUEST -> RUN_DONE`; Standard flow is
`REQUEST` -> `PLAN` -> `EXECUTE` -> `REVIEW` -> `CLOSE`.
For Standard work with a dedicated task branch, LeadAI switches to that branch
before appending `REQUEST`; no task event is written to the base branch before
the approved task branch is merged. If the session branch strategy does not use
task branches, LeadAI appends `REQUEST` on the current branch.
`FEEDBACK` is recorded when the user reports feedback or defects instead of
approving `CLOSE`. Keep the same `task-id` and artifact key. After `FEEDBACK`,
LeadAI asks whether to add to the current run or start a new run. For obvious
defects, add to the current run without asking. **Add to current run**: retry
`EXECUTE` -> `REVIEW-<NN>` within the last `RUN-<NN>` scope and
existing `PLAN`; updating `RUN-<NN>.md` is allowed before `CLOSE` approval.
**New run**: proceed with `RUN-<NN+1>`. Pad `event` to 8 characters and
`role` to 6 characters with trailing spaces. Preserve older event lines
even if their format differs.

LeadAI appends every `relay.log` event. Use these fixed event and role pairs:

```text
REQUEST  | LeadAI
PLAN     | PlanAI
EXECUTE  | ExecAI
REVIEW   | PlanAI
FEEDBACK | LeadAI
CLOSE    | LeadAI
RUN_DONE | LeadAI
```

`EXECUTE` marks run completion. LeadAI appends it with required `path` after
ExecAI writes `RUN-<NN>.md`. Each round `<NN>` is one `EXECUTE` followed by the
matching `REVIEW`. On `blocker`, append another `EXECUTE` under the same
`task-id` after ExecAI completes the next run.

## Round Artifacts

Store all round artifacts in `.agent-relay/runs/` using one stable
`<YYYYMMDD>-<HHMM>-<SLUG>` key:

- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-PLAN.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-RUN-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-REVIEW-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<HHMM>-<SLUG>-CLOSE.md`

`<NN>` starts at `01`. Never overwrite an older round, except when adding to
the current run after `FEEDBACK` before `CLOSE` approval; then updating
`RUN-<NN>.md` is allowed. ExecAI never writes
`CLOSE`; PlanAI writes it only when the matching review has no `blocker`.
Each `RUN` records changed files, change summary, validation, and unresolved
risks. `<YYYYMMDD>` and `<HHMM>` come from KST date and 24-hour minute at
`REQUEST` (no separators in `<HHMM>`). Example: `20260526-1430-diary-write`.
LeadAI chooses `<SLUG>` as lowercase kebab-case. `task-id` identifies log
events for one Standard work; `<YYYYMMDD>-<HHMM>-<SLUG>` identifies run
artifacts. Use one `task-id` and one artifact key together for the same
Standard work. Use the matching template in
`.agent-relay/templates/` for every round artifact.

Artifact creation and `relay.log` event append are separate required actions. A
stage is complete only after its artifact is written, its matching event is
appended, and the actor reports the appended event name plus the last matching
`relay.log` line. LeadAI appends the matching event, preferably with
`.agent-relay/protocol-guard append ...` instead of direct shell redirection.
PlanAI and ExecAI report artifact completion to LeadAI; LeadAI then appends
the corresponding event using the fixed role value for that event.

LeadAI must run `.agent-relay/protocol-guard gate ...` or read the last 50
`relay.log` lines before delegating the next stage. If the required prior event
is missing, stop delegation and instruct the responsible role to record it.

Required gates:

| Next stage | Required prior event |
| --- | --- |
| Delegate ExecAI | `PLAN` |
| Delegate PlanAI review | `EXECUTE` |
| Request user approval | `REVIEW` |
| Append final `CLOSE` event | explicit user approval |

## Standard Pipeline

1. LeadAI classifies the request without appending a Standard task event yet.
2. LeadAI applies the session Git branch strategy. If using a dedicated task
   branch, record the current base branch, create the task branch, and switch to
   it. Otherwise continue on the current branch.
3. LeadAI appends `REQUEST`.
4. PlanAI writes `PLAN`; LeadAI appends `PLAN`.
5. LeadAI verifies `PLAN` and delegates to ExecAI.
6. ExecAI implements, validates, and writes `RUN-01`; LeadAI appends `EXECUTE`.
7. LeadAI verifies `EXECUTE`, then PlanAI reviews and writes `REVIEW-01`;
   LeadAI appends `REVIEW`.
8. LeadAI verifies `REVIEW`. If there is no `blocker`, PlanAI writes `CLOSE`.
9. LeadAI reports result, validation, actionable nits or risks if present, and
   the `CLOSE` path to the user.
10. After explicit user approval, LeadAI appends the `CLOSE` event and commits
    the approved task state when appropriate. If a dedicated task branch was used,
    automatically merge it into its recorded base branch. If committing or merging
    cannot be completed cleanly, do not force it; report the blocker.
11. If the user gives feedback or reports defects instead of approval, LeadAI
   appends `FEEDBACK`. For obvious defects, add to the current run; otherwise
   ask the user to choose current run or new run. Then resume from ExecAI work.
12. After `CLOSE` approval, LeadAI may propose `.agent-relay/GUIDANCE.md` updates or
   `.agent-relay/lesson-learned/` additions. Add only items the user accepts.
13. If there is a `blocker`, LeadAI delegates to ExecAI again, then repeats `RUN-<NN>` and
    matching `REVIEW-<NN>` without user approval until `REVIEW-03`.
14. If blockers remain after `REVIEW-03`, LeadAI asks the user to choose retry,
    plan revision, limited acceptance, or stop.

For Standard work, user involvement is required only for final `CLOSE` approval,
feedback or defects before `CLOSE` approval (`FEEDBACK`), current-run vs new-run
choice after `FEEDBACK`, remaining blockers after `REVIEW-03`, or when LeadAI
determines a user decision is needed. LeadAI continues `FEEDBACK` reruns and
intermediate `RUN`/`REVIEW` rounds without further user approval until the next
`CLOSE` approval request. Successful approval includes automatic merge into the
recorded base branch only when a dedicated task branch was used; no separate
merge confirmation is requested in that case.

`Trivial` work closes with `RUN_DONE` without completion approval. If it reveals
durable guidance or reusable lessons, still add only user-accepted updates.

## Context Refresh

LeadAI asks the user whether to replace a PlanAI or ExecAI instance if:

- five or more follow-up messages have accumulated in one instance;
- the task topic clearly changes;
- responses slow down or the instance confuses earlier context.

## Prompt And Report Contract

PlanAI and ExecAI prompts must include:

- goal, scope, success criteria, validation, and exact artifact path;
- input artifact paths;
- matching event name to append;
- prohibition on out-of-scope work;
- instruction to return ambiguity to LeadAI instead of guessing;
- instruction to complete all three: write artifact, append event, and report
  the appended event name plus last matching `relay.log` line.

Reports should include artifact path, outcome, validation status, blockers or
risks, nits when applicable, appended event name, last matching `relay.log`
line for the task/event, and any user decision required. ExecAI reports must
also list items returned to LeadAI as out-of-scope. LeadAI keeps only artifact
paths and minimum decision data unless ambiguity requires more.

## User-Facing Reports

Default to short user-facing reports. For `Trivial` work, report the outcome,
key changed scope, and validation in one to three sentences. Do not list every
created or preserved file, narrate protocol steps, or add empty risk and next-step
sections unless the user asks or action is required.

For `Standard` work, a completion or approval report should expose only the
outcome, validation status, actionable blockers or risks, and the `CLOSE` path
when approval is needed. Detailed changes and evidence remain in artifacts unless
the user requests them.

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
`runs/`, and non-Agent-Relay instructions in `CLAUDE.md`. After a successful
update, LeadAI appends `REQUEST -> RUN_DONE` to `relay.log` with the
before/after `VERSION` in `summary`. Bootstrap and update are recording targets,
not excluded meta work.
