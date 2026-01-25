#!/usr/bin/env python3
"""Generate dynamic GitHub Actions matrix with automatic shard calculation.

This script discovers resources/datasources/actions/list-resources/ephemerals in each
service and automatically determines the optimal number of shards using Rendezvous
(Highest Random Weight) distribution for stable, deterministic assignment.

Usage:
    ./generate-test-matrix.py <category> <service1> [service2 ...]

Categories:
    resources, datasources, actions, list-resources, ephemerals

Examples:
    ./generate-test-matrix.py resources device_management groups users
    ./generate-test-matrix.py datasources device_management utility
    ./generate-test-matrix.py list-resources device_management
    ./generate-test-matrix.py ephemerals multitenant_management windows_autopilot_device_csv_import

Output:
    JSON matrix suitable for GitHub Actions matrix strategy
"""

import os
import sys
import json
import hashlib
import struct
from pathlib import Path
from typing import List, Dict

# Configuration
RESOURCES_PER_SHARD = int(os.environ.get('RESOURCES_PER_SHARD', '5'))
MIN_RESOURCES_FOR_SHARDING = int(os.environ.get('MIN_RESOURCES_FOR_SHARDING', '10'))
RENDEZVOUS_SEED = os.environ.get('RENDEZVOUS_SEED', 'terraform-provider-microsoft365')


def discover_resources(service: str, category: str) -> List[str]:
    """Discover all resource/datasource/action/list-resource/ephemeral packages in a service.

    Args:
        service: Service name (e.g., 'device_management')
        category: Category type ('resources', 'datasources', 'actions', 'list-resources', 'ephemerals')

    Returns:
        List of package paths relative to the service directory
    """
    base_path = Path(f"./internal/services/{category}/{service}")

    if not base_path.exists():
        return []

    packages = []

    # Determine marker file based on category
    marker_files = {
        'resources': ['resource.go'],
        'datasources': ['datasource.go'],
        'actions': ['action.go'],
        'list-resources': ['list_resource.go'],
        'ephemerals': ['ephemeral_resource.go']
    }

    markers = marker_files.get(category, ['resource.go', 'datasource.go', 'action.go'])

    # Discover across graph API versions
    for graph_version_dir in base_path.iterdir():
        if not graph_version_dir.is_dir():
            continue

        # Handle both graph_* versioned and non-versioned structures
        # For ephemerals: some are at service/graph_beta/name, others at service/name
        is_graph_version = graph_version_dir.name.startswith('graph_')

        if is_graph_version:
            # Find all packages within graph version directory
            for resource_dir in graph_version_dir.iterdir():
                if not resource_dir.is_dir():
                    continue

                # Check for marker files
                has_marker = any((resource_dir / marker).exists() for marker in markers)

                if has_marker:
                    # Store as relative path from base_path
                    rel_path = str(resource_dir.relative_to(base_path))
                    packages.append(rel_path)
        else:
            # Non-versioned structure (e.g., ephemerals/windows_autopilot_device_csv_import)
            # Check if this directory itself has a marker file
            has_marker = any((graph_version_dir / marker).exists() for marker in markers)

            if has_marker:
                # Store as relative path from base_path
                rel_path = str(graph_version_dir.relative_to(base_path))
                packages.append(rel_path)

    return sorted(packages)


def rendezvous_hash(item: str, shard_index: int, seed: str) -> int:
    """Compute Rendezvous (HRW) hash for item-shard pair.

    Uses SHA-256 to generate deterministic weight. Higher weight = higher priority.

    Args:
        item: Item to hash (resource package path)
        shard_index: Shard index (0-based)
        seed: Seed string for deterministic hashing

    Returns:
        64-bit unsigned integer weight
    """
    # Format matches Go implementation: "item:shard_N:seed"
    input_str = f"{item}:shard_{shard_index}:{seed}"
    hash_bytes = hashlib.sha256(input_str.encode('utf-8')).digest()

    # Extract first 8 bytes as uint64 (big-endian, matching Go)
    weight = struct.unpack('>Q', hash_bytes[:8])[0]
    return weight


def distribute_rendezvous(items: List[str], shard_count: int, seed: str) -> Dict[int, List[str]]:
    """Distribute items across shards using Rendezvous (HRW) algorithm.

    Each item is assigned to the shard with the highest hash weight.
    This ensures minimal reassignment when shard count changes.

    Args:
        items: List of items to distribute (resource paths)
        shard_count: Number of shards to distribute across
        seed: Seed for deterministic distribution

    Returns:
        Dictionary mapping shard_index -> list of items
    """
    if shard_count <= 0:
        shard_count = 1

    # Initialize empty shards
    shards = {i: [] for i in range(shard_count)}

    # Assign each item to shard with highest weight
    for item in items:
        highest_weight = 0
        selected_shard = 0

        # Compute weight for this item against all shards
        for shard_idx in range(shard_count):
            weight = rendezvous_hash(item, shard_idx, seed)
            if weight > highest_weight:
                highest_weight = weight
                selected_shard = shard_idx

        shards[selected_shard].append(item)

    return shards


def calculate_shard_count(resource_count: int) -> int:
    """Calculate optimal shard count based on resource count.

    Args:
        resource_count: Number of resources in the service

    Returns:
        Optimal number of shards (minimum 1)
    """
    if resource_count < MIN_RESOURCES_FOR_SHARDING:
        return 1

    # Calculate shards needed: ceil(resource_count / RESOURCES_PER_SHARD)
    shard_count = (resource_count + RESOURCES_PER_SHARD - 1) // RESOURCES_PER_SHARD
    return max(1, shard_count)


def generate_matrix_for_service(service: str, category: str, runner: str = "ubuntu-24.04-arm") -> List[Dict]:
    """Generate matrix entries for a single service.

    Args:
        service: Service name
        category: Category type ('resources', 'datasources', 'actions')
        runner: GitHub Actions runner type

    Returns:
        List of matrix entry dictionaries
    """
    # Discover all resources in this service
    resources = discover_resources(service, category)
    resource_count = len(resources)

    if resource_count == 0:
        # No resources found - skip this service entirely
        return []

    # Calculate optimal shard count
    shard_count = calculate_shard_count(resource_count)

    # Distribute resources using Rendezvous
    shards = distribute_rendezvous(resources, shard_count, RENDEZVOUS_SEED)

    # Generate matrix entries
    matrix_entries = []
    for shard_idx in range(shard_count):
        shard_resources = shards[shard_idx]

        entry = {
            "service": service,
            "runner": runner,
            "shard_index": shard_idx,
            "total_shards": shard_count,
            "resource_count": len(shard_resources),
            "resources": ",".join(shard_resources)  # Pass to runner as comma-separated
        }

        # Add human-readable shard label
        if shard_count == 1:
            entry["shard_label"] = ""
        else:
            entry["shard_label"] = f" ({shard_idx + 1}/{shard_count})"

        matrix_entries.append(entry)

    return matrix_entries


def generate_matrix(category: str, services: List[str]) -> Dict:
    """Generate complete matrix for all services.

    Args:
        category: Category type ('resources', 'datasources', 'actions')
        services: List of service names to process

    Returns:
        Matrix dictionary suitable for fromJSON() in GitHub Actions
    """
    all_entries = []

    for service in services:
        entries = generate_matrix_for_service(service, category)
        all_entries.extend(entries)

    return {"include": all_entries}


def main():
    """Main entry point."""
    if len(sys.argv) < 3:
        print("Usage: generate-test-matrix.py <category> <service1> [service2 ...]", file=sys.stderr)
        print("\nCategories: resources, datasources, actions, list-resources, ephemerals", file=sys.stderr)
        print("\nEnvironment variables:", file=sys.stderr)
        print(f"  RESOURCES_PER_SHARD         Target resources per shard (default: {RESOURCES_PER_SHARD})", file=sys.stderr)
        print(f"  MIN_RESOURCES_FOR_SHARDING  Min resources before sharding (default: {MIN_RESOURCES_FOR_SHARDING})", file=sys.stderr)
        print(f"  RENDEZVOUS_SEED             Seed for distribution (default: {RENDEZVOUS_SEED})", file=sys.stderr)
        sys.exit(1)

    category = sys.argv[1]
    services = sys.argv[2:]

    valid_categories = ['resources', 'datasources', 'actions', 'list-resources', 'ephemerals']
    if category not in valid_categories:
        print(f"Error: Invalid category '{category}'", file=sys.stderr)
        print(f"Valid categories: {', '.join(valid_categories)}", file=sys.stderr)
        sys.exit(1)

    # Print diagnostic info to stderr (won't pollute JSON output)
    print(f"[generate-test-matrix] Category: {category}", file=sys.stderr)
    print(f"[generate-test-matrix] Services: {', '.join(services)}", file=sys.stderr)
    print(f"[generate-test-matrix] Resources per shard: {RESOURCES_PER_SHARD}", file=sys.stderr)
    print(f"[generate-test-matrix] Min resources for sharding: {MIN_RESOURCES_FOR_SHARDING}", file=sys.stderr)
    print(f"[generate-test-matrix] Rendezvous seed: {RENDEZVOUS_SEED}", file=sys.stderr)

    # Generate matrix
    matrix = generate_matrix(category, services)

    # Print summary to stderr
    print(f"\n[generate-test-matrix] Generated {len(matrix['include'])} matrix entries:", file=sys.stderr)
    for entry in matrix['include']:
        shard_info = f"shard {entry['shard_index'] + 1}/{entry['total_shards']}" if entry['total_shards'] > 1 else "no sharding"
        print(f"  - {entry['service']}: {entry['resource_count']} resources ({shard_info})", file=sys.stderr)

    # Output JSON to stdout for GitHub Actions
    print(json.dumps(matrix))


if __name__ == "__main__":
    main()
