#!/usr/bin/env bash
# Stream command output and runner capacity; preserve the command's exit status.
# Linux CI only: GNU time records elapsed time and peak process RSS.
# Trap callbacks are invoked indirectly.
# shellcheck disable=SC2329
set -uo pipefail

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <report-directory> <command> [arguments...]" >&2
  exit 2
fi
report_dir="$1"
shift
mkdir -p "$report_dir" || exit 2
interval="${RESOURCE_MONITOR_INTERVAL_SECONDS:-15}"
if [[ ! "$interval" =~ ^[1-9][0-9]*$ ]]; then
  echo 'RESOURCE_MONITOR_INTERVAL_SECONDS must be a positive integer' >&2
  exit 2
fi
# Python otherwise buffers stdout when Actions captures it through a pipe.
export PYTHONUNBUFFERED=1

# Hosted Ubuntu runners use cgroup v2. Inspect both the command's cgroup and
# its ancestors: the effective memory limit may be imposed by a parent.
cgroup_dir='/sys/fs/cgroup'
if [[ -r /proc/self/cgroup ]]; then
  cgroup_path=$(awk -F: '$1 == "0" {print $3}' /proc/self/cgroup)
  if [[ -n "$cgroup_path" && -d "/sys/fs/cgroup$cgroup_path" ]]; then
    cgroup_dir="/sys/fs/cgroup${cgroup_path%/}"
  fi
fi

sample_cgroups() {
  local directory="$cgroup_dir" metric
  while [[ "$directory" == /sys/fs/cgroup* ]]; do
    for metric in memory.current memory.peak memory.max memory.swap.current memory.swap.max memory.events cpu.max cpu.stat; do
      if [[ -r "$directory/$metric" ]]; then
        printf '%s: ' "$directory/$metric"
        tr '\n' ' ' < "$directory/$metric" || true
        printf '\n'
      fi
    done
    [[ "$directory" == /sys/fs/cgroup ]] && break
    directory="${directory%/*}"
  done
}

sample_resources() {
  printf '\n[RESOURCE MONITOR] '
  date -u '+%Y-%m-%dT%H:%M:%SZ' || true
  uptime || true
  free -m || true
  df -Pk . "${RUNNER_TEMP:-/tmp}" || true
  df -Pi . "${RUNNER_TEMP:-/tmp}" || true
  # The second row reports current CPU use, swapping, and I/O wait over 1s.
  vmstat -w 1 2 || true
  # Process names only: command arguments may contain credentials.
  ps -eo pid,ppid,comm,%cpu,%mem,rss,stat,etime --sort=-rss | head -n 12 || true
  local pressure
  for pressure in /proc/pressure/cpu /proc/pressure/memory /proc/pressure/io; do
    if [[ -r "$pressure" ]]; then
      printf '%s:\n' "$pressure"
      cat "$pressure" || true
    fi
  done
  sample_cgroups
}

kernel_oom_evidence() {
  printf '\n[KERNEL OOM EVIDENCE]\n'
  local kernel_log
  if kernel_log=$(dmesg --ctime 2>/dev/null); then
    printf '%s\n' "$kernel_log" | grep -Ei 'out of memory|oom-kill|killed process|memory cgroup out of memory' | tail -n 20 || true
  else
    echo 'Kernel log unavailable to this user; inspect memory.events oom/oom_kill counters above.'
  fi
}

{
  git rev-parse HEAD
  go version
  go env GOCACHE GOMODCACHE
  golangci-lint version 2>/dev/null || true
  golangci-lint cache status 2>/dev/null || true
  printf 'GOMAXPROCS=%s GOMEMLIMIT=%s GOGC=%s PYTHONUNBUFFERED=%s\n' \
    "${GOMAXPROCS:-unset}" "${GOMEMLIMIT:-unset}" "${GOGC:-unset}" "$PYTHONUNBUFFERED"
  sample_resources
} | tee "$report_dir/environment.log"

monitor_pid=''
command_pid=''
record_status() {
  printf 'Command exit status: %s\n' "$1" | tee "$report_dir/exit-status.log"
  if (( $1 == 143 )); then
    echo 'Exit 143 indicates SIGTERM. Check runner shutdown/cancellation messages and resource samples; SIGTERM alone does not establish OOM.' | tee -a "$report_dir/exit-status.log"
  fi
}
cleanup() {
  trap '' TERM INT
  if [[ -n "$monitor_pid" ]]; then
    kill "$monitor_pid" 2>/dev/null || true
    wait "$monitor_pid" 2>/dev/null || true
  fi
  { sample_resources; kernel_oom_evidence; } 2>&1 | tee -a "$report_dir/resources.log"
  if [[ -f "$report_dir/time.log" ]]; then
    cat "$report_dir/time.log"
  fi
}
handle_signal() {
  local signal="$1" status="$2"
  trap '' TERM INT
  printf '[RESOURCE MONITOR] Received SIG%s at %s\n' "$signal" "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" | tee "$report_dir/termination.log"
  if [[ -n "$command_pid" ]]; then
    # The command runs in a separate session so compilation/test descendants
    # receive termination too. Bound cleanup for commands that ignore signals.
    kill -s "$signal" -- "-$command_pid" 2>/dev/null || true
    for _ in 1 2 3; do
      kill -0 -- "-$command_pid" 2>/dev/null || break
      sleep 1
    done
    kill -KILL -- "-$command_pid" 2>/dev/null || true
    wait "$command_pid" 2>/dev/null || true
  fi
  record_status "$status"
  exit "$status"
}
trap cleanup EXIT
trap 'handle_signal TERM 143' TERM
trap 'handle_signal INT 130' INT
(
  sleep_pid=''
  trap 'if [[ -n "$sleep_pid" ]]; then kill "$sleep_pid" 2>/dev/null || true; wait "$sleep_pid" 2>/dev/null || true; fi; exit 0' TERM INT
  while true; do
    sample_resources 2>&1 | tee -a "$report_dir/resources.log"
    sleep "$interval" &
    sleep_pid=$!
    wait "$sleep_pid"
    sleep_pid=''
  done
) &
monitor_pid=$!

# GNU time -v includes the full argv. Use explicit fields to avoid disclosing
# arguments that may contain credentials in stdout or uploaded diagnostics.
time_format='User time (seconds): %U
System time (seconds): %S
Percent of CPU this job got: %P
Elapsed (wall clock) time (seconds): %e
Maximum resident set size (kbytes): %M
Major page faults: %F
Minor page faults: %R
Voluntary context switches: %w
Involuntary context switches: %c
Swaps: %W
File system inputs: %I
File system outputs: %O
Exit status: %x'
setsid /usr/bin/time -f "$time_format" -o "$report_dir/time.log" -- "$@" > >(tee "$report_dir/command.log") 2>&1 &
command_pid=$!
wait "$command_pid"
command_status=$?
record_status "$command_status"
exit "$command_status"
