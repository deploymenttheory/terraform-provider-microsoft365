# Query group policy definition metadata
data "microsoft365_utility_group_policy_value_reference" "rdp_allow" {
  policy_name = "Allow users to connect remotely by using Remote Desktop Services"
}

# Output the policy details
output "rdp_policy_info" {
  value = {
    definitions_found = length(data.microsoft365_utility_group_policy_value_reference.rdp_allow.definitions)
    first_definition = {
      id            = data.microsoft365_utility_group_policy_value_reference.rdp_allow.definitions[0].id
      display_name  = data.microsoft365_utility_group_policy_value_reference.rdp_allow.definitions[0].display_name
      class_type    = data.microsoft365_utility_group_policy_value_reference.rdp_allow.definitions[0].class_type
      category_path = data.microsoft365_utility_group_policy_value_reference.rdp_allow.definitions[0].category_path
      policy_type   = data.microsoft365_utility_group_policy_value_reference.rdp_allow.definitions[0].policy_type
    }
  }
}

