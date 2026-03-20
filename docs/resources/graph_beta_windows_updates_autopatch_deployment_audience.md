---
page_title: "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience Resource - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Manages a Windows Update deployment audience in Microsoft 365, including its members and exclusions. Uses the /admin/windows/updates/deploymentAudiences endpoint to create the audience container and the updateAudienceById action to manage members and exclusions in a single resource. See the Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deploymentaudience?view=graph-rest-beta for more information.
---

# microsoft365_graph_beta_windows_updates_autopatch_deployment_audience (Resource)

Manages a Windows Update deployment audience in Microsoft 365, including its members and exclusions. Uses the `/admin/windows/updates/deploymentAudiences` endpoint to create the audience container and the `updateAudienceById` action to manage members and exclusions in a single resource. See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deploymentaudience?view=graph-rest-beta) for more information.

## Microsoft Documentation

- [deploymentAudience resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deploymentaudience?view=graph-rest-beta)
- [Get deploymentAudience](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-get?view=graph-rest-beta)
- [Update deploymentAudience](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-update?view=graph-rest-beta)
- [updateAudience](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-updateaudience?view=graph-rest-beta)
- [updateAudienceById](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-updateaudiencebyid?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Important Notes

- **Singleton Resource**: This is a singleton resource — one deployment audience exists per tenant for Windows Autopatch. It cannot be created or deleted via the API, only updated.
- **Member Types**: The `member_type` field specifies the type of members in the audience. Currently only `updatableAssetGroup` is supported.
- **Members and Exclusions**: Members are updatable asset groups included in the deployment audience. Exclusions are groups explicitly excluded from deployments. Exclusions take precedence over members.
- **Update Behavior**: Updates use the `updateAudienceById` API to add or remove members and exclusions.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Minimal

```terraform
# Minimal example — creates an empty deployment audience with no members or exclusions.
# This is the singleton autopatch deployment audience resource for the tenant.

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "example" {
  timeouts = {
    create = "5m"
    read   = "5m"
    delete = "10m"
  }
}
```

### With Members

```terraform
# Deployment audience with members — adds updatable asset groups as members.
# A 30-second time_sleep ensures the groups have fully propagated before assignment.

resource "microsoft365_graph_beta_groups_group" "member_group_1" {
  display_name     = "autopatch-audience-member-1"
  mail_enabled     = false
  mail_nickname    = "autopatch-member-1"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "member_group_2" {
  display_name     = "autopatch-audience-member-2"
  mail_enabled     = false
  mail_nickname    = "autopatch-member-2"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.member_group_1,
    microsoft365_graph_beta_groups_group.member_group_2
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "example" {
  depends_on = [time_sleep.wait_for_groups]

  member_type = "updatableAssetGroup"

  members = [
    microsoft365_graph_beta_groups_group.member_group_1.id,
    microsoft365_graph_beta_groups_group.member_group_2.id
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "10m"
  }
}
```

### With Members and Exclusions

```terraform
# Deployment audience with members and exclusions — adds updatable asset groups
# as members and specifies exclusion groups. Exclusions take precedence over members.

resource "microsoft365_graph_beta_groups_group" "member_group_1" {
  display_name     = "autopatch-audience-member-1"
  mail_enabled     = false
  mail_nickname    = "autopatch-member-1"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "member_group_2" {
  display_name     = "autopatch-audience-member-2"
  mail_enabled     = false
  mail_nickname    = "autopatch-member-2"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "exclusion_group" {
  display_name     = "autopatch-audience-exclusion"
  mail_enabled     = false
  mail_nickname    = "autopatch-exclusion"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.member_group_1,
    microsoft365_graph_beta_groups_group.member_group_2,
    microsoft365_graph_beta_groups_group.exclusion_group
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "example" {
  depends_on = [time_sleep.wait_for_groups]

  member_type = "updatableAssetGroup"

  members = [
    microsoft365_graph_beta_groups_group.member_group_1.id,
    microsoft365_graph_beta_groups_group.member_group_2.id
  ]

  exclusions = [
    microsoft365_graph_beta_groups_group.exclusion_group.id
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "10m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `exclusions` (Set of String) Set of device or Entra group IDs to exclude from the deployment audience.
- `member_type` (String) The type of members in this audience. All members and exclusions must be of the same type. Valid values are: `azureADDevice`, `updatableAssetGroup`. Defaults to `azureADDevice`.
- `members` (Set of String) Set of device or Entra group IDs to include in the deployment audience.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for the deployment audience.

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
# Import the singleton deployment audience resource
terraform import microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.example "default"
```
