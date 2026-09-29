#!/usr/bin/env bash
# Stream command output and runner capacity; preserve the command's exit status.
# Linux CI only: /usr/bin/time -v records elapsed time and peak process RSS.
set -uo pipefail

report_dir="$1"
shift
mkdir -p "$report_dir"

sample_resources() {
  date -u '+%Y-%m-%dT%H:%M:%SZ'
  free -m
  df -Pk . "${RUNNER_TEMP:-/tmp}"
  # Process names only: command arguments may contain credentials.
  ps -eo comm,rss --sort=-rss | head -n 12 || true
}

{
  git rev-parse HEAD
  go version
  go env GOCACHE GOMODCACHE
  golangci-lint version 2>/dev/null || true
  golangci-lint cache status 2>/dev/null || true
  sample_resources
} | tee "$report_dir/environment.log"

monitor_pid=''
cleanup() {
  if [[ -n "$monitor_pid" ]]; then
    kill "$monitor_pid" 2>/dev/null || true
    wait "$monitor_pid" 2>/dev/null || true
  fi
  sample_resources >> "$report_dir/resources.log" 2>&1
}
trap cleanup EXIT
(
  sleep_pid=''
  trap 'if [[ -n "$sleep_pid" ]]; then kill "$sleep_pid" 2>/dev/null || true; fi; exit 0' TERM
  while true; do
    sample_resources
    sleep 15 &
    sleep_pid=$!
    wait "$sleep_pid"
  done
) > >(tee "$report_dir/resources.log") 2>&1 &
monitor_pid=$!

/usr/bin/time -v "$@" 2>&1 | tee "$report_dir/command.log"
command_status=${PIPESTATUS[0]}
printf 'Command exit status: %s\n' "$command_status" | tee "$report_dir/exit-status.log"
exit "$command_status"
