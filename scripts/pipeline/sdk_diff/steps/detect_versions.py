#!/usr/bin/env python3
"""Detect current and latest SDK versions.

Outputs:
- current-msgraph-sdk: Current msgraph-sdk-go version
- current-msgraph-beta-sdk: Current msgraph-beta-sdk-go version
- latest-msgraph-sdk: Latest available msgraph-sdk-go version
- latest-msgraph-beta-sdk: Latest available msgraph-beta-sdk-go version
"""

import argparse
import json
import sys
import urllib.error
from pathlib import Path

# Add lib directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "lib"))
# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from github_api import get_latest_release, get_sdk_repo_name  # pylint: disable=import-error

# noqa: E402
from version_detector import (  # pylint: disable=import-error
    format_version_display,
    get_msgraph_versions,
)

# pylint: enable=wrong-import-position


def main():
    """Detect current and latest SDK versions and output comparison."""
    parser = argparse.ArgumentParser(description="Detect SDK versions")
    parser.add_argument(
        "--repo-path", type=Path, default=Path.cwd(), help="Path to repository root"
    )
    parser.add_argument("--output-file", help="GitHub Actions output file")

    args = parser.parse_args()

    print("📦 Detecting current SDK versions...")
    go_mod_path = args.repo_path / "go.mod"
    current_versions = get_msgraph_versions(go_mod_path)

    print("\n✅ Current versions:")
    print(format_version_display(current_versions))

    print("\n🔍 Getting latest SDK versions from GitHub...")
    latest_versions = {}

    for sdk_name, current_version in current_versions.items():
        try:
            repo = get_sdk_repo_name(sdk_name)
            release = get_latest_release(repo)
            latest_versions[sdk_name] = release["tag"]

            if release["tag"] == current_version:
                status = "✅ (up to date)"
            else:
                status = f"⚠️  (update available: {current_version} → {release['tag']})"

            print(f"  - {sdk_name}: {release['tag']} {status}")

        except (
            RuntimeError,
            urllib.error.URLError,
            json.JSONDecodeError,
            KeyError,
            ValueError,
        ) as e:
            print(f"  ❌ Failed to get latest for {sdk_name}: {e}")
            latest_versions[sdk_name] = current_version

    if args.output_file:
        with open(args.output_file, "a", encoding="utf-8") as f:
            f.write(f"current-msgraph-sdk={current_versions.get('msgraph-sdk-go', '')}\n")
            f.write(f"current-msgraph-beta-sdk={current_versions.get('msgraph-beta-sdk-go', '')}\n")
            f.write(f"latest-msgraph-sdk={latest_versions.get('msgraph-sdk-go', '')}\n")
            f.write(f"latest-msgraph-beta-sdk={latest_versions.get('msgraph-beta-sdk-go', '')}\n")

            has_updates = any(
                latest_versions.get(sdk) != current_versions.get(sdk)
                for sdk in current_versions.keys()
            )
            f.write(f"has-updates={'true' if has_updates else 'false'}\n")

    print("\n✅ Version detection complete")


if __name__ == "__main__":
    main()
