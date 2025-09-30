resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "allow_users_to_contact_microsoft_for_feedback_and_support" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example_with_assignments.id
  policy_name                   = "Allow users to contact Microsoft for feedback and support"
  class_type                    = "machine"
  category_path                 = "\\OneDrive"
  enabled                       = true

  values = [
    {
      value = true # Send Feedback - 1st gui value
    },
    {
      value = true # Receive user satisfication surveys - 2nd gui value
    },
    {
      value = false # Contact OneDrive Support​ - 3rd gui value
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
