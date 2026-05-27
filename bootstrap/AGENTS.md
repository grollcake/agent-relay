# Agent Instructions

<agent-relay-rules>

## Agent Relay

This project follows Agent Relay. See `.agent-relay/PROTOCOL.md`.

1. When joining or resuming work, read and follow `.agent-relay/PROTOCOL.md`.
2. Before recordable work, read `.agent-relay/GUIDANCE.md`, `.agent-relay/LESSON-LEARNED.md`, and existing records under `.agent-relay/lesson-learned/`.
3. At session start, ask whether to use a Git branch strategy for this Agent Relay session: always use branches, do not use branches, or ask per task.
4. Use only the `LeadAI`, `PlanAI`, and `ExecAI` roles, and route PlanAI/ExecAI communication through LeadAI; LeadAI delegates in the background, remains immediately available for user requests, and must not block on polling or sleep while waiting for delegated work.
5. For `Standard` work, follow `REQUEST -> PLAN -> EXECUTE -> REVIEW -> CLOSE`; LeadAI must not append `CLOSE` or close the task until the user explicitly approves completion.
6. LeadAI appends every `.agent-relay/relay.log` event. PlanAI and ExecAI notify LeadAI when their artifacts are complete.
7. Use only these event/role pairs: `REQUEST|LeadAI`, `PLAN|PlanAI`, `EXECUTE|ExecAI`, `REVIEW|PlanAI`, `FEEDBACK|LeadAI`, `CLOSE|LeadAI`, `RUN_DONE|LeadAI`.
8. Keep `relay.log` append-only, never overwrite older `.agent-relay/runs/` rounds, update `GUIDANCE.md` and `lesson-learned/` only after user approval, and never store secrets or sensitive data in `.agent-relay/`.

</agent-relay-rules>
