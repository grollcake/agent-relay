# TASK_BEGIN Template

Use this before starting a meaningful task that changes files, runs tests, investigates a non-trivial issue, or may need recovery if the session stops.

```text
<timestamp> | TASK_BEGIN | agent=<agent> | task=<MMDD-xxx> | summary=<task summary>
```

Fields:

- `timestamp`: ISO-8601 with timezone, e.g. `2026-04-28T16:20:00+09:00`
- `agent`: session agent name, e.g. `Codex(GPT-5.5)`
- `task`: generated at task start, in `MMDD-xxx` format, e.g. `0428-qmx`
- `summary`: short task summary, not the user's full prompt

Example:

```text
2026-04-28T16:20:00+09:00 | TASK_BEGIN | agent=Codex(GPT-5.5) | task=0428-qmx | summary="Align relay event names"
```
