"""Run the workflow's report script with read-only and writable PR contexts."""

import json
import re
import subprocess
import unittest
from pathlib import Path

import yaml


class CoverageReportingTests(unittest.TestCase):
    """Exercise coverage reporting behavior."""

    def run_report(self, *, fork=False, actor="contributor", existing=False, draft=False):
        """Run report."""
        workflow = Path(__file__).resolve().parents[4] / ".github/workflows/pr-tests.yml"
        steps = yaml.safe_load(workflow.read_text(encoding="utf-8"))["jobs"]["unit-tests"]["steps"]
        script = next(
            step["with"]["script"]
            for step in steps
            if step["name"] == "💬 Comment coverage summary on PR"
        )
        enforce_outputs = "steps.enforce-coverage.outputs"
        calculate_outputs = "steps.calculate-coverage.outputs"
        values = {
            "steps.enforce-coverage.outputs.coverage-pct": "" if draft else "71.9",
            "steps.calculate-coverage.outputs.coverage-pct": "55.3",
            f"{enforce_outputs}.total-lines || {calculate_outputs}.total-lines": "121",
            f"{enforce_outputs}.covered-lines || {calculate_outputs}.covered-lines": "87",
            "needs.plan-tests.outputs.service-domains": "device_management",
            "steps.pr-state.outputs.is-draft": "true" if draft else "false",
            "github.repository": "example/provider",
            "github.event.pull_request.number": "123",
        }
        script = re.sub(r"\$\{\{\s*(.*?)\s*\}\}", lambda match: values[match.group(1)], script)
        context = {
            "actor": actor,
            "repo": {"owner": "example", "repo": "provider"},
            "issue": {"number": 123},
            "payload": {
                "repository": {"full_name": "example/provider"},
                "pull_request": {
                    "head": {
                        "repo": {
                            "full_name": "contributor/provider" if fork else "example/provider"
                        }
                    }
                },
            },
        }
        comments = (
            [{"id": 1, "user": {"type": "Bot"}, "body": "Unit Test Coverage Report"}]
            if existing
            else []
        )
        harness = f"""
const calls = [];
let summary = '';
const core = {{
  info: () => {{}},
  summary: {{addRaw: text => {{ summary = text; return {{write: async () => {{}}}}; }}}}
}};
const github = {{rest: {{issues: {{
  listComments: async () => {{ calls.push('list'); return {{data: {json.dumps(comments)}}}; }},
  createComment: async () => {{ calls.push('create'); }},
  updateComment: async () => {{ calls.push('update'); }}
}}}}}};
const AsyncFunction = Object.getPrototypeOf(async function(){{}}).constructor;
new AsyncFunction('github', 'context', 'core', {json.dumps(script)})(
  github, {json.dumps(context)}, core)
  .then(() => console.log(JSON.stringify({{calls, summary}})))
  .catch(error => {{ console.error(error); process.exitCode = 1; }});
"""
        result = subprocess.run(["node", "-e", harness], check=True, capture_output=True, text=True)
        return json.loads(result.stdout)

    def test_fork_publishes_job_summary_without_github_api_calls(self):
        """Verify fork publishes job summary without github api calls."""
        result = self.run_report(fork=True)
        self.assertEqual(result["calls"], [])
        self.assertIn("Coverage Passed", result["summary"])
        self.assertIn("Changed Lines (Codecov Patch):** 71.9%", result["summary"])
        self.assertNotIn("55.3%", result["summary"])

    def test_dependabot_publishes_job_summary_without_github_api_calls(self):
        """Verify dependabot publishes job summary without github api calls."""
        result = self.run_report(actor="dependabot[bot]")
        self.assertEqual(result["calls"], [])
        self.assertIn("71.9%", result["summary"])

    def test_same_repository_keeps_creating_and_updating_comments(self):
        """Verify same repository keeps creating and updating comments."""
        for existing, expected in [(False, "create"), (True, "update")]:
            with self.subTest(existing=existing):
                result = self.run_report(existing=existing)
                self.assertEqual(result["calls"], ["list", expected])
                self.assertIn("Coverage Passed", result["summary"])

    def test_draft_package_coverage_remains_informational(self):
        """Verify draft package coverage remains informational."""
        result = self.run_report(fork=True, draft=True)
        self.assertEqual(result["calls"], [])
        self.assertIn("Coverage is informational only", result["summary"])
        self.assertIn("Coverage for Changed Packages:** 55.3%", result["summary"])


if __name__ == "__main__":
    unittest.main()
