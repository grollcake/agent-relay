Follow `AGENTS.md`.

This project uses Agent Relay.

Read `AGENTS.md` first. Then choose this session's agent name and use it consistently in relay entries.

Use `<agent>(<llm-or-version>[, <role>])`, for example `Codex(GPT-5.5)`.

Then read:

1. `README.md`
2. `.agent-relay/GUIDANCE.md`
3. `.agent-relay/PROTOCOL.md`
4. `.agent-relay/INDEX.md`
5. the last 50 lines of `.agent-relay/relay.log`
6. relevant handoff files, if referenced

After meaningful file, test, or issue changes, append one `TASK_DONE` entry to `.agent-relay/relay.log`.

When the user gives durable guidance, constraints, preferences, conventions, security rules, or "do not" rules, update `.agent-relay/GUIDANCE.md`.
