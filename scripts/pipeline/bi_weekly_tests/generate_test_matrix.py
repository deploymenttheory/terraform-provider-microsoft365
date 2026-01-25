#!/usr/bin/env python3
"""Generate dynamic GitHub Actions matrix with automatic shard calculation.

This script discovers resources/datasources/actions/list-resources/ephemerals
in each service and automatically determines the optimal number of shards using
Rendezvous (Highest Random Weight) distribution for stable, deterministic
assignment.

Usage:
    ./generate_test_matrix.py <configuration_block_type> <service1> [service2 ...]

Configuration Block Types:
    resources, datasources, actions, list-resources, ephemerals

Examples:
    ./generate_test_matrix.py resources device_management groups users
    ./generate_test_matrix.py datasources device_management utility
    ./generate_test_matrix.py list-resources device_management
    ./generate_test_matrix.py ephemerals multitenant_management

Output:
    JSON matrix suitable for GitHub Actions matrix strategy
"""

import os
import sys
import json
import hashlib
import struct
from pathlib import Path
from typing import Dict, List

RESOURCES_PER_SHARD = int(os.environ.get('RESOURCES_PER_SHARD', '5'))
MIN_RESOURCES_FOR_SHARDING = int(
    os.environ.get('MIN_RESOURCES_FOR_SHARDING', '10')
)
RENDEZVOUS_SEED = os.environ.get(
    'RENDEZVOUS_SEED',
    'terraform-provider-microsoft365'
)


def discover_resources(service: str, configuration_block_type: str) -> List[str]:
    """Discover all packages in a service for the given configuration_block_type.

    Args:
        service: Service name (e.g., 'device_management')
        configuration_block_type: Category type (resources, datasources, actions, etc.)

    Returns:
        List of package paths relative to the service directory
    """
    base_path = Path(f"./internal/services/{configuration_block_type}/{service}")

    if not base_path.exists():
        return []

    packages = []

    marker_files = {
        'resources': ['resource.go'],
        'datasources': ['datasource.go'],
        'actions': ['action.go'],
        'list-resources': ['list_resource.go'],
        'ephemerals': ['ephemeral_resource.go']
    }

    markers = marker_files.get(
        configuration_block_type,
        ['resource.go', 'datasource.go', 'action.go']
    )

    for graph_version_dir in base_path.iterdir():
        if not graph_version_dir.is_dir():
            continue

        # Why check for graph_? Ephemerals can exist at service/name
        # or service/graph_beta/name unlike other types
        is_graph_version = graph_version_dir.name.startswith('graph_')

        if is_graph_version:
            for resource_dir in graph_version_dir.iterdir():
                if not resource_dir.is_dir():
                    continue

                has_marker = any(
                    (resource_dir / marker).exists()
                    for marker in markers
                )

                if has_marker:
                    rel_path = str(resource_dir.relative_to(base_path))
                    packages.append(rel_path)
        else:
            has_marker = any(
                (graph_version_dir / marker).exists()
                for marker in markers
            )

            if has_marker:
                rel_path = str(graph_version_dir.relative_to(base_path))
                packages.append(rel_path)

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
    hash_bytes = hashlib.sha256(input_str.encode('utf-8')).digest()

    # Why big-endian? Matches Go's binary.BigEndian.Uint64
    weight = struct.unpack('>Q', hash_bytes[:8])[0]
    return weight


def distribute_rendezvous(
    items: List[str],
    shard_count: int,
    seed: str
) -> Dict[int, List[str]]:
    """Distribute items using Rendezvous (HRW) algorithm.

    Why Rendezvous? It minimizes item reassignment when shard count
    changes. Only ~1/n items move when adding a new shard, unlike
    modulo-based methods which can reassign most items.

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

    for item in items:
        highest_weight = 0
        selected_shard = 0

        for shard_idx in range(shard_count):
            weight = rendezvous_hash(item, shard_idx, seed)
            if weight > highest_weight:
                highest_weight = weight
                selected_shard = shard_idx

        shards[selected_shard].append(item)

    return shards


def calculate_shard_count(resource_count: int) -> int:
    """Calculate optimal shard count.

    Why threshold? Small services don't benefit from sharding overhead.
    Only shard when we have enough resources to justify parallel runners.

    Args:
        resource_count: Number of resources in the service

    Returns:
        Optimal number of shards (minimum 1)
    """
    if resource_count < MIN_RESOURCES_FOR_SHARDING:
        return 1

    # Why ceiling division? Ensures we don't exceed target per shard
    shard_count = (
        (resource_count + RESOURCES_PER_SHARD - 1) // RESOURCES_PER_SHARD
    )
    return max(1, shard_count)


def generate_matrix_for_service(
    service: str,
    configuration_block_type: str,
    runner: str = "ubuntu-24.04-arm"
) -> List[Dict]:
    """Generate matrix entries for a single service.

    Args:
        service: Service name
        configuration_block_type: Category type
        runner: GitHub Actions runner type

    Returns:
        List of matrix entry dictionaries
    """
    resources = discover_resources(service, configuration_block_type)
    resource_count = len(resources)

    if resource_count == 0:
        # Why skip? Avoids creating empty jobs in GitHub Actions
        return []

    shard_count = calculate_shard_count(resource_count)
    shards = distribute_rendezvous(resources, shard_count, RENDEZVOUS_SEED)

    matrix_entries = []
    for shard_idx in range(shard_count):
        shard_resources = shards[shard_idx]

        entry = {
            "service": service,
            "runner": runner,
            "shard_index": shard_idx,
            "total_shards": shard_count,
            "resource_count": len(shard_resources),
            # Why comma-separated? Easy to parse in bash
            "resources": ",".join(shard_resources)
        }

        # Why shard_label? GitHub Actions UI shows job names clearly
        if shard_count == 1:
            entry["shard_label"] = ""
        else:
            entry["shard_label"] = f" ({shard_idx + 1}/{shard_count})"

        matrix_entries.append(entry)

    return matrix_entries


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

    return {"include": all_entries}


def main():
    """Main entry point."""
    if len(sys.argv) < 3:
        print(
            "Usage: generate_test_matrix.py <configuration_block_type> <service1> "
            "[service2 ...]",
            file=sys.stderr
        )
        print(
            "\nConfiguration Block Types: resources, datasources, actions, "
            "list-resources, ephemerals",
            file=sys.stderr
        )
        print("\nEnvironment variables:", file=sys.stderr)
        print(
            f"  RESOURCES_PER_SHARD         "
            f"Target per shard (default: {RESOURCES_PER_SHARD})",
            file=sys.stderr
        )
        print(
            f"  MIN_RESOURCES_FOR_SHARDING  "
            f"Min before sharding (default: {MIN_RESOURCES_FOR_SHARDING})",
            file=sys.stderr
        )
        print(
            f"  RENDEZVOUS_SEED             "
            f"Seed for distribution (default: {RENDEZVOUS_SEED})",
            file=sys.stderr
        )
        sys.exit(1)

    configuration_block_type = sys.argv[1]
    services = sys.argv[2:]

    valid_configuration_block_types = [
        'resources',
        'datasources',
        'actions',
        'list-resources',
        'ephemerals'
    ]
    if configuration_block_type not in valid_configuration_block_types:
        print(
            f"Error: Invalid configuration_block_type "
            f"'{configuration_block_type}'",
            file=sys.stderr
        )
        print(
            f"Valid configuration block types: "
            f"{', '.join(valid_configuration_block_types)}",
            file=sys.stderr
        )
        sys.exit(1)

    # Why stderr? Keeps diagnostic output separate from JSON stdout
    print(
        f"[generate_test_matrix] Configuration block type: "
        f"{configuration_block_type}",
        file=sys.stderr
    )
    print(
        f"[generate_test_matrix] Services: {', '.join(services)}",
        file=sys.stderr
    )
    print(
        f"[generate_test_matrix] Resources per shard: "
        f"{RESOURCES_PER_SHARD}",
        file=sys.stderr
    )
    print(
        f"[generate_test_matrix] Min resources for sharding: "
        f"{MIN_RESOURCES_FOR_SHARDING}",
        file=sys.stderr
    )
    print(
        f"[generate_test_matrix] Rendezvous seed: {RENDEZVOUS_SEED}",
        file=sys.stderr
    )

    matrix = generate_matrix(configuration_block_type, services)

    print(
        f"\n[generate_test_matrix] Generated {len(matrix['include'])} "
        f"matrix entries:",
        file=sys.stderr
    )
    for entry in matrix['include']:
        if entry['total_shards'] > 1:
            shard_info = (
                f"shard {entry['shard_index'] + 1}/{entry['total_shards']}"
            )
        else:
            shard_info = "no sharding"

        print(
            f"  - {entry['service']}: {entry['resource_count']} "
            f"resources ({shard_info})",
            file=sys.stderr
        )

    # Why stdout? GitHub Actions reads matrix from stdout
    print(json.dumps(matrix))


if __name__ == "__main__":
    main()
