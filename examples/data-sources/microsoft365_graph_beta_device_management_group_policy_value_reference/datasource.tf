# Query group policy definition metadata
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "rdp_allow" {
  policy_name = "Allow users to connect remotely by using Remote Desktop Services"
}

# Output all available attributes from the data source
output "rdp_policy_full_details" {
  value = {
    # Number of definitions found
    definitions_count = length(data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions)

    # Example of all definitions and their attributes
    definitions = [
      for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions : {
        # Definition identification
        id           = def.id
        display_name = def.display_name

        # Policy classification
        class_type    = def.class_type
        category_path = def.category_path
        policy_type   = def.policy_type

        # Policy documentation
        explain_text = def.explain_text
        supported_on = def.supported_on

        # Presentation count (handle null)
        presentations_count = try(length(def.presentations), 0)

        # All presentations with complete attributes (handle null)
        presentations = try([
          for pres in def.presentations : {
            id       = pres.id
            label    = pres.label
            type     = pres.type
            required = pres.required
          }
        ], [])
      }
    ]
  }
}

# Example output for the first definition only (simplified)
output "rdp_policy_first_definition" {
  value = length(data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions) > 0 ? {
    id               = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].id
    display_name     = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].display_name
    class_type       = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].class_type
    category_path    = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].category_path
    policy_type      = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].policy_type
    explain_text     = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].explain_text
    supported_on     = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].supported_on
    presentations    = data.microsoft365_graph_beta_device_management_group_policy_value_reference.rdp_allow.definitions[0].presentations
  } : null
}
