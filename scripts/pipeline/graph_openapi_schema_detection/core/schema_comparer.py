"""Schema comparison for detecting OpenAPI changes."""

from typing import Dict, Any, List, Tuple, TYPE_CHECKING

if TYPE_CHECKING:
    from ..models import SchemaChange, PropertyChange, ParseStatistics
    from .progress_reporter import ProgressReporter
    from .schema_parser import SchemaParser


class SchemaComparer:
    """Compares two versions of OpenAPI schemas."""
    
    def __init__(self, parser: 'SchemaParser', reporter: 'ProgressReporter'):
        """Initialize schema comparer.
        
        Args:
            parser: Schema parser instance
            reporter: Progress reporter
        """
        self.parser = parser
        self.reporter = reporter
    
    def compare_schemas(
        self,
        old_schemas: Dict[str, Any],
        new_schemas: Dict[str, Any]
    ) -> Tuple[List['SchemaChange'], 'ParseStatistics']:
        """Compare all schemas and detect changes.
        
        Args:
            old_schemas: Previous version schemas
            new_schemas: New version schemas
            
        Returns:
            Tuple of (list of schema changes, statistics)
        """
        from models import SchemaChange, ParseStatistics  # type: ignore
        
        self.reporter.info("🔬 Comparing schemas...")
        
        changes = []
        stats = ParseStatistics()
        
        # Get all schema names
        old_names = set(old_schemas.keys())
        new_names = set(new_schemas.keys())
        
        stats.schemas_added = len(new_names - old_names)
        stats.schemas_removed = len(old_names - new_names)
        
        # Compare common schemas
        common_names = old_names & new_names
        stats.total_schemas_compared = len(common_names)
        
        for schema_name in sorted(common_names):
            # Only compare microsoft.graph.* schemas
            if not schema_name.startswith('microsoft.graph.'):
                continue
            
            schema_change = self.compare_model(
                schema_name,
                old_schemas[schema_name],
                new_schemas[schema_name],
                stats
            )
            
            if schema_change.has_changes:
                changes.append(schema_change)
                stats.schemas_with_changes += 1
        
        self.reporter.info(f"   Found {len(changes)} schema(s) with changes")
        
        return changes, stats
    
    def compare_model(
        self,
        model_name: str,
        old_schema: Dict[str, Any],
        new_schema: Dict[str, Any],
        stats: 'ParseStatistics'
    ) -> 'SchemaChange':
        """Compare single model between versions.
        
        Args:
            model_name: Schema name
            old_schema: Previous schema definition
            new_schema: New schema definition
            stats: Statistics object to update
            
        Returns:
            SchemaChange with all detected changes
        """
        from models import SchemaChange  # type: ignore
        
        # Extract properties from both versions
        old_data = self.parser.extract_model_properties(old_schema)
        new_data = self.parser.extract_model_properties(new_schema)
        
        old_props = old_data['properties']
        new_props = new_data['properties']
        parent_schemas = new_data.get('parent_schemas', [])
        schema_description = new_data.get('description', '')
        discriminator_value = new_data.get('discriminator_value')
        discriminator_property = new_data.get('discriminator_property')
        
        # Detect property changes
        added, removed, type_changed, required_changed, nullable_changed = self.detect_property_changes(
            old_props,
            new_props
        )
        
        # Update statistics
        stats.properties_added += len(added)
        stats.properties_removed += len(removed)
        stats.type_changes += len(type_changed)
        stats.required_changes += len(required_changed)
        stats.nullable_changes += len(nullable_changed)
        
        # Create file path for compatibility with provider filter
        # microsoft.graph.user → models/user.go
        simple_name = model_name.split('.')[-1]
        file_path = f"models/{simple_name}.go"
        
        return SchemaChange(
            schema_name=model_name,
            file_path=file_path,
            added_properties=added,
            removed_properties=removed,
            type_changed_properties=type_changed,
            required_changed_properties=required_changed,
            nullable_changed_properties=nullable_changed,
            parent_schemas=parent_schemas,
            description=schema_description,
            discriminator_value=discriminator_value,
            discriminator_property=discriminator_property
        )
    
    def detect_property_changes(
        self,
        old_props: Dict[str, Any],
        new_props: Dict[str, Any]
    ) -> Tuple[List['PropertyChange'], List['PropertyChange'], List['PropertyChange'], 
               List['PropertyChange'], List['PropertyChange']]:
        """Detect added, removed, type changes, required changes.
        
        Args:
            old_props: Previous properties
            new_props: New properties
            
        Returns:
            Tuple of (added, removed, type_changed, required_changed, nullable_changed)
        """
        from models import PropertyChange  # type: ignore
        
        old_names = set(old_props.keys())
        new_names = set(new_props.keys())
        
        # Added properties
        added = []
        for name in (new_names - old_names):
            prop = new_props[name]
            added.append(PropertyChange(
                property_name=name,
                change_type='added',
                new_type=prop['type'],
                new_required=prop['required'],
                new_nullable=prop.get('nullable'),
                description=prop.get('description'),
                enum_values=prop.get('enum'),
                format=prop.get('format'),
                pattern=prop.get('pattern'),
                min_length=prop.get('minLength'),
                max_length=prop.get('maxLength'),
                minimum=prop.get('minimum'),
                maximum=prop.get('maximum'),
                default=prop.get('default'),
                example=prop.get('example'),
                deprecated=prop.get('deprecated', False),
                read_only=prop.get('readOnly', False),
                write_only=prop.get('writeOnly', False),
                is_navigation_property=prop.get('is_navigation_property', False),
            ))
        
        # Removed properties
        removed = []
        for name in (old_names - new_names):
            old_prop = old_props[name]
            removed.append(PropertyChange(
                property_name=name,
                change_type='removed',
                old_type=old_prop['type'],
                old_required=old_prop['required'],
                old_nullable=old_prop.get('nullable'),
                enum_values=old_prop.get('enum'),  # Show enum values that were removed
                is_navigation_property=old_prop.get('is_navigation_property', False),
            ))
        
        # Check common properties for changes
        type_changed = []
        required_changed = []
        nullable_changed = []
        
        for name in (old_names & new_names):
            old_prop = old_props[name]
            new_prop = new_props[name]
            
            # Type changes
            if old_prop['type'] != new_prop['type']:
                type_changed.append(PropertyChange(
                    property_name=name,
                    change_type='type_changed',
                    old_type=old_prop['type'],
                    new_type=new_prop['type'],
                    old_required=old_prop['required'],
                    new_required=new_prop['required'],
                    old_description=old_prop.get('description'),
                    new_description=new_prop.get('description'),
                    enum_values=new_prop.get('enum'),
                    is_navigation_property=new_prop.get('is_navigation_property', False),
                ))
            
            # Required changes
            if old_prop['required'] != new_prop['required']:
                required_changed.append(PropertyChange(
                    property_name=name,
                    change_type='required_changed',
                    old_required=old_prop['required'],
                    new_required=new_prop['required'],
                    old_type=old_prop['type'],
                    new_type=new_prop['type'],
                    description=new_prop.get('description'),
                ))
            
            # Nullable changes
            old_nullable = old_prop.get('nullable', False)
            new_nullable = new_prop.get('nullable', False)
            if old_nullable != new_nullable:
                nullable_changed.append(PropertyChange(
                    property_name=name,
                    change_type='nullable_changed',
                    old_nullable=old_nullable,
                    new_nullable=new_nullable,
                    old_type=old_prop['type'],
                    new_type=new_prop['type'],
                    description=new_prop.get('description'),
                ))
            
            # Enum changes (useful for detecting new allowed values)
            old_enum = old_prop.get('enum')
            new_enum = new_prop.get('enum')
            if old_enum != new_enum and (old_enum is not None or new_enum is not None):
                # Consider enum changes as metadata changes (not breaking)
                # But track them for schema validation updates
                pass  # Could add enum_changed list if needed
        
        return added, removed, type_changed, required_changed, nullable_changed
