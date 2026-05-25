# Agent Instructions

## Agent Relay

This project follows Agent Relay. See `.agent-relay/PROTOCOL.md`.

1. When joining or resuming work, read and follow `.agent-relay/PROTOCOL.md`.
2. Before starting recordable work, read `.agent-relay/GUIDANCE.md`, `.agent-relay/LESSON-LEARNED.md`, and existing records under `.agent-relay/lesson-learned/`. Clearly excluded Q&A may be answered without this check.
3. Use the `PM`, `Planner`, and `Executor` roles and keep Planner and Executor communication routed through the PM.
4. Keep `.agent-relay/relay.log` append-only and store round artifacts in `.agent-relay/runs/` without overwriting older rounds.
5. The PM must not append the final `DONE` event or close the task until the user explicitly approves completion.
6. If a Planner or Executor instance has five or more follow-up messages, the task topic changes, or the instance slows down or confuses context, the PM asks the user whether to replace that role instance.
7. Update `.agent-relay/GUIDANCE.md` and `.agent-relay/lesson-learned/` only after user approval.
8. Never store secrets, credentials, customer data, or sensitive operational information in `.agent-relay/`.
