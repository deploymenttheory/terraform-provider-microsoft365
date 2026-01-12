#!/usr/bin/env python3
"""
Orchestrates the complete schema change detection workflow:
1. Analyzes provider code to identify used SDK models
2. Detects SDK schema changes
3. Filters changes to only those affecting used models
4. Creates GitHub issues for relevant changes

This ensures only actionable schema changes are reported.
"""

import argparse
import subprocess
import sys
from pathlib import Path
from typing import Optional


class SchemaDetectionOrchestrator:
    """Orchestrates the three-step schema detection workflow."""
    
    def __init__(self, provider_path: Path, scripts_path: Path):
        self.provider_path = provider_path
        self.scripts_path = scripts_path
        self.usage_file = Path("provider_model_usage.json")
        self.results_file = Path("schema_changes.json")
    
    def run_step1_analyze_usage(self, verbose: bool = False) -> bool:
        """Step 1: Analyze provider code to identify used models.
        
        Args:
            verbose: Show detailed output
            
        Returns:
            True if step succeeded
        """
        print("\n" + "=" * 70)
        print("📋 STEP 1: Analyzing Provider Model Usage")
        print("=" * 70)
        
        script = self.scripts_path / "analyze_provider_model_usage.py"
        cmd = [
            sys.executable,
            str(script),
            "--provider-path", str(self.provider_path),
            "--output", str(self.usage_file),
        ]
        
        if verbose:
            cmd.append("--verbose")
        
        try:
            result = subprocess.run(cmd, check=True, capture_output=False, text=True)
            print(f"\n✅ Step 1 complete: {self.usage_file}")
            return True
        except subprocess.CalledProcessError as e:
            print(f"\n❌ Step 1 failed: {e}")
            return False
    
    def run_step2_detect_changes(self, pr_number: Optional[int] = None,
                                 current_version: Optional[str] = None,
                                 new_version: Optional[str] = None,
                                 dry_run: bool = False) -> bool:
        """Step 2: Detect SDK schema changes.
        
        Args:
            pr_number: PR number to analyze
            current_version: Current SDK version
            new_version: New SDK version
            dry_run: Don't create issues (for testing)
            
        Returns:
            True if step succeeded
        """
        print("\n" + "=" * 70)
        print("📋 STEP 2: Detecting SDK Schema Changes")
        print("=" * 70)
        
        script = self.scripts_path / "kiota_graph_sdk_schema_change_detector.py"
        cmd = [
            sys.executable,
            str(script),
            "--save-results", str(self.results_file),
            "--filter-by-usage", str(self.usage_file),
        ]
        
        if pr_number:
            cmd.extend(["--pr-number", str(pr_number)])
        elif current_version and new_version:
            cmd.extend(["--current", current_version, "--new", new_version])
        else:
            print("❌ Must provide either --pr-number or --current and --new")
            return False
        
        if dry_run:
            cmd.append("--dry-run")
        
        try:
            result = subprocess.run(cmd, check=True, capture_output=False, text=True)
            print(f"\n✅ Step 2 complete: {self.results_file}")
            return True
        except subprocess.CalledProcessError as e:
            print(f"\n❌ Step 2 failed: {e}")
            return False
    
    def run_complete_workflow(self, pr_number: Optional[int] = None,
                            current_version: Optional[str] = None,
                            new_version: Optional[str] = None,
                            dry_run: bool = False,
                            verbose: bool = False,
                            skip_usage_analysis: bool = False) -> bool:
        """Run the complete three-step workflow.
        
        Args:
            pr_number: PR number to analyze
            current_version: Current SDK version
            new_version: New SDK version
            dry_run: Don't create issues (for testing)
            verbose: Show detailed output
            skip_usage_analysis: Skip step 1 (use existing usage file)
            
        Returns:
            True if workflow succeeded
        """
        print("\n" + "🔄" * 35)
        print("🚀 SCHEMA DETECTION ORCHESTRATION WORKFLOW")
        print("🔄" * 35)
        
        # Step 1: Analyze provider usage (unless skipped)
        if not skip_usage_analysis:
            if not self.run_step1_analyze_usage(verbose):
                return False
        else:
            print(f"\n⏭️  Skipping Step 1 (using existing {self.usage_file})")
            if not self.usage_file.exists():
                print(f"❌ Usage file not found: {self.usage_file}")
                print("   Run without --skip-usage-analysis first")
                return False
        
        # Step 2: Detect changes and create filtered issues
        if not self.run_step2_detect_changes(pr_number, current_version, 
                                            new_version, dry_run):
            return False
        
        # Success!
        print("\n" + "🎉" * 35)
        print("✅ WORKFLOW COMPLETE!")
        print("🎉" * 35)
        print(f"\n📊 Results:")
        print(f"   • Model usage: {self.usage_file}")
        print(f"   • Schema changes: {self.results_file}")
        
        if not dry_run:
            print(f"   • GitHub issue created (see above for details)")
        else:
            print(f"   • Dry run mode: No issue created")
        
        return True


def main() -> None:
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Orchestrate schema change detection workflow",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    
    # Analysis options
    parser.add_argument(
        "--pr-number",
        type=int,
        help="Dependabot PR number to analyze"
    )
    parser.add_argument(
        "--current",
        type=str,
        help="Current SDK version (auto-detected if not provided)"
    )
    parser.add_argument(
        "--new",
        type=str,
        help="New SDK version (required if --pr-number not provided)"
    )
    
    # Workflow options
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Analyze without creating GitHub issues"
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Show detailed output"
    )
    parser.add_argument(
        "--skip-usage-analysis",
        action="store_true",
        help="Skip provider usage analysis (use existing results)"
    )
    
    # Path options
    parser.add_argument(
        "--provider-path",
        type=Path,
        default=Path(__file__).parent.parent.parent.parent,
        help="Path to provider repository"
    )
    
    args = parser.parse_args()
    
    # Validate arguments
    if not args.pr_number and not (args.current and args.new):
        parser.error("Either --pr-number or both --current and --new must be provided")
    
    # Setup paths
    scripts_path = Path(__file__).parent
    
    # Create orchestrator
    orchestrator = SchemaDetectionOrchestrator(
        provider_path=args.provider_path,
        scripts_path=scripts_path
    )
    
    # Run workflow
    try:
        success = orchestrator.run_complete_workflow(
            pr_number=args.pr_number,
            current_version=args.current,
            new_version=args.new,
            dry_run=args.dry_run,
            verbose=args.verbose,
            skip_usage_analysis=args.skip_usage_analysis
        )
        sys.exit(0 if success else 1)
        
    except KeyboardInterrupt:
        print("\n❌ Interrupted by user", file=sys.stderr)
        sys.exit(130)
    except Exception as e:
        print(f"\n❌ Workflow failed: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
