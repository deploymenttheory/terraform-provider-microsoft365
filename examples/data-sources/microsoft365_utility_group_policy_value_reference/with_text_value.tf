# Example: Using datasource to discover policy metadata for a text value

# Query the policy definition
data "microsoft365_utility_group_policy_value_reference" "vhd_sddl" {
  policy_name = "Attached VHD SDDL"
}

# Filter for the FSLogix Profile Containers machine policy
locals {
  fslogix_sddl_policy = [
    for def in data.microsoft365_utility_group_policy_value_reference.vhd_sddl.definitions :
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
    display_name   = local.fslogix_sddl_policy.display_name
    class_type     = local.fslogix_sddl_policy.class_type
    category_path  = local.fslogix_sddl_policy.category_path
    explain_text   = local.fslogix_sddl_policy.explain_text
    presentations  = length(local.fslogix_sddl_policy.presentations)
  }
}

