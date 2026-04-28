# Task Done Entry Template

Use this format for `.agent-relay/relay.log`.

```text
<timestamp> | TASK_DONE  | agent=<agent> | task=<MMDD-xxx> | result=<success|partial|failed|blocked> | files=<summary> | tests=<summary>
```

## Fields

- `timestamp`: ISO-8601 with timezone
- `agent`: display name in `<agent>(<llm-or-version>[, <role>])` format, e.g. `Codex(GPT-5.5)`; if unknown, use only the agent name, e.g. `Codex`; use `Human` only for direct human work
- `task`: task key generated at `TASK_BEGIN`, in `MMDD-xxx` format, e.g. `0428-qmx`
- `result`: task outcome; use `success`, `partial`, `failed`, or `blocked`
- `files`: changed files or `none`
- `tests`: tests run or `not run`

## Example

```text
2026-04-28T16:40:00+09:00 | TASK_DONE  | agent=Codex(GPT-5.5) | task=0428-qmx | result=success | files=src/cube/state.ts, tests/cube-state.test.ts | tests=npm test passed
```
