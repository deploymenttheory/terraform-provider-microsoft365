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

