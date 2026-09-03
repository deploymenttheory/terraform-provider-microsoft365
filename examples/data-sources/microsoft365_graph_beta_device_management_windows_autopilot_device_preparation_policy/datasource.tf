# Look up a Windows Autopilot device preparation policy by its Intune object id
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

# Look up a Windows Autopilot device preparation policy by its name
data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "by_name" {
  name = "Autopilot Device Preparation - Cloud PC"
}

# Use the looked up policy id when provisioning Windows 365 Cloud PCs. This allows the
# provisioning policy to be managed in a different Terraform state than the device
# preparation policy itself.
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "example" {
  display_name             = "Cloud PC - Autopilot device preparation"
  description              = "Provisioning policy using an existing device preparation policy"
  image_id                 = "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365"
  image_type               = "gallery"
  provisioning_type        = "dedicated"
  cloud_pc_naming_template = "CPC-%USERNAME:4%-%RAND:5%"

  domain_join_configurations = [
    {
      domain_join_type = "azureADJoin"
      region_group     = "europeUnion"
      region_name      = "automatic"
    }
  ]

  autopilot_configuration = {
    device_preparation_profile_id = data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.by_name.id
  }
}

output "device_preparation_policy_id" {
  value       = data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.by_name.id
  description = "ID of the existing Windows Autopilot device preparation policy"
}

output "device_preparation_policy_description" {
  value       = data.microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.by_id.description
  description = "Description of the existing Windows Autopilot device preparation policy"
}
