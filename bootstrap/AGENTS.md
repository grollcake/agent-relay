# Agent Instructions

## Agent Relay

This project follows Agent Relay. See `.agent-relay/PROTOCOL.md`.

1. When joining or resuming work, follow the Agent Relay read order in `.agent-relay/PROTOCOL.md`.
2. Choose one session agent name and use it consistently in `relay.log`.
3. For meaningful work, append `TASK_BEGIN` to `.agent-relay/relay.log` before starting and `TASK_DONE` after finishing.
4. Create a handoff only when the next agent cannot continue from the issue and relay log alone within 5 minutes.
5. When creating a handoff, put long context in a handoff file and link it from `relay.log` with `path=`.
6. Update `.agent-relay/GUIDANCE.md` only for durable user instructions, constraints, preferences, conventions, security rules, or "do not" rules.
7. Never store secrets, credentials, customer data, or sensitive operational information in `.agent-relay/`.
