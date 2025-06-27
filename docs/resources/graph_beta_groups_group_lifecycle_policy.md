---
page_title: "microsoft365_graph_beta_groups_group_lifecycle_policy Resource - terraform-provider-microsoft365"
subcategory: "Groups"
description: |-
    Manages group lifecycle policies for Microsoft 365 groups using the /groupLifecyclePolicies endpoint. This resource enables administrators to set expiration periods for groups, requiring owners to renew them within specified time intervals. When a group reaches its expiration, it can be renewed to extend the expiration date, or if not renewed, it expires and is deleted with a 30-day restoration window.
---

# microsoft365_graph_beta_groups_group_lifecycle_policy (Resource)

Manages group lifecycle policies for Microsoft 365 groups using the `/groupLifecyclePolicies` endpoint. This resource enables administrators to set expiration periods for groups, requiring owners to renew them within specified time intervals. When a group reaches its expiration, it can be renewed to extend the expiration date, or if not renewed, it expires and is deleted with a 30-day restoration window.

## Microsoft Documentation

- [Group lifecycle policy resource type](https://learn.microsoft.com/en-us/graph/api/resources/grouplifecyclepolicy?view=graph-rest-beta)
- [Create groupLifecyclePolicy](https://learn.microsoft.com/en-us/graph/api/group-post-grouplifecyclepolicies?view=graph-rest-beta)
- [Update groupLifecyclePolicy](https://learn.microsoft.com/en-us/graph/api/group-update-grouplifecyclepolicy?view=graph-rest-beta)
- [Delete groupLifecyclePolicy](https://learn.microsoft.com/en-us/graph/api/group-delete-grouplifecyclepolicy?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Group.ReadWrite.All`, `Directory.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.18.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example configurations for Microsoft 365 Group Lifecycle Policy
# This resource manages group lifecycle policies that set expiration periods for Microsoft 365 groups

# Scenario 1: Default lifecycle policy for all Microsoft 365 groups
# This policy applies to all groups and sets a 180-day expiration period
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "default_policy" {
  # Number of days before a group expires and needs to be renewed
  group_lifetime_in_days = 180

  # Apply to all Microsoft 365 groups
  managed_group_types = "All"

  # Optional: List of email addresses to send notifications for groups without owners
  # Multiple email addresses can be defined by separating with semicolons
  alternate_notification_emails = "admin@example.com;notifications@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 2: Short-term project groups policy
# This policy is for temporary project groups that should expire quickly
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "project_groups_policy" {
  # Short expiration period for project groups (90 days)
  group_lifetime_in_days = 90

  # Apply to selected groups only (requires manual assignment)
  managed_group_types = "Selected"

  # Multiple notification emails for project management
  alternate_notification_emails = "pm@example.com;project-admin@example.com;it-support@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 3: Long-term department groups policy
# This policy is for permanent department groups with longer expiration
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "department_groups_policy" {
  # Longer expiration period for department groups (365 days)
  group_lifetime_in_days = 365

  # Apply to selected groups only
  managed_group_types = "Selected"

  # Single notification email for department heads
  alternate_notification_emails = "dept-heads@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 4: Disabled lifecycle policy
# This policy effectively disables group lifecycle management
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "disabled_policy" {
  # Set a very long expiration period (10 years) to effectively disable
  group_lifetime_in_days = 3650

  # Don't apply to any groups
  managed_group_types = "None"

  # No notification emails needed since policy is disabled
  # alternate_notification_emails is omitted (optional field)

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Scenario 5: Compliance-focused policy
# This policy ensures regular review of groups for compliance purposes
resource "microsoft365_graph_beta_groups_group_lifecycle_policy" "compliance_policy" {
  # 6-month expiration for compliance review cycles
  group_lifetime_in_days = 180

  # Apply to all groups for comprehensive compliance
  managed_group_types = "All"

  # Multiple stakeholders for compliance notifications
  alternate_notification_emails = "compliance@example.com;security@example.com;legal@example.com;admin@example.com"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group_lifetime_in_days` (Number) Number of days before a group expires and needs to be renewed. Once renewed, the group expiration is extended by the number of days defined.
- `managed_group_types` (String) The group type for which the expiration policy applies. Possible values are **All**, **Selected** or **None**.

### Optional

- `alternate_notification_emails` (String) List of email address to send notifications for groups without owners. Multiple email address can be defined by separating email address with a semicolon.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) A unique identifier for a policy. Read-only.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Lifecycle Policies**: Used to manage the lifecycle of Microsoft 365 groups, including expiration and renewal.
- **Notification Recipients**: The `alternate_notification_emails` attribute allows specifying additional recipients for expiration notifications.
- **Policy Scope**: Policies can be applied to all groups or a subset using `group_lifetime_in_days` and `managed_group_types`.
- **Renewal**: Groups can be renewed via API or user action if they are about to expire.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import scripts for Microsoft 365 Group Lifecycle Policy
# Replace {policy_id} with the actual policy ID from your Microsoft 365 tenant

# Import the default policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_policy.default_policy {default_policy_id}

# Import the project groups policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_policy.project_groups_policy {project_policy_id}

# Import the department groups policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_policy.department_groups_policy {department_policy_id}

# Import the disabled policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_policy.disabled_policy {disabled_policy_id}

# Import the compliance policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_policy.compliance_policy {compliance_policy_id}

# Note: You can import individual policies as needed. Not all policies need to be imported.
# To find policy IDs, you can use Microsoft Graph API or PowerShell:
# Get-MgGroupLifecyclePolicy
``` 