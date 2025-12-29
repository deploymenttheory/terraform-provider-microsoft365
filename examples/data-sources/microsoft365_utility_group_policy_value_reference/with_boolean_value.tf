# Example: Using datasource to discover policy metadata for a boolean value

# Query the policy definition
data "microsoft365_utility_group_policy_value_reference" "fslogix_enable" {
  policy_name = "Enable Profile Containers"
}

# Filter for the machine-level policy in the FSLogix category
locals {
  fslogix_machine_policy = [
    for def in data.microsoft365_utility_group_policy_value_reference.fslogix_enable.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "FSLogix\\Profile Containers")
  ][0]
}

# Create a group policy configuration
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "fslogix_config" {
  display_name = "FSLogix Profile Container Configuration"
  description  = "Enables FSLogix Profile Containers"
}

# Create the boolean value using discovered metadata
resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "enable_profile_containers" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.fslogix_config.id
  
  # Use the discovered metadata from the datasource
  policy_name   = local.fslogix_machine_policy.display_name
  class_type    = local.fslogix_machine_policy.class_type
  category_path = local.fslogix_machine_policy.category_path
  enabled       = true

  # This policy has a single boolean value
  values = [
    {
      value = true # Enable Profile Containers
    }
  ]
  
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Output the discovered policy details
output "fslogix_policy_metadata" {
  value = {
    display_name  = local.fslogix_machine_policy.display_name
    class_type    = local.fslogix_machine_policy.class_type
    category_path = local.fslogix_machine_policy.category_path
    policy_type   = local.fslogix_machine_policy.policy_type
    explain_text  = local.fslogix_machine_policy.explain_text
    presentations = length(local.fslogix_machine_policy.presentations)
  }
}

