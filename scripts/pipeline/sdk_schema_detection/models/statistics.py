"""Statistics models for schema detection."""

from dataclasses import dataclass


@dataclass
class ParseStatistics:
    """Statistics from parsing diff for diagnostic purposes."""
    total_files_in_diff: int = 0
    files_with_changes: int = 0
    files_without_changes: int = 0
    total_lines_processed: int = 0
    added_lines_processed: int = 0
    removed_lines_processed: int = 0
    
    # Struct tracking
    struct_fields_added: int = 0
    struct_fields_removed: int = 0
    
    # Interface tracking
    interface_methods_added: int = 0
    interface_methods_removed: int = 0
    
    # Embedded types tracking
    embedded_types_added: int = 0
    embedded_types_removed: int = 0
    
    # Filtering reasons
    lines_filtered_comments: int = 0
    lines_filtered_declarations: int = 0
    lines_filtered_func_impl: int = 0
    lines_filtered_no_match: int = 0
    lines_filtered_unexported: int = 0
    
    def get_summary(self) -> str:
        """Get a human-readable summary."""
        parts = [
            "  📊 Parsing Statistics:",
            f"     Files in diff: {self.total_files_in_diff}",
            f"     Files with changes: {self.files_with_changes}",
            f"     Files without changes: {self.files_without_changes}",
            "",
            "  📝 Lines Processed:",
            f"     Total change lines: {self.added_lines_processed + self.removed_lines_processed}",
            f"     Added lines (+): {self.added_lines_processed}",
            f"     Removed lines (-): {self.removed_lines_processed}",
            "",
            "  🔧 Changes Detected:",
            f"     Struct fields added: {self.struct_fields_added}",
            f"     Struct fields removed: {self.struct_fields_removed}",
            f"     Interface methods added: {self.interface_methods_added}",
            f"     Interface methods removed: {self.interface_methods_removed}",
            f"     Embedded types added: {self.embedded_types_added}",
            f"     Embedded types removed: {self.embedded_types_removed}",
            "",
            "  🔍 Filtering Breakdown:",
            f"     Comments (//): {self.lines_filtered_comments}",
            f"     Type/package/import declarations: {self.lines_filtered_declarations}",
            f"     Function implementations (func): {self.lines_filtered_func_impl}",
            f"     Unexported fields (lowercase): {self.lines_filtered_unexported}",
            f"     No regex match: {self.lines_filtered_no_match}",
        ]
        return "\n".join(parts)
