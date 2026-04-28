# Guidance

Use this file for durable project guidance that should survive across agent sessions.

Do not use this file for progress tracking, current task status, temporary plans, or next steps. Use `relay.log` and handoff files for that.

Update this file when the user gives a durable instruction, constraint, preference, convention, security rule, or "do not" rule that future agents should keep following.

Do not update this file for one-off task details, transient debugging notes, temporary plans, or work progress.

Examples to record:

- "Always keep the public API backward compatible."
- "Do not introduce new runtime dependencies without asking."
- "Use pnpm for package scripts in this repository."
- "Never store customer data in fixtures."

Examples not to record:

- "Fix this failing test."
- "Try the other implementation first."
- "The current debugging hypothesis is a race condition."
- "Next, update the button copy."

## Stable Project Context

- <long-lived project purpose or domain context>
- <important architectural or product context that is unlikely to change often>

## User Instructions

- <standing user preference or instruction>

## Constraints

- <hard technical, product, compatibility, or operational constraint>

## Do Not

- <thing future agents must not do>

## Security And Privacy

- <data, credential, privacy, or operational safety rule>

## Conventions

- <naming, style, testing, branching, or workflow convention>
