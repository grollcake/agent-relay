# Agent Instructions

<memento-rules>

## Memento AI

This project follows Memento AI. Read `.memento/PROTOCOL.md`, then read
only your role file:

- Director: `.memento/DIRECTOR.md`
- Planner: `.memento/PLANNER.md`
- Executor: `.memento/EXECUTOR.md`

1. At the start of each recordable phase, read `.memento/GUIDANCE.md`, the
   `.memento/LESSON-LEARNED.md` index, and only matching lesson records.
2. At every session start, Director detects Codex or Claude Code, runs
   `memento models get <platform>` and `memento models list <platform>`, and
   asks the user to choose Director, Planner, and Executor models, showing the
   previous choices as defaults. Save choices with `memento models set`.
3. Director then asks whether to use a Git branch strategy: always use
   branches, do not use branches, or ask per task.
4. Route Planner and Executor through Director. Director delegates in the
   background when possible, returns a short status immediately, and remains
   available to the user while delegated work runs.
5. `REVIEW` is evidence, not approval; only explicit user approval can close
   Standard work, and approval requests include 3-5 user manual checks.
6. Director passes user-reported defects to Executor; Executor records evidence
   before fixing and a self smoke test after fixing.
7. Keep older `.memento/runs/` rounds append-only; update `GUIDANCE.md` and
   `lesson-learned/` only after user approval; never store secrets or sensitive
   data in `.memento/`.
8. Memento AI uses the native `.memento/bin/memento[.exe]` Go binary
   and does not require a specific shell. Git-integrated commands require
   `git` to be available on `PATH`.

</memento-rules>
