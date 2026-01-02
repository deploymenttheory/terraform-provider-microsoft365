# Example: Full dependency chain showing datasource -> configuration -> boolean value

# Step 1: Query the policy definition to discover metadata
data "microsoft365_graph_beta_device_management_group_policy_value_reference" "onedrive_feedback" {
  policy_name = "Allow users to contact Microsoft for feedback and support"
}

# Extract the machine-level policy details
locals {
  feedback_policy = [
    for def in data.microsoft365_graph_beta_device_management_group_policy_value_reference.onedrive_feedback.definitions :
    def if def.class_type == "machine" && contains(def.category_path, "OneDrive")
  ][0]

  # Number of boolean presentations (checkboxes) for this policy
  presentation_count = length(local.feedback_policy.presentations)
}

# Step 2: Create the group policy configuration
resource "microsoft365_graph_beta_device_management_group_policy_configuration" "onedrive_config" {
  display_name = "OneDrive Feedback Configuration"
  description  = "Configure OneDrive user feedback and support options"
}

# Step 3: Create the boolean value using discovered metadata
resource "microsoft365_graph_beta_device_management_group_policy_boolean_value" "onedrive_feedback_settings" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.onedrive_config.id

  # Use the discovered metadata from the datasource
  policy_name   = local.feedback_policy.display_name
  class_type    = local.feedback_policy.class_type
  category_path = local.feedback_policy.category_path
  enabled       = true

  # This policy has 3 boolean values (Send Feedback, Surveys, Contact Support)
  values = [
    {
      value = true # Send Feedback
    },
    {
      value = true # Receive user satisfaction surveys
    },
    {
      value = false # Contact OneDrive Support
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
output "policy_metadata" {
  value = {
    display_name       = local.feedback_policy.display_name
    class_type         = local.feedback_policy.class_type
    category_path      = local.feedback_policy.category_path
    policy_type        = local.feedback_policy.policy_type
    presentation_count = local.presentation_count
    presentations = [
      for pres in local.feedback_policy.presentations : {
        label = pres.label
        type  = pres.presentation_type
      }
    ]
  }
}

