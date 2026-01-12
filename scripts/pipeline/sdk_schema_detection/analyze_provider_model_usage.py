#!/usr/bin/env python3
"""
Static code analysis to identify which Microsoft Graph SDK models
are actually used in the Terraform provider's resources and datasources.
"""

import argparse
import json
import re
from dataclasses import dataclass, asdict
from pathlib import Path
from typing import Dict, List


@dataclass
class ModelUsage:
    """Represents usage of a model in the provider."""
    model_name: str
    model_file: str
    used_in_files: List[str]
    usage_count: int


@dataclass
class AnalysisResult:
    """Complete analysis result."""
    total_files_scanned: int
    total_models_found: int
    models: List[ModelUsage]
    summary: Dict[str, int]


class RegexPatterns:
    """Compiled regex patterns for static analysis."""
    # Match main SDK model imports: graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
    SDK_IMPORT_MAIN = re.compile(
        r'graphmodels\s+"github\.com/microsoftgraph/msgraph-beta-sdk-go/models"',
        re.MULTILINE
    )
    
    # Match subpackage SDK model imports: models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/subpackage"
    SDK_IMPORT_SUBPACKAGE = re.compile(
        r'models\s+"github\.com/microsoftgraph/msgraph-beta-sdk-go/models/(security|networkaccess|externalconnectors|teamsadministration)"',
        re.MULTILINE
    )
    
    # Match model type references using graphmodels alias (e.g., graphmodels.User, graphmodels.Userable)
    MODEL_REFERENCE_MAIN = re.compile(
        r'\bgraphmodels\.([A-Z]\w+able?)\b'
    )
    
    # Match model type references using models alias for subpackages (e.g., models.Alert)
    MODEL_REFERENCE_SUBPACKAGE = re.compile(
        r'\bmodels\.([A-Z]\w+able?)\b'
    )
    
    # Match New constructors (e.g., graphmodels.NewDeviceCategory(), models.NewAlert())
    MODEL_CONSTRUCTOR = re.compile(
        r'\b(?:graphmodels|models)\.New([A-Z]\w+)\('
    )
    
    # Match model package in file path
    MODEL_FILE_PATH = re.compile(
        r'models/(?:(security|networkaccess|externalconnectors|teamsadministration)/)?(\w+)\.go$'
    )


class ProviderModelAnalyzer:
    """Analyzes provider code to identify used SDK models."""
    
    def __init__(self, provider_path: Path):
        self.provider_path = provider_path
        self.model_usage: Dict[str, ModelUsage] = {}
        
    def analyze(self) -> AnalysisResult:
        """Perform complete analysis of provider code."""
        print("🔍 Analyzing provider code for SDK model usage...")
        
        # Scan provider Go files
        go_files = self._find_go_files()
        print(f"  Found {len(go_files)} Go files to analyze")
        
        for go_file in go_files:
            self._analyze_file(go_file)
        
        # Build result
        models_list = sorted(
            self.model_usage.values(),
            key=lambda m: m.usage_count,
            reverse=True
        )
        
        summary = {
            "unique_models": len(self.model_usage),
            "total_references": sum(m.usage_count for m in models_list),
        }
        
        result = AnalysisResult(
            total_files_scanned=len(go_files),
            total_models_found=len(self.model_usage),
            models=models_list,
            summary=summary
        )
        
        print("\n✅ Analysis complete!")
        print(f"  Scanned: {result.total_files_scanned} files")
        print(f"  Found: {result.total_models_found} unique models")
        print(f"  Total references: {summary['total_references']}")
        
        return result
    
    def _find_go_files(self) -> List[Path]:
        """Find all Go files in resources, datasources, and related directories."""
        paths_to_scan = [
            self.provider_path / "internal" / "services" / "resources",
            self.provider_path / "internal" / "services" / "datasources",
            self.provider_path / "internal" / "services" / "ephemerals",
            self.provider_path / "internal" / "services" / "actions",
            self.provider_path / "internal" / "services" / "common",
            self.provider_path / "internal" / "client",
        ]
        
        go_files = []
        for path in paths_to_scan:
            if path.exists():
                go_files.extend(path.rglob("*.go"))
        
        return go_files
    
    def _analyze_file(self, file_path: Path) -> None:
        """Analyze a single Go file for model usage."""
        try:
            content = file_path.read_text(encoding='utf-8')
        except (OSError, UnicodeDecodeError) as e:
            print(f"  ⚠️  Failed to read {file_path}: {e}")
            return
        
        # Determine which import pattern is used
        has_main_import = RegexPatterns.SDK_IMPORT_MAIN.search(content)
        subpackage_match = RegexPatterns.SDK_IMPORT_SUBPACKAGE.search(content)
        
        # Track model references
        models_found = set()
        
        # Find graphmodels.* references (main package)
        if has_main_import:
            for match in RegexPatterns.MODEL_REFERENCE_MAIN.finditer(content):
                model_name = match.group(1)
                base_name = model_name[:-4] if model_name.endswith('able') else model_name
                model_file = f"models/{self._to_snake_case(base_name)}.go"
                models_found.add((base_name, model_file))
            
            # Also check for New* constructors
            for match in RegexPatterns.MODEL_CONSTRUCTOR.finditer(content):
                if 'graphmodels.New' in match.group(0):
                    model_name = match.group(1)
                    model_file = f"models/{self._to_snake_case(model_name)}.go"
                    models_found.add((model_name, model_file))
        
        # Find models.* references (subpackage)
        if subpackage_match:
            subpackage = subpackage_match.group(1)
            for match in RegexPatterns.MODEL_REFERENCE_SUBPACKAGE.finditer(content):
                model_name = match.group(1)
                base_name = model_name[:-4] if model_name.endswith('able') else model_name
                model_file = f"models/{subpackage}/{self._to_snake_case(base_name)}.go"
                models_found.add((base_name, model_file))
            
            # Also check for New* constructors
            for match in RegexPatterns.MODEL_CONSTRUCTOR.finditer(content):
                if 'models.New' in match.group(0):
                    model_name = match.group(1)
                    model_file = f"models/{subpackage}/{self._to_snake_case(model_name)}.go"
                    models_found.add((model_name, model_file))
        
        # Record all found models
        relative_path = str(file_path.relative_to(self.provider_path))
        for model_name, model_file in models_found:
            key = model_file
            if key not in self.model_usage:
                self.model_usage[key] = ModelUsage(
                    model_name=model_name,
                    model_file=model_file,
                    used_in_files=[],
                    usage_count=0
                )
            
            # Add file reference if not already present
            if relative_path not in self.model_usage[key].used_in_files:
                self.model_usage[key].used_in_files.append(relative_path)
            
            self.model_usage[key].usage_count += 1
    
    @staticmethod
    def _to_snake_case(name: str) -> str:
        """Convert PascalCase to snake_case."""
        # Insert underscore before uppercase letters (except first)
        s1 = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
        return re.sub('([a-z0-9])([A-Z])', r'\1_\2', s1).lower()


def save_results(result: AnalysisResult, output_file: Path) -> None:
    """Save analysis results to JSON file."""
    data = asdict(result)
    
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2)
    
    print(f"\n💾 Results saved to: {output_file}")


def print_summary(result: AnalysisResult, verbose: bool = False) -> None:
    """Print analysis summary."""
    print("\n" + "=" * 60)
    print("📊 PROVIDER MODEL USAGE SUMMARY")
    print("=" * 60)
    print(f"Files scanned: {result.total_files_scanned}")
    print(f"Unique models used: {result.total_models_found}")
    print(f"Total references: {result.summary['total_references']}")
    
    if verbose and result.models:
        print("\n📋 Top 20 Most Used Models:")
        for model in result.models[:20]:
            print(f"  • {model.model_name} ({model.usage_count} refs)")
            print(f"    File: {model.model_file}")
            if model.used_in_files:
                print(f"    Used in: {len(model.used_in_files)} file(s)")


def main() -> None:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Analyze provider code to identify used SDK models"
    )
    parser.add_argument(
        "--provider-path",
        type=Path,
        default=Path(__file__).parent.parent.parent,
        help="Path to provider repository"
    )
    parser.add_argument(
        "--output",
        type=Path,
        default=Path("provider_model_usage.json"),
        help="Output JSON file"
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Print detailed summary"
    )
    
    args = parser.parse_args()
    
    try:
        # Validate provider path
        if not args.provider_path.exists():
            raise ValueError(f"Provider path does not exist: {args.provider_path}")
        
        # Run analysis
        analyzer = ProviderModelAnalyzer(args.provider_path)
        result = analyzer.analyze()
        
        # Save results
        save_results(result, args.output)
        
        # Print summary
        print_summary(result, verbose=args.verbose)
        
    except (ValueError, OSError) as e:
        print(f"\n❌ Error: {e}")
        raise SystemExit(1) from e


if __name__ == "__main__":
    main()
