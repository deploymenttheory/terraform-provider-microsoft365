#!/usr/bin/env python3
"""Compile the selected tests' dependencies without executing tests or provider code."""

import argparse
import os
import subprocess

BATCH_SIZE = 64 if os.name == "nt" else 256
IMPORTS_TEMPLATE = (
    '{{join .Imports "\\n"}}{{"\\n"}}{{join .TestImports "\\n"}}{{"\\n"}}{{join .XTestImports'
    ' "\\n"}}'
)


def prepare_go_cache(packages, race=False):
    """Prime flag-aware Go build artifacts for only the selected packages' imports."""
    if not packages:
        raise ValueError("Cache preparation requires selected packages")
    flags = ["-mod=readonly", "-p=1"] + (["-race"] if race else [])
    dependencies = set()
    for start in range(0, len(packages), BATCH_SIZE):
        result = subprocess.run(
            [
                "go",
                "list",
                *flags,
                "-f",
                IMPORTS_TEMPLATE,
                *[f"./{package}" for package in packages[start : start + BATCH_SIZE]],
            ],
            check=True,
            stdout=subprocess.PIPE,
            text=True,
            encoding="utf-8",
        )
        dependencies.update(result.stdout.split())
    ordered = sorted(dependencies - {"C", "unsafe"})
    print(
        f"Compiling {len(ordered)} dependency packages ({'race' if race else 'normal'}).",
        flush=True,
    )
    for start in range(0, len(ordered), BATCH_SIZE):
        subprocess.run(["go", "build", *flags, *ordered[start : start + BATCH_SIZE]], check=True)


def main():
    """Prepare dependency artifacts; actual tests remain the downstream workers' job."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--packages", required=True, help="Exact planned test package directories")
    parser.add_argument(
        "--race", action="store_true", help="Compile race-instrumented dependencies"
    )
    args = parser.parse_args()
    prepare_go_cache(args.packages.split(), args.race)


if __name__ == "__main__":
    main()
