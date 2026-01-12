"""Go struct/interface parsing from diffs."""

from pathlib import Path
from typing import Dict, List, Optional, Tuple, TYPE_CHECKING

from regex_patterns import RegexPatterns  # type: ignore
from models import (  # type: ignore
    ModelChange,
    FieldChange,
    MethodChange,
    EmbeddedTypeChange,
    ParseStatistics,
)

if TYPE_CHECKING:
    from core.progress_reporter import ProgressReporter


class StructParser:
    """Parses Go model changes from diff text (structs, interfaces, embedded types)."""

    def __init__(self, reporter: 'ProgressReporter'):
        """Initialize parser.
        
        Args:
            reporter: Progress reporter
        """
        self.reporter = reporter
        self.stats = ParseStatistics()
        self.in_interface_context = False  # Track if we're parsing inside an interface
    
    @property
    def statistics(self) -> ParseStatistics:
        """Get parsing statistics."""
        return self.stats

    def parse_diff(self, diff_text: str) -> List[ModelChange]:
        """Parse Go model changes from diff text.
        
        Args:
            diff_text: Unified diff text
            
        Returns:
            List of ModelChange objects
        """
        model_changes: Dict[str, ModelChange] = {}
        current_file = None
        current_model = None
        self.stats = ParseStatistics()  # Reset stats for new parse
        self.in_interface_context = False

        lines = diff_text.split('\n')
        self.stats.total_lines_processed = len(lines)

        for i, line in enumerate(lines):
            # Check for file header
            if self._is_file_header(line):
                filename = self._extract_filename(line)
                if filename and (line.startswith('diff --git') or 
                               (line.startswith('+++') and current_file != filename)):
                    current_file = filename
                    current_model = self._filename_to_model_name(filename)
                    self.in_interface_context = False  # Reset for new file
                    
                    if current_file not in model_changes:
                        model_changes[current_file] = ModelChange(
                            file_path=current_file,
                            model_name=current_model
                        )
                        self.stats.total_files_in_diff += 1
                continue

            if not current_file or not line.strip():
                continue

            # Track added/removed lines
            if line.startswith('+') and not line.startswith('+++'):
                self.stats.added_lines_processed += 1
            elif line.startswith('-') and not line.startswith('---'):
                self.stats.removed_lines_processed += 1

            # Detect context switches (struct vs interface)
            self._update_context(line)

            # Parse changes based on context
            self._parse_line_change(line, i, current_file, model_changes)

        # Calculate final statistics
        result = [change for change in model_changes.values() if change.has_changes]
        files_without_changes = [change for change in model_changes.values() if not change.has_changes]
        
        # Update statistics
        self.stats.files_with_changes = len(result)
        self.stats.files_without_changes = len(files_without_changes)
        
        for change in result:
            self.stats.struct_fields_added += len(change.added_fields)
            self.stats.struct_fields_removed += len(change.removed_fields)
            self.stats.interface_methods_added += len(change.added_methods)
            self.stats.interface_methods_removed += len(change.removed_methods)
            self.stats.embedded_types_added += len(change.added_embedded_types)
            self.stats.embedded_types_removed += len(change.removed_embedded_types)
        
        self.reporter.print_parse_summary(result, self.stats, files_without_changes)
        
        return result

    def _update_context(self, line: str):
        """Update parsing context based on type declarations.
        
        Args:
            line: Current line being processed
        """
        cleaned = line.lstrip('+-').strip()
        
        # Check for interface declaration
        if RegexPatterns.GO_TYPE_INTERFACE.search(cleaned):
            self.in_interface_context = True
        # Check for struct declaration
        elif RegexPatterns.GO_TYPE_STRUCT.search(cleaned):
            self.in_interface_context = False
        # Closing brace resets context
        elif cleaned == '}':
            self.in_interface_context = False

    def _parse_line_change(self, line: str, line_number: int, current_file: str,
                          model_changes: Dict[str, ModelChange]):
        """Parse a line for any type of change (field, method, embedded type).
        
        Args:
            line: Line from diff
            line_number: Line number in diff
            current_file: Current file being processed
            model_changes: Dictionary of model changes being built
        """
        if line.startswith('+') and not line.startswith('+++'):
            self._parse_added_line(line[1:].strip(), line_number, current_file, model_changes)
        elif line.startswith('-') and not line.startswith('---'):
            self._parse_removed_line(line[1:].strip(), line_number, current_file, model_changes)

    def _parse_added_line(self, line: str, line_number: int, current_file: str,
                         model_changes: Dict[str, ModelChange]):
        """Parse an added line (+).
        
        Args:
            line: Cleaned line content
            line_number: Line number in diff
            current_file: Current file
            model_changes: Model changes dictionary
        """
        # Try embedded type first (works for both interface and struct)
        embedded_info = self._parse_embedded_type(line)
        if embedded_info:
            context = 'interface' if self.in_interface_context else 'struct'
            embedded_change = EmbeddedTypeChange(
                type_name=embedded_info,
                change_type='added',
                context=context,
                line_number=line_number
            )
            model_changes[current_file].added_embedded_types.append(embedded_change)
            self.stats.embedded_types_added += 1
            return

        # If in interface, try to parse as method
        if self.in_interface_context:
            method_info = self._parse_interface_method(line)
            if method_info:
                method_change = MethodChange(
                    method_name=method_info[0],
                    parameters=method_info[1],
                    return_type=method_info[2],
                    change_type='added',
                    line_number=line_number
                )
                model_changes[current_file].added_methods.append(method_change)
                self.stats.interface_methods_added += 1
                return

        # Otherwise, try to parse as struct field
        field_info = self._parse_field_line(line)
        if field_info:
            field_change = FieldChange(
                field_name=field_info[0],
                field_type=field_info[1],
                change_type='added',
                line_number=line_number
            )
            model_changes[current_file].added_fields.append(field_change)
            self.stats.struct_fields_added += 1

    def _parse_removed_line(self, line: str, line_number: int, current_file: str,
                           model_changes: Dict[str, ModelChange]):
        """Parse a removed line (-).
        
        Args:
            line: Cleaned line content
            line_number: Line number in diff
            current_file: Current file
            model_changes: Model changes dictionary
        """
        # Try embedded type first
        embedded_info = self._parse_embedded_type(line)
        if embedded_info:
            context = 'interface' if self.in_interface_context else 'struct'
            embedded_change = EmbeddedTypeChange(
                type_name=embedded_info,
                change_type='removed',
                context=context,
                line_number=line_number
            )
            model_changes[current_file].removed_embedded_types.append(embedded_change)
            self.stats.embedded_types_removed += 1
            return

        # If in interface, try to parse as method
        if self.in_interface_context:
            method_info = self._parse_interface_method(line)
            if method_info:
                method_change = MethodChange(
                    method_name=method_info[0],
                    parameters=method_info[1],
                    return_type=method_info[2],
                    change_type='removed',
                    line_number=line_number
                )
                model_changes[current_file].removed_methods.append(method_change)
                self.stats.interface_methods_removed += 1
                return

        # Otherwise, try to parse as struct field
        field_info = self._parse_field_line(line)
        if field_info:
            field_change = FieldChange(
                field_name=field_info[0],
                field_type=field_info[1],
                change_type='removed',
                line_number=line_number
            )
            model_changes[current_file].removed_fields.append(field_change)
            self.stats.struct_fields_removed += 1

    def _is_file_header(self, line: str) -> bool:
        """Check if line is a file header."""
        return line.startswith('diff --git') or line.startswith('+++') or line.startswith('---')

    def _extract_filename(self, line: str) -> Optional[str]:
        """Extract filename from diff header line."""
        match = RegexPatterns.MODEL_FILE_PATH.search(line)
        return match.group(0) if match else None

    def _filename_to_model_name(self, filename: str) -> str:
        """Convert filename to model name (snake_case to PascalCase)."""
        file_stem = Path(filename).stem
        return ''.join(word.capitalize() for word in file_stem.split('_'))

    def _parse_interface_method(self, line: str) -> Optional[Tuple[str, str, str]]:
        """Parse an interface method declaration.
        
        Args:
            line: Line of Go code
            
        Returns:
            Tuple of (method_name, parameters, return_type) or None
        """
        if not line or line.startswith('//') or line.startswith('}') or line.startswith('{'):
            return None
        
        # Skip function implementations (have body indicators)
        if line.startswith('func (') and '{' in line:
            self.stats.lines_filtered_func_impl += 1
            return None

        match = RegexPatterns.GO_INTERFACE_METHOD.match(line)
        if match:
            method_name = match.group(1)
            parameters = match.group(2) if match.group(2) else ""
            return_type = match.group(3) if match.group(3) else ""
            
            # Only track exported methods (uppercase first letter)
            if method_name and method_name[0].isupper():
                return (method_name, parameters.strip(), return_type.strip().strip('()'))
            else:
                self.stats.lines_filtered_unexported += 1
        
        return None

    def _parse_embedded_type(self, line: str) -> Optional[str]:
        """Parse an embedded type (interface or struct).
        
        Args:
            line: Line of Go code
            
        Returns:
            Type name or None
        """
        if not line or line.startswith('//') or line.startswith('}') or line.startswith('{'):
            return None
        
        # Skip type declarations and function implementations
        if line.startswith('type ') or line.startswith('func '):
            return None

        match = RegexPatterns.GO_EMBEDDED_TYPE.match(line)
        if match:
            type_name = match.group(1)
            # Check if it looks like a type name (not a field with type)
            # Embedded types are just the type name, no field name before it
            if type_name and (type_name[0].isupper() or '.' in type_name):
                return type_name
        
        return None

    def _parse_field_line(self, line: str) -> Optional[Tuple[str, str]]:
        """Parse a Go struct field line and track filtering reasons.
        
        Args:
            line: Line of Go code
            
        Returns:
            Tuple of (field_name, field_type) or None
        """
        if not line:
            return None
            
        # Track why lines are filtered
        if line.startswith('//'):
            self.stats.lines_filtered_comments += 1
            return None
        
        if line.startswith('type ') or line.startswith('package ') or line.startswith('import '):
            self.stats.lines_filtered_declarations += 1
            return None
        
        if line.startswith('func '):
            self.stats.lines_filtered_func_impl += 1
            return None
            
        if line.startswith('}') or line.startswith('{'):
            return None

        match = RegexPatterns.GO_STRUCT_FIELD.match(line)
        if match:
            field_name = match.group(1)
            field_type = match.group(2)

            if field_name[0].isupper():  # Go exported field
                return (field_name, field_type)
            else:
                # Unexported field (starts with lowercase)
                self.stats.lines_filtered_unexported += 1
                return None
        else:
            # Line didn't match the field pattern
            self.stats.lines_filtered_no_match += 1

        return None
