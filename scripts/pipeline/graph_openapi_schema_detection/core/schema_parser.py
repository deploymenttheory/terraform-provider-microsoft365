"""OpenAPI schema parsing."""

import yaml
import time
from typing import Dict, Any, TYPE_CHECKING

if TYPE_CHECKING:
    from .progress_reporter import ProgressReporter


class SchemaParser:
    """Parses OpenAPI schemas efficiently."""
    
    def __init__(self, reporter: 'ProgressReporter'):
        """Initialize schema parser.
        
        Args:
            reporter: Progress reporter
        """
        self.reporter = reporter
    
    def extract_schemas_section(self, spec_content: str) -> str:
        """Extract just 'components.schemas' section (~10MB vs 60MB).
        
        Args:
            spec_content: Full OpenAPI spec content
            
        Returns:
            Just the schemas section as YAML string
        """
        self.reporter.info("📄 Extracting schemas section...")
        
        lines = spec_content.splitlines()
        schemas_start = None
        schemas_end = None
        
        # Find the components.schemas section
        for i, line in enumerate(lines):
            if '  schemas:' in line and schemas_start is None:
                schemas_start = i
                self.reporter.info(f"   Found schemas at line {i:,}")
            elif schemas_start and line.startswith('  ') and not line.startswith('    ') and line.strip() and schemas_end is None:
                # Found next top-level component section
                if any(keyword in line for keyword in ['responses:', 'parameters:', 'securitySchemes:', 'requestBodies:', 'headers:']):
                    schemas_end = i
                    self.reporter.info(f"   Schemas end at line {i:,}")
                    break
        
        if not schemas_start:
            raise ValueError("Could not find schemas section in OpenAPI spec")
        
        if not schemas_end:
            schemas_end = len(lines)
        
        # Extract just the schemas
        schemas_lines = ['schemas:'] + lines[schemas_start+1:schemas_end]
        schemas_content = '\n'.join(schemas_lines)
        
        self.reporter.info(f"   Extracted {len(schemas_lines):,} lines")
        return schemas_content
    
    def parse_schemas(self, schemas_content: str) -> Dict[str, Any]:
        """Parse YAML schemas section.
        
        Args:
            schemas_content: YAML content containing schemas
            
        Returns:
            Dictionary of schemas
        """
        self.reporter.info("🔍 Parsing schemas YAML...")
        start = time.time()
        
        try:
            parsed = yaml.safe_load(schemas_content)
            elapsed = time.time() - start
            
            if not isinstance(parsed, dict) or 'schemas' not in parsed:
                raise ValueError("Invalid schemas structure")
            
            schemas = parsed['schemas']
            self.reporter.info(f"   Parsed {len(schemas):,} schemas in {elapsed:.2f}s")
            
            return schemas
            
        except yaml.YAMLError as e:
            self.reporter.error(f"YAML parsing failed: {e}")
            raise
    
    def extract_model_properties(self, schema: Dict[str, Any]) -> Dict[str, Any]:
        """Extract properties with types, required status, nullable.
        
        Args:
            schema: Schema definition
            
        Returns:
            Dictionary with property details
        """
        properties = schema.get('properties', {})
        required = schema.get('required', [])
        all_of = schema.get('allOf', [])
        
        # Handle inheritance (allOf)
        inherited_props = {}
        if all_of:
            for item in all_of:
                if 'properties' in item:
                    inherited_props.update(item['properties'])
        
        # Merge inherited and direct properties
        all_props = {**inherited_props, **properties}
        
        # Build detailed property info with rich metadata
        property_details = {}
        for prop_name, prop_def in all_props.items():
            prop_type = prop_def.get('type')
            if not prop_type:
                # Check for $ref
                if '$ref' in prop_def:
                    prop_type = prop_def['$ref'].split('/')[-1]
                else:
                    prop_type = 'unknown'
            
            property_details[prop_name] = {
                # Core fields
                'type': prop_type,
                'required': prop_name in required,
                'nullable': prop_def.get('nullable', False),
                'description': prop_def.get('description', ''),
                
                # Validation metadata
                'enum': prop_def.get('enum'),
                'format': prop_def.get('format'),
                'pattern': prop_def.get('pattern'),
                'minLength': prop_def.get('minLength'),
                'maxLength': prop_def.get('maxLength'),
                'minimum': prop_def.get('minimum'),
                'maximum': prop_def.get('maximum'),
                
                # Other metadata
                'default': prop_def.get('default'),
                'example': prop_def.get('example'),
                'deprecated': prop_def.get('deprecated', False),
                'readOnly': prop_def.get('readOnly', False),
                'writeOnly': prop_def.get('writeOnly', False),
            }
        
        return property_details
