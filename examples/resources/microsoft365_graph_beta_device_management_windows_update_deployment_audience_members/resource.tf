terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

provider "microsoft365" {
  cloud = "public"
}

# Example 1: Create an audience container
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "example" {
  timeouts = {
    create = "10m"
    read   = "5m"
    delete = "5m"
  }
}

# Example 2: Populate the audience with device members and exclusions
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members" "example_devices" {
  audience_id = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.example.id
  member_type = "azureADDevice"

  members = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002"
  ]

  exclusions = [
    "00000000-0000-0000-0000-000000000003"
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 3: Use updatable asset groups (for group-based targeting)
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience_members" "example_groups" {
  audience_id = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.example.id
  member_type = "updatableAssetGroup"

  members = [
    "00000000-0000-0000-0000-000000000004",
    "00000000-0000-0000-0000-000000000005"
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
