# Repository Instructions

## Environment And Editing

- Detect the current OS and shell; do not assume this repository is tied to one
  platform.
- Treat repository text as UTF-8. In PowerShell, specify `-Encoding utf8` when
  reading text and verify Korean text is not corrupted.
- On Windows, prefer PowerShell. Use Bash only for POSIX scripts, and preserve
  variables, quotes, paths, and encoding across the PowerShell-to-Bash boundary.
- Use `apply_patch` for repository edits. After changes, inspect the result and
  run `git diff --check`.
