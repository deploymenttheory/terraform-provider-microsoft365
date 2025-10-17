#!/bin/bash
# Import ID format: {customSecurityAttributeDefinitionId}/{id}
# The format includes both the parent definition ID and the value identifier

# Example 1: Import an allowed value
terraform import microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.engineering_project_alpha "Engineering_Project/Alpha"

# Note: The import ID format is {attributeSet}_{attributeName}/{id}
# For example: Engineering_Project/Alpine
