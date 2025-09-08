---
page_title: "microsoft365_graph_beta_device_management_group_policy_category Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  This data source retrieves comprehensive information about a specific Group Policy setting by performing three sequential Microsoft Graph API calls:
  GET /deviceManagement/groupPolicyCategories?$expand=parent,definitions&$select=id,displayName,isRoot,ingestionSource - Retrieves all categories with their definitionsGET /deviceManagement/groupPolicyDefinitions('{id}') - Gets detailed information about the specific policy definitionGET /deviceManagement/groupPolicyDefinitions('{id}')/presentations - Retrieves all presentation configurations for the policy
  The data source consolidates information from all three API calls into a single Terraform resource, making it easy to access category details, policy definitions, and presentation configurations (including dropdown options, text boxes, checkboxes, etc.) for a given Group Policy setting.
  Permissions
  One of the following permissions is required to call this API. To learn more, including how to choose permissions, see Permissions https://docs.microsoft.com/en-us/graph/permissions-reference.
  |Permission type|Permissions (from least to most privileged)|
  |:---|:---|
  |Delegated (work or school account)|DeviceManagementConfiguration.Read.All, DeviceManagementConfiguration.ReadWrite.All|
  |Delegated (personal Microsoft account)|Not supported.|
  |Application|DeviceManagementConfiguration.Read.All, DeviceManagementConfiguration.ReadWrite.All|
---

# microsoft365_graph_beta_device_management_group_policy_category (Data Source)

The Microsoft 365 Intune group policy category data source returns detailed information on every available group
policy available for configuration. This data source supports a search by group policy name and will return all matching
group policies, there settings and presentations. PResentations in this context, means the possible configuration choices
that are support for a given group policy setting.

## Microsoft Documentation

- [groupPolicyCategory resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-grouppolicy-grouppolicycategory?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.Read.All`

## Example Usage

```terraform
# Get the WSL networking configuration setting
data "microsoft365_graph_beta_device_management_group_policy_category" "wsl_networking" {
  setting_name = "Configure default networking mode" // Define the group policy item you wish to return

  timeouts = {
    read = "5m"
  }
}

# Output the complete data structure
output "wsl_networking_setting" {
  description = "Complete group policy setting information from all three API calls"
  value       = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking
}

# Access group policy top tier list 
output "category_info" {
  description = "Category information from the first API call"
  value = {
    id               = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.id
    display_name     = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.display_name
    is_root          = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.is_root
    ingestion_source = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.ingestion_source
    parent_category  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.parent
  }
}

# Access group policy definition by id (from 2nd API call)
output "definition_info" {
  description = "Detailed policy definition from the second API call"
  value = {
    id                       = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.id
    display_name             = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.display_name
    explain_text             = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.explain_text
    category_path            = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.category_path
    class_type               = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.class_type
    policy_type              = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.policy_type
    version                  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.version
    has_related_definitions  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.has_related_definitions
    group_policy_category_id = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.group_policy_category_id
    min_device_csp_version   = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.min_device_csp_version
    min_user_csp_version     = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.min_user_csp_version
    supported_on             = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.supported_on
    last_modified_date_time  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.last_modified_date_time
  }
}

# Access presentation information (from 3rd API call) - now properly populated!
output "presentation_info" {
  description = "Presentation configuration from the third API call - dropdown with options"
  value = {
    presentation_id         = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].id
    odata_type              = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].odata_type
    label                   = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].label
    required                = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].required
    last_modified_date_time = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].last_modified_date_time

    # Dropdown-specific properties.
    default_item      = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].default_item
    available_options = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].items
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `setting_name` (String) The display name of the Group Policy setting to search for (case-insensitive)

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `category` (Attributes) The Group Policy category information from the first API call (see [below for nested schema](#nestedatt--category))
- `definition` (Attributes) The detailed Group Policy definition information from the second API call (see [below for nested schema](#nestedatt--definition))
- `id` (String) The unique identifier for this data source
- `presentations` (Attributes List) The list of presentations associated with the group policy definition from the third API call (see [below for nested schema](#nestedatt--presentations))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--category"></a>
### Nested Schema for `category`

Read-Only:

- `display_name` (String) The display name of the category
- `id` (String) The unique identifier of the category
- `ingestion_source` (String) The source of the category (e.g., builtIn, custom)
- `is_root` (Boolean) Indicates if the category is a root category
- `parent` (Attributes) The parent category if this is not a root category (see [below for nested schema](#nestedatt--category--parent))

<a id="nestedatt--category--parent"></a>
### Nested Schema for `category.parent`

Read-Only:

- `display_name` (String) The display name of the parent category
- `id` (String) The unique identifier of the parent category
- `is_root` (Boolean) Indicates if the parent category is a root category



<a id="nestedatt--definition"></a>
### Nested Schema for `definition`

Read-Only:

- `category_path` (String) The category path of the definition
- `class_type` (String) The class type of the definition (e.g., machine, user)
- `display_name` (String) The display name of the definition
- `explain_text` (String) The explanation text for the definition
- `group_policy_category_id` (String) The ID of the group policy category this definition belongs to
- `has_related_definitions` (Boolean) Indicates if the definition has related definitions
- `id` (String) The unique identifier of the definition
- `last_modified_date_time` (String) The date and time the definition was last modified
- `min_device_csp_version` (String) The minimum device CSP version required
- `min_user_csp_version` (String) The minimum user CSP version required
- `policy_type` (String) The policy type of the definition
- `supported_on` (String) The supported platforms for the definition
- `version` (String) The version of the definition


<a id="nestedatt--presentations"></a>
### Nested Schema for `presentations`

Read-Only:

- `default_checked` (Boolean) Whether the checkbox is checked by default
- `default_decimal_value` (Number) The default decimal value for decimal text box presentations
- `default_item` (Attributes) The default item for dropdown list presentations (see [below for nested schema](#nestedatt--presentations--default_item))
- `default_value` (String) The default value for text box presentations
- `explicit_value` (Boolean) Whether the user must specify the registry subkey value and name for list box presentations
- `id` (String) The ID of the presentation
- `items` (Attributes List) The list of items for dropdown list presentations (see [below for nested schema](#nestedatt--presentations--items))
- `label` (String) The localized text label for the presentation
- `last_modified_date_time` (String) The date and time the entity was last modified
- `max_length` (Number) The maximum length for text box presentations
- `max_value` (Number) The maximum value for decimal text box presentations
- `min_value` (Number) The minimum value for decimal text box presentations
- `odata_type` (String) The OData type of the presentation (e.g., #microsoft.graph.groupPolicyPresentationDropdownList)
- `required` (Boolean) Whether a value is required for the parameter box (if applicable)
- `spin` (Boolean) Whether spin controls are enabled for decimal text box presentations
- `spin_step` (Number) The spin step for decimal text box presentations
- `value_prefix` (String) The value prefix for list box presentations

<a id="nestedatt--presentations--default_item"></a>
### Nested Schema for `presentations.default_item`

Read-Only:

- `display_name` (String) The display name of the default item
- `value` (String) The value of the default item


<a id="nestedatt--presentations--items"></a>
### Nested Schema for `presentations.items`

Read-Only:

- `display_name` (String) The display name of the item
- `value` (String) The value of the item