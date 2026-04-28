# Agent Relay Protocol

Agent Relay is a vendor-neutral, file-based relay protocol for heterogeneous AI coding agents.

It exists to help agents continue each other's work across:

- sessions
- models
- harnesses
- machines
- roles

Agent Relay is intentionally small.

## Core Principle

Read when joining.
Record when finishing.
Handoff only when needed.

## Directory

Agent Relay files live under:

`.agent-relay/`

Recommended structure:

```text
.agent-relay/
├── PROTOCOL.md
├── INDEX.md
├── GUIDANCE.md
├── relay.log
├── handoff/
└── templates/
```

## Rule 1. Read Context When Joining

Agents must read recent relay context when starting or resuming work.

Read `AGENTS.md` first. Then choose a session agent name before reading or writing relay entries, and use it consistently for the session.

Use:

```text
<agent>(<llm-or-version>[, <role>])
```

Examples:

```text
Codex(GPT-5.5)
Claude Code(Claude Sonnet 4.5)
Cursor(GPT-5.5, Reviewer)
```

If the LLM or version is unknown, use only the agent name without parentheses, for example `Codex`, `Claude Code`, or `Cursor`.

Use `Human` only for work performed directly by a human.

Read relay context when:

- starting a new session;
- switching agents, models, or harnesses;
- resuming after a long pause;
- taking over from a handoff;
- the current context is unclear.

Recommended read order:

1. `AGENTS.md`
2. choose the session agent name
3. `README.md`
4. `.agent-relay/GUIDANCE.md`
5. `.agent-relay/PROTOCOL.md`
6. `.agent-relay/INDEX.md`
7. last 50 lines of `.agent-relay/relay.log`
8. relevant handoff files under `.agent-relay/handoff/`

Do not re-read `relay.log` before every user message in the same active session.

Within one continuous session, keep working from the current conversation context.

## Rule 2. Record Results After Meaningful Work

After meaningful work, append a short event to `.agent-relay/relay.log`.

Meaningful work includes:

- code changes
- file changes
- test execution
- debugging attempts
- architecture or design decisions
- issue or task updates
- work that another agent may need later

Do not log simple Q&A, brainstorming, or explanation-only conversations unless the user asks to preserve the context.

## Rule 2a. Preserve Durable Guidance

Update `.agent-relay/GUIDANCE.md` when the user gives a durable instruction, constraint, preference, convention, security rule, or "do not" rule that future agents should keep following.

Do not update `GUIDANCE.md` for one-off task details, transient debugging notes, temporary plans, or work progress.

When updating `GUIDANCE.md`, keep the entry short and stable. If the guidance came from a task conversation, also record the meaningful work result in `relay.log` when appropriate.

## Rule 3. Handoff Only When Needed

Create a handoff file only when the next agent cannot continue from the issue and relay log alone.

Create a handoff when:

- another agent or role must continue;
- the result is `partial`, `failed`, or `blocked`;
- tests, review, or analysis must be performed by another role;
- important context cannot fit in a short relay log event;
- the next agent would not be able to continue within 5 minutes.

Do not create a handoff when:

- the task is fully completed;
- the result can be summarized in one short relay log entry;
- there is no next owner.

## Relay Log

`relay.log` is append-only.

Never edit existing lines in place.

If a correction is needed, append a new `CORRECTION` event.

Recommended event types:

- `SESSION_START`
- `WORK_START`
- `WORK_END`
- `NOTE`
- `HANDOFF`
- `HANDOFF_RECEIVED`
- `HANDOFF_CLOSED`
- `CORRECTION`

Keep entries short.

If the context is long, create a separate markdown file under `.agent-relay/handoff/` and reference it from `relay.log`.

Use the session agent name in every `agent=`, `from=`, and `to=` field that names an agent or role.

When creating a handoff file, append a `HANDOFF` event to `relay.log` with `path=<handoff-file>`.

Example:

```text
2026-04-28T16:40:00+09:00 | HANDOFF | agent=Codex(GPT-5.5) | result=blocked | path=.agent-relay/handoff/20260428-1640-codex-auth-refresh.md | summary="Auth refresh behavior needs follow-up analysis"
```

## Minimal WORK_END Format

```text
<timestamp> | WORK_END | agent=<agent> | result=<success|partial|failed|blocked> | files=<summary> | tests=<summary>
```

Use `agent=<agent>(<llm-or-version>[, <role>])`, for example `Codex(GPT-5.5)` or `Cursor(GPT-5.5, Reviewer)`. If the LLM or version is unknown, use only the agent name, for example `Codex`. Use `agent=Human` only for work performed directly by a human.

Example:

```text
2026-04-28T16:40:00+09:00 | WORK_END | agent=Codex(GPT-5.5) | result=success | files=src/cube/state.ts, tests/cube-state.test.ts | tests=npm test passed
```

## Handoff File Naming

By default, do not assume the next agent is known. Use the current session agent and topic:

```text
.agent-relay/handoff/YYYYMMDD-HHMM-<from-agent>-<topic>.md
```

Example:

```text
.agent-relay/handoff/20260428-1640-codex-cube-rotation.md
```

If the next agent is already known, `to-<target-agent>` may be included:

```text
.agent-relay/handoff/YYYYMMDD-HHMM-<from-agent>-to-<target-agent>-<topic>.md
```

## Safety

Never store the following in `.agent-relay/`:

- API keys
- tokens
- passwords
- private credentials
- customer data
- sensitive internal information
- production secrets

## Git Policy

Commit `.agent-relay/` to Git.

Agent Relay is project handoff state. The next agent should receive the same protocol, index, relay log, and handoff files from the repository.

Do not add `.agent-relay/` to `.gitignore`. If an ignore rule already exists, remove it.

Never store secrets or sensitive data in `.agent-relay/`.
