"""Regression coverage for interrupted acceptance runs; no external writes."""

import contextlib
import importlib.util
import io
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

SCRIPTS = Path(__file__).resolve().parents[1]


def load(name):
    spec = importlib.util.spec_from_file_location(name, SCRIPTS / f"{name}.py")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


runner = load("run_tests")
merger = load("merge_test_results")
detector = load("detect_job_failures")
issues = load("manage_test_issues")
job_issues = load("manage_job_failure_issues")


def result(package="./internal/one", name="TestSame", **extra):
    return dict(package=package, test_name=name, package_complete=True,
                configuration_block_type="resources", service="example", **extra)


def issue(number=1, package="./internal/one", name="TestSame"):
    return {"number": number, "title": f"Bug: {name} Failing",
            "body": f"<!-- acceptance-test: {json.dumps({'package': package, 'test_name': name})} -->"}


def job(conclusion="failure", step="cancelled"):
    return {"id": 42, "name": "Resource Tests", "conclusion": conclusion,
            "steps": [{"name": "Set up Go", "conclusion": "success", "number": 1},
                      {"name": "Run resource tests", "conclusion": step, "number": 2}],
            "check_run_url": "https://api.github.com/repos/o/r/check-runs/99",
            "html_url": "https://github.com/o/r/actions/runs/1/job/42"}


def output(package, body, exit_code=0):
    end = f"[COMMAND] Exit={exit_code} package={package} elapsed=1.23s\n" if exit_code is not None else ""
    return f"[COMMAND] Started package={package} pid=123\n{body}{end}"


class TemporaryTests(unittest.TestCase):
    def setUp(self):
        directory = tempfile.TemporaryDirectory()
        self.addCleanup(directory.cleanup)
        previous = Path.cwd()
        os.chdir(directory.name)
        self.addCleanup(os.chdir, previous)
        stdout = contextlib.redirect_stdout(io.StringIO())
        stdout.__enter__()
        self.addCleanup(stdout.__exit__, None, None, None)

    def write(self, path, data):
        path = Path(path)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(json.dumps(data))


class ParserTests(TemporaryTests):
    def parse(self, content, packages):
        Path("output.log").write_text(content)
        runner.parse_test_results("output.log", "resources", "example", packages)
        return [json.loads(Path(name).read_text()) for name in
                ("test-failures.json", "test-successes.json", "test-run.json")]

    def test_same_name_in_different_packages_and_panic_keep_identity(self):
        content = output("./one", "=== RUN   TestSame\n--- PASS: TestSame (0.1s)\n")
        content += output("./two", "=== RUN   TestSame\npanic: failed\nFAIL github.com/example/two\n", 1)
        failures, successes, manifest = self.parse(content, ["./one", "./two"])
        self.assertEqual([(x["package"], x["test_name"]) for x in failures], [("./two", "TestSame")])
        self.assertEqual(successes[0]["package"], "./one")
        self.assertTrue(manifest["complete"])

    def test_interrupted_package_keeps_pass_but_cannot_resolve(self):
        content = output("./one", "=== RUN   TestSame\n--- PASS: TestSame (0.1s)\n", None)
        failures, successes, manifest = self.parse(content, ["./one", "./two"])
        self.assertEqual(failures, [])
        self.assertFalse(successes[0]["package_complete"])
        self.assertFalse(manifest["complete"])

    def test_failed_package_pass_is_not_resolution_evidence(self):
        _, successes, _ = self.parse(output("./one", "--- PASS: TestSame (0.1s)\n", 1), ["./one"])
        self.assertFalse(successes[0]["package_complete"])

    def test_build_failure_without_test_outcomes_is_incomplete(self):
        failures, successes, manifest = self.parse(
            output("./one", "FAIL github.com/example/one [build failed]\n", 1), ["./one"])
        self.assertEqual((failures, successes), ([], []))
        self.assertFalse(manifest["complete"])

    def test_skips_are_neither_passes_nor_panics(self):
        body = "=== RUN   TestSkip\n--- SKIP: TestSkip (0.0s)\nFAIL github.com/example/one\n"
        failures, successes, _ = self.parse(output("./one", body, 1), ["./one"])
        self.assertEqual((failures, successes), ([], []))

    def test_unframed_legacy_output_cannot_resolve(self):
        _, successes, manifest = self.parse("--- PASS: TestSame (0.1s)\n", ["./one"])
        self.assertEqual(successes[0]["package"], "")
        self.assertFalse(successes[0]["package_complete"])
        self.assertFalse(manifest["complete"])

    def test_go_children_inherit_memory_budget(self):
        with patch.dict(os.environ, {"GOMEMLIMIT": "10GiB"}):
            status = runner.run_command([sys.executable, "-c", "import os; print(os.environ['GOMEMLIMIT'])"], "child.log")
        self.assertEqual(status, 0)
        self.assertIn("\n10GiB\n", Path("child.log").read_text())


class MergeTests(TemporaryTests):
    def setUp(self):
        super().setUp()
        self.write("jobs.json", [job("success", "success")])
        self.artifact = Path("artifacts/test-results-one")

    def artifact_files(self):
        self.write(self.artifact / "test-failures.json", [])
        self.write(self.artifact / "test-successes.json", [result()])
        self.write(self.artifact / "test-run.json", {"schema_version": 1, "complete": True,
                   "run_id": os.environ.get("GITHUB_RUN_ID", ""), "run_attempt": os.environ.get("GITHUB_RUN_ATTEMPT", "")})

    def aggregate(self, **kwargs):
        return merger.aggregate(Path("artifacts"), Path("jobs.json"), **kwargs)

    def test_missing_artifacts_cli_always_writes_outputs(self):
        proc = subprocess.run([sys.executable, str(SCRIPTS / "merge_test_results.py"),
                              "missing", "failures.json", "passes.json", "jobs.json", "status.json"], capture_output=True)
        self.assertEqual(proc.returncode, 0, proc.stderr)
        self.assertEqual(json.loads(Path("failures.json").read_text()), [])
        self.assertEqual(json.loads(Path("passes.json").read_text()), [])
        self.assertFalse(json.loads(Path("status.json").read_text())["complete"])

    def test_complete_artifacts_and_expected_credential_skip(self):
        self.artifact_files()
        self.write("jobs.json", [job("success", "success"), job("success", "skipped")])
        _, successes, status = self.aggregate()
        self.assertTrue(status["complete"], status)
        self.assertEqual(status["expected_artifacts"], 1)
        self.assertEqual(len(successes), 1)

    def test_entire_missing_shard_blocks_resolution(self):
        self.artifact_files()
        self.write("jobs.json", [job("success", "success"), job("success", "success")])
        self.assertFalse(self.aggregate()[2]["complete"])

    def test_corrupt_missing_or_invalid_failure_array_blocks_resolution(self):
        for value in ("{", "{}", '["wrong"]', '[{"test_name": "TestSame"}]', None):
            with self.subTest(value=value):
                self.artifact_files()
                path = self.artifact / "test-failures.json"
                path.unlink() if value is None else path.write_text(value)
                _, successes, status = self.aggregate()
                self.assertEqual(len(successes), 1)
                self.assertFalse(status["complete"])

    def test_incomplete_or_stale_manifest_blocks_resolution(self):
        for manifest in ({"schema_version": 1, "complete": False},
                         {"schema_version": 1, "complete": True, "run_id": "old", "run_attempt": "1"}):
            self.artifact_files()
            self.write(self.artifact / "test-run.json", manifest)
            with patch.dict(os.environ, {"GITHUB_RUN_ID": "123", "GITHUB_RUN_ATTEMPT": "2"}):
                self.assertFalse(self.aggregate()[2]["complete"])

    def test_job_inventory_failure_and_failed_download_block_resolution(self):
        self.artifact_files()
        self.assertFalse(self.aggregate(download_ok=False)[2]["complete"])
        Path("jobs.json").unlink()
        self.assertFalse(self.aggregate()[2]["complete"])

    def test_interrupted_job_blocks_resolution(self):
        self.artifact_files()
        self.write("jobs.json", [job()])
        self.assertFalse(self.aggregate()[2]["complete"])


class IssueTests(TemporaryTests):
    def process(self, failures, successes, open_issues, complete=True):
        self.write("failures.json", failures)
        self.write("successes.json", successes)
        self.write("status.json", {"complete": complete})
        with patch.object(issues, "get_all_open_test_issues", return_value=open_issues), \
             patch.object(issues, "run_gh_command", return_value="url") as gh, \
             patch.object(issues.time, "sleep"):
            issues.process_test_failures("o", "r", "1", "failures.json", "successes.json", "status.json")
        return [call.args[0] for call in gh.call_args_list]

    def test_only_matching_package_and_test_explicit_pass_closes(self):
        calls = self.process([], [result()], [issue(1), issue(2, package="./internal/two"), issue(3, name="TestOther")])
        self.assertEqual([c[2] for c in calls if c[:2] == ["issue", "close"]], ["1"])

    def test_legacy_title_without_package_is_not_closed(self):
        self.assertEqual(self.process([], [result()], [{"number": 1, "title": "Bug: TestSame Failing"}]), [])

    def test_incomplete_reporting_or_package_cannot_close(self):
        self.assertEqual(self.process([], [result()], [issue()], complete=False), [])
        passed = result()
        passed["package_complete"] = False
        self.assertEqual(self.process([], [passed], [issue()]), [])

    def test_failure_overrides_pass_and_updates_matching_issue_only(self):
        calls = self.process([result(context="failed")], [result()], [issue(1, package="./internal/two"), issue(2)])
        self.assertEqual([c[2] for c in calls if c[:2] == ["issue", "comment"]], ["2"])
        self.assertFalse(any(c[:2] == ["issue", "close"] for c in calls))

    def test_new_issue_contains_identity_even_with_same_name_elsewhere(self):
        calls = self.process([result(context="failed")], [], [issue(package="./internal/two")])
        create = next(c for c in calls if c[:2] == ["issue", "create"])
        body = create[create.index("--body") + 1]
        self.assertEqual(issues.issue_identity({"body": body}), ("./internal/one", "TestSame"))

    def test_no_results_do_not_close_anything(self):
        self.assertEqual(self.process([], [], [issue()]), [])

    def test_closure_helper_rejects_no_explicit_pass(self):
        with patch.object(issues, "run_gh_command") as gh:
            issues.close_resolved_issue("o", "r", "1", "date", "run", "url", False)
        gh.assert_not_called()


class ReportingIntegrationTests(TemporaryTests):
    def test_parser_to_aggregation_to_issue_updates(self):
        artifact = Path("artifacts/test-results-example")
        artifact.mkdir(parents=True)
        previous = Path.cwd()
        os.chdir(artifact)
        try:
            content = output("./one", "=== RUN   TestSame\n--- PASS: TestSame (0.1s)\n")
            content += output("./two", "=== RUN   TestSame\n    example_test.go:7: broken\n--- FAIL: TestSame (0.1s)\n", 1)
            Path("output.log").write_text(content)
            runner.parse_test_results("output.log", "resources", "example", ["./one", "./two"])
        finally:
            os.chdir(previous)
        self.write("jobs.json", [job("success", "failure"), job("success", "skipped")])
        failures, successes, status = merger.aggregate(Path("artifacts"), Path("jobs.json"))
        self.assertTrue(status["complete"], status)
        for name, data in (("failures", failures), ("successes", successes), ("status", status)):
            self.write(f"{name}.json", data)
        with patch.object(issues, "get_all_open_test_issues", return_value=[issue(1, "./one"), issue(2, "./two"), issue(3, "./unrun")]), \
             patch.object(issues, "run_gh_command", return_value="url") as gh, patch.object(issues.time, "sleep"):
            issues.process_test_failures("o", "r", "1", "failures.json", "successes.json", "status.json")
        calls = [call.args[0] for call in gh.call_args_list]
        self.assertEqual([call[2] for call in calls if call[:2] == ["issue", "close"]], ["1"])
        self.assertEqual([call[2] for call in calls if call[:2] == ["issue", "comment"]], ["2"])


class DetectionTests(TemporaryTests):
    def test_cancelled_step_preserved_under_failed_job(self):
        failure = detector.analyze_job_failure(job())
        self.assertEqual(failure["failed_step"], "Run resource tests")
        self.assertEqual(failure["step_number"], 2)
        self.assertNotEqual(failure["failure_type"], "out_of_memory")

    def test_positive_oom_delta_proves_oom_and_does_not_double_count(self):
        log = "\n".join(f"{path}/memory.events: low 0 oom 0 oom_kill {value}" for value in (0, 1)
                        for path in ("/sys/fs/cgroup", "/sys/fs/cgroup/job"))
        failure = detector.analyze_job_failure(job(), log)
        self.assertEqual(failure["failure_type"], "out_of_memory")
        self.assertIn("increased by 1", failure["evidence"][0])

    def test_old_oom_counter_signal_or_step_name_does_not_prove_oom(self):
        for log in ("exit=143", "/sys/fs/cgroup/memory.events: oom_kill 1",
                    "/sys/fs/cgroup/memory.events: oom_kill 1\n/sys/fs/cgroup/memory.events: oom_kill 1"):
            candidate = job()
            candidate["steps"][1]["name"] = "Monitor memory and OOM"
            self.assertNotEqual(detector.analyze_job_failure(candidate, log)["failure_type"], "out_of_memory")

    def test_runner_loss_annotation_preserves_uncertainty(self):
        failure = detector.analyze_job_failure(job(), annotations=[{"message": "The hosted runner lost communication with the server."}])
        self.assertEqual(failure["failure_type"], "runner_failure")
        self.assertIn("unconfirmed", failure["evidence"][0])

    def test_timeout_and_cancellation_remain_distinct(self):
        self.assertEqual(detector.analyze_job_failure(job("timed_out"))["failure_type"], "timeout")
        self.assertEqual(detector.analyze_job_failure(job("cancelled"))["failure_type"], "cancelled")
        self.assertEqual(detector.analyze_job_failure(job(), annotations=[{"message": "The job exceeded the maximum execution time"}])["failure_type"], "timeout")

    def test_credential_skip_is_not_an_incident(self):
        self.assertIsNone(detector.analyze_job_failure(job("success", "skipped")))
        self.assertIsNone(detector.analyze_job_failure(job("skipped", "skipped")))

    def test_multiple_job_pages_and_current_attempt(self):
        with patch.dict(os.environ, {"GITHUB_RUN_ATTEMPT": "2"}), \
             patch.object(detector, "run_gh_command", return_value=json.dumps([{"jobs": [job()]}, {"jobs": [job()]}])) as gh:
            self.assertEqual(len(detector.get_workflow_jobs("o", "r", "1")), 2)
        self.assertIn("/attempts/2/jobs?per_page=100", gh.call_args.args[0][1])
        self.assertIn("--slurp", gh.call_args.args[0])

    def test_inventory_error_does_not_emit_false_empty_success(self):
        self.write("job-failures.json", [])
        self.write("workflow-jobs.json", [])
        with patch.object(detector, "get_workflow_jobs", side_effect=subprocess.CalledProcessError(1, "gh")):
            with self.assertRaises(subprocess.CalledProcessError):
                detector.detect_job_failures("o", "r", "1", "job-failures.json")
        self.assertFalse(Path("job-failures.json").exists())
        self.assertFalse(Path("workflow-jobs.json").exists())

    def test_missing_logs_still_reports_cancelled_step_and_annotation(self):
        with patch.object(detector, "get_workflow_jobs", return_value=[job()]), \
             patch.object(detector, "run_gh_command", side_effect=[subprocess.CalledProcessError(1, "gh"),
                 json.dumps([[{"message": "The runner lost communication with the server"}]])]):
            detector.detect_job_failures("o", "r", "1", "failures.json")
        failure = json.loads(Path("failures.json").read_text())[0]
        self.assertEqual(failure["failure_type"], "runner_failure")
        self.assertEqual(failure["failed_step"], "Run resource tests")
        self.assertIn("Job logs unavailable", failure["evidence"])

    def test_issue_body_and_recurrence_include_evidence(self):
        failure = detector.analyze_job_failure(job())
        with patch.object(job_issues, "run_gh_command", return_value="url") as gh:
            job_issues.create_job_failure_issue("o", "r", failure, "date", "1", "url")
            job_issues.update_job_failure_issue("o", "r", "7", failure, "date", "1", "url")
        bodies = [call.args[0][call.args[0].index("--body") + 1] for call in gh.call_args_list if "--body" in call.args[0]]
        self.assertEqual(len(bodies), 2)
        for body in bodies:
            self.assertIn(failure["evidence"][0], body)


if __name__ == "__main__":
    unittest.main()
