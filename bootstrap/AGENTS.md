# Agent Instructions

<agent-relay-rules>

## Agent Relay

This project follows Agent Relay. See `.agent-relay/PROTOCOL.md`.

1. When joining or resuming work, read and follow `.agent-relay/PROTOCOL.md`.
2. At the start of each recordable phase, the responsible role reads `.agent-relay/GUIDANCE.md` and the `.agent-relay/LESSON-LEARNED.md` index, then reads only lesson records whose `Applies When` or `Trigger / Symptom` matches its current scope; do not rely only on lessons selected by an earlier role.
3. At session start, ask whether to use a Git branch strategy for this Agent Relay session: always use branches, do not use branches, or ask per task.
4. Use only the `LeadAI`, `PlanAI`, and `ExecAI` roles, and route PlanAI/ExecAI communication through LeadAI; LeadAI delegates in the background, reads only the `LeadAI Brief` by default before ExecAI delegation, immediately returns control to the user with a short status, remains available for new user requests, and must not block on polling, sleep, or delegated work.
5. For `Standard` work, follow `REQUEST -> PLANNED -> EXECUTED -> REVIEW -> CLOSE`; `REVIEW` is evidence, not approval, only the user can approve `CLOSE`, and `FEEDBACK` after any review is a normal pipeline step.
6. LeadAI appends every `.agent-relay/relay.log` event. PlanAI and ExecAI notify LeadAI when their artifacts are complete.
7. Use only these event/role pairs: `REQUEST|LeadAI`, `PLANNED|PlanAI`, `EXECUTED|ExecAI`, `REVIEW|PlanAI`, `FEEDBACK|LeadAI`, `CLOSE|LeadAI`, `RUN_DONE|LeadAI`.
8. Keep `relay.log` append-only, never overwrite older `.agent-relay/runs/` rounds, update `GUIDANCE.md` and `lesson-learned/` only after user approval, and never store secrets or sensitive data in `.agent-relay/`.

</agent-relay-rules>
