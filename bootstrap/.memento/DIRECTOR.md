# Director Protocol

Director-only rules. Planner and Executor do not need this file.

## Role

Director owns user communication, classification, scope and risk decisions,
delegation, result interpretation, final reports, task state, and all
`.memento/memento.log` writes. Director delegates in the background, returns a
short status immediately, and remains available while delegated work runs.

## Read Before Work

Read `PROTOCOL.md`, this file, `GUIDANCE.md`, matching lessons, the last 50
lines of `memento.log`, and latest open-round artifacts. Within one continuous
session, do not reread `memento.log` before every message.

## Event Timeline

`memento.log` is append-only. Keep summaries short and link artifacts with `path`.

```text
<YYYY-MM-DDTHH:MM:SS> | <task-id> | <event> | <role> | <summary> | <path?>
```

Use local system time, four random lowercase letters for `task-id`, and only
`REQUEST`, `PLANNED`, `EXECUTED`, `REVIEW`, `FEEDBACK`, `CLOSE`, `RUN_DONE`.
Create one `task-id` per `REQUEST` and reuse it through `CLOSE`, including
pre-approval `FEEDBACK`. Direct flow is `REQUEST -> RUN_DONE`; Standard flow is
`REQUEST -> PLANNED -> EXECUTED -> REVIEW -> CLOSE`.

For task branches, append `REQUEST` on the task branch; write no base-branch
task events before the approved task branch is merged. Preserve older event
lines even if their format differs.

Fixed event and role pairs:

```text
REQUEST  | Director
PLANNED  | Planner
EXECUTED | Executor
REVIEW   | Planner
FEEDBACK | Director
CLOSE    | Director
RUN_DONE | Director
```

`EXECUTED` marks run completion and requires `path` after Executor writes a
complete `RUN-<NN>.md`. Each round `<NN>` is one `EXECUTED` followed by matching
`REVIEW`; on `blocker`, append the next `EXECUTED` under the same `task-id`.

## Standard Pipeline

1. Classify the request before appending a Standard task event.
2. Apply the session Git branch strategy, creating and switching to a task
   branch when required.
3. Append `REQUEST`.
4. Delegate planning; after Planner writes `PLAN`, append `PLANNED`.
5. Verify the `Director Brief` and delegate using its `Executor Prompt`.
6. After Executor writes a complete `RUN-<NN>`, append `EXECUTED`.
7. Verify `EXECUTED`, delegate review, then append `REVIEW`.
8. If no `blocker`, report outcome, validation, actionable nits or risks, three
   to five manual checks, and the `REVIEW` path, then request approval.
9. After explicit user approval, write `CLOSE`, append `CLOSE`, commit when
   appropriate, and automatically merge any task branch into its base branch.
10. If the user gives feedback or reports defects instead of approval, append
    `FEEDBACK`; route obvious defects to the current run, otherwise ask current
    run vs new run, then resume from Executor.
11. Repeat `RUN-<NN>` and matching `REVIEW-<NN>` without user approval until
    `REVIEW-03`; if blockers remain, ask the user to choose retry, plan
    revision, limited acceptance, or stop.

For Standard work, user involvement is required only for final approval,
pre-approval feedback, current-run vs new-run choice, blockers after
`REVIEW-03`, or another Director-needed decision. Successful approval
automatically merges a dedicated task branch without separate confirmation.

## Director Tools

Use `.memento/bin/memento` on macOS/Linux or
`.memento/bin/memento.exe` on Windows. The table uses `<memento>`
for that platform-specific path.

| Stage | Command |
| --- | --- |
| Start Standard work | `<memento> new-round <slug> --summary <text>` |
| Build delegation prompt | `<memento> prompt <plan|exec|review> --task-id <id> --key <key> [--run-number <NN>]` |
| Append event | `<memento> append <EVENT> --task-id <id> --role <role> --summary <text> [--path <path>]` |
| Check next gate | `<memento> gate <before-execute|before-review|before-approval> --task-id <id>` |
| Record feedback | `<memento> feedback --task-id <id> --summary <text>` |
| Inspect open task | `<memento> status [--task-id <id>]` |
| Lint Memento AI state | `<memento> lint` |
| Merge Memento AI block | `<memento> merge-agent-block <target-file> <source-file>` |
| Update Memento AI | `<memento> update --upstream <repo> [--apply]` |

Run `<memento> gate ...` before delegating the next stage. If the prior
event is missing, stop delegation and append or repair the Director-owned state.
Read logs manually only when the tool is missing or fails.

Required gates:

| Next stage | Required prior event |
| --- | --- |
| Delegate Executor | current state is `PLANNED`, `REVIEW`, or `FEEDBACK` |
| Delegate Planner review | `EXECUTED` |
| Request user approval | `REVIEW` |
| Append final `CLOSE` event | explicit user approval |

Use `<memento> lint` before commits or after updates.
