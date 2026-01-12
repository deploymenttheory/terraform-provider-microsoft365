"""Progress reporting for schema detection."""

import sys
from typing import List, TYPE_CHECKING

if TYPE_CHECKING:
    from models import ModelChange, ParseStatistics


class ProgressReporter:
    """Handles all user-facing output and progress reporting."""

    def __init__(self, verbose: bool = True):
        """Initialize the reporter.
        
        Args:
            verbose: If True, show detailed progress messages
        """
        self.verbose = verbose

    def section(self, message: str):
        """Print a section header."""
        if self.verbose:
            print(f"\n{message}")

    def info(self, message: str, indent: int = 0):
        """Print an info message."""
        if self.verbose:
            prefix = "  " * indent
            print(f"{prefix}{message}")

    def success(self, message: str):
        """Print a success message."""
        print(f"✅ {message}")

    def warning(self, message: str):
        """Print a warning message."""
        print(f"⚠️  {message}", file=sys.stderr)

    def error(self, message: str):
        """Print an error message."""
        print(f"❌ {message}", file=sys.stderr)

    def print_parse_summary(self, model_changes: List['ModelChange'], stats: 'ParseStatistics' = None,
                           files_without_changes: List['ModelChange'] = None):
        """Print summary of parsed model changes with statistics.
        
        Args:
            model_changes: List of model changes detected
            stats: Optional parsing statistics for diagnostics
            files_without_changes: Optional list of files that had no field changes
        """
        self.info("\n📋 Parse Summary:", indent=1)
        
        if model_changes:
            self.info(f"✓ Found {len(model_changes)} model(s) with changes:\n", indent=1)
            
            for change in model_changes:
                self.info(f"📄 {change.model_name} ({change.file_path})", indent=2)
                self.info(change.change_summary, indent=3)

                # Show struct fields
                if change.added_fields:
                    self.info("Added struct fields:", indent=3)
                    for fld in change.added_fields[:5]:
                        self.info(f"+ {fld.field_name}: {fld.field_type}", indent=4)
                    if len(change.added_fields) > 5:
                        self.info(f"... and {len(change.added_fields) - 5} more", indent=4)

                if change.removed_fields:
                    self.info("Removed struct fields:", indent=3)
                    for fld in change.removed_fields[:5]:
                        self.info(f"- {fld.field_name}: {fld.field_type}", indent=4)
                    if len(change.removed_fields) > 5:
                        self.info(f"... and {len(change.removed_fields) - 5} more", indent=4)

                # Show interface methods
                if change.added_methods:
                    self.info("Added interface methods:", indent=3)
                    for method in change.added_methods[:5]:
                        self.info(f"+ {method.signature}", indent=4)
                    if len(change.added_methods) > 5:
                        self.info(f"... and {len(change.added_methods) - 5} more", indent=4)

                if change.removed_methods:
                    self.info("Removed interface methods:", indent=3)
                    for method in change.removed_methods[:5]:
                        self.info(f"- {method.signature}", indent=4)
                    if len(change.removed_methods) > 5:
                        self.info(f"... and {len(change.removed_methods) - 5} more", indent=4)

                # Show embedded types
                if change.added_embedded_types:
                    self.info("Added embedded types:", indent=3)
                    for emb in change.added_embedded_types[:5]:
                        self.info(f"+ {emb.type_name} ({emb.context})", indent=4)
                    if len(change.added_embedded_types) > 5:
                        self.info(f"... and {len(change.added_embedded_types) - 5} more", indent=4)

                if change.removed_embedded_types:
                    self.info("Removed embedded types:", indent=3)
                    for emb in change.removed_embedded_types[:5]:
                        self.info(f"- {emb.type_name} ({emb.context})", indent=4)
                    if len(change.removed_embedded_types) > 5:
                        self.info(f"... and {len(change.removed_embedded_types) - 5} more", indent=4)

                print()
        else:
            self.info("ℹ️  No changes detected in diff", indent=1)
        
        # Print detailed statistics if provided
        if stats:
            print()
            print(stats.get_summary())
            
            # Explain why files were filtered
            if stats.files_without_changes > 0:
                print()
                self.info(f"ℹ️  {stats.files_without_changes} file(s) had changes but no detectable model modifications.", indent=1)
                self.info("   Possible reasons:", indent=1)
                self.info("   • Only comments, imports, or package declarations changed", indent=1)
                self.info("   • Method implementations (func body) changed", indent=1)
                self.info("   • Type aliases or constants changed", indent=1)
                self.info("   • Only unexported (lowercase) fields/methods changed", indent=1)
                self.info("   • Changes didn't match expected patterns", indent=1)
                
                # Show examples of files without field changes
                if files_without_changes:
                    print()
                    self.info("📁 Examples of files without field changes (showing up to 10):", indent=1)
                    for change in files_without_changes[:10]:
                        self.info(f"   • {change.model_name} ({change.file_path})", indent=1)
                    if len(files_without_changes) > 10:
                        self.info(f"   ... and {len(files_without_changes) - 10} more", indent=1)

    def print_dry_run_issue(self, title: str, body: str):
        """Print issue content for dry run."""
        print("\n🔍 DRY RUN: Would create issue with following content:")
        print("=" * 80)
        print(f"Title: {title}\n")
        print(body)
        print("=" * 80)
