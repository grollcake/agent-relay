# Agent Instructions

This project uses Agent Relay.

Agent Relay is a lightweight protocol for helping heterogeneous AI agents continue each other's work across sessions, models, and tools.

## Core Rules

1. Read context when joining or resuming work.
2. Record the result after meaningful work.
3. Create a handoff only when the next agent or role cannot continue from the issue and log alone.

## When Starting a New Session

Read `AGENTS.md` first. Then choose this session's agent name and use it consistently in relay entries.

Use `<agent>(<llm-or-version>[, <role>])`, for example `Codex(GPT-5.5)` or `Claude Code(Claude Sonnet 4.5, Reviewer)`.
Use `Human` only for work performed directly by a human.

Then read:

1. `README.md`
2. `.agent-relay/GUIDANCE.md`
3. `.agent-relay/PROTOCOL.md`
4. `.agent-relay/INDEX.md`
5. the last 50 lines of `.agent-relay/relay.log`
6. any handoff file referenced by the latest log, issue, or user request

Do not re-read the relay log before every user message in the same active session.

## During Work

Use the current conversation context during one continuous session.

Do not store secrets, tokens, passwords, customer data, or sensitive internal information in `.agent-relay/`.

Update `.agent-relay/GUIDANCE.md` when the user gives a durable instruction, constraint, preference, convention, security rule, or "do not" rule that future agents should keep following.

Do not update `GUIDANCE.md` for one-off task details, transient debugging notes, temporary plans, or work progress.

## After Meaningful Work

Append one result entry to `.agent-relay/relay.log`.

Meaningful work includes:

- code changes
- document changes
- configuration changes
- test execution
- architecture decisions
- issue updates
- work that another agent may need to continue later

For simple explanations, brainstorming, or questions, no relay log is required unless the user asks to preserve the context.

## Handoff

Create a handoff only when:

- another agent or role must continue the work;
- the result is partial, failed, or blocked;
- tests, review, or analysis must be performed by another role;
- important context is too long for a short relay log entry;
- the next agent would not be able to continue within 5 minutes from the issue and log alone.
