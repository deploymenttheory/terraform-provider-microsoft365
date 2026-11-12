#!/usr/bin/env python3
"""
update_graph_service_permission_doc_templates.py

Reads the ReadPermissions / WritePermissions fields from each Go service
constructor and synchronises the **Required:** permission list in the
corresponding Terraform documentation template (.md.tmpl file).

Source of truth
---------------
The Go constructor function (New*()) in each service directory — the same
file that update_graph_service_permissions.py writes to.  The script reads
the current state of those files, so it should be run *after*
``update_graph_service_permissions.py --apply``.

Template update logic
---------------------
* Resources / actions     : Required = union(ReadPermissions, WritePermissions)
* Datasources / list-resources / ephemerals : Required = ReadPermissions only

The **Optional:** block is never modified by this script.

Safety behaviour
----------------
* Dry-run by default — pass ``--apply`` to write changes.
* Only updates templates where a **Required:** block is found.
* Skips resources with no Go constructor file or no permissions.

Usage::

    # Preview what would change (dry-run)
    python3 update_graph_service_permission_doc_templates.py

    # Write changes to template files
    python3 update_graph_service_permission_doc_templates.py --apply

    # Limit to one service domain or category
    python3 update_graph_service_permission_doc_templates.py --domain applications --apply
    python3 update_graph_service_permission_doc_templates.py --category resources --apply

Author:  Deployment Theory
Version: 1.0
"""

from __future__ import annotations

import argparse
import re
import sys
from dataclasses import dataclass, field
from functools import partial
from pathlib import Path

# ---------------------------------------------------------------------------
# Extend sys.path so the sibling scripts in the same directory are importable.
# ---------------------------------------------------------------------------
_SCRIPT_DIR = Path(__file__).resolve().parent
if str(_SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(_SCRIPT_DIR))

# pylint: disable=wrong-import-position
from get_utilised_graph_api_endpoints import (  # noqa: E402
    SERVICES_SUBPATH,
    ResourceEndpoints,
    scan_services,
)
from permission_update_common import (
    build_update_parser,
)
from permission_update_common import find_constructor_file as _find_constructor_file
from permission_update_common import find_repo_root as _find_repo_root
from permission_update_common import scan_update_domains, update_scanned_resources

# pylint: enable=wrong-import-position

# =============================================================================
# CONSTANTS
# =============================================================================

# Service categories that carry both ReadPermissions and WritePermissions.
# For these, the template Required list is the union of both sets.
_WRITE_CATEGORIES: frozenset[str] = frozenset({"resources", "actions"})

# Maps Go service category name → subdirectory name under templates/
_CATEGORY_TEMPLATE_DIR: dict[str, str] = {
    "resources": "resources",
    "datasources": "data-sources",
    "actions": "actions",
    "list-resources": "list-resources",
    "ephemerals": "ephemeral-resources",
}

# Maps the Go API-version directory name → template filename prefix.
# graph_v1.0 resources use the bare "graph_" prefix in their template names.
_API_DIR_TO_PREFIX: dict[str, str] = {
    "graph_beta": "graph_beta_",
    "graph_v1.0": "graph_",
}

# Subpath within the repo root where documentation templates live.
_TEMPLATES_SUBPATH: str = "templates"

# Regex: ReadPermissions []string{...} block in a Go constructor
_RE_READ_BLOCK = re.compile(
    r"ReadPermissions:\s*\[\]string\{([^}]*)\}",
    re.DOTALL,
)

# Regex: WritePermissions []string{...} block in a Go constructor
_RE_WRITE_BLOCK = re.compile(
    r"WritePermissions:\s*\[\]string\{([^}]*)\}",
    re.DOTALL,
)

# Regex: individual quoted permission string inside a []string{} block
_RE_PERM_ITEM = re.compile(r'"([^"]+)"')

# Regex: locate any New*() constructor function signature

# Regex: the **Required:** header and its ``- `permission` `` list lines.
# Captures only the lines that follow the header so the blank line before
# **Optional:** is preserved unchanged.
_RE_REQUIRED_BLOCK = re.compile(
    r"(?P<header>\*\*Required:\*\*\n)" + r"(?P<items>(?:- `[^\n]+`\n)+)",
    re.MULTILINE,
)

# =============================================================================
# DATA CLASSES
# =============================================================================


@dataclass
class _TemplateResult:
    """Outcome of processing one resource's documentation template."""

    domain: str
    resource: ResourceEndpoints
    go_file: Path | None
    template_file: Path | None
    required_perms: list[str] = field(default_factory=list)
    status: str = "pending"


# =============================================================================
# GO CONSTRUCTOR READER
# =============================================================================


def _read_go_permissions(go_file: Path) -> tuple[list[str], list[str]]:
    """
    Extract ReadPermissions and WritePermissions from a Go constructor file.

    Args:
        go_file: Path to the Go source file.

    Returns:
        ``(sorted_read_perms, sorted_write_perms)`` lists, either may be empty.
    """
    try:
        text = go_file.read_text(encoding="utf-8")
    except OSError:
        return [], []

    read_match = _RE_READ_BLOCK.search(text)
    write_match = _RE_WRITE_BLOCK.search(text)

    # Strip trailing commas that may appear inside string literals due to
    # copy-paste errors in the Go source (e.g. "Perm.ReadWrite.All,").
    read_perms = (
        sorted(p.rstrip(",") for p in _RE_PERM_ITEM.findall(read_match.group(1)))
        if read_match
        else []
    )
    write_perms = (
        sorted(p.rstrip(",") for p in _RE_PERM_ITEM.findall(write_match.group(1)))
        if write_match
        else []
    )
    return read_perms, write_perms


# =============================================================================
# TEMPLATE FILE DISCOVERY
# =============================================================================


def _find_template(templates_root: Path, resource: ResourceEndpoints) -> Path | None:
    """
    Locate the ``.md.tmpl`` documentation template for *resource*.

    The template filename is derived from the resource's service category,
    domain, API version directory, and resource name.  A fallback with the
    alternative API prefix is attempted when the primary name does not exist
    (handles cases where a v1.0 resource's template was named with the
    ``graph_beta_`` prefix).

    Args:
        templates_root: Path to the ``templates/`` directory.
        resource:       Resource scan result.

    Returns:
        Path to the template file, or ``None`` if not found.
    """
    parts = Path(resource.source_dir).parts
    # Expected layout: ('internal', 'services', category, domain, api_dir, resource_name)
    if len(parts) < 6:
        return None

    category, domain, api_dir, resource_name = parts[2], parts[3], parts[4], parts[5]

    template_subdir = _CATEGORY_TEMPLATE_DIR.get(category)
    if template_subdir is None:
        return None

    template_dir = templates_root / template_subdir
    api_prefix = _API_DIR_TO_PREFIX.get(api_dir, "graph_")

    # Primary candidate
    primary = template_dir / f"{api_prefix}{domain}_{resource_name}.md.tmpl"
    if primary.exists():
        return primary

    # Fallback: alternative API prefix (handles naming inconsistencies)
    alt_prefix = "graph_" if api_prefix == "graph_beta_" else "graph_beta_"
    fallback = template_dir / f"{alt_prefix}{domain}_{resource_name}.md.tmpl"
    if fallback.exists():
        return fallback

    return None


# =============================================================================
# TEMPLATE UPDATER
# =============================================================================


def _make_required_items(perms: list[str]) -> str:
    """
    Render the permission list as Markdown "- `perm`" lines.

    Args:
        perms: Sorted permission name strings.

    Returns:
        Multi-line string where every permission occupies one line ending
        with ``\\n``.
    """
    return "".join(f"- `{p}`\n" for p in perms)


def _update_template(source: str, required_perms: list[str]) -> str | None:
    """
    Replace the **Required:** permission list in a template.

    Locates the ``**Required:**`` block and rewrites its "- `...`" item
    lines.  The ``**Optional:**`` section and surrounding blank lines are left
    untouched.

    Args:
        source:         Full template source text.
        required_perms: New sorted list of required permission strings.

    Returns:
        Updated template string, or ``None`` when no change is needed or the
        block cannot be found.
    """
    match = _RE_REQUIRED_BLOCK.search(source)
    if not match:
        return None
    if not required_perms:
        return None

    new_items = _make_required_items(required_perms)
    new_block = match.group("header") + new_items
    updated = source[: match.start()] + new_block + source[match.end() :]
    return updated if updated != source else None


def _update_template_file(
    result: _TemplateResult,
    apply_changes: bool,
) -> None:
    """
    Update (or preview) the Required permissions block in the template.

    Mutates ``result.status`` in-place to reflect the outcome.

    Args:
        result:        Template result object (template_file must be set).
        apply_changes: When ``True``, write the updated file to disk.
    """
    if result.template_file is None:
        result.status = "no template"
        return
    if not result.required_perms:
        result.status = "no perms"
        return

    try:
        original = result.template_file.read_text(encoding="utf-8")
    except OSError as exc:
        result.status = f"read error: {exc}"
        return

    updated = _update_template(original, result.required_perms)
    if updated is None:
        result.status = "no change"
        return

    if apply_changes:
        try:
            result.template_file.write_text(updated, encoding="utf-8")
        except OSError as exc:
            result.status = f"write error: {exc}"
            return
        result.status = "updated"
    else:
        result.status = "would update"


# =============================================================================
# OUTPUT
# =============================================================================


def _fmt_perms(perms: list[str]) -> str:
    """Return a comma-separated permission list, or ``'(none)'`` if empty."""
    return ", ".join(perms) if perms else "(none)"


def _print_result(result: _TemplateResult) -> None:
    """
    Print a one-line summary row for a template update result.

    Args:
        result: Completed template result.
    """
    tmpl_label = result.template_file.name if result.template_file else "(no template)"
    name = f"{result.domain}/{result.resource.resource_name}"
    label = f"{name}  ({result.resource.category})  {tmpl_label}"
    print(f"  [{result.status:^14}] {label}")
    print(f"               Required: {_fmt_perms(result.required_perms)}")


# =============================================================================
# ORCHESTRATION
# =============================================================================


def _process_resource(
    resource: ResourceEndpoints,
    domain: str,
    repo_root: Path,
    apply_changes: bool,
) -> _TemplateResult:
    """
    Read Go permissions, find the template, and update the Required block.

    Args:
        resource:      Resource scan result.
        domain:        Service domain name (used for display).
        repo_root:     Repository root path.
        apply_changes: Whether to write changes to disk.

    Returns:
        Populated :class:`_TemplateResult` describing the outcome.
    """
    templates_root = repo_root / _TEMPLATES_SUBPATH
    resource_dir = repo_root / resource.source_dir
    go_file = _find_constructor_file(resource_dir)

    result = _TemplateResult(
        domain=domain,
        resource=resource,
        go_file=go_file,
        template_file=_find_template(templates_root, resource),
    )

    if go_file is None:
        result.status = "no go file"
        return result

    read_perms, write_perms = _read_go_permissions(go_file)

    if resource.category in _WRITE_CATEGORIES:
        required = sorted(set(read_perms) | set(write_perms))
    else:
        required = sorted(read_perms)

    result.required_perms = required
    _update_template_file(result, apply_changes)
    return result


def _run(args: argparse.Namespace) -> int:
    """
    Scan services, read Go permissions, and update documentation templates.

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
        print(f"error: services directory not found: {services_root}", file=sys.stderr)
        return 1

    templates_root = repo_root / _TEMPLATES_SUBPATH
    if not templates_root.is_dir():
        print(f"error: templates directory not found: {templates_root}", file=sys.stderr)
        return 1

    domains = scan_update_domains(scan_services, services_root, args)
    return update_scanned_resources(
        args,
        domains,
        partial(_process_resource, repo_root=repo_root, apply_changes=args.apply),
        _print_result,
        skip_unmatched=False,
    )


# =============================================================================
# CLI
# =============================================================================


def _parse_args() -> argparse.Namespace:
    """Parse the shared update flags with this tool's description and output target."""
    return build_update_parser(
        "Synchronise the **Required:** permission list in Terraform documentation templates from"
        " the Go service constructor files.",
        "template",
        metadata=False,
    ).parse_args()


def main() -> None:
    """Entry point: parse arguments and run the template update pipeline."""
    sys.exit(_run(_parse_args()))


if __name__ == "__main__":
    main()
