# Agent Instructions

<agent-relay-rules>

## Agent Relay

This project follows Agent Relay. Read `.agent-relay/PROTOCOL.md`, then read
only your role file:

- Director: `.agent-relay/DIRECTOR.md`
- Planner: `.agent-relay/PLANNER.md`
- Executor: `.agent-relay/EXECUTOR.md`

1. At the start of each recordable phase, read `.agent-relay/GUIDANCE.md`, the
   `.agent-relay/LESSON-LEARNED.md` index, and only matching lesson records.
2. At session start, Director asks whether to use a Git branch strategy: always
   use branches, do not use branches, or ask per task.
3. Route Planner and Executor through Director. Director delegates in the
   background when possible, returns a short status immediately, and remains
   available to the user while delegated work runs.
4. `REVIEW` is evidence, not approval; only explicit user approval can close
   Standard work, and approval requests include 3-5 user manual checks.
5. Director passes user-reported defects to Executor; Executor records evidence
   before fixing and a self smoke test after fixing.
6. Keep older `.agent-relay/runs/` rounds append-only; update `GUIDANCE.md` and
   `lesson-learned/` only after user approval; never store secrets or sensitive
   data in `.agent-relay/`.
7. Agent Relay scripts require POSIX `sh`. On Windows, run them in Git for
   Windows Git Bash; native PowerShell and cmd are not supported.

</agent-relay-rules>
