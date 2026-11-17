"""Verify declarative constructor additions to the provider's aggregate registries."""

import json
import re
import subprocess
from pathlib import Path

from lib.git_operations import get_merge_base
from lib.partitions import MODULE, without_comments


REGISTRIES = {
    'resources.go': ('Resources', 'resource.Resource', 'resources'),
    'datasources.go': ('DataSources', 'datasource.DataSource', 'datasources'),
    'actions.go': ('Actions', 'action.Action', 'actions'),
    'list_resources.go': ('ListResources', 'list.ListResource', 'list-resources'),
    'ephemeral_resources.go': ('EphemeralResources', 'ephemeral.EphemeralResource', 'ephemerals'),
}


def _registry_shape(source, method, interface):
    """Read literal imports and a single return of constructor references, rejecting calls or logic."""
    source = without_comments(source)
    imports = re.search(r'\bimport\s*\((?P<body>[^()]*)\)', source)
    return_type = r'\[\]func\(\)\s*' + re.escape(interface)
    registry = re.search(
        rf'\bfunc\s+\(p\s+\*M365Provider\)\s+{method}\s*\((?:ctx|_)\s+context\.Context\)'
        rf'\s*{return_type}\s*\{{\s*return\s+{return_type}\s*\{{'
        r'(?P<entries>(?:\s*\w+\.\w+\s*,)*)\s*\}\s*\}',
        source,
    )
    if not imports or not registry or imports.end() > registry.start():
        return None
    imported = re.findall(r'(?:([A-Za-z_]\w*)\s+)?"([A-Za-z0-9_./-]+)"', imports['body'])
    remainder = re.sub(r'(?:[A-Za-z_]\w*\s+)?"[A-Za-z0-9_./-]+"', '', imports['body'])
    if remainder.strip():
        return None
    aliases = {alias or path.rsplit('/', 1)[-1]: path for alias, path in imported}
    entries = re.findall(r'(\w+)\.(\w+)', registry['entries'])
    if len(aliases) != len(imported) or len(set(entries)) != len(entries):
        return None
    # Preserve literals while ignoring layout outside the two verified regions.
    token = r'"(?:\\.|[^"\\])*"|`[^`]*`|\x27(?:\\.|[^\x27\\])*\x27|[^\s]'
    outline = tuple(
        tuple(re.findall(token, text)) for text in (
            source[:imports.start()], source[imports.end():registry.start('entries')],
            source[registry.end('entries'):],
        )
    )
    return aliases, entries, outline


def registry_only_additions(name, before, after):
    """Accept only new literal imports and constructor references with existing code unchanged."""
    prefix = 'internal/provider/'
    if not name.startswith(prefix) or name[len(prefix):] not in REGISTRIES:
        return False
    method, interface, service = REGISTRIES[name[len(prefix):]]
    old = _registry_shape(before, method, interface)
    new = _registry_shape(after, method, interface)
    if old is None or new is None or old[2] != new[2]:
        return False
    old_imports, old_entries, _ = old
    imports, entries, _ = new
    if not old_imports.items() <= imports.items():
        return False
    remaining = iter(entries)
    if not all(any(entry == existing for entry in remaining) for existing in old_entries):
        return False
    added = set(entries) - set(old_entries)
    if not added or not (imports.keys() - old_imports.keys()) <= {alias for alias, _ in added}:
        return False
    for alias, constructor in added:
        path = imports.get(alias, '')
        if (
            not path.startswith(f'{MODULE}/internal/services/{service}/')
            or any(part in {'.', '..', 'mocks', 'tests'} for part in path.split('/'))
            or not re.fullmatch(r'New[A-Z]\w*' + interface.split('.')[-1], constructor)
        ):
            return False
    return True


def verified_registry_additions(added, base_ref, head):
    """Bind the registry proof to the same Git merge base and head as patch coverage."""
    candidates = [name for name in added if name in {'internal/provider/' + key for key in REGISTRIES}]
    if not candidates:
        return set()
    base = get_merge_base(base_ref)
    verified = set()
    for name in candidates:
        before, after = (
            subprocess.run(['git', 'show', f'{ref}:{name}'], check=True, capture_output=True, text=True).stdout
            for ref in (base, head)
        )
        if registry_only_additions(name, before, after):
            verified.add(name)
    return verified


def verified_declaration_only_files(added, head):
    """Use Go's parser to prove that missing files have no functions or variable initializers."""
    sources = {
        name: subprocess.run(
            ["git", "show", f"{head}:{name}"], check=True, capture_output=True, text=True
        ).stdout
        for name in added
    }
    if not sources:
        return set()
    helper = Path(__file__).resolve().parents[1] / "tools/declarations/main.go"
    result = subprocess.run(
        ["go", "run", str(helper)], input=json.dumps(sources), check=True,
        capture_output=True, text=True,
    )
    verified = json.loads(result.stdout)
    if not isinstance(verified, list) or any(not isinstance(name, str) or name not in sources for name in verified):
        raise ValueError("Invalid declaration-only file proof")
    return set(verified)
