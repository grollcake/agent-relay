Follow `AGENTS.md`.

This project uses Agent Relay.

Read `AGENTS.md` first. Then choose this session's agent name and use it consistently in relay entries.

Use `<agent>(<llm-or-version>[, <role>])`, for example `Claude Code(Claude Sonnet 4.5)`.

Then read:

1. `README.md`
2. `.agent-relay/GUIDANCE.md`
3. `.agent-relay/PROTOCOL.md`
4. `.agent-relay/INDEX.md`
5. the last 50 lines of `.agent-relay/relay.log`
6. relevant handoff files, if referenced

Do not re-read relay context before every message in the same active session.

After meaningful work, append a short result entry to `.agent-relay/relay.log`.

When the user gives durable guidance, constraints, preferences, conventions, security rules, or "do not" rules, update `.agent-relay/GUIDANCE.md`.

Create a handoff only when the next agent cannot continue from the issue and relay log alone.
