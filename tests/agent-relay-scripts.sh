#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
TMP_ROOT=$(mktemp -d)

cleanup() {
  cleanup_status=$?
  trap - EXIT HUP INT TERM
  rm -rf "$TMP_ROOT"
  exit "$cleanup_status"
}

trap cleanup EXIT HUP INT TERM

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

expect_failure() {
  label="$1"
  shift
  if "$@" >"$TMP_ROOT/output" 2>&1; then
    fail "$label unexpectedly succeeded"
  fi
}

new_project() {
  project="$1"
  mkdir -p "$project"
  cp -R "$ROOT/bootstrap/.agent-relay" "$project/.agent-relay"
  cp "$ROOT/bootstrap/AGENTS.md" "$project/AGENTS.md"
  chmod +x "$project/.agent-relay/scripts/"*
  printf '%s\n' \
    '2026-07-11T10:00:00 | boot | REQUEST  | Director | Bootstrap Agent Relay' \
    '2026-07-11T10:00:00 | boot | RUN_DONE | Director | Agent Relay initialized' \
    > "$project/.agent-relay/relay.log"
}

write_plan() {
  write_path="$1"
  write_task="$2"
  printf '%s\n' \
    '# PLAN: Test flow' \
    "Task ID: $write_task" \
    'Date: 2026-07-11' \
    'Planner: test' \
    '## Director Brief' \
    '- Goal: validate the flow' \
    '## Success Criteria' \
    '- The test passes.' \
    '## Validation' \
    '- Run the test.' > "$write_path"
}

write_run() {
  write_path="$1"
  write_task="$2"
  write_round="$3"
  printf '%s\n' \
    "# RUN-$write_round: Test flow" \
    "Task ID: $write_task" \
    'Date: 2026-07-11' \
    'Executor: test' \
    'Status: complete' \
    '## Validation' \
    '- Test passed.' \
    '## Success Criteria Status' \
    '- Flow: met' > "$write_path"
}

write_review() {
  write_path="$1"
  write_task="$2"
  write_round="$3"
  write_result="$4"
  printf '%s\n' \
    "# REVIEW-$write_round: Test flow" \
    "Task ID: $write_task" \
    'Date: 2026-07-11' \
    'Planner: test' \
    "Result: $write_result" \
    '## Suggested User Checks' \
    '- Inspect the output.' \
    '## Evidence Reviewed' \
    '- Test output.' > "$write_path"
}

write_close() {
  write_path="$1"
  write_task="$2"
  printf '%s\n' \
    '# CLOSE: Test flow' \
    "Task ID: $write_task" \
    'Date: 2026-07-11' \
    'Director: test' \
    'Approved By: User' \
    '## Acceptance' \
    '- Accepted.' \
    '## Validation Summary' \
    '- Test passed.' > "$write_path"
}

event_role() {
  case "$1" in
    REQUEST|FEEDBACK|CLOSE|RUN_DONE) printf 'Director\n' ;;
    PLANNED|REVIEW) printf 'Planner\n' ;;
    EXECUTED) printf 'Executor\n' ;;
  esac
}

transition_allowed() {
  case "$1:$2" in
    START:REQUEST|\
    REQUEST:PLANNED|REQUEST:RUN_DONE|\
    PLANNED:EXECUTED|\
    EXECUTED:REVIEW|\
    REVIEW:EXECUTED|REVIEW:FEEDBACK|REVIEW:CLOSE|\
    FEEDBACK:EXECUTED)
      return 0
      ;;
  esac
  return 1
}

test_templates_are_rejected() {
  project="$TMP_ROOT/templates"
  new_project "$project"
  cd "$project"

  cp .agent-relay/templates/plan.md .agent-relay/runs/20260711-1000-probe-PLAN.md
  cp .agent-relay/templates/run.md .agent-relay/runs/20260711-1000-probe-RUN-01.md
  cp .agent-relay/templates/review.md .agent-relay/runs/20260711-1000-probe-REVIEW-01.md
  cp .agent-relay/templates/close.md .agent-relay/runs/20260711-1000-probe-CLOSE.md

  expect_failure "PLAN template" .agent-relay/scripts/artifact-check PLANNED \
    .agent-relay/runs/20260711-1000-probe-PLAN.md
  expect_failure "RUN template" .agent-relay/scripts/artifact-check EXECUTED \
    .agent-relay/runs/20260711-1000-probe-RUN-01.md
  expect_failure "REVIEW template" .agent-relay/scripts/artifact-check REVIEW \
    .agent-relay/runs/20260711-1000-probe-REVIEW-01.md
  expect_failure "CLOSE template" .agent-relay/scripts/artifact-check CLOSE \
    .agent-relay/runs/20260711-1000-probe-CLOSE.md
}

test_complete_flow_passes() {
  project="$TMP_ROOT/complete"
  new_project "$project"
  cd "$project"

  round=$(.agent-relay/scripts/director-tool new-round valid-flow --summary 'Validate complete flow')
  eval "$round"

  plan_path=".agent-relay/runs/$key-PLAN.md"
  run_path=".agent-relay/runs/$key-RUN-01.md"
  review_path=".agent-relay/runs/$key-REVIEW-01.md"
  close_path=".agent-relay/runs/$key-CLOSE.md"

  printf '%s\n' \
    '# PLAN: Valid flow' \
    "Task ID: $task_id" \
    'Date: 2026-07-11' \
    'Planner: test' \
    '## Director Brief' \
    '- Goal: validate the flow' \
    '## Success Criteria' \
    '- The test passes.' \
    '## Validation' \
    '- Run the test.' > "$plan_path"
  .agent-relay/scripts/director-tool append PLANNED --task-id "$task_id" \
    --role Planner --summary 'Plan complete' --path "$plan_path" >/dev/null

  printf '%s\n' \
    '# RUN-01: Valid flow' \
    "Task ID: $task_id" \
    'Date: 2026-07-11' \
    'Executor: test' \
    'Status: complete' \
    '## Validation' \
    '- Test passed.' \
    '## Success Criteria Status' \
    '- Flow: met' > "$run_path"
  .agent-relay/scripts/director-tool append EXECUTED --task-id "$task_id" \
    --role Executor --summary 'Run complete' --path "$run_path" >/dev/null

  printf '%s\n' \
    '# REVIEW-01: Valid flow' \
    "Task ID: $task_id" \
    'Date: 2026-07-11' \
    'Planner: test' \
    'Result: ready-for-user-decision' \
    '## Suggested User Checks' \
    '- Inspect the output.' \
    '## Evidence Reviewed' \
    '- Test output.' > "$review_path"
  .agent-relay/scripts/director-tool append REVIEW --task-id "$task_id" \
    --role Planner --summary 'Review complete' --path "$review_path" >/dev/null

  printf '%s\n' \
    '# CLOSE: Valid flow' \
    "Task ID: $task_id" \
    'Date: 2026-07-11' \
    'Director: test' \
    'Approved By: User' \
    '## Acceptance' \
    '- Accepted.' \
    '## Validation Summary' \
    '- Test passed.' > "$close_path"
  .agent-relay/scripts/director-tool append CLOSE --task-id "$task_id" \
    --role Director --summary 'Flow closed' --path "$close_path" >/dev/null

  .agent-relay/scripts/relay-lint >"$TMP_ROOT/lint-complete" 2>&1 || \
    fail "complete flow did not pass relay-lint"
}

test_lint_does_not_execute_paths() {
  project="$TMP_ROOT/injection"
  new_project "$project"
  cd "$project"

  marker="$TMP_ROOT/command-ran"
  printf '%s\n' \
    "2026-07-11T10:01:00 | evil | REQUEST  | Director | Probe | \$(touch $marker)" \
    '2026-07-11T10:01:01 | evil | RUN_DONE | Director | Probe complete' \
    >> .agent-relay/relay.log

  expect_failure "malicious path lint" .agent-relay/scripts/relay-lint
  [ ! -e "$marker" ] || fail "relay-lint executed a path from relay.log"
}

test_log_delimiters_are_rejected() {
  project="$TMP_ROOT/delimiters"
  new_project "$project"
  cd "$project"

  expect_failure "summary delimiter" .agent-relay/scripts/director-tool new-round \
    bad-summary --summary 'invalid | summary'
}

test_transition_matrix() {
  project="$TMP_ROOT/transitions"
  new_project "$project"
  cd "$project"

  for prior in START REQUEST PLANNED EXECUTED REVIEW FEEDBACK CLOSE RUN_DONE; do
    for next in REQUEST PLANNED EXECUTED REVIEW FEEDBACK CLOSE RUN_DONE; do
      transition_allowed "$prior" "$next" && continue

      printf '%s\n' \
        '2026-07-11T10:00:00 | boot | REQUEST  | Director | Bootstrap Agent Relay' \
        '2026-07-11T10:00:00 | boot | RUN_DONE | Director | Agent Relay initialized' \
        > .agent-relay/relay.log
      if [ "$prior" != "START" ]; then
        role=$(event_role "$prior")
        printf '2026-07-11T10:01:00 | matr | %-8s | %-8s | Matrix state\n' \
          "$prior" "$role" >> .agent-relay/relay.log
      fi

      role=$(event_role "$next")
      expect_failure "$prior -> $next" .agent-relay/scripts/director-tool append "$next" \
        --task-id matr --role "$role" --summary 'Invalid transition probe'
    done
  done

  printf '%s\n' \
    '2026-07-11T10:00:00 | boot | REQUEST  | Director | Bootstrap Agent Relay' \
    '2026-07-11T10:00:00 | boot | RUN_DONE | Director | Agent Relay initialized' \
    > .agent-relay/relay.log
  .agent-relay/scripts/director-tool append REQUEST --task-id drct --role Director \
    --summary 'Direct flow' >/dev/null
  .agent-relay/scripts/director-tool append RUN_DONE --task-id drct --role Director \
    --summary 'Direct flow complete' >/dev/null

  round=$(.agent-relay/scripts/director-tool new-round retry-flow --summary 'Retry flow')
  eval "$round"
  plan_path=".agent-relay/runs/$key-PLAN.md"
  write_plan "$plan_path" "$task_id"
  .agent-relay/scripts/director-tool append PLANNED --task-id "$task_id" \
    --role Planner --summary 'Plan complete' --path "$plan_path" >/dev/null

  run_path=".agent-relay/runs/$key-RUN-01.md"
  review_path=".agent-relay/runs/$key-REVIEW-01.md"
  write_run "$run_path" "$task_id" 01
  write_review "$review_path" "$task_id" 01 blockers
  .agent-relay/scripts/director-tool append EXECUTED --task-id "$task_id" \
    --role Executor --summary 'Run 01 complete' --path "$run_path" >/dev/null
  .agent-relay/scripts/director-tool append REVIEW --task-id "$task_id" \
    --role Planner --summary 'Review 01 complete' --path "$review_path" >/dev/null

  run_path=".agent-relay/runs/$key-RUN-02.md"
  review_path=".agent-relay/runs/$key-REVIEW-02.md"
  write_run "$run_path" "$task_id" 02
  write_review "$review_path" "$task_id" 02 ready-for-user-decision
  .agent-relay/scripts/director-tool append EXECUTED --task-id "$task_id" \
    --role Executor --summary 'Run 02 complete' --path "$run_path" >/dev/null
  .agent-relay/scripts/director-tool append REVIEW --task-id "$task_id" \
    --role Planner --summary 'Review 02 complete' --path "$review_path" >/dev/null

  .agent-relay/scripts/director-tool feedback --task-id "$task_id" \
    --summary 'User feedback' >/dev/null
  run_path=".agent-relay/runs/$key-RUN-03.md"
  review_path=".agent-relay/runs/$key-REVIEW-03.md"
  close_path=".agent-relay/runs/$key-CLOSE.md"
  write_run "$run_path" "$task_id" 03
  write_review "$review_path" "$task_id" 03 ready-for-user-decision
  write_close "$close_path" "$task_id"
  .agent-relay/scripts/director-tool append EXECUTED --task-id "$task_id" \
    --role Executor --summary 'Run 03 complete' --path "$run_path" >/dev/null
  .agent-relay/scripts/director-tool append REVIEW --task-id "$task_id" \
    --role Planner --summary 'Review 03 complete' --path "$review_path" >/dev/null
  .agent-relay/scripts/director-tool append CLOSE --task-id "$task_id" \
    --role Director --summary 'Retry flow closed' --path "$close_path" >/dev/null

  .agent-relay/scripts/relay-lint >"$TMP_ROOT/lint-transitions" 2>&1 || \
    fail "allowed transition flows did not pass relay-lint"
}

test_artifact_consistency() {
  project="$TMP_ROOT/consistency"
  new_project "$project"
  cd "$project"

  round=$(.agent-relay/scripts/director-tool new-round consistent-flow --summary 'Consistency flow')
  eval "$round"
  plan_path=".agent-relay/runs/$key-PLAN.md"

  write_plan "$plan_path" "$task_id"
  grep -v '^Date:' "$plan_path" > "$plan_path.tmp"
  mv "$plan_path.tmp" "$plan_path"
  expect_failure "PLAN missing required date" .agent-relay/scripts/director-tool append PLANNED \
    --task-id "$task_id" --role Planner --summary 'Missing date' --path "$plan_path"

  write_plan "$plan_path" wrong
  expect_failure "PLAN task-id mismatch" .agent-relay/scripts/director-tool append PLANNED \
    --task-id "$task_id" --role Planner --summary 'Wrong task' --path "$plan_path"

  write_plan "$plan_path" "$task_id"
  .agent-relay/scripts/director-tool append PLANNED --task-id "$task_id" \
    --role Planner --summary 'Plan complete' --path "$plan_path" >/dev/null

  wrong_run=".agent-relay/runs/20260711-1000-wrong-key-RUN-01.md"
  write_run "$wrong_run" "$task_id" 01
  expect_failure "RUN key mismatch" .agent-relay/scripts/director-tool append EXECUTED \
    --task-id "$task_id" --role Executor --summary 'Wrong key' --path "$wrong_run"

  run_path=".agent-relay/runs/$key-RUN-01.md"
  write_run "$run_path" "$task_id" 01
  .agent-relay/scripts/director-tool append EXECUTED --task-id "$task_id" \
    --role Executor --summary 'Run complete' --path "$run_path" >/dev/null

  wrong_review=".agent-relay/runs/$key-REVIEW-02.md"
  write_review "$wrong_review" "$task_id" 02 ready-for-user-decision
  expect_failure "REVIEW round mismatch" .agent-relay/scripts/director-tool append REVIEW \
    --task-id "$task_id" --role Planner --summary 'Wrong round' --path "$wrong_review"

  review_path=".agent-relay/runs/$key-REVIEW-01.md"
  write_review "$review_path" "$task_id" 01 ready-for-user-decision
  .agent-relay/scripts/director-tool append REVIEW --task-id "$task_id" \
    --role Planner --summary 'Review complete' --path "$review_path" >/dev/null

  close_path=".agent-relay/runs/$key-CLOSE.md"
  write_close "$close_path" wrong
  expect_failure "CLOSE task-id mismatch" .agent-relay/scripts/director-tool append CLOSE \
    --task-id "$task_id" --role Director --summary 'Wrong task' --path "$close_path"
  write_close "$close_path" "$task_id"
  .agent-relay/scripts/director-tool append CLOSE --task-id "$task_id" \
    --role Director --summary 'Consistency flow closed' --path "$close_path" >/dev/null

  cp .agent-relay/relay.log "$TMP_ROOT/consistent-relay.log"
  sed "s@$run_path@$wrong_run@" .agent-relay/relay.log > .agent-relay/relay.log.tmp
  mv .agent-relay/relay.log.tmp .agent-relay/relay.log
  expect_failure "lint artifact key mismatch" .agent-relay/scripts/relay-lint
  cp "$TMP_ROOT/consistent-relay.log" .agent-relay/relay.log

  sed 's/^Task ID:.*/Task ID: tampered/' "$review_path" > "$review_path.tmp"
  mv "$review_path.tmp" "$review_path"
  expect_failure "lint task-id mismatch" .agent-relay/scripts/relay-lint
}

test_update_preserves_project_state() {
  project="$TMP_ROOT/update"
  backup="$TMP_ROOT/update-backup"
  new_project "$project"
  mkdir -p "$backup"

  printf '0.9.0\n' > "$project/.agent-relay/VERSION"
  printf 'custom guidance\n' > "$project/.agent-relay/GUIDANCE.md"
  printf 'custom lesson index\n' > "$project/.agent-relay/LESSON-LEARNED.md"
  printf 'custom lesson\n' > "$project/.agent-relay/lesson-learned/custom.md"
  printf 'custom run\n' > "$project/.agent-relay/runs/preserved.md"
  printf 'stale managed protocol\n' > "$project/.agent-relay/PROTOCOL.md"

  cp "$project/.agent-relay/GUIDANCE.md" "$backup/GUIDANCE.md"
  cp "$project/.agent-relay/LESSON-LEARNED.md" "$backup/LESSON-LEARNED.md"
  cp "$project/.agent-relay/lesson-learned/custom.md" "$backup/custom.md"
  cp "$project/.agent-relay/runs/preserved.md" "$backup/preserved.md"
  cp "$project/.agent-relay/relay.log" "$backup/relay.log"
  before_lines=$(wc -l < "$backup/relay.log" | tr -d ' ')

  cd "$project"
  .agent-relay/scripts/update-agent-relay --upstream "$ROOT" --apply \
    > "$TMP_ROOT/update-output" 2>&1 || fail "update apply failed"

  cmp "$backup/GUIDANCE.md" .agent-relay/GUIDANCE.md || fail "GUIDANCE.md changed"
  cmp "$backup/LESSON-LEARNED.md" .agent-relay/LESSON-LEARNED.md || \
    fail "LESSON-LEARNED.md changed"
  cmp "$backup/custom.md" .agent-relay/lesson-learned/custom.md || \
    fail "lesson record changed"
  cmp "$backup/preserved.md" .agent-relay/runs/preserved.md || fail "run artifact changed"
  sed -n "1,${before_lines}p" .agent-relay/relay.log > "$backup/relay-prefix.log"
  cmp "$backup/relay.log" "$backup/relay-prefix.log" || fail "existing relay.log lines changed"
  after_lines=$(wc -l < .agent-relay/relay.log | tr -d ' ')
  [ "$after_lines" -eq $((before_lines + 2)) ] || fail "update did not append exactly two events"
  cmp "$ROOT/bootstrap/.agent-relay/PROTOCOL.md" .agent-relay/PROTOCOL.md || \
    fail "managed protocol was not updated"
  .agent-relay/scripts/relay-lint > "$TMP_ROOT/lint-update" 2>&1 || \
    fail "updated project did not pass relay-lint"
}

test_templates_are_rejected
test_complete_flow_passes
test_lint_does_not_execute_paths
test_log_delimiters_are_rejected
test_transition_matrix
test_artifact_consistency
test_update_preserves_project_state

printf 'agent-relay script tests passed\n'
