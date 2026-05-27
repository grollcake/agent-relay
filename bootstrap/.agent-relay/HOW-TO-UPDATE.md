# How To Update Agent Relay

Use this when the user asks to update, sync, or refresh Agent Relay in a project
that already has `.agent-relay/`.

## Steps

1. Read `.agent-relay/VERSION`; stop and report if it is missing.
2. Fetch or copy the latest upstream `main` from `https://github.com/grollcake/agent-relay` into a temporary location.
3. Run `.agent-relay/scripts/update-agent-relay --upstream <repo>` for a dry-run and inspect managed files for local customizations.
4. If managed files can be replaced, run `.agent-relay/scripts/update-agent-relay --upstream <repo> --apply`.
5. If updating manually, use `.agent-relay/scripts/merge-agent-block` for `AGENTS.md` and `CLAUDE.md` blocks.
6. Run `.agent-relay/scripts/relay-lint` after the update.
7. Record `REQUEST -> RUN_DONE` in `relay.log` with the before/after `VERSION` in the summary if the update was manual. The update script records this automatically.

## Preserve

Do not overwrite:

- non-Agent-Relay instructions in `AGENTS.md` or `CLAUDE.md`
- `.agent-relay/GUIDANCE.md`
- `.agent-relay/LESSON-LEARNED.md`
- `.agent-relay/lesson-learned/`
- `.agent-relay/relay.log`
- `.agent-relay/runs/`

If an Agent Relay block contains project-specific instructions that cannot be
separated safely, leave the file unchanged and report the conflict.

## Report

Keep the final report short: version change, updated categories, and any manual
checks or conflicts.
