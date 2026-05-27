# Agent Instructions

<agent-relay-rules>

## Agent Relay

This project follows Agent Relay. See `.agent-relay/PROTOCOL.md`.

1. When joining or resuming work, read and follow `.agent-relay/PROTOCOL.md`.
2. At the start of each recordable phase, the responsible role reads `.agent-relay/GUIDANCE.md` and the `.agent-relay/LESSON-LEARNED.md` index, then reads only lesson records whose `Applies When` or `Trigger / Symptom` matches its current scope; do not rely only on lessons selected by an earlier role.
3. At session start, ask whether to use a Git branch strategy for this Agent Relay session: always use branches, do not use branches, or ask per task.
4. Route Planner/Executor through Director. For `Standard` work, Director promptly delegates each phase in the background using only explicit constraints and source paths, reads only the `Director Brief` by default before Executor delegation, returns a short status immediately, and never blocks on delegated work, polling, or sleep.
5. For `Standard` work, follow `REQUEST -> PLANNED -> EXECUTED -> REVIEW -> CLOSE`; `REVIEW` is evidence, not approval, only the user can approve `CLOSE`, and approval requests include 3-5 user manual checks.
6. Director appends every `.agent-relay/relay.log` event. Planner and Executor notify Director when their artifacts are complete.
7. `.agent-relay/scripts/director-tool` is Director-owned; delegates must not use it to mutate relay state. Director uses it for routine task creation, event append, gate checks, status, and subagent prompts. Use `relay-lint` before commits or after updates. Use only these event/role pairs: `REQUEST|Director`, `PLANNED|Planner`, `EXECUTED|Executor`, `REVIEW|Planner`, `FEEDBACK|Director`, `CLOSE|Director`, `RUN_DONE|Director`.
8. Keep `relay.log` append-only, never overwrite older `.agent-relay/runs/` rounds, update `GUIDANCE.md` and `lesson-learned/` only after user approval, and never store secrets or sensitive data in `.agent-relay/`.

</agent-relay-rules>
