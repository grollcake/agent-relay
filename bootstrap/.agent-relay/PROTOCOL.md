# Agent Relay Protocol

Minimum rules for installed projects. Background, bootstrap, and template
details belong in the repository README and guide.

## Roles

- `LeadAI`: user communication, classification, scope and risk decisions,
  delegation, result interpretation, final report, and every `relay.log` append.
  Not a passive relay. LeadAI delegates to PlanAI and ExecAI in the background
  and immediately returns control to the user with a short status. LeadAI remains
  available for new user requests while delegated work runs.
- `PlanAI`: writes `PLAN` and reviews `RUN`. A `REVIEW` is evidence for the
  user's decision, not approval; PlanAI never approves completion or writes `CLOSE`.
- `ExecAI`: implements `PLAN`, validates work, and writes `RUN`. Returns
  ambiguity to LeadAI without expanding scope.

PlanAI and ExecAI communicate only through LeadAI.

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
guidance, LeadAI completes this check before classifying it as `Direct` or
`Standard`. Before their delegated phases, PlanAI and ExecAI each repeat the
index selection for their own potentially expanded scope rather than relying
only on lessons selected by an earlier role.
Within one continuous session, do not reread `relay.log` before every message.

## Work Classes

- Excluded from records: simple Q&A, short explanation, or brainstorming.
- `Direct`: minor localized edit, Agent Relay bootstrap, or Agent Relay update
  sync. LeadAI records `REQUEST -> RUN_DONE`; no completion approval is required.
  Bootstrap and update are not excluded as meta work.
- `Standard`: multi-file work, design judgment, or work needing verification.
  After classification, follow the session Git branch strategy before appending
  `REQUEST`. If a dedicated task branch is used, keep the full workflow, records,
  artifacts, and implementation on that branch through approval, then commit the
  approved state and automatically merge it into the recorded base branch. If no
  task branch is used, keep the workflow on the current branch. Use
  PlanAI -> ExecAI -> PlanAI review, preferably in the background. LeadAI stays
  responsible for user replies during delegation. After each delegation, LeadAI
  must respond to the user instead of waiting for completion; do not block on
  polling, sleep, or delegated work.

At session start, LeadAI asks whether to use a Git branch strategy for this
Agent Relay session: always use branches, do not use branches, or ask per task.

## Event Timeline

`relay.log` is append-only. Keep summaries short and link artifacts with `path`.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

Use local system time in `YYYY-MM-DDTHH:MM:SS` format, four random lowercase letters for `task-id`, and
only `REQUEST`, `PLANNED`, `EXECUTED`, `REVIEW`, `FEEDBACK`, `CLOSE`, `RUN_DONE`.
LeadAI creates one `task-id` at `REQUEST` and reuses it through `CLOSE` for the same
Standard work, including any `FEEDBACK` before `CLOSE` approval. Use a new
`task-id` for each new `REQUEST`. LeadAI direct flow is
`REQUEST -> RUN_DONE`; Standard flow is
`REQUEST` -> `PLANNED` -> `EXECUTED` -> `REVIEW` -> `CLOSE`.
For Standard work with a dedicated task branch, LeadAI switches to that branch
before appending `REQUEST`; no task event is written to the base branch before
the approved task branch is merged. If the session branch strategy does not use
task branches, LeadAI appends `REQUEST` on the current branch.
`FEEDBACK` is recorded when the user reports feedback or defects instead of
approving `CLOSE`. Keep the same `task-id` and artifact key. After `FEEDBACK`,
LeadAI asks whether to add to the current run or start a new run. For obvious
defects, add to the current run without asking. **Add to current run**: retry
`EXECUTED` -> `REVIEW-<NN>` within the last `RUN-<NN>` scope and
existing `PLAN`; updating `RUN-<NN>.md` is allowed before `CLOSE` approval.
**New run**: proceed with `RUN-<NN+1>`. Pad `event` to 8 characters and
`role` to 6 characters with trailing spaces. Preserve older event lines
even if their format differs.

LeadAI appends every `relay.log` event. Use these fixed event and role pairs:

```text
REQUEST  | LeadAI
PLANNED  | PlanAI
EXECUTED | ExecAI
REVIEW   | PlanAI
FEEDBACK | LeadAI
CLOSE    | LeadAI
RUN_DONE | LeadAI
```

`EXECUTED` marks run completion. LeadAI appends it with required `path` after
ExecAI writes `RUN-<NN>.md`. Each round `<NN>` is one `EXECUTED` followed by the
matching `REVIEW`. On `blocker`, append another `EXECUTED` under the same
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
`CLOSE`; LeadAI writes it only after explicit user approval.
Each `RUN` records changed files, change summary, validation, and unresolved
risks. `<YYYYMMDD>` and `<HHMM>` come from the local system date and 24-hour minute at
`REQUEST` (no separators in `<HHMM>`). Example: `20260526-1430-diary-write`.
LeadAI chooses `<SLUG>` as lowercase kebab-case. `task-id` identifies log
events for one Standard work; `<YYYYMMDD>-<HHMM>-<SLUG>` identifies run
artifacts. Use one `task-id` and one artifact key together for the same
Standard work. Use the matching template in
`.agent-relay/templates/` for every round artifact.

Artifact creation and `relay.log` event append are separate required actions. A
stage is complete only after its artifact is written and LeadAI appends and
verifies its matching event. PlanAI and ExecAI notify LeadAI when their
artifacts are complete and provide a suggested event summary; they never append
`relay.log` or claim that an event has been appended. LeadAI appends the matching
event, preferably with `.agent-relay/protocol-guard append ...` instead of direct
shell redirection, then verifies the appended line before the next delegation.

LeadAI must run `.agent-relay/protocol-guard gate ...` or read the last 50
`relay.log` lines before delegating the next stage. If the required prior event
is missing, stop delegation and append or repair the LeadAI-owned log action.

Required gates:

| Next stage | Required prior event |
| --- | --- |
| Delegate ExecAI | `PLANNED` |
| Delegate PlanAI review | `EXECUTED` |
| Request user approval | `REVIEW` |
| Append final `CLOSE` event | explicit user approval |

## Standard Pipeline

1. LeadAI classifies the request without appending a Standard task event yet.
2. LeadAI applies the session Git branch strategy. If using a dedicated task
   branch, record the current base branch, create the task branch, and switch to
   it. Otherwise continue on the current branch.
3. LeadAI appends `REQUEST`.
4. PlanAI writes `PLAN` with a top `LeadAI Brief`; LeadAI appends `PLANNED`.
5. LeadAI reads only the `LeadAI Brief` by default, verifies it is complete,
   and delegates to ExecAI using its `ExecAI Prompt`.
6. ExecAI implements, validates, and writes `RUN-01`; LeadAI appends `EXECUTED`.
7. LeadAI verifies `EXECUTED`, then PlanAI reviews and writes `REVIEW-01`;
   LeadAI appends `REVIEW`.
8. LeadAI verifies `REVIEW`. If there is no `blocker`, the work is ready for a
   user decision; this is not approval.
9. LeadAI reports result, validation, actionable nits or risks if present, and
   the `REVIEW` path to the user, then requests approval.
10. After explicit user approval, LeadAI writes `CLOSE`, appends the `CLOSE`
    event, and commits
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

`REVIEW` is evidence for the user, not approval. Only the user may approve
`CLOSE`. `FEEDBACK` after any `REVIEW` round is a normal pipeline step, not an
exception or reversal of a completed task.

`Direct` work closes with `RUN_DONE` without completion approval. If it reveals
durable guidance or reusable lessons, still add only user-accepted updates.

## Context Refresh

LeadAI asks the user whether to replace a PlanAI or ExecAI instance if:

- five or more follow-up messages have accumulated in one instance;
- the task topic clearly changes;
- responses slow down or the instance confuses earlier context.

## Prompt And Report Contract

Delegation prompts pass explicit requirements, source and artifact paths, and
the event LeadAI will append after completion. PlanAI derives criteria,
validation, and risks in `PLAN`; ExecAI and review use `PLAN` without LeadAI
restating it. Delegates do not expand scope, guess through ambiguity, or write
`relay.log`; they notify LeadAI with completion and a suggested summary.

Reports should include artifact path, outcome, validation status, blockers or
risks, nits when applicable, requested event name and suggested summary for
LeadAI to append, and any user decision required. ExecAI reports must
also list items returned to LeadAI as out-of-scope. LeadAI keeps only artifact
paths and minimum decision data unless ambiguity requires more. PlanAI must put
a short `LeadAI Brief` at the top of each `PLAN` with goal, scope, success
criteria, risks, required checks, and a minimal `ExecAI Prompt`. LeadAI reads
only that brief by default before delegating to ExecAI. LeadAI may read the full
`PLAN` only when the brief is missing, incomplete, inconsistent, high-risk, or a
user decision requires detailed inspection.

## User-Facing Reports

Default to short user-facing reports. For `Direct` work, report the outcome,
key changed scope, and validation in one to three sentences. Do not list every
created or preserved file, narrate protocol steps, or add empty risk and next-step
sections unless the user asks or action is required.

For `Standard` work, an approval request should expose only the outcome,
validation status, actionable blockers or risks, and the `REVIEW` path.
Detailed changes and evidence remain in artifacts unless
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
preserve `GUIDANCE.md`, `LESSON-LEARNED.md`, `lesson-learned/`, `relay.log`,
`runs/`, and non-Agent-Relay instructions in `CLAUDE.md`. After a successful
update, LeadAI appends `REQUEST -> RUN_DONE` to `relay.log` with the
before/after `VERSION` in `summary`. Bootstrap and update are recording targets,
not excluded meta work.
