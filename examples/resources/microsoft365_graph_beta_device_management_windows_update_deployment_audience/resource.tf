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

# Example: Create a deployment audience container
# This creates an empty audience that can be populated with members using the
# microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience_members resource
resource "microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience" "example" {
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Output the audience ID for use with the members resource
output "audience_id" {
  value = microsoft365_graph_beta_device_management_windows_autopatch_deployment_audience.example.id
}
