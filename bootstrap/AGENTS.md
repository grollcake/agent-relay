# Agent Instructions

<agent-relay-rules>

## Agent Relay

This project follows Agent Relay. See `.agent-relay/PROTOCOL.md`.

1. When joining or resuming work, read and follow `.agent-relay/PROTOCOL.md`.
2. Before starting recordable work, read `.agent-relay/GUIDANCE.md`, `.agent-relay/LESSON-LEARNED.md`, and existing records under `.agent-relay/lesson-learned/`. Clearly excluded Q&A may be answered without this check.
3. Use the `Leader`, `Planner`, and `Runner` roles and keep Planner and Runner communication routed through the Leader.
4. Keep `.agent-relay/relay.log` append-only and store round artifacts in `.agent-relay/runs/` without overwriting older rounds.
5. For `Standard` work, create a task branch before recording `REQUEST`, keep all task records and changes on that branch, then commit and automatically merge it after approved `DONE`.
6. The Leader must not append the final `DONE` event or close the task until the user explicitly approves completion.
7. If a Planner or Runner instance has five or more follow-up messages, the task topic changes, or the instance slows down or confuses context, the Leader asks the user whether to replace that role instance.
8. Update `.agent-relay/GUIDANCE.md` and `.agent-relay/lesson-learned/` only after user approval.
9. Never store secrets, credentials, customer data, or sensitive operational information in `.agent-relay/`.

</agent-relay-rules>
