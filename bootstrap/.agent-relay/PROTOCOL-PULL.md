# Agent Relay Protocol

Agent Relay is a file-based collaboration protocol for a `PM`, `Planner`, and
`Executor` team. This file contains the minimum rules for installed projects.
Background, installation, and template details belong in the repository README
and guide.

## Roles

- `PM`: user communication, request classification, scope and risk decisions,
  delegation, result interpretation, and final reporting.
- `Planner`: writes `PLAN` artifacts and reviews implementation evidence.
- `Executor`: implements the plan, validates the work, and reports ambiguity to
  the PM without expanding scope.

Planner and Executor communicate only through the PM.

## Required Preflight

Before starting recordable work, PM, Planner, and Executor must read:

1. `.agent-relay/GUIDANCE.md`
2. `.agent-relay/LESSON-LEARNED.md`
3. existing records under `.agent-relay/lesson-learned/`

Clearly excluded requests may be answered without this check. If the request
requires file changes, investigation, design judgment, or project-specific
guidance, complete this check before classifying it as `Trivial` or `Standard`.

## Read Order

When joining, resuming, switching agents, or taking over an open round, read:

1. `AGENTS.md`
2. `.agent-relay/PROTOCOL.md`
3. `.agent-relay/GUIDANCE.md`
4. `.agent-relay/LESSON-LEARNED.md`
5. existing records under `.agent-relay/lesson-learned/`
6. the last 50 lines of `.agent-relay/relay.log`
7. latest open-round `PLAN`, `RUN`, and `REVIEW` artifacts, if any

Within one continuous session, do not reread `relay.log` before every message.
Still rerun the required preflight before recordable work.

## Work Classes

| Class | Scope | Handling |
| --- | --- | --- |
| Excluded from records | Simple Q&A, short explanation, brainstorming | Respond without an event |
| `Trivial` | Minor text/config edit or obvious localized edit | PM handles directly and records `REQUEST -> RUN_DONE`; no completion approval is required |
| `Standard` | Multi-file implementation, design judgment, or work needing implementation verification | Planner -> Executor -> Planner review |

Delegate `Standard` work in the background when the active tool supports it.
The PM remains responsible for user communication. Do not poll or sleep only to
wait for completion.

## Event Timeline

`relay.log` is append-only:

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

- `timestamp`: KST, `YYYY-MM-DDTHH:MM:SS`
- `task-id`: four random lowercase ASCII letters
- events: `REQUEST`, `PLAN`, `RUN`, `REVIEW`, `DONE`, `RUN_DONE`
- PM direct flow: `REQUEST -> RUN_DONE`
- Standard flow: `REQUEST -> PLAN -> RUN -> REVIEW -> DONE`
- spaces around `role` are for alignment only

Keep summaries short. Put implementation details and review findings in round
artifacts and link them with `path`. Preserve old event lines even if they used
an older format.

## Round Artifacts

All round artifacts live in `.agent-relay/runs/` and use one stable
`<YYYYMMDD>-<SLUG>` key per task.

| Artifact | Path | Author |
| --- | --- | --- |
| Plan | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-PLAN.md` | Planner |
| Submission | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-RUN-<NN>.md` | Executor |
| Review | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-REVIEW-<NN>.md` | Planner |
| Acceptance | `.agent-relay/runs/<YYYYMMDD>-<SLUG>-DONE.md` | Planner |

- `<NN>` starts at `01`.
- Never overwrite an older round artifact.
- Executor never writes `DONE`.
- Planner writes `DONE` only when the matching review has no `blocker`.
- PM appends the final `DONE` event only after explicit user approval.
- Each `RUN` records changed files, summary, validation, and unresolved risks.
- PM chooses `<SLUG>` as lowercase kebab-case.
- `task-id` identifies `relay.log` events; `<SLUG>` identifies run artifacts.
- Use the matching template in `.agent-relay/templates/` for every round artifact.

## Standard Pipeline

1. PM classifies the request.
2. Planner writes `PLAN`.
3. Executor implements, validates, and writes `RUN-01`.
4. Planner reviews and writes `REVIEW-01`.
5. If there is no `blocker`, Planner writes `DONE`.
6. PM reports result, nits, residual risks, and `DONE` path to the user.
7. After explicit user approval, PM appends the `DONE` event.
8. After `DONE` approval, PM may propose `.agent-relay/GUIDANCE.md` updates or
   `.agent-relay/lesson-learned/` additions. Add only items the user accepts.
9. If there is a `blocker`, Executor writes the next `RUN` and Planner writes
   the matching `REVIEW`.
10. If blockers remain after `REVIEW-03`, PM asks the user to choose retry, plan
   revision, limited acceptance, or stop.

For `Trivial` work, PM records `REQUEST -> RUN_DONE` and may close the task
without completion approval. If the task reveals durable guidance or reusable
lessons, PM must still propose those updates and add only user-accepted items.

## Context Refresh

Planner and Executor should keep the same context when possible. PM must ask the
user whether to replace a role instance if:

- five or more follow-up messages have accumulated in one instance;
- the task topic clearly changes;
- responses become noticeably slow or the instance appears to confuse earlier
  context.

## Prompt And Report Contract

Planner and Executor prompts must include:

- goal;
- relevant files or investigation scope;
- artifact type and exact output path;
- success criteria and validation;
- prohibition on out-of-scope work;
- instruction to return ambiguity to the PM instead of guessing;
- required input artifact paths.

PM keeps only artifact paths and minimum decision data unless ambiguity requires
more: outcome, validation status, blocker count and summary, residual risk, and
required user decision.

Planner reports: artifact path, decision outcome, blocker summary, nit summary,
and whether a user decision is required.

Executor reports: `RUN` path, change summary, validation result, unresolved
risks, and out-of-scope items returned to the PM.

## Guidance, Lessons, And Security

Use `.agent-relay/GUIDANCE.md` only for durable instructions, constraints,
preferences, conventions, security rules, and prohibitions.

Use `.agent-relay/lesson-learned/` only for reusable mistakes, solutions, and
validation knowledge from completed work. Use
`.agent-relay/templates/lesson-learned.md` and save records as
`.agent-relay/lesson-learned/<YYYYMMDD>-<slug>.md`.

After `DONE` approval, PM may propose guidance or lesson updates. Add only items
the user accepts.

Do not store task progress, temporary plans, debugging notes, or one-off
requests in guidance or lesson records. Keep task records in `relay.log` and
round artifacts.

Never store API keys, tokens, passwords, private credentials, customer data,
personal information, sensitive internal information, or production secrets
anywhere under `.agent-relay/`.

## Git And Updates

Commit `.agent-relay/` to Git. Do not add it to `.gitignore`.

When updating Agent Relay, preserve project-specific state:

- `.agent-relay/GUIDANCE.md`
- `.agent-relay/LESSON-LEARNED.md`
- `.agent-relay/lesson-learned/`
- `.agent-relay/relay.log`
- `.agent-relay/runs/`
- non-Agent-Relay instructions in tool instruction files
