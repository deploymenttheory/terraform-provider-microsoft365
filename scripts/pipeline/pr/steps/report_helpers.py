#!/usr/bin/env python3
"""Report package-local helpers in changed service packages. Report-only: never fails the job."""

import argparse
import os
import subprocess
import sys
import tempfile
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from lib.git_operations import get_changed_packages, get_merge_base  # noqa: E402

# pylint: enable=wrong-import-position

SERVICES_PREFIX = "internal/services/"
SCANNER = "./tools/helperscan"
SKIPPED = "No changed Go packages under internal/services; helper report skipped."


def service_packages(packages):
    """Keep the changed packages that live under internal/services."""
    return [package for package in packages if package.startswith(SERVICES_PREFIX)]


def append_summary(summary_file, markdown):
    """Append Markdown to the GitHub step summary when one is configured."""
    if not summary_file:
        return
    with open(summary_file, "a", encoding="utf-8", newline="\n") as handle:
        handle.write(markdown.rstrip("\n") + "\n")


def scan(packages, summary_file, runner=subprocess.run):
    """Build the scanner once, print the text report and add the Markdown report to the summary."""
    with tempfile.TemporaryDirectory() as directory:
        binary = Path(directory) / ("helperscan.exe" if os.name == "nt" else "helperscan")
        runner(["go", "build", "-o", str(binary), SCANNER], check=True)
        text = runner(
            [str(binary), "-format", "text", *packages], check=True, capture_output=True, text=True
        )
        print(text.stdout, end="", flush=True)
        if text.stderr:
            print(text.stderr, end="", file=sys.stderr, flush=True)
        markdown = runner(
            [str(binary), "-format", "markdown", *packages],
            check=True,
            capture_output=True,
            text=True,
        )
        append_summary(summary_file, markdown.stdout)


def report_helpers(base_ref, summary_file=None, runner=subprocess.run):
    """Run the helper report on changed service packages; any failure becomes a warning."""
    try:
        packages = service_packages(get_changed_packages(get_merge_base(base_ref)))
        if not packages:
            print(SKIPPED, flush=True)
            append_summary(summary_file, f"## Helper report\n\n{SKIPPED}\n")
            return 0
        print(f"Helper report for {len(packages)} changed service package(s)", flush=True)
        scan(packages, summary_file, runner)
    # The report is informational, so no failure of git, go or the scanner may fail the job.
    except Exception as error:  # pylint: disable=broad-exception-caught
        print(f"::warning title=Helper report::The report could not run: {error}", flush=True)
    return 0


def main():
    """Parse arguments and run the report."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-ref", required=True)
    args = parser.parse_args()
    return report_helpers(args.base_ref, os.environ.get("GITHUB_STEP_SUMMARY"))


if __name__ == "__main__":
    sys.exit(main())
