#!/usr/bin/env python3
"""
get_graph_endpoint_permissions.py

Discovers required Microsoft Graph API permissions for a specific endpoint by
fetching the Microsoft Graph PowerShell SDK command metadata directly from GitHub
and parsing it locally — no PowerShell modules required.

Optionally compares these requirements against an existing enterprise application's
granted permissions using the Microsoft Graph Python SDK (msgraph-sdk) and
Azure Identity (azure-identity) for OAuth 2.0 client credentials authentication.

The script provides two modes of operation:
  1. Discovery Mode:  Shows what permissions are required for a given endpoint.
  2. Evaluation Mode: Compares an app's current permissions against the
                      requirements and reports any gaps with remediation guidance.

Permission metadata is cached to the system temp directory for 24 hours to avoid
repeated GitHub downloads. Use --refresh-metadata to force a fresh download.

Usage:
    python3 get_graph_endpoint_permissions.py --uri "/groups/{id}" --method GET

    python3 get_graph_endpoint_permissions.py --uri "/groups" --method POST \\
        --app-id-to-be-evaluated "00000000-0000-0000-0000-000000000001" \\
        --tenant-id             "00000000-0000-0000-0000-000000000002" \\
        --client-secret         "your-secret-here"

Author:    Deployment Theory
Version:   1.0
Reference: https://learn.microsoft.com/en-us/graph/permissions-reference

Metadata source: https://github.com/microsoftgraph/msgraph-sdk-powershell
"""

from __future__ import annotations

import argparse
import asyncio
import json
import os
import re
import sys
import textwrap
import time
from pathlib import Path
from typing import Any

import httpx
from azure.identity import ClientSecretCredential
from msgraph import GraphServiceClient
from msgraph.generated.models.app_role_assignment import AppRoleAssignment
from msgraph.generated.models.directory_role import DirectoryRole
from msgraph.generated.models.service_principal import ServicePrincipal
from msgraph.generated.service_principals.service_principals_request_builder import (
    ServicePrincipalsRequestBuilder,
)

# =============================================================================
# CONSTANTS
# =============================================================================

METADATA_URLS: list[str] = [
    (
        "https://raw.githubusercontent.com/microsoftgraph/msgraph-sdk-powershell/main/"
        "src/Authentication/Authentication/custom/common/MgCommandMetadata.json"
    ),
    (
        "https://raw.githubusercontent.com/microsoftgraph/msgraph-sdk-powershell/dev/"
        "src/Authentication/Authentication/custom/common/MgCommandMetadata.json"
    ),
]

CACHE_DIR_NAME: str = "MgGraphPermissionsRaw"
CACHE_FILE_NAME: str = "MgCommandMetadata.json"
CACHE_TTL_HOURS: int = 24
MICROSOFT_GRAPH_APP_ID: str = "00000003-0000-0000-c000-000000000000"

# ANSI colour codes for terminal output
_CYAN: str = "\033[96m"
_GREEN: str = "\033[92m"
_YELLOW: str = "\033[93m"
_RED: str = "\033[91m"
_WHITE: str = "\033[97m"
_RESET: str = "\033[0m"


# =============================================================================
# LOGGING HELPERS
# =============================================================================


def _indent(level: int) -> str:
    """Return a string of spaces for the requested indentation level."""
    return "  " * (level + 1)


def log_header(message: str) -> None:
    """Print a major section header in cyan."""
    print(f"\n{_CYAN}{'=' * 40}{_RESET}")
    print(f"{_CYAN} {message}{_RESET}")
    print(f"{_CYAN}{'=' * 40}{_RESET}")


def log_section(message: str) -> None:
    """Print a sub-section label in yellow."""
    print(f"\n{_YELLOW}  {message}{_RESET}")
    print(f"{_YELLOW}  {'-' * len(message)}{_RESET}")


def log_info(message: str, indent: int = 0) -> None:
    """Print an informational message in white."""
    print(f"{_WHITE}{_indent(indent)}{message}{_RESET}")


def log_success(message: str, indent: int = 0) -> None:
    """Print a success message in green with a check mark."""
    print(f"{_GREEN}{_indent(indent)}✓ {message}{_RESET}")


def log_warning(message: str, indent: int = 0) -> None:
    """Print a warning message in yellow with a warning symbol."""
    print(f"{_YELLOW}{_indent(indent)}⚠ {message}{_RESET}")


def log_error(message: str, indent: int = 0) -> None:
    """Print an error message in red with an x symbol."""
    print(f"{_RED}{_indent(indent)}✗ {message}{_RESET}")


# =============================================================================
# METADATA FUNCTIONS
# =============================================================================


def _get_cache_path() -> Path:
    """Resolve the platform-appropriate cache file path."""
    tmp = Path(
        os.environ.get("TEMP")
        or os.environ.get("TMPDIR")
        or "/tmp"
    )
    return tmp / CACHE_DIR_NAME / CACHE_FILE_NAME


async def get_graph_command_metadata(force_refresh: bool = False) -> list[dict[str, Any]]:
    """
    Download or retrieve the cached Microsoft Graph SDK command metadata.

    The metadata maps Graph API URIs and HTTP methods to their required permissions.
    It is sourced from the public Microsoft Graph PowerShell SDK GitHub repository
    and cached locally in the system temp directory for ``CACHE_TTL_HOURS`` hours.

    Args:
        force_refresh: When True, bypass the local cache and re-download.

    Returns:
        Parsed list of command metadata entries.

    Raises:
        RuntimeError: If the metadata cannot be downloaded from any known URL.
    """
    cache_path = _get_cache_path()

    if not force_refresh and cache_path.exists():
        age_hours = (time.time() - cache_path.stat().st_mtime) / 3600
        if age_hours < CACHE_TTL_HOURS:
            log_success(
                f"Using cached metadata (age: {age_hours:.1f}h, path: {cache_path})",
                indent=1,
            )
            try:
                data: list[dict[str, Any]] = json.loads(cache_path.read_text(encoding="utf-8"))
                return data
            except (json.JSONDecodeError, OSError):
                log_warning("Cache file is corrupted, re-downloading...", indent=1)
        else:
            log_info(f"Cache is stale ({age_hours:.1f}h old), re-downloading...", indent=1)

    cache_path.parent.mkdir(parents=True, exist_ok=True)

    async with httpx.AsyncClient(timeout=60.0) as http_client:
        for url in METADATA_URLS:
            try:
                log_info("Downloading Graph command metadata...", indent=1)
                log_info(f"Source: {url}", indent=2)
                response = await http_client.get(url)
                response.raise_for_status()
                entries: list[dict[str, Any]] = response.json()
                cache_path.write_text(json.dumps(entries), encoding="utf-8")
                log_success(
                    f"Metadata downloaded and cached ({len(entries):,} entries)",
                    indent=1,
                )
                return entries
            except (httpx.HTTPError, ValueError) as exc:
                log_warning(f"Failed to download from {url} — {exc}", indent=1)

    raise RuntimeError(
        "Could not download Graph command metadata from any known URL. "
        "Check network connectivity."
    )


def normalize_uri_template(uri: str) -> str:
    """
    Normalise a Graph URI template for reliable equality comparison.

    Replaces all ``{param}`` tokens with ``{*}`` and lowercases the result,
    enabling matching regardless of parameter naming differences (e.g.
    ``{id}`` vs ``{application-id}``).

    Args:
        uri: The URI template string to normalise.

    Returns:
        Lowercased URI with all ``{...}`` placeholders replaced by ``{*}``.
    """
    return re.sub(r"\{[^}]+\}", "{*}", uri.rstrip("/").lower())


def find_endpoint_permissions(
    metadata: list[dict[str, Any]],
    uri: str,
    method: str,
    api_version: str,
) -> list[dict[str, Any]]:
    """
    Return all metadata entries matching the supplied URI, method, and API version.

    Both the user-supplied URI and each metadata entry's URI are normalised via
    :func:`normalize_uri_template` before comparison, so parameter name differences
    are handled transparently.

    Args:
        metadata:    Full list of command metadata entries.
        uri:         Graph API URI to match.
        method:      HTTP method filter. Empty string matches all methods.
        api_version: API version filter (``"v1.0"`` or ``"beta"``).

    Returns:
        List of matching metadata entry dicts.
    """
    normalized_input = normalize_uri_template(uri)
    matched: list[dict[str, Any]] = []

    for entry in metadata:
        if entry.get("ApiVersion") != api_version:
            continue
        if method and entry.get("Method", "").upper() != method.upper():
            continue
        if normalize_uri_template(entry.get("Uri", "")) == normalized_input:
            matched.append(entry)

    return matched


def get_permission_name(perm_obj: Any) -> str | None:
    """
    Extract a permission name string from a permission object or plain string.

    The metadata JSON represents permissions as plain strings (older SDK versions)
    or as objects with a ``Name`` property (newer SDK versions). This function
    normalises both forms to a plain string.

    Args:
        perm_obj: A permission entry from the metadata (string or dict).

    Returns:
        The permission name string, or ``None`` if it cannot be determined.
    """
    if isinstance(perm_obj, str):
        return perm_obj
    if isinstance(perm_obj, dict):
        return perm_obj.get("Name") or perm_obj.get("name")
    name = getattr(perm_obj, "Name", None) or getattr(perm_obj, "name", None)
    return str(name) if name is not None else None


def is_application_permission(name: str) -> bool:
    """
    Return True if the permission is an Application (not Delegated) permission.

    Application permissions end in ``.All`` and do not contain ``AccessAsUser``.

    Args:
        name: The permission name string to test.

    Returns:
        True for application permissions, False for delegated permissions.
    """
    return name.endswith(".All") and "AccessAsUser" not in name


# =============================================================================
# GRAPH SDK EVALUATION FUNCTIONS
# =============================================================================


def _make_graph_client(
    tenant_id: str,
    client_id: str,
    client_secret: str,
) -> GraphServiceClient:
    """
    Build an authenticated Microsoft Graph SDK client using client credentials.

    Args:
        tenant_id:     Entra ID tenant ID.
        client_id:     Application (client) ID.
        client_secret: Client secret for the application.

    Returns:
        A configured :class:`GraphServiceClient` instance.
    """
    credential = ClientSecretCredential(
        tenant_id=tenant_id,
        client_id=client_id,
        client_secret=client_secret,
    )
    return GraphServiceClient(
        credentials=credential,
        scopes=["https://graph.microsoft.com/.default"],
    )


async def get_service_principal_by_app_id(
    client: GraphServiceClient,
    app_id: str,
) -> ServicePrincipal:
    """
    Retrieve a service principal object by its application (client) ID.

    Calls ``GET /servicePrincipals?$filter=appId eq '{app_id}'`` and returns
    the first matching result.

    Args:
        client: Authenticated Graph SDK client.
        app_id: The application (client) ID to look up.

    Returns:
        The matching :class:`ServicePrincipal` object.

    Raises:
        ValueError:   If no matching service principal is found.
        RuntimeError: If the Graph API call fails.
    """
    log_info(f"Retrieving service principal for AppId: {app_id}...", indent=1)

    QueryParams = ServicePrincipalsRequestBuilder.ServicePrincipalsRequestBuilderGetQueryParameters
    RequestConfig = (
        ServicePrincipalsRequestBuilder.ServicePrincipalsRequestBuilderGetRequestConfiguration
    )

    request_config = RequestConfig(
        query_parameters=QueryParams(
            filter=f"appId eq '{app_id}'",
            select=["id", "appId", "displayName"],
        )
    )

    try:
        result = await client.service_principals.get(request_configuration=request_config)
    except Exception as exc:
        raise RuntimeError(f"Failed to retrieve service principal: {exc}") from exc

    if not result or not result.value:
        raise ValueError(f"No service principal found with appId '{app_id}'")

    return result.value[0]


async def get_app_role_assignments(
    client: GraphServiceClient,
    sp_id: str,
) -> list[AppRoleAssignment]:
    """
    Return all app role assignments granted to a service principal.

    Calls ``GET /servicePrincipals/{id}/appRoleAssignments``. The returned
    assignment objects contain ``app_role_id`` GUIDs that must be resolved
    against the Microsoft Graph resource SP's ``app_roles`` collection to
    obtain human-readable permission names.

    Args:
        client: Authenticated Graph SDK client.
        sp_id:  Object ID of the service principal.

    Returns:
        List of :class:`AppRoleAssignment` objects.

    Raises:
        RuntimeError: If the Graph API call fails.
    """
    log_info(f"Retrieving app role assignments for service principal '{sp_id}'...", indent=1)

    try:
        result = (
            await client.service_principals.by_service_principal_id(sp_id)
            .app_role_assignments.get()
        )
    except Exception as exc:
        raise RuntimeError(f"Failed to retrieve app role assignments: {exc}") from exc

    return result.value if result and result.value else []


async def get_microsoft_graph_service_principal(
    client: GraphServiceClient,
) -> ServicePrincipal:
    """
    Retrieve the Microsoft Graph first-party service principal.

    Calls ``GET /servicePrincipals?$filter=appId eq '{MICROSOFT_GRAPH_APP_ID}'``.
    The returned object contains the ``app_roles`` collection used to resolve
    app role assignment GUIDs back to human-readable permission names.

    Args:
        client: Authenticated Graph SDK client.

    Returns:
        The Microsoft Graph first-party :class:`ServicePrincipal` object.

    Raises:
        RuntimeError: If the service principal cannot be retrieved.
    """
    log_info("Retrieving Microsoft Graph first-party service principal...", indent=1)

    QueryParams = ServicePrincipalsRequestBuilder.ServicePrincipalsRequestBuilderGetQueryParameters
    RequestConfig = (
        ServicePrincipalsRequestBuilder.ServicePrincipalsRequestBuilderGetRequestConfiguration
    )

    request_config = RequestConfig(
        query_parameters=QueryParams(
            filter=f"appId eq '{MICROSOFT_GRAPH_APP_ID}'",
            select=["id", "appId", "appRoles"],
        )
    )

    try:
        result = await client.service_principals.get(request_configuration=request_config)
    except Exception as exc:
        raise RuntimeError(f"Failed to retrieve Microsoft Graph SP: {exc}") from exc

    if not result or not result.value:
        raise RuntimeError("Microsoft Graph first-party service principal not found in tenant")

    return result.value[0]


async def get_service_principal_directory_roles(
    client: GraphServiceClient,
    sp_id: str,
) -> list[DirectoryRole]:
    """
    Return the directory roles currently assigned to a service principal.

    Calls ``GET /servicePrincipals/{id}/memberOf/microsoft.graph.directoryRole``
    which returns only directory role memberships (filtering out groups and other
    directory objects automatically via the type-cast OData path).

    Args:
        client: Authenticated Graph SDK client.
        sp_id:  Object ID of the service principal.

    Returns:
        List of :class:`DirectoryRole` objects.

    Raises:
        RuntimeError: If the Graph API call fails.
    """
    log_info("Retrieving directory role memberships...", indent=1)

    try:
        result = (
            await client.service_principals.by_service_principal_id(sp_id)
            .member_of.graph_directory_role.get()
        )
    except Exception as exc:
        raise RuntimeError(f"Failed to retrieve directory roles: {exc}") from exc

    return result.value if result and result.value else []


# =============================================================================
# MAIN EXECUTION HELPERS
# =============================================================================

# Collected permission data returned from _collect_permissions.
# Tuple layout: (app_perms, del_perms, perm_uses)
_PermData = tuple[list[str], list[str], dict[str, list[str]]]


def _collect_permissions(matched_commands: list[dict[str, Any]]) -> _PermData:
    """
    Iterate matched metadata commands and categorise all unique permissions.

    Prints the per-command analysis block as a side effect.

    Args:
        matched_commands: List of metadata entries returned by
                          :func:`find_endpoint_permissions`.

    Returns:
        A three-tuple of ``(app_perms, del_perms, perm_uses)`` where
        ``app_perms`` and ``del_perms`` are deduplicated sorted-ready lists
        and ``perm_uses`` maps each permission name to the command tags that
        require it.
    """
    app_perms: list[str] = []
    del_perms: list[str] = []
    perm_uses: dict[str, list[str]] = {}

    log_section("Analysing Commands and Permissions")

    for cmd in matched_commands:
        print()
        log_info(f"Command: {cmd.get('Command', 'N/A')}", indent=1)
        log_info(f"Method:  {cmd.get('Method', 'N/A')}", indent=1)
        log_info(f"URI:     {cmd.get('Uri', 'N/A')}", indent=1)

        permissions: list[Any] = cmd.get("Permissions") or []
        if not permissions:
            log_warning("No permissions documented in metadata for this command", indent=2)
            continue

        log_info("Permissions:", indent=1)
        tag = f"{cmd.get('Method', '')} {cmd.get('Uri', '')}"
        for perm_obj in permissions:
            perm = get_permission_name(perm_obj)
            if not perm:
                continue
            perm_uses.setdefault(perm, []).append(tag)
            if is_application_permission(perm):
                if perm not in app_perms:
                    app_perms.append(perm)
                log_info(f"[App] {perm}", indent=2)
            else:
                if perm not in del_perms:
                    del_perms.append(perm)
                log_info(f"[Del] {perm}", indent=2)

    return app_perms, del_perms, perm_uses


def _print_permission_summary(
    app_perms: list[str],
    del_perms: list[str],
    perm_uses: dict[str, list[str]],
) -> None:
    """
    Print the Permission Requirements Summary block.

    Args:
        app_perms: Deduplicated list of application permission names.
        del_perms: Deduplicated list of delegated permission names.
        perm_uses: Map of permission name → list of command tags that use it.
    """
    log_header("Permission Requirements Summary")

    if app_perms:
        log_section("Application Permissions (for service principals / daemon apps)")
        for perm in sorted(app_perms):
            log_success(perm, indent=1)
            for use in sorted(set(perm_uses.get(perm, []))):
                log_info(f"→ {use}", indent=2)

    if del_perms:
        log_section("Delegated Permissions (for user context / interactive apps)")
        for perm in sorted(del_perms):
            log_info(perm, indent=1)
            for use in sorted(set(perm_uses.get(perm, []))):
                log_info(f"→ {use}", indent=2)

    if not app_perms and not del_perms:
        log_warning(
            "No permissions were documented in the metadata for the matched command(s)",
            indent=1,
        )


def _compare_and_print_permissions(
    app_perms: list[str],
    granted_perms: list[str],
) -> tuple[list[str], list[str]]:
    """
    Compare required application permissions against those already granted and print results.

    Args:
        app_perms:     Required application permissions (from metadata discovery).
        granted_perms: Permissions already granted to the service principal.

    Returns:
        A two-tuple of ``(covered_perms, missing_perms)``.
    """
    missing: list[str] = []
    covered: list[str] = []
    for required in sorted(app_perms):
        if required in granted_perms:
            covered.append(required)
            log_success(f"{required}  [GRANTED]", indent=1)
        else:
            missing.append(required)
            log_error(f"{required}  [MISSING]", indent=1)
    return covered, missing


def _resolve_granted_permissions(
    assignments: list[AppRoleAssignment],
    mg_app_roles: list[Any],
) -> list[str]:
    """
    Resolve app role assignment GUIDs to human-readable permission names.

    Args:
        assignments:  List of :class:`AppRoleAssignment` objects.
        mg_app_roles: ``app_roles`` collection from the Microsoft Graph SP.

    Returns:
        List of resolved permission name strings.
    """
    granted: list[str] = []
    for assignment in assignments:
        role = next(
            (r for r in mg_app_roles if r.id == assignment.app_role_id),
            None,
        )
        if role and role.value:
            granted.append(role.value)
    return granted


def _print_coverage_summary(
    display_name: str,
    app_id: str,
    app_perms: list[str],
    covered_perms: list[str],
    missing_perms: list[str],
) -> None:
    """
    Print the Evaluation Summary / Coverage Status block.

    Args:
        display_name:  Display name of the service principal.
        app_id:        Application (client) ID.
        app_perms:     All required application permissions.
        covered_perms: Permissions that are already granted.
        missing_perms: Permissions that are not yet granted.
    """
    log_header("Evaluation Summary")
    log_section("Coverage Status")
    log_info(f"Service Principal:              {display_name}", indent=1)
    log_info(f"AppId:                          {app_id}", indent=1)
    log_info(f"Total Required App Permissions: {len(app_perms)}", indent=1)
    log_success(f"Covered: {len(covered_perms)}", indent=1)

    if missing_perms:
        log_error(f"Missing: {len(missing_perms)}", indent=1)
    else:
        log_success("Missing: 0", indent=1)

    if missing_perms:
        log_section("Missing Permissions — Remediation")
        log_info("The following application permissions must be granted:", indent=1)
        for perm in missing_perms:
            log_error(perm, indent=2)
        print()
        log_info("To grant these permissions:", indent=1)
        log_info("1. Open Azure Portal → Entra ID → Enterprise Applications", indent=2)
        log_info(f"2. Locate: {display_name}  (AppId: {app_id})", indent=2)
        log_info("3. Navigate to: API Permissions", indent=2)
        log_info("4. Add each missing Microsoft Graph Application permission listed above",
                 indent=2)
        log_info("5. Click 'Grant admin consent'", indent=2)
    else:
        print()
        log_success(
            "All required application permissions are granted — no action needed",
            indent=1,
        )


def _check_role_assignable_groups(
    directory_roles: list[DirectoryRole],
    uri: str,
    method: str,
) -> None:
    """
    Print the special requirement note for role-assignable group creation.

    Only outputs anything when the URI targets groups and the method is POST
    (or unfiltered), since ``isAssignableToRole=true`` requires a specific
    directory role beyond the standard Graph application permissions.

    Args:
        directory_roles: Directory roles currently held by the service principal.
        uri:             The Graph API URI being evaluated.
        method:          The HTTP method being evaluated.
    """
    if "groups" not in uri.lower() or method.upper() not in ("POST", ""):
        return

    log_section("Special Requirement: Role-Assignable Groups")
    if not directory_roles:
        log_warning(
            "Creating groups with isAssignableToRole=true also requires one of:",
            indent=1,
        )
        log_info("- Privileged Role Administrator (directory role)", indent=2)
        log_info("- Global Administrator (directory role)", indent=2)
        return

    required_names = {"Privileged Role Administrator", "Global Administrator"}
    has_required = any((r.display_name or "") in required_names for r in directory_roles)
    if has_required:
        log_success(
            "Service principal holds the required directory role for "
            "role-assignable group creation",
            indent=1,
        )
    else:
        log_warning(
            "Service principal has directory roles but none required for "
            "isAssignableToRole groups",
            indent=1,
        )
        log_info("Required: Privileged Role Administrator OR Global Administrator", indent=2)


async def _run_evaluation_mode(
    client: GraphServiceClient,
    app_id: str,
    app_perms: list[str],
    uri: str,
    method: str,
) -> None:
    """
    Perform evaluation mode: compare required permissions against those granted to an app.

    Args:
        client:    Authenticated Graph SDK client.
        app_id:    Application (client) ID of the enterprise app to evaluate.
        app_perms: Required application permissions discovered from the metadata.
        uri:       The Graph API URI being evaluated (used for group role check).
        method:    The HTTP method being evaluated (used for group role check).
    """
    sp = await get_service_principal_by_app_id(client, app_id)
    display_name = sp.display_name or "Unknown"
    sp_id = sp.id or ""

    log_success(f"Found service principal: {display_name}", indent=1)
    log_info(f"Object ID: {sp_id}", indent=2)

    assignments = await get_app_role_assignments(client, sp_id)
    log_success(f"Retrieved {len(assignments)} app role assignment(s)", indent=1)

    mg_sp_roles = (await get_microsoft_graph_service_principal(client)).app_roles or []
    granted_perms = _resolve_granted_permissions(assignments, mg_sp_roles)
    log_success(f"Resolved {len(granted_perms)} granted permission name(s)", indent=1)

    log_section("Permission Comparison")
    covered_perms, missing_perms = _compare_and_print_permissions(app_perms, granted_perms)

    log_section("Directory Role Assignments")
    directory_roles = await get_service_principal_directory_roles(client, sp_id)
    if not directory_roles:
        log_error(
            "No directory roles currently assigned to this service principal",
            indent=1,
        )
    else:
        for role in directory_roles:
            log_success(role.display_name or "Unknown role", indent=1)

    _print_coverage_summary(display_name, app_id, app_perms, covered_perms, missing_perms)
    _check_role_assignable_groups(directory_roles, uri, method)


def _print_discovery_recommendations(
    app_perms: list[str],
    del_perms: list[str],
    uri: str,
    method: str,
) -> None:
    """
    Print the Recommendations block shown in discovery-only mode.

    Args:
        app_perms: Required application permissions.
        del_perms: Required delegated permissions.
        uri:       The Graph API URI (used for group role note).
        method:    The HTTP method (used for group role note).
    """
    log_header("Recommendations")

    if app_perms:
        log_section("For Service Principal (Application Authentication)")
        log_info(
            "Grant these Microsoft Graph Application permissions in Azure Portal:",
            indent=1,
        )
        for perm in sorted(app_perms):
            log_info(f"  • {perm}", indent=1)
        print()
        log_info("Admin consent is required for all application permissions", indent=1)

        if "groups" in uri.lower() and method.upper() in ("POST", ""):
            print()
            log_warning(
                "For groups with isAssignableToRole=true, the service principal also needs:",
                indent=1,
            )
            log_info(
                "  • Privileged Role Administrator  OR  Global Administrator  (directory role)",
                indent=1,
            )

    if del_perms:
        log_section("For Delegated (User Context) Authentication")
        log_info("Grant these Microsoft Graph Delegated permissions:", indent=1)
        for perm in sorted(del_perms):
            log_info(f"  • {perm}", indent=1)


def _print_documentation_references() -> None:
    """Print the Documentation References footer block."""
    log_header("Documentation References")
    log_info("Microsoft Graph Permissions Reference:", indent=1)
    log_info("  https://learn.microsoft.com/en-us/graph/permissions-reference", indent=1)
    print()
    log_info("Microsoft Graph API Reference:", indent=1)
    log_info("  https://learn.microsoft.com/en-us/graph/api/overview", indent=1)
    print()
    log_info("Graph Explorer (interactive permission discovery):", indent=1)
    log_info("  https://developer.microsoft.com/en-us/graph/graph-explorer", indent=1)
    print()


# =============================================================================
# MAIN EXECUTION
# =============================================================================


async def _run(args: argparse.Namespace) -> int:
    """
    Orchestrate discovery and optional evaluation of Graph endpoint permissions.

    Args:
        args: Parsed command-line arguments.

    Returns:
        Exit code (0 for success, 1 for failure).
    """
    # Step 1: Load metadata
    log_header("Loading Graph Command Metadata")
    metadata = await get_graph_command_metadata(force_refresh=args.refresh_metadata)
    if not metadata:
        log_error("Metadata loaded but contained no entries")
        return 1
    log_info(f"Total command entries available: {len(metadata):,}", indent=1)

    # Step 2: Discover matching commands
    log_header("Discovering Commands for URI")
    log_info(f"URI:         {args.uri}", indent=1)
    log_info(f"API Version: {args.api_version}", indent=1)
    log_info(f"Method:      {args.method or '(all methods)'}", indent=1)
    log_info(f"Normalised:  {normalize_uri_template(args.uri)}", indent=1)

    matched = find_endpoint_permissions(metadata, args.uri, args.method, args.api_version)
    if not matched:
        log_error(f"No commands found for: {args.method or 'ANY'} {args.uri} ({args.api_version})")
        print()
        log_warning("Troubleshooting suggestions:", indent=1)
        log_info("- Use {id} or any {placeholder} for resource identifiers: /groups/{id}", indent=2)
        log_info("- Try the alternate API version with --api-version v1.0 or beta", indent=2)
        log_info("- Some deep nested paths are not indexed in the metadata", indent=2)
        log_info("- Use --refresh-metadata to download the latest SDK metadata", indent=2)
        return 1
    log_success(f"Found {len(matched)} matching command(s)", indent=1)

    # Steps 3 + 4: Collect permissions and print summary
    app_perms, del_perms, perm_uses = _collect_permissions(matched)
    _print_permission_summary(app_perms, del_perms, perm_uses)

    # Step 5: Evaluation mode or discovery recommendations
    if args.app_id_to_be_evaluated:
        log_header("Evaluating Enterprise Application")
        graph_client = _make_graph_client(
            tenant_id=args.tenant_id,
            client_id=args.app_id_to_be_evaluated,
            client_secret=args.client_secret,
        )
        await _run_evaluation_mode(
            graph_client, args.app_id_to_be_evaluated, app_perms, args.uri, args.method
        )
    else:
        _print_discovery_recommendations(app_perms, del_perms, args.uri, args.method)

    _print_documentation_references()
    return 0


def _parse_args() -> argparse.Namespace:
    """Parse and validate command-line arguments."""
    parser = argparse.ArgumentParser(
        description="Discover Microsoft Graph API permissions for a specific endpoint.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=textwrap.dedent("""\
            Examples:
              %(prog)s --uri "/groups/{id}" --method GET
              %(prog)s --uri "/applications/{id}" --method PATCH \\
                  --app-id-to-be-evaluated "00000000-0000-0000-0000-000000000001" \\
                  --tenant-id             "00000000-0000-0000-0000-000000000002" \\
                  --client-secret         "your-secret"
        """),
    )

    parser.add_argument(
        "--uri",
        required=True,
        metavar="URI",
        help="Graph API URI to check (e.g. '/groups', '/groups/{id}').",
    )
    parser.add_argument(
        "--method",
        default="",
        choices=["GET", "POST", "PATCH", "PUT", "DELETE", ""],
        metavar="METHOD",
        help="HTTP method filter (GET, POST, PATCH, PUT, DELETE). Omit for all methods.",
    )
    parser.add_argument(
        "--api-version",
        default="beta",
        choices=["v1.0", "beta"],
        dest="api_version",
        help="API version to check against (default: beta).",
    )
    parser.add_argument(
        "--app-id-to-be-evaluated",
        dest="app_id_to_be_evaluated",
        metavar="APP_ID",
        help="Application (Client) ID of the enterprise app to evaluate.",
    )
    parser.add_argument(
        "--tenant-id",
        dest="tenant_id",
        metavar="TENANT_ID",
        help="Entra ID tenant ID (required with --app-id-to-be-evaluated).",
    )
    parser.add_argument(
        "--client-secret",
        dest="client_secret",
        metavar="SECRET",
        help="Client secret (required with --app-id-to-be-evaluated).",
    )
    parser.add_argument(
        "--refresh-metadata",
        action="store_true",
        dest="refresh_metadata",
        help="Force a fresh download of the Graph SDK command metadata from GitHub.",
    )

    args = parser.parse_args()

    # Cross-validate evaluation parameters
    if args.app_id_to_be_evaluated:
        if not args.tenant_id:
            parser.error("--tenant-id is required when --app-id-to-be-evaluated is provided")
        if not args.client_secret:
            parser.error("--client-secret is required when --app-id-to-be-evaluated is provided")

    return args


def main() -> None:
    """Entry point: parse arguments and run the async core logic."""
    args = _parse_args()
    exit_code = asyncio.run(_run(args))
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
