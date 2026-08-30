# How To Update Memento AI

Use this when the user asks to update, sync, or refresh Memento AI in a project
that already has `.memento/`.

## Steps

1. Read `.memento/VERSION`; stop and report if it is missing.
2. Fetch or copy the latest upstream `main` from `https://github.com/grollcake/memento` into a temporary location.
3. Select and checksum-verify the current platform binary under upstream `bootstrap/.memento/bin/<os>-<arch>/`.
4. Use the installed binary as `<memento>` when available. For legacy script-only installs and all Windows updates, run the new upstream binary from its upstream or a temporary path while the current directory is the target project.
5. Run `<memento> update --upstream <repo>` for a dry-run and inspect managed files for local customizations.
6. If managed files can be replaced, run `<memento> update --upstream <repo> --apply`.
7. If updating manually, use `<memento> merge-agent-block` for `AGENTS.md` and `CLAUDE.md` blocks.
8. Run the installed `<memento> lint` after the update.
9. Record `REQUEST -> RUN_DONE` in `memento.log` with the before/after `VERSION` in the summary if the update was manual. The update command records this automatically.

## Preserve

Do not overwrite:

- non-Memento AI instructions in `AGENTS.md` or `CLAUDE.md`
- `.memento/GUIDANCE.md`
- `.memento/LESSON-LEARNED.md`
- `.memento/lesson-learned/`
- `.memento/memento.log`
- `.memento/runs/`

If a Memento AI block contains project-specific instructions that cannot be
separated safely, leave the file unchanged and report the conflict.

## Report

Keep the final report short: version change, updated categories, and any manual
checks or conflicts.
