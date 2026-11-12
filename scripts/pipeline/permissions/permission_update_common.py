"""Shared discovery and dry-run CLI behavior for permission update tools."""

import argparse
import re
from pathlib import Path


def find_repo_root(start: Path) -> Path:
    """Walk up to go.mod, retaining the current directory as the fallback."""
    current = start.resolve()
    for _ in range(10):
        if (current / "go.mod").exists():
            return current
        parent = current.parent
        if parent == current:
            break
        current = parent
    return Path.cwd()


def find_constructor_file(resource_dir: Path) -> Path | None:
    """Find the first non-test Go file containing a service constructor."""
    pattern = re.compile(r"\bfunc\s+New\w+\s*\([^)]*\)\s+\w+(?:\.\w+)?\s*\{", re.MULTILINE)
    for go_file in sorted(resource_dir.glob("*.go")):
        if go_file.name.endswith("_test.go"):
            continue
        try:
            text = go_file.read_text(encoding="utf-8")
        except OSError:
            continue
        if pattern.search(text):
            return go_file
    return None


def build_update_parser(description: str, target: str, metadata: bool = False):
    """Keep the source/template update tools' shared flags and dry-run defaults aligned."""
    parser = argparse.ArgumentParser(
        description=description,
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=(
            "Examples:\n"
            "  %(prog)s                          # dry-run all domains\n"
            "  %(prog)s --apply                  # apply all changes\n"
            "  %(prog)s --domain applications --apply\n"
            "  %(prog)s --category resources --apply\n"
        ),
    )
    parser.add_argument(
        "--root", metavar="PATH", help="Repository root (default: auto-detected via go.mod)."
    )
    parser.add_argument(
        "--apply", action="store_true", help=f"Write changes to {target} files (default: dry-run)."
    )
    parser.add_argument(
        "--domain", metavar="DOMAIN", help="Limit to this service domain (e.g. device_management)."
    )
    parser.add_argument(
        "--category",
        metavar="CATEGORY",
        choices=sorted({"actions", "datasources", "ephemerals", "list-resources", "resources"}),
        help="Limit to this service category.",
    )
    if metadata:
        parser.add_argument(
            "--refresh-metadata",
            action="store_true",
            dest="refresh_metadata",
            help="Force re-download of Graph command metadata.",
        )
    return parser


def scan_update_domains(scan, services_root: Path, args):
    """Scan the requested domains and display the selected apply/dry-run mode."""
    print("Scanning Go source tree …", flush=True)
    domains = scan(services_root, category_filter=args.category, domain_filter=args.domain)
    total = sum(len(domain.resources) for domain in domains)
    print(f"  {len(domains)} domains, {total} resources found.")
    mode = "APPLY" if args.apply else "DRY-RUN"
    print(f"\n  Mode: {mode}\n")
    return domains


def update_scanned_resources(args, domains, process, print_result, skip_unmatched=False) -> int:
    """Apply the selected updater and report its changed, skipped, and failed counts."""
    changed = errors = skipped = 0
    for domain in domains:
        for resource in domain.resources:
            if skip_unmatched and not resource.endpoints:
                skipped += 1
                continue
            result = process(resource, domain.domain)
            print_result(result)
            if result.status in ("updated", "would update"):
                changed += 1
            elif "error" in result.status:
                errors += 1
            else:
                skipped += 1
    print(f"\n  Changed: {changed}  |  Skipped/no-change: {skipped}  |  Errors: {errors}")
    if not args.apply:
        print("\n  (dry-run — re-run with --apply to write changes)")
    return 0 if errors == 0 else 1
