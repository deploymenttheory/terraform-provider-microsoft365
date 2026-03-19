---
page_title: "microsoft365_graph_beta_windows_updates_update_policy Resource - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Manages a Windows Update policy using the /admin/windows/updates/updatePolicies endpoint. An update policy serves as a container for compliance changes (content approvals) that define which updates should be deployed to devices. This resource is a prerequisite for creating content approvals.
---

# microsoft365_graph_beta_windows_updates_update_policy (Resource)

Manages a Windows Update policy using the `/admin/windows/updates/updatePolicies` endpoint. An update policy serves as a container for compliance changes (content approvals) that define which updates should be deployed to devices. This resource is a prerequisite for creating content approvals.

## Microsoft Documentation

- [updatePolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-updatepolicy?view=graph-rest-beta)
- [Create updatePolicy](https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-post-updatepolicies?view=graph-rest-beta)
- [Get updatePolicy](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-get?view=graph-rest-beta)
- [Update updatePolicy](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-update?view=graph-rest-beta)
- [Delete updatePolicy](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-delete?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Important Behavioral Notes

### Compliance Change Rules Immutability

The `compliance_change_rules` block and its nested attributes (`filter_type`, `duration_before_deployment_start`) **cannot be modified** after the resource is created. Any attempt to change these values will trigger a resource replacement (destroy and recreate). This is enforced through the `RequiresReplace` plan modifier.

### Filter Type Restrictions

Despite what the Microsoft Graph API documentation may suggest, only `driverUpdateFilter` is valid for the `filter_type` attribute. Using `windowsUpdateFilter` will result in a schema validation error from the API. This has been validated through extensive API testing.

### Duration Format and Validation

Duration values must be specified in ISO 8601 duration format using day notation (e.g., `P7D` for 7 days):

- **`duration_before_deployment_start`**: Valid range is P1D to P30D (1 to 30 days)
- **`duration_between_offers`**: Valid range is P1D to P30D (1 to 30 days)

The provider automatically handles SDK normalization issues where the Microsoft Graph SDK converts day-based durations to week-based formats (e.g., P7D → P1W) by denormalizing them back to the original day format to prevent spurious diffs.

### Devices Per Offer Validation

The `devices_per_offer` attribute must be a positive integer (minimum value of 1). Testing has confirmed that values up to 1,000,000 are accepted by the API.

### Compliance Changes Field

The `compliance_changes` field is **write-only** and is not returned by the API after creation. It will always show as changed during import operations and must be added to `ImportStateVerifyIgnore` in tests.

### Update Operation Restrictions

When updating an existing update policy, only the `deployment_settings` block can be modified. The API does not accept changes to `audience`, `compliance_changes`, or `compliance_change_rules` in PATCH operations, even though the documentation may suggest otherwise.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Scenario 1: Minimal Update Policy

```terraform
# ==============================================================================
# Minimal Update Policy
# ==============================================================================
# Creates a Windows Update policy with only the required fields.
# This minimal configuration creates a policy with an audience reference
# and enables compliance changes, but does not configure any compliance
# change rules or deployment settings.

resource "microsoft365_graph_beta_windows_updates_update_policy" "minimal" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true
}
```

### Scenario 2: Update Policy with Compliance Change Rules

```terraform
# ==============================================================================
# Update Policy with Compliance Change Rules
# ==============================================================================
# Creates a Windows Update policy with compliance change rules that automatically
# approve driver updates after a 7-day waiting period. The compliance change rules
# cannot be modified after creation - any changes will require resource replacement.

data "microsoft365_graph_beta_windows_updates_deployment_audience" "example" {
  id = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "with_compliance_rules" {
  audience_id        = data.microsoft365_graph_beta_windows_updates_deployment_audience.example.id
  compliance_changes = true

  compliance_change_rules = [
    {
      content_filter = {
        filter_type = "driverUpdateFilter"
      }
      duration_before_deployment_start = "P7D"
    }
  ]
}
```

### Scenario 3: Update Policy with Gradual Rollout

```terraform
# ==============================================================================
# Update Policy with Gradual Rollout Settings
# ==============================================================================
# Creates a Windows Update policy with deployment settings that control how
# updates are rolled out to devices. The gradual rollout settings specify that
# updates should be offered to 1000 devices at a time, with 1 day between each
# batch of offers.

resource "microsoft365_graph_beta_windows_updates_update_policy" "with_gradual_rollout" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true

  deployment_settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P1D"
        devices_per_offer       = 1000
      }
    }
  }
}
```

### Scenario 4: Update Policy with Scheduled Start Time

```terraform
# ==============================================================================
# Update Policy with Scheduled Start Time
# ==============================================================================
# Creates a Windows Update policy with a specific start date/time for the
# deployment schedule. This allows you to control when updates begin rolling
# out to devices. The gradual rollout settings control the pace of deployment.

resource "microsoft365_graph_beta_windows_updates_update_policy" "with_scheduled_start" {
  audience_id        = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
  compliance_changes = true

  deployment_settings = {
    schedule = {
      start_date_time = "2026-04-01T00:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P2D"
        devices_per_offer       = 500
      }
    }
  }
}
```

### Scenario 5: Complete Update Policy Configuration

```terraform
# ==============================================================================
# Complete Update Policy Configuration
# ==============================================================================
# Creates a comprehensive Windows Update policy with all available options:
# - Compliance change rules for automatic driver update approval
# - Deployment settings with scheduled start time
# - Gradual rollout configuration for controlled deployment
#
# Note: Compliance change rules cannot be modified after creation. Changes to
# filter_type or duration_before_deployment_start will require resource replacement.

data "microsoft365_graph_beta_windows_updates_deployment_audience" "production" {
  id = "8c4eb1eb-d7a3-4633-8e2f-f926e82df08e"
}

resource "microsoft365_graph_beta_windows_updates_update_policy" "complete" {
  audience_id        = data.microsoft365_graph_beta_windows_updates_deployment_audience.production.id
  compliance_changes = true

  compliance_change_rules = [
    {
      content_filter = {
        filter_type = "driverUpdateFilter"
      }
      duration_before_deployment_start = "P14D"
    }
  ]

  deployment_settings = {
    schedule = {
      start_date_time = "2026-04-15T02:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P7D"
        devices_per_offer       = 2000
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `audience_id` (String) The ID of the deployment audience to target with this policy.
- `compliance_changes` (Boolean) Enable compliance changes (content approvals) for this policy. Must be set to `true` to create content approvals.

### Optional

- `compliance_change_rules` (Attributes Set) Rules for governing the automatic creation of compliance changes. Cannot be updated after creation - changes require resource replacement. (see [below for nested schema](#nestedatt--compliance_change_rules))
- `deployment_settings` (Attributes) Settings for governing how to deploy content. (see [below for nested schema](#nestedatt--deployment_settings))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time when the update policy was created. Read-only.
- `id` (String) The unique identifier for the update policy.

<a id="nestedatt--compliance_change_rules"></a>
### Nested Schema for `compliance_change_rules`

Optional:

- `content_filter` (Attributes) The content filter for the compliance change rule. (see [below for nested schema](#nestedatt--compliance_change_rules--content_filter))
- `duration_before_deployment_start` (String) The duration before deployment starts (ISO 8601 duration format). Valid range: P1D to P30D (1 to 30 days). Cannot be updated after creation. Examples: 'P7D' for 7 days, 'P14D' for 14 days.

Read-Only:

- `created_date_time` (String) The date and time when the rule was created. Read-only.
- `last_evaluated_date_time` (String) The date and time when the rule was last evaluated. Read-only.
- `last_modified_date_time` (String) The date and time when the rule was last modified. Read-only.

<a id="nestedatt--compliance_change_rules--content_filter"></a>
### Nested Schema for `compliance_change_rules.content_filter`

Required:

- `filter_type` (String) The type of content filter. Only `driverUpdateFilter` is supported. Note: `windowsUpdateFilter` is not valid despite appearing in API documentation.



<a id="nestedatt--deployment_settings"></a>
### Nested Schema for `deployment_settings`

Optional:

- `schedule` (Attributes) Settings for the schedule of the deployment. (see [below for nested schema](#nestedatt--deployment_settings--schedule))

<a id="nestedatt--deployment_settings--schedule"></a>
### Nested Schema for `deployment_settings.schedule`

Optional:

- `gradual_rollout` (Attributes) Settings for gradual rollout. (see [below for nested schema](#nestedatt--deployment_settings--schedule--gradual_rollout))
- `start_date_time` (String) The start date and time for the deployment (ISO 8601 format).

<a id="nestedatt--deployment_settings--schedule--gradual_rollout"></a>
### Nested Schema for `deployment_settings.schedule.gradual_rollout`

Required:

- `devices_per_offer` (Number) The number of devices to offer the update to in each batch. Must be a positive integer (minimum 1).
- `duration_between_offers` (String) The duration between offers (ISO 8601 duration format). Valid range: P1D to P30D (1 to 30 days). Examples: 'P1D' for 1 day, 'P7D' for 7 days.




<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import an existing Windows Update policy using its ID
terraform import microsoft365_graph_beta_windows_updates_update_policy.example a1b2c3d4-1234-5678-abcd-a1b2c3d4e5f6
```
