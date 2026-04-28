# Work End Entry Template

Use this format for `.agent-relay/relay.log`.

```text
<timestamp> | WORK_END | agent=<agent> | result=<success|partial|failed|blocked> | files=<summary> | tests=<summary>
```

## Fields

- `timestamp`: ISO-8601 with timezone
- `agent`: display name in `<agent>(<llm-or-version>[, <role>])` format, e.g. `Codex(GPT-5.5)`; if unknown, use only the agent name, e.g. `Codex`; use `Human` only for direct human work
- `result`: `success`, `partial`, `failed`, or `blocked`
- `files`: changed files or `none`
- `tests`: tests run or `not run`

## Example

```text
2026-04-28T16:40:00+09:00 | WORK_END | agent=Codex(GPT-5.5) | result=success | files=src/cube/state.ts, tests/cube-state.test.ts | tests=npm test passed
```
