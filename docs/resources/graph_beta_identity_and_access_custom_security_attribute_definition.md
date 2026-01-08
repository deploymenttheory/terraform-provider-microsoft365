---
page_title: "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft Entra custom security attribute definitions using the /directory/customSecurityAttributeDefinitions endpoint. Custom security attribute definitions define the structure and behavior of custom security attributes that can be assigned to users, groups, and other directory objects.
  Note: Custom security attribute definitions cannot be deleted once created. When removed from Terraform configuration, the resource will be deactivated by setting its status to 'Deprecated' and then removed from Terraform state. The attribute definition will remain in Microsoft Entra in a deprecated state.
---

# microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition (Resource)

Manages Microsoft Entra custom security attribute definitions using the `/directory/customSecurityAttributeDefinitions` endpoint. Custom security attribute definitions define the structure and behavior of custom security attributes that can be assigned to users, groups, and other directory objects.

**Note:** Custom security attribute definitions cannot be deleted once created. When removed from Terraform configuration, the resource will be deactivated by setting its status to 'Deprecated' and then removed from Terraform state. The attribute definition will remain in Microsoft Entra in a deprecated state.

## Microsoft Documentation

- [customSecurityAttributeDefinition resource type](https://learn.microsoft.com/en-us/graph/api/resources/customsecurityattributedefinition?view=graph-rest-beta)
- [Create customSecurityAttributeDefinition](https://learn.microsoft.com/en-us/graph/api/directory-post-customsecurityattributedefinitions?view=graph-rest-beta&tabs=http)
- [Update customSecurityAttributeDefinition](https://learn.microsoft.com/en-us/graph/api/customsecurityattributedefinition-update?view=graph-rest-beta&tabs=http)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `CustomSecAttributeDefinition.Read.All`, `CustomSecAttributeDefinition.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Simple String Attribute (Most Common Use Case)
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "project_name" {
  attribute_set               = "Engineering"
  name                        = "ProjectName"
  description                 = "Name of the project the user is assigned to"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

# Example 2: Integer Attribute
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "cost_center" {
  attribute_set               = microsoft365_graph_beta_identity_and_access_attribute_set.example.id
  name                        = "CostCenter"
  description                 = "Cost center number for budget tracking"
  type                        = "Integer"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false
}

# Example 3: Boolean Attribute
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "security_clearance" {
  attribute_set               = "Security"
  name                        = "HasClearance"
  description                 = "Indicates if user has security clearance"
  type                        = "Boolean"
  status                      = "Available"
  is_collection               = false # Must be false for Boolean type
  is_searchable               = true
  use_pre_defined_values_only = false # Must be false for Boolean type
}

# Example 4: Collection Attribute (Multiple Values)
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "skills" {
  attribute_set               = "HumanResources"
  name                        = "Skills"
  description                 = "Skills and competencies of the employee"
  type                        = "String"
  status                      = "Available"
  is_collection               = true # Allows multiple values
  is_searchable               = true
  use_pre_defined_values_only = false
}

# Example 5: Attribute with Predefined Values Only
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "department" {
  attribute_set               = "Organization"
  name                        = "Department"
  description                 = "Department assignment with predefined values"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = true # Only predefined values allowed

  # Note: After creating this definition, you would need to add allowed values
  # using the microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value resource
}

# Example 6: Multi-Collection Attribute for Team Assignments
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "team_assignments" {
  attribute_set               = "Engineering"
  name                        = "Teams"
  description                 = "Engineering teams the user belongs to"
  type                        = "String"
  status                      = "Available"
  is_collection               = true
  is_searchable               = true
  use_pre_defined_values_only = true # Restrict to predefined team names

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "15m"
  }
}

# Example 7: Deprecated Attribute (Inactive State)
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "legacy_code" {
  attribute_set               = "Legacy"
  name                        = "OldSystemCode"
  description                 = "Legacy system identifier (deprecated)"
  type                        = "String"
  status                      = "Deprecated" # Marks as inactive
  is_collection               = false
  is_searchable               = false
  use_pre_defined_values_only = false
}

# Example 8: Employee Classification with Limited Search
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "employee_classification" {
  attribute_set               = "HumanResources"
  name                        = "Classification"
  description                 = "Employee classification level"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = false # Not searchable for privacy
  use_pre_defined_values_only = true
}

# Example 9: Attribute for Compliance Tracking
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "compliance_training_date" {
  attribute_set               = "Compliance"
  name                        = "LastTrainingDate"
  description                 = "Date of last compliance training completion (format: YYYY-MM-DD)"
  type                        = "String"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

# Example 10: Office Location Tracking
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "office_locations" {
  attribute_set               = "Facilities"
  name                        = "Locations"
  description                 = "Office locations where the employee works"
  type                        = "String"
  status                      = "Available"
  is_collection               = true # Employee can work from multiple locations
  is_searchable               = true
  use_pre_defined_values_only = true

  # After creation, add specific office locations as allowed values
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `attribute_set` (String) Name of the attribute set. Case insensitive.
- `is_collection` (Boolean) Indicates whether multiple values can be assigned to the custom security attribute. Cannot be changed later. If type is set to Boolean, isCollection cannot be set to true.
- `is_searchable` (Boolean) Indicates whether custom security attribute values are indexed for searching on objects that are assigned attribute values. Cannot be changed later.
- `name` (String) Name of the custom security attribute. Must be unique within an attribute set. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case insensitive.
- `status` (String) Specifies whether the custom security attribute is active or deactivated. Acceptable values are: Available and Deprecated. Can be changed later.
- `type` (String) Data type for the custom security attribute values. Supported types are: Boolean, Integer, and String. Cannot be changed later.
- `use_pre_defined_values_only` (Boolean) Indicates whether only predefined values can be assigned to the custom security attribute. If set to false, free-form values are allowed. Can later be changed from true to false, but cannot be changed from false to true. If type is set to Boolean, usePreDefinedValuesOnly cannot be set to true.

### Optional

- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) Identifier of the custom security attribute, which is a combination of the attribute set name and the custom security attribute name separated by an underscore (attributeSet_name). The id property is auto generated and cannot be set. Case insensitive. Inherited from entity.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

You can define up to 500 active objects in a tenant. The customSecurityAttributeDefinition object can't be renamed or deleted, but it can be deactivated by using the Update customSecurityAttributeDefinition operation. Must be part of an attributeSet.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import ID format: {attribute_set}/{name}
terraform import microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.example "Engineering/ExampleCustomAttribute"
``` 