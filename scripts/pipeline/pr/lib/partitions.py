"""Partition Go jobs using the acceptance pipeline's bounded Rendezvous sharder."""

import ast
import re
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[2] / "bi_weekly_tests"))
sys.path.insert(0, str(Path(__file__).resolve().parents[2]))
# Direct script execution requires sibling search paths before these imports.
# pylint: disable=wrong-import-position
from generate_test_factories import without_comments
from generate_test_matrix import RENDEZVOUS_SEED, calculate_shard_count, distribute_rendezvous

# pylint: enable=wrong-import-position

MODULE = "github.com/deploymenttheory/terraform-provider-microsoft365"


def _go_import_paths(source):
    """Read Go import literals after removing comments without joining adjacent tokens."""
    literal = r'"(?:\\.|[^"\\\n])*"|`[^`]*`'
    block = rf'\((?:{literal}|[^()"`])*\)'
    declarations = re.finditer(
        rf"(?:{literal})|\bimport\s*(?P<imports>{block}|(?:(?:\w+|\.)\s*)?(?:{literal}))",
        without_comments(source),
    )
    for match in declarations:
        declaration = match.group("imports")
        if declaration is None:
            continue
        for value in re.findall(literal, declaration):
            try:
                yield value[1:-1].replace("\r", "") if value[0] == "`" else ast.literal_eval(value)
            except (ValueError, SyntaxError) as error:
                raise ValueError("Cannot verify a Go import literal") from error


def _local_import_graph(root):
    """Read production and test imports without executing or compiling the provider."""
    imports = {}
    test_imports = {}
    tests = set()
    source_files = list(root.glob("*.go")) + list((root / "internal").rglob("*.go"))
    for directory in sorted({path.parent for path in source_files}):
        package = directory.relative_to(root).as_posix()
        sources = list(directory.glob("*.go"))
        if any(path.name.endswith("_test.go") for path in sources):
            tests.add(package)
        imports[package] = set()
        test_imports[package] = set()
        for path in sources:
            source = path.read_text(encoding="utf-8")
            target = test_imports if path.name.endswith("_test.go") else imports
            target[package].update(
                "." if name == MODULE else name.removeprefix(MODULE + "/")
                for name in _go_import_paths(source)
                if name == MODULE or name.startswith(MODULE + "/")
            )
    return imports, test_imports, tests


def production_import_closure(root=Path(".")):
    """Prove which local packages are reachable from the provider's runtime entry points."""
    root = root.resolve()
    module = (root / "go.mod").read_text(encoding="utf-8")
    if not re.search(r"(?m)^module\s+" + re.escape(MODULE) + r"\s*$", module):
        raise ValueError("Cannot verify the provider module identity")
    if not (root / "main.go").is_file():
        raise ValueError("Provider main entry point is missing")
    if not any(
        not path.name.endswith("_test.go") for path in (root / "internal/provider").glob("*.go")
    ):
        raise ValueError("Provider runtime registry is missing")
    for target in re.findall(r'=>\s*("[^"\n]+"|\S+)', without_comments(module)):
        target = target.strip('"')
        if (target.startswith((".", "/", "\\")) or Path(target).is_absolute()) and (
            root / target
        ).resolve().is_relative_to(root):
            raise ValueError("Cannot verify local module replacements inside the provider tree")
    imports, _, _ = _local_import_graph(root)
    pending = [".", "internal/provider"]
    reachable = set()
    while pending:
        package = pending.pop()
        if package in reachable:
            continue
        if package not in imports:
            raise ValueError(f"Missing runtime package import evidence: {package}")
        reachable.add(package)
        pending.extend(imports[package] - reachable)
    return reachable


def affected_test_packages(changed, root=Path("."), production_changed=None):
    """Include tests whose local import graph depends on a changed package."""
    imports, test_imports, tests = _local_import_graph(root)
    affected = set(changed if production_changed is None else production_changed)
    while True:
        expanded = affected | {package for package, deps in imports.items() if deps & affected}
        if expanded == affected:
            break
        affected = expanded
    dependent_tests = {package for package in tests if test_imports[package] & affected}
    return sorted(((set(changed) | affected) & tests) | dependent_tests)


def fixture_packages(files):
    """Resolve changed internal fixtures and embedded data to their nearest Go package."""
    owners = set()
    for name in files:
        path = Path(name)
        if path.suffix == ".go" or not path.parts or path.parts[0] != "internal":
            continue
        for directory in path.parents:
            if not directory.parts or directory.parts[0] != "internal":
                break
            if any(directory.glob("*.go")):
                owners.add(directory.as_posix())
                break
    return owners


def partition_packages(packages):
    """Assign every selected package exactly once to a bounded standard-runner job."""
    packages = sorted(set(packages))
    if not packages:
        return {"include": []}
    weights = dict.fromkeys(packages, 1)
    count = calculate_shard_count(len(packages))
    while count <= len(packages):
        try:
            shards = distribute_rendezvous(packages, count, RENDEZVOUS_SEED, weights)
            break
        except ValueError:
            count += 1
    occupied = [items for items in shards.values() if items]
    if len(occupied) > 256:
        raise ValueError("Go job matrix exceeds the GitHub Actions 256-job limit")
    return {
        "include": [
            {"shard_index": index, "packages": " ".join(items)}
            for index, items in enumerate(occupied)
        ]
    }


def partition_lint_packages(packages):
    """Keep the full production registry on a dedicated standard-runner worker."""
    packages = set(packages)
    registry = "internal/provider"
    matrix = partition_packages(packages - {registry})
    for entry in matrix["include"]:
        entry.update({"label": f"Partition {entry['shard_index']}", "gomemlimit": "", "gogc": ""})
    if registry in packages:
        matrix["include"].append(
            {
                "shard_index": len(matrix["include"]),
                "packages": registry,
                "label": "Provider Registry",
                "gomemlimit": "20GiB",
                "gogc": "100",
            }
        )
    if len(matrix["include"]) > 256:
        raise ValueError("Go job matrix exceeds the GitHub Actions 256-job limit")
    return matrix


def select_lint_packages(changed, requested=None):
    """Allow manual diagnostic runs to narrow, but never expand, the changed package set."""
    if requested is None:
        return sorted(set(changed))
    selected = set(requested.split())
    if not selected or not selected.issubset(changed):
        raise ValueError("Lint selection must contain only packages changed against the baseline")
    return sorted(selected)
