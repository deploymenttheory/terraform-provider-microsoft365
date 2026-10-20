"""Exercise the real tool wrapper and missing/incomplete evidence handling."""

import importlib.util
import io
import json
import os
from pathlib import Path
import signal
import shutil
import subprocess
import sys
import tempfile
import time
import unittest

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT))


def load(name):
    spec = importlib.util.spec_from_file_location(name, ROOT / f"{name}.py")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


class ToolTests(unittest.TestCase):
    def test_start_event_streams_while_tool_output_is_buffered(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            events = root / "events.jsonl"
            output = io.StringIO()
            env = {
                **os.environ, "OOM_DISCOVERY_REPORT": directory,
                "OOM_DISCOVERY_EVENT_LOG": str(events),
            }
            with load("stream_tools").ToolEventStream(events, output):
                process = subprocess.Popen(
                    [sys.executable, str(ROOT / "tool_exec.py"), sys.executable,
                     "-c", "import time; time.sleep(2)"],
                    env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                )
                try:
                    deadline = time.monotonic() + 1.5
                    while '"event": "start"' not in output.getvalue() and time.monotonic() < deadline:
                        time.sleep(0.01)
                    self.assertIn('"event": "start"', output.getvalue())
                    self.assertIsNone(process.poll(), "Start event must arrive before the tool exits")
                finally:
                    _, stderr = process.communicate(timeout=10)
            self.assertEqual(stderr, b"", "Relay must not duplicate records through buffered stderr")
            records = [json.loads(line.split("[GO TOOL] ", 1)[1]) for line in output.getvalue().splitlines()]
            self.assertEqual([record["event"] for record in records], ["start", "end"])

    def run_tool(self, script):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        env = {
            **os.environ,
            "OOM_DISCOVERY_REPORT": directory.name,
            "TOOLEXEC_IMPORTPATH": "example.test/package",
            "M365_CLIENT_SECRET": "must-not-be-logged",
        }
        result = subprocess.run(
            [sys.executable, str(ROOT / "tool_exec.py"), sys.executable, "-c", script],
            env=env, capture_output=True, text=True,
        )
        records = [
            json.loads(line)
            for path in (Path(directory.name) / "tools").glob("*.jsonl")
            for line in path.read_text().splitlines()
        ]
        self.assertNotIn("must-not-be-logged", result.stdout + result.stderr + json.dumps(records))
        self.assertEqual(records[0]["package"], "example.test/package")
        self.assertGreater(records[1]["peak_rss_kib"], 0)
        return result, records

    def test_failure_and_output_preserved_with_memory_attribution(self):
        result, records = self.run_tool('print("child output"); raise SystemExit(7)')
        self.assertEqual(result.returncode, 7)
        self.assertEqual(result.stdout, "child output\n")
        self.assertEqual(records[-1]["exit_code"], 7)
        self.assertEqual(records[0]["pid"], records[1]["pid"])

    def test_signal_preserved(self):
        for signum in (signal.SIGTERM, signal.SIGKILL):
            with self.subTest(signum=signum):
                result, records = self.run_tool(f"import os; os.kill(os.getpid(), {signum})")
                self.assertEqual(result.returncode, -signum)
                self.assertEqual(records[-1]["exit_code"], -signum)


class EvidenceTests(unittest.TestCase):
    def test_streamed_evidence_survives_missing_artifact_and_killed_tool(self):
        text = '\n'.join([
            '2026-10-20T00:00:00Z [EXPERIMENT] {"name":"baseline"}',
            '2026-10-20T00:00:00Z Mem: 15947 1000 1000 0 13947 14000',
            '2026-10-20T00:00:00Z /sys/fs/cgroup/test/memory.events: oom_kill 0',
            '2026-10-20T00:00:00Z [GO TOOL] {"event":"start","invocation":"1","tool":"vet","package":"example"}',
            '2026-10-20T00:00:00Z Mem: 15947 15000 100 0 847 500',
            '2026-10-20T00:00:00Z Command exit status: 143',
            '2026-10-20T00:00:00Z /sys/fs/cgroup/test/memory.events: oom_kill 1',
        ])
        result = load("summarize").summarize_log(text)
        self.assertEqual(result["exit_code"], 143)
        self.assertEqual(result["oom_kill_delta"], 1)
        self.assertEqual(result["unfinished_tools"][0]["package"], "example")
        self.assertFalse(result["candidate_pass"])

    def test_replay_oom_does_not_contaminate_completed_build(self):
        text = '\n'.join([
            '[EXPERIMENT RESULT] {"exit_code":0}',
            'Mem: 15947 1000 1000 0 13947 14000',
            '/sys/fs/cgroup/test/memory.events: oom_kill 0',
            '[VET REPLAY BEGIN] {}',
            'Mem: 15947 15000 100 0 847 500',
            '/sys/fs/cgroup/test/memory.events: oom_kill 1',
        ])
        result = load("summarize").summarize_log(text)
        self.assertTrue(result["candidate_pass"])
        self.assertEqual(result["oom_kill_delta"], 0)

    def test_missing_results_are_not_success(self):
        with tempfile.TemporaryDirectory() as directory:
            result = load("summarize").summarize(Path(directory))
            self.assertFalse(result["candidate_pass"])
            self.assertFalse(result["evidence_complete"])
            self.assertIsNone(result["oom_kill_delta"])

    def test_oom_delta_and_memory_headroom(self):
        with tempfile.TemporaryDirectory() as directory:
            report = Path(directory)
            (report / "experiment.json").write_text('{"exit_code": 0}')
            (report / "environment.log").write_text(
                "Mem: 15947 1000 1000 0 13947 14000\n"
                "/sys/fs/cgroup/test/memory.events: low 0 oom 0 oom_kill 2\n"
            )
            (report / "resources.log").write_text(
                "Mem: 15947 13000 1000 0 1947 2500\n"
                "/sys/fs/cgroup/test/memory.events: low 0 oom 0 oom_kill 2\n"
            )
            result = load("summarize").summarize(report)
            self.assertEqual(result["min_available_mib"], 2500)
            self.assertEqual(result["oom_kill_delta"], 0)
            self.assertTrue(result["candidate_pass"])
            with (report / "resources.log").open("a") as stream:
                stream.write("/sys/fs/cgroup/test/memory.events: low 0 oom 0 oom_kill 3\n")
            self.assertFalse(load("summarize").summarize(report)["candidate_pass"])

    def test_command_builds_but_never_runs_test_binary(self):
        runner = load("run")
        for variant in runner.VARIANTS:
            command = runner.build_command("./example", Path("/tmp/report"), True, variant)
            self.assertIn("-c", command)
            self.assertIn("-race", command)
            self.assertEqual(command[-1], "./example")
            self.assertEqual("-vet=off" in command, variant in ("no-vet", "combined"))


@unittest.skipUnless(shutil.which("go"), "Go required for compile-only integration check")
class GoIntegrationTests(unittest.TestCase):
    def test_compile_only_invokes_vet_and_does_not_execute_init_or_tests(self):
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "go.mod").write_text("module example.test/discovery\n\ngo 1.23.0\n")
            (root / "example_test.go").write_text(
                'package discovery\n'
                'import ("fmt"; "testing")\n'
                'func init() { panic("TEST_BINARY_MUST_NOT_EXECUTE") }\n'
                'func TestExample(t *testing.T) { fmt.Printf("%d", "incorrect type") }\n'
            )
            (root / "go-work").mkdir()
            env = {
                **os.environ, "GOTOOLCHAIN": "local", "GOWORK": "off",
                "GOMAXPROCS": "1", "GOCACHE": str(root / "build-cache"),
                "GOTMPDIR": str(root / "go-work"),
            }
            for variant, expected in (("baseline", 1), ("no-vet", 0)):
                report = root / variant
                report.mkdir()
                env["OOM_DISCOVERY_REPORT"] = str(report)
                command = load("run").build_command(".", report, False, variant)
                result = subprocess.run(
                    command, cwd=root, env=env, text=True, capture_output=True, timeout=120,
                )
                self.assertEqual(result.returncode, expected, result.stdout + result.stderr)
                self.assertNotIn("panic: TEST_BINARY_MUST_NOT_EXECUTE", result.stdout + result.stderr)
                evidence = load("summarize").summarize(report)
                if variant == "baseline":
                    self.assertGreater(evidence["vet_invocations"], 0)
                    configs = list((report / "tools").glob("*.vet-config.json"))
                    self.assertTrue(configs, "Vet configuration must be retained for attribution/replay")
                else:
                    self.assertEqual(evidence["vet_invocations"], 0)
                    self.assertTrue((report / "discovery.test").is_file())


if __name__ == "__main__":
    unittest.main()
