"""Run the Linux monitor against real commands, including quiet work and SIGTERM."""

import os
import selectors
import signal
import subprocess
import sys
import tempfile
import time
import unittest
from pathlib import Path

MONITOR = Path(__file__).resolve().parents[1] / 'steps/monitor_resources.sh'


@unittest.skipUnless(sys.platform.startswith('linux'), 'Hosted Ubuntu monitor requires Linux')
class ResourceMonitorTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.addCleanup(self.directory.cleanup)
        self.report = Path(self.directory.name) / 'diagnostics'
        self.environment = {**os.environ, 'RESOURCE_MONITOR_INTERVAL_SECONDS': '1'}

    def command(self, script, *args):
        return ['bash', str(MONITOR), str(self.report), sys.executable, '-c', script, *args]

    def wait_for_output(self, process, marker, timeout=15):
        output = b''
        deadline = time.monotonic() + timeout
        with selectors.DefaultSelector() as selector:
            selector.register(process.stdout, selectors.EVENT_READ)
            while time.monotonic() < deadline:
                if selector.select(timeout=0.2):
                    chunk = os.read(process.stdout.fileno(), 65536)
                    output += chunk
                    if marker in output:
                        return output.decode()
                    if not chunk:
                        break
        self.fail(f'Missing live output {marker!r}: {output.decode()}')

    def stop(self, process):
        if process.poll() is None:
            process.terminate()
        try:
            process.communicate(timeout=12)
        except subprocess.TimeoutExpired:
            process.kill()
            process.communicate(timeout=5)

    def test_quiet_command_streams_resources_and_unbuffered_python(self):
        process = subprocess.Popen(
            self.command('import time; print("QUIET_WORK_STARTED"); time.sleep(4)', 'sensitive-argument-canary'),
            stdout=subprocess.PIPE, stderr=subprocess.STDOUT, env=self.environment,
        )
        self.addCleanup(self.stop, process)
        output = self.wait_for_output(process, b'QUIET_WORK_STARTED')
        self.assertIsNone(process.poll(), 'Output must arrive while the command is still running')
        output += self.wait_for_output(process, b'[RESOURCE MONITOR]')
        self.assertIsNone(process.poll(), 'Resource sample must arrive during quiet work')
        remaining, _ = process.communicate(timeout=15)
        output += remaining.decode()
        self.assertEqual(process.returncode, 0, output)
        self.assertIn('QUIET_WORK_STARTED', (self.report / 'command.log').read_text())
        samples = (self.report / 'resources.log').read_text()
        self.assertGreaterEqual(samples.count('[RESOURCE MONITOR]'), 2)
        self.assertIn('Mem:', samples)
        self.assertIn('Filesystem', samples)
        self.assertIn('Maximum resident set size', (self.report / 'time.log').read_text())
        self.assertNotIn('sensitive-argument-canary', output)
        self.assertNotIn('sensitive-argument-canary', samples)
        self.assertNotIn('sensitive-argument-canary', (self.report / 'time.log').read_text())
        if Path('/sys/fs/cgroup/memory.events').exists():
            self.assertIn('memory.events:', samples)
            self.assertIn('oom_kill', samples)

    def test_command_failure_and_signal_exit_status_are_preserved(self):
        for script, expected in [('raise SystemExit(7)', 7), ('import os, signal; os.kill(os.getpid(), signal.SIGTERM)', 143)]:
            with self.subTest(expected=expected):
                result = subprocess.run(
                    self.command(script), capture_output=True, text=True,
                    env=self.environment, timeout=15,
                )
                self.assertEqual(result.returncode, expected, result.stdout + result.stderr)
                self.assertIn(f'Command exit status: {expected}', (self.report / 'exit-status.log').read_text())
                self.assertIn('[KERNEL OOM EVIDENCE]', result.stdout)

    def test_sigterm_to_monitor_terminates_command_descendants(self):
        child_pid_file = Path(self.directory.name) / 'child.pid'
        script = (
            'import subprocess, sys, time; from pathlib import Path; '
            'child = subprocess.Popen([sys.executable, "-c", "import time; time.sleep(60)"]); '
            'Path(sys.argv[1]).write_text(str(child.pid)); '
            'print("COMMAND_TREE_STARTED"); time.sleep(60)'
        )
        process = subprocess.Popen(
            self.command(script, str(child_pid_file)), stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT, env=self.environment,
        )
        self.addCleanup(self.stop, process)
        self.wait_for_output(process, b'COMMAND_TREE_STARTED')
        child_pid = int(child_pid_file.read_text())
        process.send_signal(signal.SIGTERM)
        output, _ = process.communicate(timeout=15)
        self.assertEqual(process.returncode, 143, output.decode())
        self.assertIn('Received SIGTERM', (self.report / 'termination.log').read_text())
        self.assertIn('Command exit status: 143', (self.report / 'exit-status.log').read_text())
        self.assertIn('[RESOURCE MONITOR]', output.decode(), 'Final snapshot must be streamed')
        child_status = Path(f'/proc/{child_pid}/stat')
        if child_status.exists():
            self.assertEqual(child_status.read_text().split(') ', 1)[1].split()[0], 'Z', 'Child must be exited, not still running')
