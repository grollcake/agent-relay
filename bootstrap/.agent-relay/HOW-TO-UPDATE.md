# How To Update Agent Relay

Use this when the user asks to update, sync, or refresh Agent Relay in a project
that already has `.agent-relay/`.

## Steps

1. Read `.agent-relay/VERSION`; stop and report if it is missing.
2. Fetch or copy the latest upstream `main` from `https://github.com/grollcake/agent-relay` into a temporary location.
3. Select and checksum-verify the current platform binary under upstream `bootstrap/.agent-relay/bin/<os>-<arch>/`.
4. Use the installed binary as `<agent-relay>` when available. For legacy script-only installs and all Windows updates, run the new upstream binary from its upstream or a temporary path while the current directory is the target project.
5. Run `<agent-relay> update --upstream <repo>` for a dry-run and inspect managed files for local customizations.
6. If managed files can be replaced, run `<agent-relay> update --upstream <repo> --apply`.
7. If updating manually, use `<agent-relay> merge-agent-block` for `AGENTS.md` and `CLAUDE.md` blocks.
8. Run the installed `<agent-relay> lint` after the update.
9. Record `REQUEST -> RUN_DONE` in `relay.log` with the before/after `VERSION` in the summary if the update was manual. The update command records this automatically.

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
