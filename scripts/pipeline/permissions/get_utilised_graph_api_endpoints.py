#!/usr/bin/env python3
"""
get_utilised_graph_api_endpoints.py

Scans the Terraform provider Go source tree under internal/services and extracts
every Microsoft Graph API endpoint (URI + HTTP method + API version) referenced in
each service resource, grouped by service domain.

The script covers all five service categories:
  - actions
  - datasources
  - ephemerals
  - list-resources
  - resources

Two extraction strategies are applied to each Go file:
  1. SDK method-chain calls  — ``r.client.PathSeg1().PathSeg2().ByXxxId(id).Post(ctx, ...)``
  2. Custom request helpers  — explicit ``Endpoint:`` fields in request config structs
     (PostRequest, GetRequest, PatchRequest, PutRequest, DeleteRequest, GetRequestByResourceId)

Hard-delete verification helpers (``ExecuteDeleteWithVerification``) also imply
a ``DELETE /directory/deletedItems/{id}`` call, which is added automatically.

Output:
  - A JSON document written to ``--output`` (or stdout with ``--json-stdout``)
    suitable for piping to ``get_graph_endpoint_permissions.py``
  - A human-readable summary table printed to stdout

Usage:
    # Scan from current directory (must be the repository root)
    python3 get_utilised_graph_api_endpoints.py

    # Specify repo root explicitly
    python3 get_utilised_graph_api_endpoints.py --root /path/to/repo

    # Write JSON to a file
    python3 get_utilised_graph_api_endpoints.py --output endpoints.json

    # Filter to one service domain
    python3 get_utilised_graph_api_endpoints.py --domain device_management

    # Filter to one category
    python3 get_utilised_graph_api_endpoints.py --category resources

    # Print JSON to stdout instead of the summary table
    python3 get_utilised_graph_api_endpoints.py --json-stdout

Author:  Deployment Theory
Version: 1.0
"""

from __future__ import annotations

import argparse
import datetime
import json
import re
import sys
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

# =============================================================================
# CONSTANTS
# =============================================================================

SERVICES_SUBPATH: str = "internal/services"

CATEGORIES: frozenset[str] = frozenset({
    "actions",
    "datasources",
    "ephemerals",
    "list-resources",
    "resources",
})

API_VERSION_DIR_MAP: dict[str, str] = {
    "graph_beta": "beta",
    "graph_v1.0": "v1.0",
}

# Go SDK method names that represent HTTP verbs (mixed case as in Go source)
_HTTP_VERB_SET: frozenset[str] = frozenset({"Get", "Post", "Patch", "Put", "Delete"})

# Go client infrastructure methods that are NOT URL path segments
_INFRA_METHODS: frozenset[str] = frozenset({
    "GetAdapter",
    "GetRequestAdapter",
    "GetBackingStore",
    "GetAdditionalData",
    "GetChangeTracker",
})

# Custom request helper function name → HTTP method
_CUSTOM_FUNC_METHOD: dict[str, str] = {
    "GetRequest": "GET",
    "GetRequestByResourceId": "GET",
    "GetRequestWithPagination": "GET",
    "PostRequest": "POST",
    "PatchRequest": "PATCH",
    "PutRequest": "PUT",
    "DeleteRequest": "DELETE",
}

# Compiled regex: client chain leading up to an HTTP verb call
# Matches: {receiver}.client.{chain}.{Verb}(ctx
# Receiver can be r, d, a, or any single/multi-char lowercase identifier
_RE_CLIENT_CHAIN = re.compile(
    r'(?:[a-z]\w*\.)?client\.'
    r'(.+?)'
    r'\.(Get|Post|Patch|Put|Delete)\s*\(ctx'
)

# Compiled regex: custom request function call site
_RE_CUSTOM_FUNC = re.compile(
    r'\b(Get|Post|Patch|Put|Delete)Request'
    r'(?:ByResourceId|WithPagination|WithFilter)?\s*\('
)

# Compiled regex: Endpoint field in a request config struct
_RE_ENDPOINT_FIELD = re.compile(r'\bEndpoint:\s*"([^"]+)"')

# Compiled regex: EndpointSuffix field
_RE_ENDPOINT_SUFFIX = re.compile(r'\bEndpointSuffix:\s*"(/[^"]+)"')

# Compiled regex: ResourceIDPattern field (any value means an ID segment exists)
_RE_RESOURCE_ID = re.compile(r'\bResourceID:\s*\S')

# Compiled regex: API version constant or literal in config struct
_RE_API_VERSION = re.compile(
    r'GraphAPI(?P<const>Beta|V1)\b'
    r'|APIVersion:\s*"(?P<literal>beta|v1\.0)"'
)

# Compiled regex: ResourcePath struct field (used in list-resources)
_RE_RESOURCE_PATH = re.compile(r'\bResourcePath:\s*"(/[^"]+)"')

# Compiled regex: hard/soft delete verification helpers
_RE_HARD_DELETE = re.compile(r'\bExecuteDeleteWithVerification\b|\bExecuteHardDelete\b')

# =============================================================================
# DATA CLASSES
# =============================================================================


@dataclass
class GraphEndpoint:
    """A single Graph API endpoint reference discovered in Go source code."""

    uri: str
    method: str
    api_version: str
    source_file: str


@dataclass
class ResourceEndpoints:
    """All Graph API endpoints discovered for one Terraform resource."""

    resource_name: str
    category: str
    api_version: str
    source_dir: str
    endpoints: list[GraphEndpoint] = field(default_factory=list)


@dataclass
class ServiceDomainResult:
    """Aggregated endpoint data for one service domain (e.g. ``device_management``)."""

    domain: str
    resources: list[ResourceEndpoints] = field(default_factory=list)


@dataclass
class ScanResult:
    """Top-level output from a full scan of ``internal/services``."""

    scan_root: str
    generated_at: str
    statistics: dict[str, int]
    service_domains: list[ServiceDomainResult]


# =============================================================================
# GO SOURCE PARSING — CHAIN ANALYSIS
# =============================================================================


def _join_chain_lines(source: str) -> str:
    """
    Join multi-line Go method-chain continuations into single logical lines.

    In Go a method chain may span multiple lines by ending each line with
    a ``.`` token so the parser knows the expression continues. This function
    strips inline ``//`` comments and concatenates such continuation lines.

    Args:
        source: Raw Go source text.

    Returns:
        Source text with chain continuation lines merged into single lines.
    """
    lines = source.splitlines()
    joined: list[str] = []
    idx = 0
    while idx < len(lines):
        line = re.sub(r'//.*$', '', lines[idx]).rstrip()
        while line.endswith('.') and (idx + 1) < len(lines):
            idx += 1
            next_line = re.sub(r'//.*$', '', lines[idx]).rstrip()
            line = line + next_line.lstrip()
        joined.append(line)
        idx += 1
    return '\n'.join(joined)


def _extract_chain_method_names(chain: str) -> list[str]:
    """
    Extract top-level method names from a Go SDK method-chain string.

    Uses a parenthesis-depth tracker so that nested calls inside argument
    lists (e.g. ``ByDeviceCategoryId(obj.ID.ValueString())``) are not
    mistakenly collected as chain path segments.

    Args:
        chain: The chain portion between ``.client.`` and ``.HttpMethod(ctx``.

    Returns:
        Ordered list of top-level method names in the chain.
    """
    methods: list[str] = []
    current: list[str] = []
    depth = 0
    # Start True so the very first method name (no leading '.') is captured
    after_dot = True

    for char in chain:
        if depth == 0:
            if char == '.':
                current = []
                after_dot = True
            elif after_dot and (char.isalnum() or char == '_'):
                current.append(char)
            elif after_dot and char == '(':
                name = ''.join(current)
                if name:
                    methods.append(name)
                current = []
                after_dot = False
                depth = 1
            elif char == '(':
                depth = 1
                after_dot = False
                current = []
            else:
                after_dot = False
                current = []
        else:
            if char == '(':
                depth += 1
            elif char == ')':
                depth -= 1

    return methods


def _chain_methods_to_url_segments(method_names: list[str]) -> list[str]:
    """
    Convert Go SDK chain method names to URL path segments.

    Conversion rules:

    * ``ByXxxId(...)`` / ``ByXxx(...)``  → ``{id}``
    * ``Get``, ``Post``, ``Patch``, ``Put``, ``Delete`` → skipped (HTTP verbs)
    * Infrastructure methods (e.g. ``GetAdapter``) → skipped
    * Any other ``TitleCase`` method → first character lowercased (camelCase segment)

    Args:
        method_names: Top-level method names extracted from the SDK chain.

    Returns:
        List of URL path segments (without leading ``/``).
    """
    segments: list[str] = []
    for name in method_names:
        if name in _HTTP_VERB_SET or name in _INFRA_METHODS:
            continue
        if re.match(r'^By[A-Z]', name):
            segments.append('{id}')
        elif name and name[0].isupper():
            segments.append(name[0].lower() + name[1:])
    return segments


# =============================================================================
# GO SOURCE PARSING — ENDPOINT EXTRACTION
# =============================================================================


def _extract_sdk_endpoints(
    joined_source: str,
    api_version: str,
    filepath: str,
) -> list[GraphEndpoint]:
    """
    Extract Graph API endpoints from Microsoft Graph SDK method-chain calls.

    Detects patterns of the form::

        r.client.PathSeg1().PathSeg2().ByXxxId(id).Post(ctx, body, nil)

    Args:
        joined_source: Go source text with continuation lines merged.
        api_version:   API version inferred from the directory path.
        filepath:      Source file path (stored for attribution).

    Returns:
        List of :class:`GraphEndpoint` objects discovered via SDK chains.
    """
    endpoints: list[GraphEndpoint] = []
    for match in _RE_CLIENT_CHAIN.finditer(joined_source):
        chain_str = match.group(1)
        http_method = match.group(2).upper()
        names = _extract_chain_method_names(chain_str)
        segments = _chain_methods_to_url_segments(names)
        if not segments:
            continue
        endpoints.append(GraphEndpoint(
            uri='/' + '/'.join(segments),
            method=http_method,
            api_version=api_version,
            source_file=filepath,
        ))
    return endpoints


def _parse_api_version_from_window(window: str, fallback: str) -> str:
    """
    Determine the Graph API version from a code window around a request call.

    Checks for ``GraphAPIBeta`` / ``GraphAPIV1`` constants or a literal
    ``APIVersion: "beta"`` / ``APIVersion: "v1.0"`` field.

    Args:
        window:   Multi-line string containing the request config context.
        fallback: Version to return when no explicit marker is found.

    Returns:
        ``"beta"`` or ``"v1.0"``.
    """
    version_match = _RE_API_VERSION.search(window)
    if not version_match:
        return fallback
    const = version_match.group('const')
    if const:
        return 'beta' if const == 'Beta' else 'v1.0'
    return version_match.group('literal') or fallback


def _build_custom_request_uri(base: str, window: str) -> str:
    """
    Construct a full URI from a custom request config's fields.

    Combines the base ``Endpoint`` path with an optional ``/{id}`` segment
    (when ``ResourceID:`` is present) and an optional ``EndpointSuffix``.

    Args:
        base:   The raw value of the ``Endpoint:`` field (may lack leading ``/``).
        window: Code window containing the surrounding config struct fields.

    Returns:
        Normalised URI string starting with ``/``.
    """
    uri = base if base.startswith('/') else '/' + base
    if _RE_RESOURCE_ID.search(window):
        uri = uri + '/{id}'
    suffix_match = _RE_ENDPOINT_SUFFIX.search(window)
    if suffix_match:
        uri = uri + suffix_match.group(1)
    return uri


def _extract_custom_request_endpoints(
    source: str,
    api_version: str,
    filepath: str,
) -> list[GraphEndpoint]:
    """
    Extract Graph API endpoints from explicit custom request helper calls.

    Detects patterns where ``PostRequest``, ``GetRequest``, etc. are called
    with a config struct that includes an ``Endpoint: "/path/..."`` field.

    Args:
        source:      Raw Go source text (not joined).
        api_version: Fallback API version from the directory path.
        filepath:    Source file path (stored for attribution).

    Returns:
        List of :class:`GraphEndpoint` objects discovered via custom requests.
    """
    endpoints: list[GraphEndpoint] = []
    lines = source.splitlines()

    for line_idx, line in enumerate(lines):
        func_match = _RE_CUSTOM_FUNC.search(line)
        if not func_match:
            continue

        http_method = func_match.group(1).upper()
        window_end = min(len(lines), line_idx + 25)
        window = '\n'.join(lines[line_idx:window_end])

        ep_match = _RE_ENDPOINT_FIELD.search(window)
        if not ep_match:
            continue

        uri = _build_custom_request_uri(ep_match.group(1), window)
        ep_version = _parse_api_version_from_window(window, api_version)
        endpoints.append(GraphEndpoint(
            uri=uri,
            method=http_method,
            api_version=ep_version,
            source_file=filepath,
        ))

    return endpoints


def _extract_resource_path_endpoints(
    source: str,
    api_version: str,
    filepath: str,
) -> list[GraphEndpoint]:
    """
    Extract endpoints declared via the ``ResourcePath`` struct field pattern.

    Some list-resources store their base collection path in a ``ResourcePath``
    field (e.g. ``ResourcePath: "/deviceManagement/deviceManagementScripts"``).
    A ``GET`` is inferred for the collection.

    Args:
        source:      Raw Go source text.
        api_version: API version for the inferred endpoint.
        filepath:    Source file path (stored for attribution).

    Returns:
        List of inferred :class:`GraphEndpoint` objects, or empty list.
    """
    endpoints: list[GraphEndpoint] = []
    for match in _RE_RESOURCE_PATH.finditer(source):
        endpoints.append(GraphEndpoint(
            uri=match.group(1),
            method='GET',
            api_version=api_version,
            source_file=filepath,
        ))
    return endpoints


def _extract_hard_delete_endpoints(
    source: str,
    api_version: str,
    filepath: str,
) -> list[GraphEndpoint]:
    """
    Infer hard-delete endpoints for resources that use delete-verification helpers.

    When ``ExecuteDeleteWithVerification`` or ``ExecuteHardDelete`` is present,
    the resource flow always calls ``DELETE /directory/deletedItems/{id}`` and
    ``GET /directory/deletedItems/{id}`` as part of the permanent-delete check.

    Args:
        source:      Raw Go source text.
        api_version: API version for the inferred endpoints.
        filepath:    Source file path (stored for attribution).

    Returns:
        Two inferred hard-delete :class:`GraphEndpoint` objects, or empty list.
    """
    if not _RE_HARD_DELETE.search(source):
        return []
    return [
        GraphEndpoint(
            uri='/directory/deletedItems/{id}',
            method='DELETE',
            api_version=api_version,
            source_file=filepath,
        ),
        GraphEndpoint(
            uri='/directory/deletedItems/{id}',
            method='GET',
            api_version=api_version,
            source_file=filepath,
        ),
    ]


def extract_endpoints_from_file(
    go_file: Path,
    api_version: str,
) -> list[GraphEndpoint]:
    """
    Extract all Graph API endpoints referenced in a single Go source file.

    Combines SDK chain, custom request, resource-path, and hard-delete
    detection strategies. Test files (``*_test.go``) are skipped.

    Args:
        go_file:     Path to the ``.go`` source file.
        api_version: API version inferred from the containing directory.

    Returns:
        Deduplicated list of :class:`GraphEndpoint` objects found in the file.
    """
    if go_file.name.endswith('_test.go'):
        return []

    try:
        source = go_file.read_text(encoding='utf-8')
    except OSError:
        return []

    joined = _join_chain_lines(source)
    raw: list[GraphEndpoint] = []
    fpath = str(go_file)

    raw.extend(_extract_sdk_endpoints(joined, api_version, fpath))
    raw.extend(_extract_custom_request_endpoints(source, api_version, fpath))
    raw.extend(_extract_resource_path_endpoints(source, api_version, fpath))
    raw.extend(_extract_hard_delete_endpoints(source, api_version, fpath))

    seen: set[tuple[str, str, str]] = set()
    deduped: list[GraphEndpoint] = []
    for ep in raw:
        key = (ep.uri, ep.method, ep.api_version)
        if key not in seen:
            seen.add(key)
            deduped.append(ep)

    return deduped


# =============================================================================
# DIRECTORY SCANNING
# =============================================================================


def _api_version_from_dir_name(dir_name: str) -> str | None:
    """
    Map a directory name to a Graph API version string.

    Args:
        dir_name: A single directory component (e.g. ``"graph_beta"``).

    Returns:
        ``"beta"``, ``"v1.0"``, or ``None`` if not an API version directory.
    """
    return API_VERSION_DIR_MAP.get(dir_name)


def _scan_resource_dir(
    resource_dir: Path,
    resource_name: str,
    category: str,
    api_version: str,
    repo_root: Path,
) -> ResourceEndpoints:
    """
    Scan a single resource directory and collect all its Graph API endpoints.

    Args:
        resource_dir:  Path to the resource directory.
        resource_name: Leaf directory name used as the resource identifier.
        category:      Service category (e.g. ``"resources"``, ``"datasources"``).
        api_version:   API version for this resource.
        repo_root:     Repository root for computing relative ``source_dir``.

    Returns:
        :class:`ResourceEndpoints` populated with discovered endpoints.
    """
    try:
        source_dir = str(resource_dir.relative_to(repo_root))
    except ValueError:
        source_dir = str(resource_dir)

    entry = ResourceEndpoints(
        resource_name=resource_name,
        category=category,
        api_version=api_version,
        source_dir=source_dir,
    )

    seen: set[tuple[str, str, str]] = set()
    all_eps: list[GraphEndpoint] = []

    for go_file in sorted(resource_dir.rglob('*.go')):
        for ep in extract_endpoints_from_file(go_file, api_version):
            key = (ep.uri, ep.method, ep.api_version)
            if key not in seen:
                seen.add(key)
                all_eps.append(ep)

    entry.endpoints = sorted(all_eps, key=lambda e: (e.method, e.uri))
    return entry


def _iter_resources_under_domain(
    domain_dir: Path,
    category: str,
    repo_root: Path,
) -> list[ResourceEndpoints]:
    """
    Enumerate resource directories under a domain directory.

    Handles both structures:
    - ``{domain}/{api_version_dir}/{resource}/``  (most common)
    - ``{domain}/{resource}/``                    (no explicit api version dir)

    Args:
        domain_dir: Path to the domain directory (e.g. ``…/device_management``).
        category:   Parent service category name.
        repo_root:  Repository root for relative path computation.

    Returns:
        List of :class:`ResourceEndpoints` with at least one endpoint.
    """
    results: list[ResourceEndpoints] = []
    subdirs = sorted(d for d in domain_dir.iterdir() if d.is_dir())

    # Determine if subdirs are API version directories
    api_version_dirs = [d for d in subdirs if _api_version_from_dir_name(d.name)]
    non_api_dirs = [d for d in subdirs if not _api_version_from_dir_name(d.name)]

    if api_version_dirs:
        for api_dir in api_version_dirs:
            version = API_VERSION_DIR_MAP[api_dir.name]
            for resource_dir in sorted(api_dir.iterdir()):
                if not resource_dir.is_dir():
                    continue
                entry = _scan_resource_dir(
                    resource_dir, resource_dir.name, category, version, repo_root
                )
                if entry.endpoints:
                    results.append(entry)
    else:
        # No api_version subdirectory — treat subdirs directly as resources
        for resource_dir in non_api_dirs:
            entry = _scan_resource_dir(
                resource_dir, resource_dir.name, category, 'beta', repo_root
            )
            if entry.endpoints:
                results.append(entry)

    return results


def scan_services(
    services_root: Path,
    category_filter: str | None = None,
    domain_filter: str | None = None,
) -> list[ServiceDomainResult]:
    """
    Scan the entire ``internal/services`` tree and group results by domain.

    Args:
        services_root:   Path to the ``internal/services`` directory.
        category_filter: When set, only this category is scanned.
        domain_filter:   When set, only this domain is scanned.

    Returns:
        Sorted list of :class:`ServiceDomainResult` objects, one per domain.
    """
    repo_root = services_root.parent.parent
    domain_map: dict[str, ServiceDomainResult] = {}

    for category_dir in sorted(services_root.iterdir()):
        if not category_dir.is_dir():
            continue
        category_name = category_dir.name
        if category_name not in CATEGORIES:
            continue
        if category_filter and category_name != category_filter:
            continue

        for domain_dir in sorted(category_dir.iterdir()):
            if not domain_dir.is_dir():
                continue
            domain_name = domain_dir.name
            if domain_filter and domain_name != domain_filter:
                continue

            resources = _iter_resources_under_domain(domain_dir, category_name, repo_root)
            if not resources:
                continue

            if domain_name not in domain_map:
                domain_map[domain_name] = ServiceDomainResult(domain=domain_name)
            domain_map[domain_name].resources.extend(resources)

    return sorted(domain_map.values(), key=lambda d: d.domain)


# =============================================================================
# OUTPUT HELPERS
# =============================================================================


def _compute_statistics(domains: list[ServiceDomainResult]) -> dict[str, int]:
    """
    Compute summary statistics from the scan results.

    Args:
        domains: List of :class:`ServiceDomainResult` objects.

    Returns:
        Dict with counts for domains, resources, total and unique endpoints.
    """
    total_resources = sum(len(d.resources) for d in domains)
    all_eps = [
        ep
        for d in domains
        for r in d.resources
        for ep in r.endpoints
    ]
    unique_eps = len({(e.uri, e.method, e.api_version) for e in all_eps})
    return {
        'total_domains': len(domains),
        'total_resources': total_resources,
        'total_endpoint_references': len(all_eps),
        'total_unique_endpoints': unique_eps,
    }


def _build_json_output(result: ScanResult) -> dict[str, Any]:
    """
    Convert a :class:`ScanResult` to a plain JSON-serialisable dict.

    Args:
        result: The completed scan result.

    Returns:
        JSON-serialisable dictionary representation.
    """
    domains_out: list[dict[str, Any]] = []
    for domain in result.service_domains:
        resources_out: list[dict[str, Any]] = []
        for res in domain.resources:
            resources_out.append({
                'resource_name': res.resource_name,
                'category': res.category,
                'api_version': res.api_version,
                'source_dir': res.source_dir,
                'endpoints': [
                    {
                        'uri': ep.uri,
                        'method': ep.method,
                        'api_version': ep.api_version,
                    }
                    for ep in res.endpoints
                ],
            })
        domains_out.append({
            'domain': domain.domain,
            'resources': resources_out,
        })
    return {
        'scan_root': result.scan_root,
        'generated_at': result.generated_at,
        'statistics': result.statistics,
        'service_domains': domains_out,
    }


def _print_summary(result: ScanResult) -> None:
    """
    Print a human-readable summary table of discovered endpoints to stdout.

    Args:
        result: The completed scan result.
    """
    stats = result.statistics
    print()
    print('╔══════════════════════════════════════════════════════════════════╗')
    print('║       GRAPH API ENDPOINT SCAN — RESULTS SUMMARY                 ║')
    print('╚══════════════════════════════════════════════════════════════════╝')
    print(f"  Scan root:          {result.scan_root}")
    print(f"  Generated:          {result.generated_at}")
    print(f"  Service domains:    {stats['total_domains']}")
    print(f"  Resources scanned:  {stats['total_resources']}")
    print(f"  Unique endpoints:   {stats['total_unique_endpoints']}")
    print()

    for domain in result.service_domains:
        print(f"  ┌─ {domain.domain}")
        for res in domain.resources:
            api_v = res.api_version
            print(f"  │  ├─ [{res.category}] {res.resource_name}  ({api_v})")
            for ep in res.endpoints:
                method_padded = ep.method.ljust(7)
                print(f"  │  │      {method_padded} {ep.uri}")
        print('  │')

    print()
    print(f"  Total unique endpoints: {stats['total_unique_endpoints']}")
    print()


def _print_permissions_script_calls(
    result: ScanResult,
    perms_script: Path,
) -> None:
    """
    Print the shell commands to run ``get_graph_endpoint_permissions.py``
    for every unique endpoint discovered.

    Args:
        result:       The completed scan result.
        perms_script: Path to ``get_graph_endpoint_permissions.py``.
    """
    seen: set[tuple[str, str, str]] = set()
    print()
    print('# ── get_graph_endpoint_permissions.py calls ──────────────────────')
    for domain in result.service_domains:
        for res in domain.resources:
            for ep in res.endpoints:
                key = (ep.uri, ep.method, ep.api_version)
                if key in seen:
                    continue
                seen.add(key)
                print(
                    f'python3 {perms_script} '
                    f'--uri "{ep.uri}" '
                    f'--method {ep.method} '
                    f'--api-version {ep.api_version}'
                )
    print()


# =============================================================================
# CLI
# =============================================================================


def _find_repo_root(start: Path) -> Path:
    """
    Walk up from ``start`` until a ``go.mod`` file is found.

    Falls back to the current working directory if no ``go.mod`` is found.

    Args:
        start: Directory to begin the upward search from.

    Returns:
        Path to the repository root.
    """
    current = start.resolve()
    for _ in range(10):
        if (current / 'go.mod').exists():
            return current
        parent = current.parent
        if parent == current:
            break
        current = parent
    return Path.cwd()


def _parse_args() -> argparse.Namespace:
    """Parse and validate command-line arguments."""
    parser = argparse.ArgumentParser(
        description='Scan Go sources and extract Microsoft Graph API endpoints.',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=(
            'Examples:\n'
            '  %(prog)s\n'
            '  %(prog)s --root /path/to/repo --output endpoints.json\n'
            '  %(prog)s --domain device_management --category resources\n'
            '  %(prog)s --json-stdout | jq .statistics\n'
        ),
    )
    parser.add_argument(
        '--root',
        metavar='PATH',
        help='Repository root directory (default: auto-detected via go.mod).',
    )
    parser.add_argument(
        '--output',
        metavar='FILE',
        help='Write JSON output to this file path.',
    )
    parser.add_argument(
        '--domain',
        metavar='DOMAIN',
        help='Limit scan to this service domain (e.g. device_management).',
    )
    parser.add_argument(
        '--category',
        metavar='CATEGORY',
        choices=sorted(CATEGORIES),
        help='Limit scan to this service category.',
    )
    parser.add_argument(
        '--json-stdout',
        action='store_true',
        dest='json_stdout',
        help='Write JSON to stdout instead of the summary table.',
    )
    parser.add_argument(
        '--show-commands',
        action='store_true',
        dest='show_commands',
        help='Print get_graph_endpoint_permissions.py shell commands to stdout.',
    )
    return parser.parse_args()


def main() -> None:
    """Entry point: parse arguments, scan the codebase, and output results."""
    args = _parse_args()

    # Resolve repository root
    if args.root:
        repo_root = Path(args.root).resolve()
    else:
        repo_root = _find_repo_root(Path(__file__).parent)

    services_root = repo_root / SERVICES_SUBPATH
    if not services_root.is_dir():
        print(
            f'error: services directory not found: {services_root}',
            file=sys.stderr,
        )
        sys.exit(1)

    # Run the scan
    domains = scan_services(
        services_root,
        category_filter=args.category,
        domain_filter=args.domain,
    )

    stats = _compute_statistics(domains)
    result = ScanResult(
        scan_root=str(repo_root),
        generated_at=datetime.datetime.now(datetime.timezone.utc).isoformat(),
        statistics=stats,
        service_domains=domains,
    )

    # Output
    json_data = _build_json_output(result)

    if args.json_stdout:
        print(json.dumps(json_data, indent=2))
    else:
        _print_summary(result)

    if args.output:
        output_path = Path(args.output)
        output_path.write_text(json.dumps(json_data, indent=2), encoding='utf-8')
        print(f'  JSON written to: {output_path}')

    if args.show_commands:
        perms_script = Path(__file__).parent / 'get_graph_endpoint_permissions.py'
        _print_permissions_script_calls(result, perms_script)


if __name__ == '__main__':
    main()
