#!/usr/bin/env python3
"""Generate dynamic GitHub Actions matrix with automatic shard calculation.

This script discovers resources/datasources/actions/list-resources/ephemerals
in each service and automatically determines the optimal number of shards using
Rendezvous (Highest Random Weight) distribution for stable, deterministic
assignment.

Usage:
    ./generate_test_matrix.py <configuration_block_type> [service1] [service2 ...]

Configuration Block Types:
    resources, datasources, actions, list-resources, ephemerals

Service Discovery:
    - If no services are specified, all services will be auto-discovered
    - Services with no tests are automatically skipped
    - Explicit service names can be provided to limit scope

Examples:
    # Auto-discover all services
    ./generate_test_matrix.py resources
    ./generate_test_matrix.py datasources

    # Explicit service list
    ./generate_test_matrix.py resources device_management groups users windows_updates
    ./generate_test_matrix.py datasources device_management utility
    ./generate_test_matrix.py list-resources device_management
    ./generate_test_matrix.py ephemerals multitenant_management

Output:
    JSON matrix suitable for GitHub Actions matrix strategy
"""

import hashlib
import json
import os
import re
import struct
import sys
from pathlib import Path
from typing import Dict, List

ITEMS_PER_SHARD = int(os.environ.get("ITEMS_PER_SHARD", "5"))
TESTS_PER_SHARD = int(os.environ.get("TESTS_PER_SHARD", "25"))
if ITEMS_PER_SHARD < 1 or TESTS_PER_SHARD < 1:
    raise ValueError("ITEMS_PER_SHARD and TESTS_PER_SHARD must be positive integers")
RENDEZVOUS_SEED = os.environ.get("RENDEZVOUS_SEED", "terraform-provider-microsoft365")
VERBOSE = os.environ.get("VERBOSE", "").lower() in ("true", "1", "yes")


def discover_test_names(package_dir: Path) -> List[str]:
    """Inventory top-level tests, fuzz seeds and examples without compiling providers."""
    names = set()
    for source in package_dir.glob("*_test.go"):
        # Remove comments and string literals, preserving newlines for anchors.
        content = re.sub(
            r'//[^\n]*|/\*[\s\S]*?\*/|"(?:\\.|[^"\\])*"|`[^`]*`|\'(?:\\.|[^\'\\])*\'',
            lambda match: "\n" * match.group().count("\n"),
            source.read_text(encoding="utf-8"),
        )
        names.update(
            re.findall(
                r"^func\s+((?:Test|Fuzz)(?:[A-Z_0-9]\w*)?)\s*\(\s*\w+\s+\*testing\.[TF]\s*,?\s*\)",
                content,
                re.MULTILINE,
            )
        )
        names.update(re.findall(r"^func\s+(Example\w*)\s*\(\s*\)", content, re.MULTILINE))
    return sorted(names)


def discover_services(configuration_block_type: str) -> List[str]:
    """Auto-discover all services that have tests for the given configuration_block_type.

    Args:
        configuration_block_type: Category type (resources, datasources, actions, etc.)

    Returns:
        List of service names (e.g., ['device_management', 'identity_and_access'])
    """
    base_path = Path(f"./internal/services/{configuration_block_type}")

    if not base_path.exists():
        return []

    services = []
    for service_dir in base_path.iterdir():
        if service_dir.is_dir() and not service_dir.name.startswith("."):
            # Check if this service actually has any resources
            if discover_resources(service_dir.name, configuration_block_type):
                services.append(service_dir.name)

    return sorted(services)


def _package_directories(base_path: Path, category: str):
    """Yield package directories using the service/category layout rules."""
    for version in base_path.iterdir():
        if not version.is_dir():
            continue
        if not version.name.startswith("graph_"):
            yield version
            continue
        for resource in version.iterdir():
            if not resource.is_dir():
                continue
            if category == "actions":
                yield from (action for action in resource.iterdir() if action.is_dir())
            else:
                yield resource


def discover_resources(
    service: str, configuration_block_type: str, verbose: bool = False
) -> List[str]:
    """Discover implemented test packages without including incomplete resources."""
    base_path = Path(f"./internal/services/{configuration_block_type}/{service}")
    if not base_path.exists():
        return []
    impl_file_markers = {
        "resources": ["resource.go"],
        "datasources": ["datasource.go"],
        "actions": ["action.go"],
        "list-resources": ["list_resource.go"],
        "ephemerals": ["ephemeral_resource.go"],
    }

    # Only discover packages that have test files
    # This ensures we don't include incomplete implementations
    test_file_markers = {
        "resources": ["resource_test.go"],
        "datasources": ["datasource_test.go"],
        "actions": ["action_test.go"],
        "list-resources": ["list_resource_test.go"],
        "ephemerals": ["ephemeral_resource_test.go"],
    }

    impl_markers = impl_file_markers.get(
        configuration_block_type, ["resource.go", "datasource.go", "action.go"]
    )
    test_markers = test_file_markers.get(
        configuration_block_type, ["resource_test.go", "datasource_test.go", "action_test.go"]
    )
    packages, skipped = [], []
    for directory in _package_directories(base_path, configuration_block_type):
        relative = directory.relative_to(base_path).as_posix()
        if any((directory / marker).exists() for marker in test_markers):
            packages.append(relative)
        elif verbose and any((directory / marker).exists() for marker in impl_markers):
            skipped.append(relative)
    if verbose and skipped:
        print(
            f"⚠️  [DISCOVERY] Skipped {len(skipped)} package(s) in {service} (no tests):",
            file=sys.stderr,
        )
        for package in sorted(skipped):
            print(f"    - {package}", file=sys.stderr)
    return sorted(packages)


def rendezvous_hash(item: str, shard_index: int, seed: str) -> int:
    """Compute Rendezvous (HRW) hash for item-shard pair.

    Uses SHA-256 for deterministic weight calculation.
    Higher weight = higher priority for assignment.

    Args:
        item: Item to hash (resource package path)
        shard_index: Shard index (0-based)
        seed: Seed string for deterministic hashing

    Returns:
        64-bit unsigned integer weight
    """
    # Why this format? Matches Go implementation for consistency
    input_str = f"{item}:shard_{shard_index}:{seed}"
    hash_bytes = hashlib.sha256(input_str.encode("utf-8")).digest()

    # Why big-endian? Matches Go's binary.BigEndian.Uint64
    weight = struct.unpack(">Q", hash_bytes[:8])[0]
    return weight


def distribute_rendezvous(
    items: List[str],
    shard_count: int,
    seed: str,
    weights: Dict[str, int] = None,
) -> Dict[int, List[str]]:
    """Distribute items using Rendezvous (HRW) algorithm.

    Each item prefers the eligible shard with the highest hash. Optional test
    weights enforce test and item budgets. Capacity constraints can move more
    items than unconstrained Rendezvous when the inventory changes.

    Args:
        items: List of items to distribute (resource paths)
        shard_count: Number of shards to distribute across
        seed: Seed for deterministic distribution

    Returns:
        Dictionary mapping shard_index -> list of items
    """
    if shard_count <= 0:
        shard_count = 1

    shards = {i: [] for i in range(shard_count)}

    loads = {i: 0 for i in range(shard_count)}
    for item in sorted(items, key=lambda item: (-(weights or {}).get(item, 1), item)):
        cost = (weights or {}).get(item, 1)
        candidates = [
            index
            for index in shards
            if weights is None
            or (len(shards[index]) < ITEMS_PER_SHARD and loads[index] + cost <= TESTS_PER_SHARD)
        ]
        if not candidates:
            raise ValueError("Insufficient shard capacity")
        selected_shard = max(
            candidates, key=lambda index, item=item: rendezvous_hash(item, index, seed)
        )
        shards[selected_shard].append(item)
        loads[selected_shard] += cost

    return shards


def calculate_shard_count(item_count: int) -> int:
    """Calculate shard count based on items per shard.

    Always shard consistently - every ITEMS_PER_SHARD items = 1 shard.
    The final shard may contain fewer items, which is expected and acceptable.

    Args:
        item_count: Number of test packages/items in the service
                   (resources, datasources, actions, list-resources, or ephemerals)

    Returns:
        Number of shards (minimum 1)

    Examples:
        - 1 item with 5/shard = 1 shard (1 item)
        - 5 items with 5/shard = 1 shard (5 items)
        - 6 items with 5/shard = 2 shards (5+1 items)
        - 9 items with 5/shard = 2 shards (5+4 items)
        - 11 items with 5/shard = 3 shards (5+5+1 items)
    """
    # Why ceiling division? Ensures all items are included
    # and last shard can have fewer items
    shard_count = (item_count + ITEMS_PER_SHARD - 1) // ITEMS_PER_SHARD
    return max(1, shard_count)


def _service_test_units(service: str, category: str, items: list[str]) -> dict:
    """Split oversized packages into independently filtered test batches."""
    units = {}
    for item in items:
        names = discover_test_names(Path("internal/services") / category / service / item)
        if len(names) > TESTS_PER_SHARD:
            for start in range(0, len(names), TESTS_PER_SHARD):
                batch = names[start : start + TESTS_PER_SHARD]
                units[f"{item}::batch-{start // TESTS_PER_SHARD}"] = (
                    item,
                    "^(" + "|".join(batch) + ")$",
                    len(batch),
                )
        else:
            units[item] = (item, "", max(1, len(names)))
    return units


def _distribute_test_units(units: dict) -> list:
    """Find enough bounded Rendezvous capacity for indivisible package batches."""
    weights = {identity: unit[2] for identity, unit in units.items()}
    count = max(
        calculate_shard_count(len(units)),
        (sum(weights.values()) + TESTS_PER_SHARD - 1) // TESTS_PER_SHARD,
    )
    while count <= len(units):
        try:
            shards = distribute_rendezvous(list(units), count, RENDEZVOUS_SEED, weights)
            return [shard for shard in shards.values() if shard]
        except ValueError:
            count += 1
    raise ValueError("Unable to fit test batches within shard budgets")


def _service_matrix_entry(metadata: dict, index: int, identities: list, units: dict) -> dict:
    """Describe one shard without losing its package-specific test filters."""
    return {
        **metadata,
        "shard_index": index,
        "item_count": len(identities),
        "resources": ",".join(units[identity][0] for identity in identities),
        "test_filters": json.dumps(
            {
                units[identity][0]: units[identity][1]
                for identity in identities
                if units[identity][1]
            },
            sort_keys=True,
        ),
        "test_count": sum(units[identity][2] for identity in identities),
        "shard_label": (
            "" if metadata["total_shards"] == 1 else f" ({index + 1}/{metadata['total_shards']})"
        ),
    }


def generate_matrix_for_service(
    service: str, configuration_block_type: str, runner: str = "ubuntu-24.04-arm"
) -> List[Dict]:
    """Generate bounded matrix entries for every test package in one service."""
    items = discover_resources(service, configuration_block_type, verbose=VERBOSE)
    if not items:
        return []
    units = _service_test_units(service, configuration_block_type, items)
    shards = _distribute_test_units(units)
    metadata = {"service": service, "runner": runner, "total_shards": len(shards)}
    return [
        _service_matrix_entry(metadata, index, identities, units)
        for index, identities in enumerate(shards)
    ]


def generate_matrix(configuration_block_type: str, services: List[str]) -> Dict:
    """Generate complete matrix for all services.

    Args:
        configuration_block_type: Category type
        services: List of service names to process

    Returns:
        Matrix dictionary suitable for fromJSON() in GitHub Actions
    """
    all_entries = []

    for service in services:
        entries = generate_matrix_for_service(service, configuration_block_type)
        all_entries.extend(entries)

    if len(all_entries) > 256:
        raise ValueError("Matrix exceeds 256 jobs; increase TESTS_PER_SHARD or ITEMS_PER_SHARD")

    return {"include": all_entries}


def main():
    """Main entry point."""
    if len(sys.argv) < 2:
        print(
            "Usage: generate_test_matrix.py <configuration_block_type> [service1] [service2 ...]",
            file=sys.stderr,
        )
        print(
            "\nConfiguration Block Types: resources, datasources, actions, "
            "list-resources, ephemerals",
            file=sys.stderr,
        )
        print(
            "\nService Discovery: If no services specified, all services will be auto-discovered",
            file=sys.stderr,
        )
        print("\nEnvironment variables:", file=sys.stderr)
        print(
            f"  TESTS_PER_SHARD             Test budget per shard (default: {TESTS_PER_SHARD})",
            file=sys.stderr,
        )
        print(
            f"  ITEMS_PER_SHARD             Items per shard (default: {ITEMS_PER_SHARD})",
            file=sys.stderr,
        )
        print(
            f"  RENDEZVOUS_SEED             Seed for distribution (default: {RENDEZVOUS_SEED})",
            file=sys.stderr,
        )
        print(
            "  VERBOSE                     Show packages skipped (no tests) (default: false)",
            file=sys.stderr,
        )
        sys.exit(1)

    configuration_block_type = sys.argv[1]
    services = sys.argv[2:]

    valid_configuration_block_types = [
        "resources",
        "datasources",
        "actions",
        "list-resources",
        "ephemerals",
    ]
    if configuration_block_type not in valid_configuration_block_types:
        print(
            f"Error: Invalid configuration_block_type '{configuration_block_type}'",
            file=sys.stderr,
        )
        print(
            f"Valid configuration block types: {', '.join(valid_configuration_block_types)}",
            file=sys.stderr,
        )
        sys.exit(1)

    # Auto-discover services if none specified
    if not services:
        print(
            "[generate_test_matrix] No services specified, "
            f"auto-discovering services for {configuration_block_type}...",
            file=sys.stderr,
        )
        services = discover_services(configuration_block_type)
        if not services:
            print(
                f"[generate_test_matrix] No services found for {configuration_block_type}",
                file=sys.stderr,
            )
            # Return empty matrix
            print(json.dumps({"include": []}))
            sys.exit(0)
        print(
            f"[generate_test_matrix] Discovered {len(services)} service(s): {', '.join(services)}",
            file=sys.stderr,
        )

    # Keeps diagnostic output separate from JSON stdout
    print(
        f"[generate_test_matrix] Configuration block type: {configuration_block_type}",
        file=sys.stderr,
    )
    print(f"[generate_test_matrix] Services: {', '.join(services)}", file=sys.stderr)
    print(f"[generate_test_matrix] Items per shard: {ITEMS_PER_SHARD}", file=sys.stderr)
    print(f"[generate_test_matrix] Rendezvous seed: {RENDEZVOUS_SEED}", file=sys.stderr)

    print(f"[generate_test_matrix] Tests per shard: {TESTS_PER_SHARD}", file=sys.stderr)
    matrix = generate_matrix(configuration_block_type, services)

    print(
        f"\n[generate_test_matrix] Generated {len(matrix['include'])} matrix entries:",
        file=sys.stderr,
    )
    for entry in matrix["include"]:
        if entry["total_shards"] > 1:
            shard_info = f"shard {entry['shard_index'] + 1}/{entry['total_shards']}"
        else:
            shard_info = "no sharding"

        print(
            f"  - {entry['service']}: {entry['item_count']} "
            f"items, {entry['test_count']} discovered tests ({shard_info})",
            file=sys.stderr,
        )

    # GitHub Actions reads matrix from stdout
    print(json.dumps(matrix))


if __name__ == "__main__":
    main()
