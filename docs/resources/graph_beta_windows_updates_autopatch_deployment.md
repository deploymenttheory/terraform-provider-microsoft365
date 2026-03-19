---
page_title: "microsoft365_graph_beta_windows_updates_autopatch_deployment Resource - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Manages Windows Update deployments for deploying content to a set of devices using the /admin/windows/updates/deployments endpoint. Deployments define which update content (feature or quality updates) should be deployed, to which audience, and with what settings (schedule, monitoring, etc.).
---

# microsoft365_graph_beta_windows_updates_autopatch_deployment (Resource)

Manages Windows Update deployments for deploying content to a set of devices using the `/admin/windows/updates/deployments` endpoint. Deployments define which update content (feature or quality updates) should be deployed, to which audience, and with what settings (schedule, monitoring, etc.).

## Microsoft Documentation

- [deployment resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deployment?view=graph-rest-beta)
- [Create deployment](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updates-post-deployments?view=graph-rest-beta)
- [Get deployment](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-get?view=graph-rest-beta)
- [Update deployment](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-update?view=graph-rest-beta)
- [Delete deployment](https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-delete?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Monitoring Rule Constraints

The `signal` and `action` fields within each monitoring rule must satisfy the following pairing rules, which are enforced at plan time:

| Signal | Required Action | Threshold |
|--------|----------------|-----------|
| `rollback` | `pauseDeployment` or `alertError` | Required (≥ 1) |
| `ineligible` | `offerFallback` | Must not be set |

The `offerFallback` action is only supported on feature update deployments targeting Windows 11 and automatically offers Windows 10 22H2 to ineligible devices.

## Settings Mutability

Deployment settings (`schedule`, `monitoring`) may be added to a deployment that was created without them via an in-place update. However, once settings are configured they cannot be modified — any subsequent change to `settings` will require the resource to be **replaced** (destroyed and re-created).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Scenario 1: Minimal Deployment

```terraform
# ==============================================================================
# Minimal Deployment (no settings)
# ==============================================================================
# Creates a Windows Update deployment with only the required content block.
# No schedule or monitoring rules are configured. The deployment will be
# created in a pending state and can have settings added later by updating
# the resource. Once settings are applied they cannot be modified in place —
# changes will require the resource to be replaced.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "minimal" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }
}
```

### Scenario 2: Feature Update with Rollback / Pause Deployment

```terraform
# ==============================================================================
# Feature Update — Rollback Signal with Pause Deployment Action
# ==============================================================================
# Deploys a feature update with a rate-driven gradual rollout schedule and a
# monitoring rule that pauses the deployment if too many devices roll back.
# This is the most common monitoring configuration for feature update deployments.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "rollback_pause" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2025-02-01T08:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P7D"
        devices_per_offer       = 100
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 5
          action    = "pauseDeployment"
        }
      ]
    }
  }
}
```

### Scenario 3: Feature Update with Rollback / Alert Error

```terraform
# ==============================================================================
# Feature Update — Rollback Signal with Alert Error Action
# ==============================================================================
# Deploys a feature update with a gradual rollout schedule and a monitoring
# rule that raises an alert (rather than pausing) when the rollback threshold
# is reached. Use this when you want visibility into rollback rates without
# automatically stopping the deployment.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "rollback_alert" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P7D"
        devices_per_offer       = 100
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 10
          action    = "alertError"
        }
      ]
    }
  }
}
```

### Scenario 4: Feature Update with Ineligible / Offer Fallback

```terraform
# ==============================================================================
# Feature Update — Ineligible Signal with Offer Fallback Action
# ==============================================================================
# Deploys a Windows 11 feature update with a monitoring rule that automatically
# offers Windows 10 22H2 as a fallback to devices that are ineligible for the
# Windows 11 update. This combination is only valid for Windows 11 feature
# update deployments.
#
# Note: the "ineligible" signal must always be paired with the "offerFallback"
# action, and no threshold is accepted for this combination — the fallback is
# offered to all ineligible devices unconditionally.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "ineligible_fallback" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    monitoring = {
      monitoring_rules = [
        {
          signal = "ineligible"
          action = "offerFallback"
        }
      ]
    }
  }
}
```

### Scenario 5: Feature Update with Multiple Monitoring Rules

```terraform
# ==============================================================================
# Feature Update — Multiple Monitoring Rules
# ==============================================================================
# Deploys a feature update with two monitoring rules configured simultaneously:
#   1. Pause the deployment if more than 5% of devices roll back.
#   2. Offer Windows 10 22H2 as a fallback to any device ineligible for Windows 11.
#
# Multiple monitoring rules allow you to respond to different failure signals
# independently within a single deployment.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "multiple_rules" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_update.entries[0].id
    catalog_entry_type = "featureUpdate"
  }

  settings = {
    schedule = {
      gradual_rollout = {
        duration_between_offers = "P14D"
        devices_per_offer       = 200
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 5
          action    = "pauseDeployment"
        },
        {
          signal = "ineligible"
          action = "offerFallback"
        }
      ]
    }
  }
}
```

### Scenario 6: Quality Update with Date-Driven Rollout

```terraform
# ==============================================================================
# Quality Update — Date-Driven Rollout
# ==============================================================================
# Deploys a cumulative quality (security) update using a date-driven gradual
# rollout. Devices are progressively offered the update until the specified
# end date, rather than rolling out a fixed number of devices per wave.

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "quality_update" {
  content = {
    catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_update.entries[0].id
    catalog_entry_type = "qualityUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2026-02-01T08:00:00Z"
      gradual_rollout = {
        end_date_time = "2026-03-01T08:00:00Z"
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 5
          action    = "pauseDeployment"
        }
      ]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (Attributes) Specifies what content to deploy. Cannot be changed after creation. (see [below for nested schema](#nestedatt--content))

### Optional

- `settings` (Attributes) Settings specified on the deployment governing how to deploy content. Settings may be added to a deployment that was created without them, but once configured they cannot be modified in place — changes will require the resource to be replaced. (see [below for nested schema](#nestedatt--settings))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time when the deployment was created.
- `id` (String) Unique identifier for the deployment.
- `last_modified_date_time` (String) The date and time when the deployment was last modified.

<a id="nestedatt--content"></a>
### Nested Schema for `content`

Required:

- `catalog_entry_id` (String) The ID of the catalog entry to deploy. This should reference a feature or quality update from the Windows Update catalog. Cannot be changed after creation and will require a replacement to update.
- `catalog_entry_type` (String) The type of catalog entry being deployed. Valid values are: `featureUpdate`, `qualityUpdate`.


<a id="nestedatt--settings"></a>
### Nested Schema for `settings`

Optional:

- `monitoring` (Attributes) Monitoring settings for the deployment. (see [below for nested schema](#nestedatt--settings--monitoring))
- `schedule` (Attributes) Schedule settings for the deployment. (see [below for nested schema](#nestedatt--settings--schedule))

<a id="nestedatt--settings--monitoring"></a>
### Nested Schema for `settings.monitoring`

Optional:

- `monitoring_rules` (Attributes Set) Rules for monitoring the deployment and the action that should be taken when the signal threshold is met. (see [below for nested schema](#nestedatt--settings--monitoring--monitoring_rules))

<a id="nestedatt--settings--monitoring--monitoring_rules"></a>
### Nested Schema for `settings.monitoring.monitoring_rules`

Required:

- `action` (String) The action to take when the monitoring threshold is met. Valid values are: `alertError`, `offerFallback`, `pauseDeployment`.
- `signal` (String) The signal to monitor. Valid values are: `rollback` or `ineligible`.

Optional:

- `threshold` (Number) The percentage of devices that trigger the action. Required for `rollback` signal. Must not be set when `action` is `offerFallback` (the fallback is offered to all ineligible devices unconditionally).



<a id="nestedatt--settings--schedule"></a>
### Nested Schema for `settings.schedule`

Optional:

- `gradual_rollout` (Attributes) Settings for gradual rollout of the deployment. (see [below for nested schema](#nestedatt--settings--schedule--gradual_rollout))
- `start_date_time` (String) The date and time when the deployment should start. Must be in ISO 8601 format.

<a id="nestedatt--settings--schedule--gradual_rollout"></a>
### Nested Schema for `settings.schedule.gradual_rollout`

Optional:

- `devices_per_offer` (Number) The number of devices to offer the update to in each rollout wave.
- `duration_between_offers` (String) The duration between each offer in ISO 8601 format (e.g., `P7D` for 7 days).
- `end_date_time` (String) The date and time when the gradual rollout should complete. Must be in ISO 8601 format.




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
# Import an existing Windows Update autopatch deployment by its ID
terraform import microsoft365_graph_beta_windows_updates_autopatch_deployment.example 00000000-0000-0000-0000-000000000000
```
