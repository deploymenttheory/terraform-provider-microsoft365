# Example: Discovering and selecting from multiple policy variants

# Query a policy that exists in multiple locations (Chrome, Edge, etc.)
data "microsoft365_utility_group_policy_value_reference" "home_button" {
  policy_name = "Show Home button on toolbar"
}

# Output all discovered variants
output "all_home_button_variants" {
  value = [
    for def in data.microsoft365_utility_group_policy_value_reference.home_button.definitions : {
      id            = def.id
      class_type    = def.class_type
      category_path = def.category_path
      policy_type   = def.policy_type
    }
  ]
  description = "Shows all variants of this policy across Chrome, Edge, and their default settings"
}

# Filter for specific browser and class type
locals {
  # Get Microsoft Edge machine policy
  edge_machine_policy = [
    for def in data.microsoft365_utility_group_policy_value_reference.home_button.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "Microsoft Edge\\Startup")
  ][0]
  
  # Get Google Chrome user policy
  chrome_user_policy = [
    for def in data.microsoft365_utility_group_policy_value_reference.home_button.definitions :
    def if def.class_type == "user" && contains(def.category_path, "Google\\Google Chrome\\Startup")
  ][0]
}

# Output selected variants
output "selected_policies" {
  value = {
    edge_machine = {
      display_name  = local.edge_machine_policy.display_name
      class_type    = local.edge_machine_policy.class_type
      category_path = local.edge_machine_policy.category_path
    }
    chrome_user = {
      display_name  = local.chrome_user_policy.display_name
      class_type    = local.chrome_user_policy.class_type
      category_path = local.chrome_user_policy.category_path
    }
  }
}

# Create configuration for the selected Edge policy
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "edge_config" {
  display_name = "Microsoft Edge Settings"
  description  = "Configure Edge browser home button"
}

resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "edge_home_button" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.edge_config.id

  policy_name   = local.edge_machine_policy.display_name
  class_type    = local.edge_machine_policy.class_type
  category_path = local.edge_machine_policy.category_path
  enabled       = true

  values = [
    {
      value = true # Show home button
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

