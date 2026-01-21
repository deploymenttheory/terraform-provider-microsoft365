---
page_title: "microsoft365_graph_beta_device_management_group_policy_value_reference Data Source - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Retrieves Group Policy definition metadata from Microsoft Intune using the /deviceManagement/groupPolicyDefinitions endpoint. This data source is used to discover ADMX policy details including class type, category path, and presentation configurations for policy authoring.
---

# microsoft365_graph_beta_device_management_group_policy_value_reference

Queries Microsoft Graph API for group policy definition metadata including class type, category path, presentations, and other policy details. This data source enables you to discover policy structure and configuration requirements before creating group policy values in Microsoft Intune.

## Background

Group Policy Objects (GPOs) in Microsoft Intune are based on ADMX (Administrative Template) files that define policy settings. Each policy definition contains:

- **Display Name**: The human-readable policy name as shown in management portals
- **Class Type**: Whether the policy applies to `user` or `machine` (computer) configurations
- **Category Path**: The hierarchical location of the policy in the Group Policy tree
- **Presentations**: Individual settings/checkboxes available for the policy
- **Policy Type**: Whether the policy is `admxBacked` (traditional GPO) or `admxIngested` (custom ADMX)

The same policy name may exist multiple times with different class types (user vs machine) or across different category paths (e.g., Microsoft Edge vs Google Chrome). This data source helps you identify exactly which policy variant you need for your configuration.

## Search Behavior

The data source uses **exact matching** with the following normalization:
- **Case-insensitive**: "Enable Profile Containers" matches "enable profile containers"
- **Whitespace-normalized**: Extra spaces are collapsed to single spaces

If no exact match is found, the error message provides helpful suggestions of similar policy names ranked by similarity using fuzzy matching.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.41.0-alpha | Experimental | Initial release |

## Example Usage

### Basic Policy Lookup

```terraform
# Query group policy definition metadata
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "rdp_allow" {
  policy_name = "Allow users to connect remotely by using Remote Desktop Services"
}

# Output all available attributes from the data source
output "rdp_policy_full_details" {
  value = {
    # Number of definitions found
    definitions_count = length(data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions)

    # Example of all definitions and their attributes
    definitions = [
      for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions : {
        # Definition identification
        id           = def.id
        display_name = def.display_name

        # Policy classification
        class_type    = def.class_type
        category_path = def.category_path
        policy_type   = def.policy_type

        # Policy documentation
        explain_text = def.explain_text
        supported_on = def.supported_on

        # Presentation count (handle null)
        presentations_count = try(length(def.presentations), 0)

        # All presentations with complete attributes (handle null)
        presentations = try([
          for pres in def.presentations : {
            id       = pres.id
            label    = pres.label
            type     = pres.type
            required = pres.required
          }
        ], [])
      }
    ]
  }
}

# Example output for the first definition only (simplified)
output "rdp_policy_first_definition" {
  value = length(data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions) > 0 ? {
    id            = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].id
    display_name  = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].display_name
    class_type    = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].class_type
    category_path = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].category_path
    policy_type   = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].policy_type
    explain_text  = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].explain_text
    supported_on  = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].supported_on
    presentations = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].presentations
  } : null
}
```

### Full Dependency Chain: Datasource → Configuration → Boolean Value

```terraform
# Example: Full dependency chain showing datasource -> configuration -> boolean value

# Step 1: Query the policy definition to discover metadata
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "onedrive_feedback" {
  policy_name = "Allow users to contact Microsoft for feedback and support"
}

# Extract the machine-level policy details
locals {
  feedback_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.onedrive_feedback.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "OneDrive")
  ][0]

  # Number of boolean presentations (checkboxes) for this policy
  presentation_count = length(local.feedback_policy.presentations)
}

# Step 2: Create the group policy configuration
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "onedrive_config" {
  display_name = "OneDrive Feedback Configuration"
  description  = "Configure OneDrive user feedback and support options"
}

# Step 3: Create the boolean value using discovered metadata
resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "onedrive_feedback_settings" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.onedrive_config.id

  # Use the discovered metadata from the datasource
  policy_name   = local.feedback_policy.display_name
  class_type    = local.feedback_policy.class_type
  category_path = local.feedback_policy.category_path
  enabled       = true

  # This policy has 3 boolean values (Send Feedback, Surveys, Contact Support)
  values = [
    {
      value = true # Send Feedback
    },
    {
      value = true # Receive user satisfaction surveys
    },
    {
      value = false # Contact OneDrive Support
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Output the discovered policy details
output "policy_metadata" {
  value = {
    display_name       = local.feedback_policy.display_name
    class_type         = local.feedback_policy.class_type
    category_path      = local.feedback_policy.category_path
    policy_type        = local.feedback_policy.policy_type
    presentation_count = local.presentation_count
    presentations = [
      for pres in local.feedback_policy.presentations : {
        label = pres.label
        type  = pres.presentation_type
      }
    ]
  }
}
```

### Using with Boolean Value Resource

```terraform
# Example: Using datasource to discover policy metadata for a boolean value

# Query the policy definition
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "fslogix_enable" {
  policy_name = "Enable Profile Containers"
}

# Filter for the machine-level policy in the FSLogix category
locals {
  fslogix_machine_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.fslogix_enable.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "FSLogix\\Profile Containers")
  ][0]
}

# Create a group policy configuration
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "fslogix_config" {
  display_name = "FSLogix Profile Container Configuration"
  description  = "Enables FSLogix Profile Containers"
}

# Create the boolean value using discovered metadata
resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "enable_profile_containers" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.fslogix_config.id

  # Use the discovered metadata from the datasource
  policy_name   = local.fslogix_machine_policy.display_name
  class_type    = local.fslogix_machine_policy.class_type
  category_path = local.fslogix_machine_policy.category_path
  enabled       = true

  # This policy has a single boolean value
  values = [
    {
      value = true # Enable Profile Containers
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Output the discovered policy details
output "fslogix_policy_metadata" {
  value = {
    display_name  = local.fslogix_machine_policy.display_name
    class_type    = local.fslogix_machine_policy.class_type
    category_path = local.fslogix_machine_policy.category_path
    policy_type   = local.fslogix_machine_policy.policy_type
    explain_text  = local.fslogix_machine_policy.explain_text
    presentations = length(local.fslogix_machine_policy.presentations)
  }
}
```

### Using with Text Value Resource

```terraform
# Example: Using datasource to discover policy metadata for a text value

# Query the policy definition
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "vhd_sddl" {
  policy_name = "Attached VHD SDDL"
}

# Filter for the FSLogix Profile Containers machine policy
locals {
  fslogix_sddl_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.vhd_sddl.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "FSLogix\\Profile Containers")
  ][0]
}

# Create group policy configuration
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "fslogix_config" {
  display_name = "FSLogix Configuration"
  description  = "FSLogix Profile Container settings"
}

# Create the text value using discovered metadata
resource "microsoft365_graph_beta_device_management_group_policy_text_value" "vhd_sddl" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.fslogix_config.id

  # Use discovered metadata
  policy_name   = local.fslogix_sddl_policy.display_name
  class_type    = local.fslogix_sddl_policy.class_type
  category_path = local.fslogix_sddl_policy.category_path
  enabled       = true

  # SDDL string giving Full access for admins, read/write for authenticated users
  value = "D:P(A;;FA;;;BA)(A;;FRFW;;;AU)"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Output the discovered policy information
output "sddl_policy_info" {
  value = {
    display_name  = local.fslogix_sddl_policy.display_name
    class_type    = local.fslogix_sddl_policy.class_type
    category_path = local.fslogix_sddl_policy.category_path
    explain_text  = local.fslogix_sddl_policy.explain_text
    presentations = length(local.fslogix_sddl_policy.presentations)
  }
}
```

### Using with Multi-Text Value Resource

```terraform
# Example: Using datasource to discover policy metadata for a multi-text value

# Query a policy that accepts multiple text values
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "vhd_locations" {
  policy_name = "VHD location"
}

# Filter for the FSLogix Profile Containers machine policy
locals {
  vhd_locations_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.vhd_locations.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "FSLogix\\Profile Containers")
  ][0]
}

# Create group policy configuration
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "fslogix_config" {
  display_name = "FSLogix Profile Locations"
  description  = "Configure FSLogix Profile Container storage locations"
}

# Create the multi-text value with multiple network paths
resource "microsoft365_graph_beta_device_management_group_policy_multi_text_value" "vhd_locations" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.fslogix_config.id

  # Use discovered metadata
  policy_name   = local.vhd_locations_policy.display_name
  class_type    = local.vhd_locations_policy.class_type
  category_path = local.vhd_locations_policy.category_path
  enabled       = true

  # Multiple UNC paths for profile storage (primary, secondary, tertiary)
  values = [
    "\\\\fileserver01\\profiles",
    "\\\\fileserver02\\profiles",
    "\\\\fileserver03\\profiles"
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Output the discovered policy information
output "vhd_locations_policy_info" {
  value = {
    display_name  = local.vhd_locations_policy.display_name
    class_type    = local.vhd_locations_policy.class_type
    category_path = local.vhd_locations_policy.category_path
    explain_text  = local.vhd_locations_policy.explain_text
    presentations = [
      for pres in local.vhd_locations_policy.presentations : {
        label = pres.label
        type  = pres.presentation_type
      }
    ]
  }
}
```

### Discovering and Selecting from Multiple Policy Variants

```terraform
# Example: Discovering and selecting from multiple policy variants

# Query a policy that exists in multiple locations (Chrome, Edge, etc.)
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "home_button" {
  policy_name = "Show Home button on toolbar"
}

# Output all discovered variants
output "all_home_button_variants" {
  value = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.home_button.definitions : {
      id            = def.id
      class_type    = def.class_type
      category_path = def.category_path
      policy_type   = def.policy_type
    }
  ]
  description = "Shows all variants of this policy across Chrome, Edge, and their default settings"
}

# Filter for specific browser and class type
locals {
  # Get Microsoft Edge machine policy
  edge_machine_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.home_button.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "Microsoft Edge\\Startup")
  ][0]

  # Get Google Chrome user policy
  chrome_user_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.home_button.definitions :
    def if def.class_type == "user" && contains(def.category_path, "Google\\Google Chrome\\Startup")
  ][0]
}

# Output selected variants
output "selected_policies" {
  value = {
    edge_machine = {
      display_name  = local.edge_machine_policy.display_name
      class_type    = local.edge_machine_policy.class_type
      category_path = local.edge_machine_policy.category_path
    }
    chrome_user = {
      display_name  = local.chrome_user_policy.display_name
      class_type    = local.chrome_user_policy.class_type
      category_path = local.chrome_user_policy.category_path
    }
  }
}

# Create configuration for the selected Edge policy
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "edge_config" {
  display_name = "Microsoft Edge Settings"
  description  = "Configure Edge browser home button"
}

resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "edge_home_button" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.edge_config.id

  policy_name   = local.edge_machine_policy.display_name
  class_type    = local.edge_machine_policy.class_type
  category_path = local.edge_machine_policy.category_path
  enabled       = true

  values = [
    {
      value = true # Show home button
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

### Handling Multiple Definitions (User and Machine)

Many group policies exist as both user and machine configurations. The data source returns all matching definitions with a helpful warning:

```hcl
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "edge_home_button" {
  policy_name = "Show Home button on toolbar"
}

# Filter for specific class type
locals {
  machine_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.edge_home_button.definitions : 
    def if def.class_type == "machine" && contains(def.category_path, "Microsoft Edge")
  ][0]
}

output "selected_policy" {
  value = {
    id            = local.machine_policy.id
    class_type    = local.machine_policy.class_type
    category_path = local.machine_policy.category_path
  }
}
```

### Real-World Example: Multiple Definitions

When you query a policy that exists across multiple browsers and class types, you'll receive all variants:

```terraform
╷
│ Warning: Multiple Definitions Found
│ 
│   with data.microsoft365_graph_beta_device_management_group_policy_value_reference.edge_home_button,
│   on test_group_policy_value_reference.tf line 15:
│   15: data "microsoft365_graph_beta_device_management_group_policy_value_reference" "edge_home_button" {
│ 
│ Found 8 group policy definitions with the exact name 'Show Home button on toolbar'.
│ 
│ This is common when policies exist for both User and Machine configurations.
│ 
│ Matched definitions:
│   1. "Show Home button on toolbar"
│      - Class Type: user
│      - Category: \Google\Google Chrome\Startup, Home page and New Tab page
│ 
│   2. "Show Home button on toolbar"
│      - Class Type: machine
│      - Category: \Microsoft Edge - Default Settings (users can override)\Startup, home page and new tab page
│ 
│   3. "Show Home button on toolbar"
│      - Class Type: machine
│      - Category: \Google\Google Chrome - Default Settings (users can override)\Startup, Home page and New Tab page
│ 
│   4. "Show Home button on toolbar"
│      - Class Type: machine
│      - Category: \Microsoft Edge\Startup, home page and new tab page
│ 
│   5. "Show Home button on toolbar"
│      - Class Type: user
│      - Category: \Microsoft Edge - Default Settings (users can override)\Startup, home page and new tab page
│ 
│   6. "Show Home button on toolbar"
│      - Class Type: machine
│      - Category: \Google\Google Chrome\Startup, Home page and New Tab page
│ 
│   7. "Show Home button on toolbar"
│      - Class Type: user
│      - Category: \Microsoft Edge\Startup, home page and new tab page
│ 
│   8. "Show Home button on toolbar"
│      - Class Type: user
│      - Category: \Google\Google Chrome - Default Settings (users can override)\Startup, Home page and New Tab page
│ 
│ All matching definitions are included in the results. Use the class_type, category_path, and other attributes to
│ distinguish between them in your configuration.
╵
```

### Fuzzy Matching with Suggestions

If you provide a policy name that doesn't exactly match, you'll receive helpful suggestions:

```terraform
╷
│ Error: No Exact Match Found
│ 
│   with data.microsoft365_graph_beta_device_management_group_policy_value_reference.fuzzy_test_typo,
│   on test_group_policy_fuzzy_match_demo.tf line 5:
│    5: data "microsoft365_graph_beta_device_management_group_policy_value_reference" "fuzzy_test_typo" {
│ 
│ No exact match found for policy name 'Show Home button'.
│ 
│ Did you mean one of these? (ranked by similarity):
│   1. "Show Home button on toolbar"
│ 
│ Please use the exact policy name from the list above.
╵
```


## Argument Reference

* `policy_name` - (Required) The display name of the group policy definition to search for. Requires an exact match (case-insensitive, whitespace-normalized). If no exact match is found, the error message will suggest similar policy names ranked by similarity using fuzzy matching. Example: `"Allow Cloud Policy Management"`, `"Enable Profile Containers"`.

* `timeouts` - (Optional) Timeout configuration block. See [Timeouts](#timeouts) below.

## Attributes Reference

* `id` - The computed ID of this data source operation in the format `policy_name:{policy_name}`.

* `definitions` - (Computed) List of group policy definitions matching the policy name. Multiple definitions may be returned if the same policy name exists in different categories or for different class types. Each definition contains:
  - `id` - The unique identifier (GUID) of the policy definition
  - `display_name` - The display name of the policy as shown in management portals
  - `class_type` - Whether the policy applies to `user` or `machine` configurations
  - `category_path` - The hierarchical path in the Group Policy tree (e.g., `\FSLogix\Profile Containers`)
  - `explain_text` - Detailed explanation of what the policy does
  - `supported_on` - Which operating systems/versions support this policy
  - `policy_type` - The type of policy: `admxBacked` (traditional GPO) or `admxIngested` (custom ADMX)
  - `presentations` - List of individual settings/checkboxes available for this policy:
    - `id` - The presentation GUID
    - `label` - The label text for the setting
    - `presentation_type` - The type of presentation (e.g., `groupPolicyPresentationCheckBox`, `groupPolicyPresentationText`)

## Timeouts

The `timeouts` block supports:

* `read` - (Optional) Timeout for reading data. Defaults to 3 minutes.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `policy_name` (String) The display name of the group policy definition to search for. Requires an exact match (case-insensitive, whitespace-normalized). If no exact match is found, the error message will suggest similar policy names ranked by similarity using fuzzy matching. Example: `"Allow Cloud Policy Management"`, `"Enable Profile Containers"`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `definitions` (Attributes List) List of group policy definitions matching the policy name. Multiple definitions may be returned if the same policy name exists in different categories or for different class types. (see [below for nested schema](#nestedatt--definitions))
- `id` (String) The ID of this data source operation.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--definitions"></a>
### Nested Schema for `definitions`

Read-Only:

- `category_path` (String) The full category path of the policy in the Group Policy hierarchy. Format: `\Category\Subcategory`. Example: `\FSLogix\Profile Containers`, `\Microsoft Edge\Content settings`.
- `class_type` (String) The class type of the policy. Valid values: `user` (applies to user settings), `machine` (applies to computer settings).
- `display_name` (String) The display name of the group policy definition.
- `explain_text` (String) The detailed explanation text describing what the policy does. This is the same text shown in the Group Policy Management Console.
- `id` (String) The unique identifier (GUID) of the group policy definition template.
- `policy_type` (String) The type of policy (e.g., `admxBacked`, `admxIngested`).
- `presentations` (Attributes List) List of presentations (individual settings/controls) available for this policy definition. Each presentation represents a configurable element like a checkbox, text box, dropdown, etc. (see [below for nested schema](#nestedatt--definitions--presentations))
- `supported_on` (String) The supported platforms/versions for this policy. Example: `"At least Windows 10 Server, Windows 10 or Windows 10 RT"`.

<a id="nestedatt--definitions--presentations"></a>
### Nested Schema for `definitions.presentations`

Read-Only:

- `id` (String) The unique identifier (GUID) of the presentation template.
- `label` (String) The label text displayed for this presentation in the UI. Example: `"Enable Profile Containers"`, `"Set maximum size"`.
- `required` (Boolean) Indicates whether this presentation is required when the policy is enabled. If true, the presentation must have a value set.
- `type` (String) The OData type of the presentation, indicating the control type. Examples: `#microsoft.graph.groupPolicyPresentationCheckBox` (boolean on/off), `#microsoft.graph.groupPolicyPresentationText` (text input), `#microsoft.graph.groupPolicyPresentationDecimalTextBox` (numeric input), `#microsoft.graph.groupPolicyPresentationDropdownList` (dropdown selection), `#microsoft.graph.groupPolicyPresentationListBox` (list of values).

## Use Cases

This data source supports multiple Group Policy automation scenarios:

1. **Discovery Before Configuration** - Find the exact `class_type` and `category_path` needed for a policy before creating group policy values
2. **Presentation Enumeration** - Discover which individual settings are available for a policy
3. **Policy Variant Selection** - Identify the correct policy variant when multiple exist (e.g., Chrome vs Edge, User vs Machine)
4. **Documentation Generation** - Dynamically generate documentation of available policies and their structure
5. **Configuration Validation** - Verify that policies exist before attempting to configure them

## Best Practices

1. **Use Exact Names**: Always provide the exact policy display name as it appears in Microsoft Graph API. Use the fuzzy match suggestions when you get them.

2. **Handle Multiple Results**: Many policies exist for both user and machine configurations. Always filter results by `class_type` and `category_path` to select the correct variant:
   ```hcl
   locals {
     machine_policy = [
       for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.example.definitions : 
       def if def.class_type == "machine"
     ][0]
   }
   ```

3. **Check Result Count**: Validate that `definitions` is not empty and contains the expected number of results before accessing array indices.

4. **Use Locals for Filtering**: Store filtered results in local values for reuse across multiple resources to avoid repeating filter logic.

5. **Document Your Selection**: When multiple definitions exist, add comments explaining which variant you chose and why (e.g., "Using Microsoft Edge machine policy, not Chrome variant").

## Important Notes

### Requires Microsoft Graph API Access

This data source queries the Microsoft Graph API and requires appropriate credentials and permissions. Unlike some utility data sources that work offline, this one makes real-time API calls to your Microsoft Intune tenant.

**Required Permissions:**
- `DeviceManagementConfiguration.Read.All`

### Multiple Definitions are Common

It's normal to receive multiple definitions for the same policy name. This happens because:
- Policies often exist for both **user** and **machine** configurations
- The same policy may exist across multiple products (Microsoft Edge, Google Chrome, Firefox)
- Policies may have "Default Settings (users can override)" variants

Always use `class_type` and `category_path` to select the correct definition for your needs.

### Policy Availability Depends on ADMX Files

The policies available in your tenant depend on which ADMX files have been ingested. Common sources include:
- **Windows Built-in** - Traditional Windows policies (`admxBacked`)
- **Microsoft Edge** - Edge browser policies
- **Microsoft Office** - Office application policies
- **Custom ADMX** - Uploaded custom ADMX files (`admxIngested`)

If a policy isn't found, verify that the corresponding ADMX files are loaded in your Intune tenant.

### Search Performance

The data source uses the Microsoft Graph API's `contains()` filter for efficient searching, then performs exact matching in code. Large result sets are handled efficiently through pagination using PageIterator.

## Common Policies

Here are some commonly used group policies that can be looked up with this data source:

| Policy Name | Common Class Types | Typical Category Path | Use Case |
|-------------|-------------------|----------------------|----------|
| Enable Profile Containers | machine, user | \FSLogix\Profile Containers | FSLogix profile redirection |
| Allow Cloud Policy Management | machine | \FSLogix\Profile Containers | Enable cloud-based FSLogix policy |
| Show Home button on toolbar | user, machine | \Microsoft Edge\Startup, home page... | Configure Edge home button |
| Configure Automatic Updates | machine | \Windows Components\Windows Update | Windows Update settings |
| Allow users to connect remotely by using Remote Desktop Services | machine | \Windows Components\Remote Desktop Services\... | Enable RDP connections |

## Additional Resources

- [Microsoft Graph API - Group Policy Definitions](https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicydefinition-get?view=graph-rest-beta)
- [Group Policy Overview](https://learn.microsoft.com/en-us/mem/intune/configuration/administrative-templates-windows)
- [ADMX-backed Policies in Intune](https://learn.microsoft.com/en-us/mem/intune/configuration/administrative-templates-import-custom)
- [FSLogix Group Policy Reference](https://learn.microsoft.com/en-us/fslogix/reference-configuration-settings)

