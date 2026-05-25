# Agent Relay Protocol

Minimum rules for installed projects. Background, bootstrap, and template
details belong in the repository README and guide.

## Roles

- `PM`: user communication, classification, scope and risk decisions,
  delegation, result interpretation, and final report. Not a passive relay.
- `Planner`: writes `PLAN`, reviews `RUN`, and writes `DONE` when accepted.
- `Executor`: implements `PLAN`, validates work, and writes `RUN`. Returns
  ambiguity to the PM without expanding scope.

Planner and Executor communicate only through the PM.

## Read Before Work

When joining or resuming, read `AGENTS.md`, this file, `GUIDANCE.md`,
`LESSON-LEARNED.md`, existing `lesson-learned/` records, the last 50 lines of
`relay.log`, and latest open-round artifacts if any.

Before starting recordable work, PM, Planner, and Executor must reread:

1. `.agent-relay/GUIDANCE.md`
2. `.agent-relay/LESSON-LEARNED.md`
3. existing records under `.agent-relay/lesson-learned/`

Clearly excluded requests may be answered without this check. If the request
requires file changes, investigation, design judgment, or project-specific
guidance, complete this check before classifying it as `Trivial` or `Standard`.
Within one continuous session, do not reread `relay.log` before every message.

## Work Classes

- Excluded from records: simple Q&A, short explanation, or brainstorming.
- `Trivial`: minor localized edit. PM records `REQUEST -> RUN_DONE`; no
  completion approval is required.
- `Standard`: multi-file work, design judgment, or work needing verification.
  Use Planner -> Executor -> Planner review, preferably in the background. PM
  stays responsible for user replies during delegation; do not poll or sleep to
  wait for completion.

## Event Timeline

`relay.log` is append-only. Keep summaries short and link artifacts with `path`.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

Use KST `YYYY-MM-DDTHH:MM:SS`, four random lowercase letters for `task-id`, and
only `REQUEST`, `PLAN`, `RUN_ST`, `RUN_ED`, `REVIEW`, `DONE`, `RUN_DONE`. PM
direct flow is `REQUEST -> RUN_DONE`; Standard flow is
`REQUEST -> PLAN -> RUN_ST -> RUN_ED -> REVIEW -> DONE`. Spaces around `role`
are for alignment only. Preserve older event lines even if their format differs.

`RUN_ST` marks run start. Executor appends it when starting `RUN-<NN>`. No
`path` required. `RUN_ED` marks run completion. Executor appends it with required
`path` after writing `RUN-<NN>.md`.

## Round Artifacts

Store all round artifacts in `.agent-relay/runs/` using one stable
`<YYYYMMDD>-<SLUG>` key:

- `.agent-relay/runs/<YYYYMMDD>-<SLUG>-PLAN.md`
- `.agent-relay/runs/<YYYYMMDD>-<SLUG>-RUN-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<SLUG>-REVIEW-<NN>.md`
- `.agent-relay/runs/<YYYYMMDD>-<SLUG>-DONE.md`

`<NN>` starts at `01`. Never overwrite an older round. Executor never writes
`DONE`; Planner writes it only when the matching review has no `blocker`.
Each `RUN` records changed files, change summary, validation, and unresolved
risks. PM chooses `<SLUG>` as lowercase kebab-case. `task-id` identifies log
events; `<SLUG>` identifies run artifacts. Use the matching template in
`.agent-relay/templates/` for every round artifact.

Artifact creation and `relay.log` event append are separate required actions. A
stage is complete only after both its artifact is written and its matching event
is appended. The artifact author appends the matching event: Planner appends
`PLAN` and `REVIEW`, Executor appends `RUN_ST` and `RUN_ED`, and PM appends
`REQUEST`, `RUN_DONE`, and the final `DONE` after user approval. PM verifies the
previous event exists before delegating the next stage.

## Standard Pipeline

1. PM classifies the request.
2. Planner writes `PLAN`.
3. PM delegates to Executor.
4. Executor appends `RUN_ST`, implements, validates, writes `RUN-01`, and appends
   `RUN_ED`.
5. Planner reviews and writes `REVIEW-01`.
6. If there is no `blocker`, Planner writes `DONE`.
7. PM reports result, nits, risks, and `DONE` path to the user.
8. After explicit user approval, PM appends the `DONE` event.
9. After `DONE` approval, PM may propose `.agent-relay/GUIDANCE.md` updates or
   `.agent-relay/lesson-learned/` additions. Add only items the user accepts.
10. If there is a `blocker`, Executor appends `RUN_ST` again and repeat
    `RUN-<NN>` and matching `REVIEW-<NN>`.
11. If blockers remain after `REVIEW-03`, PM asks the user to choose retry,
    plan revision, limited acceptance, or stop.

`Trivial` work closes with `RUN_DONE` without completion approval. If it reveals
durable guidance or reusable lessons, still add only user-accepted updates.

## Context Refresh

PM asks the user whether to replace a Planner or Executor instance if:

- five or more follow-up messages have accumulated in one instance;
- the task topic clearly changes;
- responses slow down or the instance confuses earlier context.

## Prompt And Report Contract

Planner and Executor prompts must include:

- goal, scope, success criteria, validation, and exact artifact path;
- input artifact paths;
- prohibition on out-of-scope work;
- instruction to return ambiguity to the PM instead of guessing.

Reports should include artifact path, outcome, validation status, blockers or
risks, nits when applicable, and any user decision required. Executor reports
must also list items returned to PM as out-of-scope. PM keeps only artifact
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
`runs/`, and non-Agent-Relay instructions in tool instruction files.
