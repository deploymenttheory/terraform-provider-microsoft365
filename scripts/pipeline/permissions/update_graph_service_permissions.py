#!/usr/bin/env python3
"""
update_graph_service_permissions.py

Derives the required Microsoft Graph API application permissions for every
endpoint used by each Terraform provider resource and updates the
``ReadPermissions`` / ``WritePermissions`` fields in the corresponding Go
source files.

The script covers all five service categories (actions, datasources,
ephemerals, list-resources, resources) for both ``graph_beta`` and
``graph_v1.0`` API versions.

Permission classification
-------------------------
ReadPermissions
    Permissions for GET operations. When both ``Foo.Read.All`` and
    ``Foo.ReadWrite.All`` would satisfy a GET, only ``Foo.Read.All``
    is kept (least-privilege for read).

WritePermissions
    Permissions for POST / PATCH / PUT / DELETE operations. Not written
    for datasources, list-resources, or ephemerals (read-only types).

Safety behaviour
----------------
* **Dry-run by default** — pass ``--apply`` to write changes.
* ReadPermissions are only updated when at least one permission is found.
* WritePermissions are only updated when write-operation permissions are
  found *and* the service category supports a WritePermissions field.

Usage::

    # Preview what would change (dry-run)
    python3 update_graph_service_permissions.py

    # Write changes to Go source files
    python3 update_graph_service_permissions.py --apply

    # Limit to one service domain or category
    python3 update_graph_service_permissions.py --domain applications --apply
    python3 update_graph_service_permissions.py --category resources --apply

    # Force re-download of Graph command metadata cache
    python3 update_graph_service_permissions.py --refresh-metadata

Author:  Deployment Theory
Version: 1.0
"""

from __future__ import annotations

import argparse
import asyncio
import re
import sys
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any

# ---------------------------------------------------------------------------
# Extend sys.path so the sibling scripts in the same directory are importable.
# ---------------------------------------------------------------------------
_SCRIPT_DIR = Path(__file__).resolve().parent
if str(_SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(_SCRIPT_DIR))

# pylint: disable=wrong-import-position
from get_graph_endpoint_permissions import (  # noqa: E402
    find_endpoint_permissions,
    get_graph_command_metadata,
    get_permission_name,
    is_application_permission,
)
from get_utilised_graph_api_endpoints import (  # noqa: E402
    ResourceEndpoints,
    SERVICES_SUBPATH,
    scan_services,
)
# pylint: enable=wrong-import-position

# =============================================================================
# CONSTANTS
# =============================================================================

# Categories whose Go structs include a WritePermissions field
_WRITE_CATEGORIES: frozenset[str] = frozenset({'resources', 'actions'})

# Per-domain allowlist of permission-name prefixes.
# Only permissions whose name starts with one of the listed prefixes are kept
# for that domain. This prevents cross-domain contamination where Graph URIs
# shared across service areas (e.g. /applications/{id}) return permissions
# from completely unrelated domains (e.g. AgentIdentity.* or DeviceManagement.*).
# Domains absent from this dict receive no filtering.
_DOMAIN_PERM_PREFIXES: dict[str, tuple[str, ...]] = {
    'agents': (
        'Agent', 'Application.', 'Directory.', 'User.',
    ),
    'applications': (
        'AppRoleAssignment.', 'Application.', 'Directory.', 'Group.', 'User.',
    ),
    'device_and_app_management': (
        'DeviceManagement', 'Directory.',
    ),
    'device_management': (
        'DeviceManagement', 'Directory.', 'Group.', 'GroupMember.', 'User.',
    ),
    'groups': (
        'Application.', 'Directory.', 'Group.', 'GroupMember.',
        'LicenseAssignment.', 'RoleManagement.', 'User.',
    ),
    'identity_and_access': (
        # CustomSecAttribute (no dot) matches both CustomSecAttributeAssignment.*
        # and CustomSecAttributeDefinition.* — the dot variant would miss them.
        'Application.', 'CustomSecAttribute', 'Directory.',
        'Policy.', 'RoleManagement.',
        # User.* needed for validate.go calls e.g. validateUserExists -> GET /users/{id}
        'User.',
    ),
    'multitenant_management': (
        'ManagedTenant',
    ),
    'users': (
        'Directory.', 'LicenseAssignment.', 'Mailbox', 'User.',
    ),
    'utility': (
        'Application.', 'Device.', 'Directory.', 'Group.', 'User.',
    ),
    'windows_365': (
        'CloudPC.', 'Directory.', 'RoleManagement.',
    ),
}

# Permission-name suffixes that indicate purely read-only access.
# These are excluded from WritePermissions — write operation blocks must only
# contain permissions that actually grant write / action capabilities.
_READ_ONLY_SUFFIXES: tuple[str, ...] = (
    '.Read.All', '.ReadBasic.All', '.Read.Basic.All',
)

# Regex: locate any New*() constructor function signature
_RE_CONSTRUCTOR = re.compile(
    r'\bfunc\s+New\w+\s*\([^)]*\)\s+\w+(?:\.\w+)?\s*\{',
    re.MULTILINE,
)

# Regex: ReadPermissions: []string{...}, — flexible leading whitespace
_RE_READ_BLOCK = re.compile(
    r'(?P<indent>[ \t]*)ReadPermissions:\s*\[\]string\{[^}]*\},?',
)

# Regex: WritePermissions: []string{...}, — flexible leading whitespace
_RE_WRITE_BLOCK = re.compile(
    r'(?P<indent>[ \t]*)WritePermissions:\s*\[\]string\{[^}]*\},?',
)

# =============================================================================
# DATA CLASSES
# =============================================================================


@dataclass
class _UpdateResult:
    """Outcome of processing one resource's Go constructor file."""

    domain: str
    resource: ResourceEndpoints
    go_file: Path | None
    read_perms: list[str] = field(default_factory=list)
    write_perms: list[str] = field(default_factory=list)
    status: str = 'pending'


# =============================================================================
# PERMISSION LOOKUP
# =============================================================================


def _get_app_perms(
    metadata: list[dict[str, Any]],
    uri: str,
    method: str,
    api_version: str,
) -> frozenset[str]:
    """
    Return application permissions for one Graph endpoint from the metadata.

    Args:
        metadata:    Full list of MgCommandMetadata entries.
        uri:         Graph URI template (e.g. ``/groups/{id}``).
        method:      HTTP method (``GET``, ``POST``, etc.).
        api_version: ``"beta"`` or ``"v1.0"``.

    Returns:
        Frozen set of application permission name strings.
    """
    matched = find_endpoint_permissions(metadata, uri, method, api_version)
    result: set[str] = set()
    for entry in matched:
        for perm_obj in entry.get('Permissions', []):
            name = get_permission_name(perm_obj)
            if name and is_application_permission(name):
                result.add(name)
    return frozenset(result)


def _is_write_capable(name: str) -> bool:
    """
    Return ``True`` if *name* is not a purely read-only permission.

    Permissions ending in ``.Read.All`` (or similar read-only suffixes) are
    excluded from WritePermissions blocks — they grant no write capability.

    Args:
        name: Permission name string to test.

    Returns:
        ``True`` when the permission may grant write or action access.
    """
    return not any(name.endswith(s) for s in _READ_ONLY_SUFFIXES)


def _filter_for_read(perms: frozenset[str]) -> frozenset[str]:
    """
    Prefer least-privilege read permissions over read-write equivalents.

    When both ``Foo.Read.All`` and ``Foo.ReadWrite.All`` are present,
    ``Foo.ReadWrite.All`` is excluded — ``Foo.Read.All`` is sufficient for
    GET operations.

    Args:
        perms: Raw permission set returned by GET-endpoint metadata lookups.

    Returns:
        Filtered set with redundant ``*.ReadWrite.*`` entries removed.
    """
    result: set[str] = set()
    for perm in perms:
        if '.ReadWrite.' in perm:
            read_only = perm.replace('.ReadWrite.', '.Read.')
            if read_only in perms:
                continue  # read-only equivalent exists; skip the ReadWrite one
        result.add(perm)
    return frozenset(result)


def _compute_permissions(
    resource: ResourceEndpoints,
    domain: str,
    metadata: list[dict[str, Any]],
) -> tuple[list[str], list[str]]:
    """
    Derive ReadPermissions and WritePermissions for a resource's endpoints.

    GET endpoints feed ReadPermissions (filtered for least privilege).
    POST / PATCH / PUT / DELETE endpoints feed WritePermissions (filtered to
    exclude read-only permission names).

    Permissions are also scoped to the service domain via
    :data:`_DOMAIN_PERM_PREFIXES` to prevent cross-domain contamination when
    Graph URIs are shared across multiple service areas.

    Args:
        resource: Resource scan result with endpoint list.
        domain:   Service domain name (e.g. ``"applications"``).
        metadata: Full Graph command metadata list.

    Returns:
        ``(sorted_read_perms, sorted_write_perms)`` string lists.
    """
    allowed = _DOMAIN_PERM_PREFIXES.get(domain)  # None → no prefix filtering
    read_raw: set[str] = set()
    write_raw: set[str] = set()

    for ep in resource.endpoints:
        perms = _get_app_perms(metadata, ep.uri, ep.method, ep.api_version)
        if allowed is not None:
            perms = frozenset(p for p in perms if any(p.startswith(pf) for pf in allowed))
        if ep.method.upper() == 'GET':
            read_raw.update(perms)
        else:
            write_raw.update(perms)

    read_perms = sorted(_filter_for_read(frozenset(read_raw)))
    write_perms = sorted(p for p in write_raw if _is_write_capable(p))
    return read_perms, write_perms


# =============================================================================
# GO SOURCE UPDATER
# =============================================================================


def _make_perms_block(field_name: str, perms: list[str], indent: str) -> str:
    """
    Format a Go ``[]string{...}`` permissions field block.

    Args:
        field_name: ``"ReadPermissions"`` or ``"WritePermissions"``.
        perms:      Sorted list of permission name strings.
        indent:     Leading whitespace for the field line (e.g. ``\\t\\t``).

    Returns:
        Formatted Go field assignment string ending with ``,``.
    """
    if not perms:
        return f'{indent}{field_name}: []string{{}},'
    item_indent = indent + '\t'
    inner = '\n'.join(f'{item_indent}"{p}",' for p in perms)
    return f'{indent}{field_name}: []string{{\n{inner}\n{indent}}},'


def _replace_first(
    source: str,
    pattern: re.Pattern[str],
    replacement: str,
) -> tuple[str, bool]:
    """
    Replace the first match of *pattern* in *source* with *replacement*.

    Args:
        source:      Input text.
        pattern:     Compiled regex to search for.
        replacement: Text to substitute in place of the match.

    Returns:
        ``(updated_source, was_changed)`` tuple.
    """
    match = pattern.search(source)
    if not match:
        return source, False
    updated = source[:match.start()] + replacement + source[match.end():]
    return updated, updated != source


def _apply_write_update(
    source: str,
    write_perms: list[str],
    category: str,
) -> tuple[str, bool]:
    """
    Update the WritePermissions block when applicable.

    Only modifies the source when *category* is in ``_WRITE_CATEGORIES`` and
    *write_perms* is non-empty.

    Args:
        source:      Current Go source text (may already have ReadPermissions
                     replaced).
        write_perms: Computed WritePermissions list.
        category:    Service category.

    Returns:
        ``(updated_source, was_changed)`` tuple.
    """
    if category not in _WRITE_CATEGORIES or not write_perms:
        return source, False
    write_match = _RE_WRITE_BLOCK.search(source)
    if not write_match:
        return source, False
    w_indent = write_match.group('indent')
    write_block = _make_perms_block('WritePermissions', write_perms, w_indent)
    return _replace_first(source, _RE_WRITE_BLOCK, write_block)


def _update_source(
    source: str,
    read_perms: list[str],
    write_perms: list[str],
    category: str,
) -> str | None:
    """
    Locate and update permissions fields in a Go constructor function.

    Returns the updated source string, or ``None`` when no changes are
    needed or the constructor / ReadPermissions field cannot be found.

    Args:
        source:      Full Go source text.
        read_perms:  New sorted ReadPermissions list.
        write_perms: New sorted WritePermissions list.
        category:    Service category (e.g. ``"resources"``).

    Returns:
        Updated source string, or ``None`` if unchanged / not applicable.
    """
    if not _RE_CONSTRUCTOR.search(source):
        return None
    if not read_perms:
        return None

    read_match = _RE_READ_BLOCK.search(source)
    if not read_match:
        return None

    indent = read_match.group('indent')
    read_block = _make_perms_block('ReadPermissions', read_perms, indent)
    updated, read_changed = _replace_first(source, _RE_READ_BLOCK, read_block)
    updated, write_changed = _apply_write_update(updated, write_perms, category)

    return updated if (read_changed or write_changed) else None


def _find_constructor_file(resource_dir: Path) -> Path | None:
    """
    Return the Go file in *resource_dir* that contains a ``New*()`` constructor.

    Test files (``*_test.go``) are skipped.

    Args:
        resource_dir: Directory to search.

    Returns:
        Path to the first matching Go file, or ``None``.
    """
    for go_file in sorted(resource_dir.glob('*.go')):
        if go_file.name.endswith('_test.go'):
            continue
        try:
            text = go_file.read_text(encoding='utf-8')
        except OSError:
            continue
        if _RE_CONSTRUCTOR.search(text):
            return go_file
    return None


def _update_go_file(
    result: _UpdateResult,
    apply_changes: bool,
) -> None:
    """
    Update (or preview) the permissions in the Go file referenced by *result*.

    Mutates ``result.status`` in-place to reflect the outcome.

    Args:
        result:        Update result object (go_file must be set).
        apply_changes: When ``True``, write the updated file to disk.
    """
    if result.go_file is None:
        result.status = 'no file'
        return

    try:
        original = result.go_file.read_text(encoding='utf-8')
    except OSError as exc:
        result.status = f'read error: {exc}'
        return

    updated = _update_source(
        original, result.read_perms, result.write_perms, result.resource.category
    )
    if updated is None:
        result.status = 'no change'
        return

    if apply_changes:
        try:
            result.go_file.write_text(updated, encoding='utf-8')
        except OSError as exc:
            result.status = f'write error: {exc}'
            return
        result.status = 'updated'
    else:
        result.status = 'would update'


# =============================================================================
# OUTPUT
# =============================================================================


def _fmt_perms(perms: list[str]) -> str:
    """Return a comma-separated permission list, or ``'(none)'`` if empty."""
    return ', '.join(perms) if perms else '(none)'


def _print_result(result: _UpdateResult) -> None:
    """
    Print a one-line summary row for a resource update result.

    Args:
        result: Completed update result.
    """
    file_label = result.go_file.name if result.go_file else '(no constructor file)'
    name = f'{result.domain}/{result.resource.resource_name}'
    label = f'{name}  ({result.resource.category})  {file_label}'
    print(f'  [{result.status:^14}] {label}')
    print(f'               R: {_fmt_perms(result.read_perms)}')
    if result.resource.category in _WRITE_CATEGORIES:
        print(f'               W: {_fmt_perms(result.write_perms)}')


# =============================================================================
# ORCHESTRATION
# =============================================================================


def _process_resource(
    resource: ResourceEndpoints,
    domain: str,
    repo_root: Path,
    metadata: list[dict[str, Any]],
    apply_changes: bool,
) -> _UpdateResult:
    """
    Compute permissions and update one resource's Go constructor file.

    Args:
        resource:      Resource scan result.
        domain:        Service domain name (used for permission filtering and display).
        repo_root:     Repository root path.
        metadata:      Graph command metadata list.
        apply_changes: Whether to write changes to disk.

    Returns:
        Populated :class:`_UpdateResult` describing the outcome.
    """
    read_perms, write_perms = _compute_permissions(resource, domain, metadata)
    resource_dir = repo_root / resource.source_dir
    go_file = _find_constructor_file(resource_dir)

    result = _UpdateResult(
        domain=domain,
        resource=resource,
        go_file=go_file,
        read_perms=read_perms,
        write_perms=write_perms,
    )
    _update_go_file(result, apply_changes)
    return result


async def _run(args: argparse.Namespace) -> int:
    """
    Load metadata, scan services, and orchestrate permission updates.

    Args:
        args: Parsed CLI arguments.

    Returns:
        ``0`` on success, ``1`` if any errors occurred.
    """
    if args.root:
        repo_root = Path(args.root).resolve()
    else:
        repo_root = _find_repo_root(Path(__file__).parent)

    services_root = repo_root / SERVICES_SUBPATH
    if not services_root.is_dir():
        print(f'error: services directory not found: {services_root}', file=sys.stderr)
        return 1

    print('Loading Graph command metadata …', flush=True)
    metadata = await get_graph_command_metadata(force_refresh=args.refresh_metadata)
    print(f'  {len(metadata):,} metadata entries loaded.')

    print('Scanning Go source tree …', flush=True)
    domains = scan_services(
        services_root,
        category_filter=args.category,
        domain_filter=args.domain,
    )
    total = sum(len(d.resources) for d in domains)
    print(f'  {len(domains)} domains, {total} resources found.')

    mode = 'APPLY' if args.apply else 'DRY-RUN'
    print(f'\n  Mode: {mode}\n')

    changed = errors = skipped = 0

    for domain_result in domains:
        for resource in domain_result.resources:
            if not resource.endpoints:
                skipped += 1
                continue
            result = _process_resource(
                resource, domain_result.domain, repo_root, metadata, args.apply,
            )
            _print_result(result)
            if result.status in ('updated', 'would update'):
                changed += 1
            elif 'error' in result.status:
                errors += 1
            else:
                skipped += 1

    print(f'\n  Changed: {changed}  |  Skipped/no-change: {skipped}  |  Errors: {errors}')
    if not args.apply:
        print('\n  (dry-run — re-run with --apply to write changes)')
    return 0 if errors == 0 else 1


# =============================================================================
# CLI
# =============================================================================


def _find_repo_root(start: Path) -> Path:
    """
    Walk up from *start* until a ``go.mod`` file is found.

    Args:
        start: Directory to begin the upward search from.

    Returns:
        Repository root path, or the current working directory as a fallback.
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
        description=(
            'Update ReadPermissions/WritePermissions fields in Go Terraform '
            'provider service source files.'
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=(
            'Examples:\n'
            '  %(prog)s                          # dry-run all domains\n'
            '  %(prog)s --apply                  # apply all changes\n'
            '  %(prog)s --domain applications --apply\n'
            '  %(prog)s --category resources --apply\n'
        ),
    )
    parser.add_argument(
        '--root', metavar='PATH',
        help='Repository root (default: auto-detected via go.mod).',
    )
    parser.add_argument(
        '--apply', action='store_true',
        help='Write changes to source files (default: dry-run).',
    )
    parser.add_argument(
        '--domain', metavar='DOMAIN',
        help='Limit to this service domain (e.g. device_management).',
    )
    parser.add_argument(
        '--category', metavar='CATEGORY',
        choices=sorted({'actions', 'datasources', 'ephemerals', 'list-resources', 'resources'}),
        help='Limit to this service category.',
    )
    parser.add_argument(
        '--refresh-metadata', action='store_true', dest='refresh_metadata',
        help='Force re-download of Graph command metadata.',
    )
    return parser.parse_args()


def main() -> None:
    """Entry point: parse arguments and run the async update pipeline."""
    sys.exit(asyncio.run(_run(_parse_args())))


if __name__ == '__main__':
    main()
