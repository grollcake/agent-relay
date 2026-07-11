# Repository Instructions

## Environment And Editing

- Detect the current OS and shell before running repository commands.
- Treat repository text as UTF-8 and verify Korean text is not corrupted.
- On Windows, use Git for Windows Git Bash as the required shell for this
  repository. Agent Relay itself is a native Go binary; Git Bash remains the
  supported operational shell for repository and Git workflows.
- Use Go 1.26 or newer for source builds and tests.
- Use `apply_patch` for repository edits. After changes, inspect the result and
  run `git diff --check`.
