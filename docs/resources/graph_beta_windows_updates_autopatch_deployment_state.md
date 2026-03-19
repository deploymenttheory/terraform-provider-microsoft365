---
page_title: "microsoft365_graph_beta_windows_updates_autopatch_deployment_state Resource - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Manages the lifecycle state of a Windows Update autopatch deployment. Using the /admin/windows/updates/deployments/{deploymentId}/update endpoint. This resource allows pausing, resuming, or archiving a deployment independently from its configuration. The deployment must be created first using the microsoft365_graph_beta_windows_updates_autopatch_deployment resource.
---

# microsoft365_graph_beta_windows_updates_autopatch_deployment_state (Resource)

Manages the lifecycle state of a Windows Update autopatch deployment. Using the `/admin/windows/updates/deployments/{deploymentId}/update` endpoint. This resource allows pausing, resuming, or archiving a deployment independently from its configuration. The deployment must be created first using the `microsoft365_graph_beta_windows_updates_autopatch_deployment` resource.

## Microsoft Documentation

- [deployment resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deployment?view=graph-rest-beta)
- [Update deployment](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-update?view=graph-rest-beta)
- [Get deployment](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-get?view=graph-rest-beta)
- [deploymentState resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deploymentstate?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Relationship to Deployment Resource

This resource manages only the lifecycle state of a deployment created by the
`microsoft365_graph_beta_windows_updates_autopatch_deployment` resource. The
`deployment_id` must reference an existing deployment.

On destroy, this resource resets the deployment state to `none` (active/offering)
rather than deleting the deployment itself. To delete the deployment, destroy the
`microsoft365_graph_beta_windows_updates_autopatch_deployment` resource.

## State Values

| `requested_value` | Description |
|-------------------|-------------|
| `none` | Active — the deployment is offering updates to devices |
| `paused` | Paused — no further devices are offered the update until resumed |
| `archived` | Archived — permanently stopped; cannot be resumed |

The `effective_value` attribute reflects the actual state reported by the Graph API
(e.g. `scheduled`, `offering`, `paused`, `faulted`, `archived`) and may differ from
`requested_value` while the service processes the state transition.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Scenario 1: Pause a Deployment

```terraform
# ==============================================================================
# Pause a Deployment
# ==============================================================================
# Pauses an active autopatch deployment. The deployment must already exist,
# managed by the microsoft365_graph_beta_windows_updates_autopatch_deployment
# resource. Pausing stops the deployment from offering the update to additional
# devices while preserving progress on devices already receiving it.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "example" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_state" "paused" {
  deployment_id   = microsoft365_graph_beta_windows_updates_autopatch_deployment.example.id
  requested_value = "paused"
}
```

### Scenario 2: Resume a Deployment

```terraform
# ==============================================================================
# Resume a Deployment (set to none / active)
# ==============================================================================
# Sets the deployment state to "none", which means the deployment is active and
# offering updates to devices. Use this to resume a previously paused deployment
# or to explicitly set a deployment into its offering state.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "example" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_state" "active" {
  deployment_id   = microsoft365_graph_beta_windows_updates_autopatch_deployment.example.id
  requested_value = "none"
}
```

### Scenario 3: Archive a Deployment

```terraform
# ==============================================================================
# Archive a Deployment
# ==============================================================================
# Archives a deployment, permanently stopping it from offering updates to any
# further devices. Archived deployments are read-only and cannot be resumed.
# Use this when a deployment is complete or no longer needed.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "example" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
    catalog_entry_type = "qualityUpdate"
  }
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_state" "archived" {
  deployment_id   = microsoft365_graph_beta_windows_updates_autopatch_deployment.example.id
  requested_value = "archived"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `deployment_id` (String) The ID of the deployment whose state is being managed.
- `requested_value` (String) The requested state of the deployment. Valid values are: `none` (active/offering), `paused`, `archived`.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `effective_value` (String) The effective state value of the deployment. Possible values: `scheduled`, `offering`, `paused`, `faulted`, `archived` (read-only).
- `id` (String) The ID of the deployment (same as deployment_id).

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
# Import an existing deployment state by the deployment ID
terraform import microsoft365_graph_beta_windows_updates_autopatch_deployment_state.example 00000000-0000-0000-0000-000000000000
```
