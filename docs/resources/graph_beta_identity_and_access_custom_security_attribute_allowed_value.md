---
page_title: "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft Entra custom security attribute allowed values using the /directory/customSecurityAttributeDefinitions/{customSecurityAttributeDefinitionId}/allowedValues endpoint. This resource is used to allowed values represent predefined values that can be assigned to custom security attributes.
  Note: You can define up to 100 allowed values per custom security attribute definition. Allowed values cannot be renamed or deleted, but they can be deactivated..
---

# microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value (Resource)

Manages Microsoft Entra custom security attribute allowed values using the `/directory/customSecurityAttributeDefinitions/{customSecurityAttributeDefinitionId}/allowedValues` endpoint. This resource is used to allowed values represent predefined values that can be assigned to custom security attributes.

**Note:** You can define up to 100 allowed values per custom security attribute definition. Allowed values cannot be renamed or deleted, but they can be deactivated..

## Microsoft Documentation

- [customSecurityAttributeAllowedValue resource type](https://learn.microsoft.com/en-us/graph/api/resources/allowedvalue?view=graph-rest-beta)
- [Create customSecurityAttributeAllowedValue](https://learn.microsoft.com/en-us/graph/api/customsecurityattributedefinition-post-allowedvalues?view=graph-rest-beta&tabs=http)
- [Update customSecurityAttributeAllowedValue](https://learn.microsoft.com/en-us/graph/api/allowedvalue-update?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `CustomSecAttributeDefinition.Read.All`
- `CustomSecAttributeDefinition.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.33.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example 1: Basic Allowed Value
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "engineering_project_alpha" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Alpha"
  is_active                               = true

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

# Example 2: Multiple Allowed Values for Department
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_engineering" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Engineering"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_sales" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Sales"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_marketing" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Marketing"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "dept_hr" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Human Resources"
  is_active                               = true
}

# Example 3: Deprecated/Inactive Allowed Value
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "legacy_dept" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.department.id
  id                                      = "Legacy Department"
  is_active                               = false # Marked as inactive
}

# Example 4: Office Locations with Allowed Values
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_seattle" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "Seattle"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_new_york" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "New York"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_london" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "London"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "location_tokyo" {
  custom_security_attribute_definition_id = microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition.office_locations.id
  id                                      = "Tokyo"
  is_active                               = true
}

# Example 5: Project Names with Allowed Values
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "project_apollo" {
  custom_security_attribute_definition_id = "Engineering_ProjectName"
  id                                      = "Apollo"
  is_active                               = true
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "project_orion" {
  custom_security_attribute_definition_id = "Engineering_ProjectName"
  id                                      = "Orion"
  is_active                               = true
}

# Example 6: Values with Special Characters
resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "classification_level_1" {
  custom_security_attribute_definition_id = "Security_Classification"
  id                                      = "Level-1"
  is_active                               = true

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "classification_level_2" {
  custom_security_attribute_definition_id = "Security_Classification"
  id                                      = "Level-2"
  is_active                               = true
}

# Note: You can define up to 100 allowed values per custom security attribute definition
# The id is the identifier for the predefined value (e.g., "Alpine", "Engineering")
# Allowed values cannot be deleted, only deactivated by setting is_active to false
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `custom_security_attribute_definition_id` (String) The ID of the custom security attribute definition that this allowed value belongs to. Format: 'attributeSet_attributeName'.
- `id` (String) Identifier for the predefined value. Can be up to 64 characters long and include Unicode characters. Can include spaces, but some special characters aren't allowed. Cannot be changed later. Case sensitive.
- `is_active` (Boolean) Indicates whether the predefined value is active or deactivated. If set to false, this predefined value cannot be assigned to any more supported directory objects. Can be changed later.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

You can define up to 100 allowed values per custom security attribute definition. Allowed values cannot be renamed or deleted, but they can be deactivated.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import ID format: {customSecurityAttributeDefinitionId}/{id}
# The format includes both the parent definition ID and the value identifier

# Example 1: Import an allowed value
terraform import microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value.engineering_project_alpha "Engineering_Project/Alpha"

# Note: The import ID format is {attributeSet}_{attributeName}/{id}
# For example: Engineering_Project/Alpine
``` 