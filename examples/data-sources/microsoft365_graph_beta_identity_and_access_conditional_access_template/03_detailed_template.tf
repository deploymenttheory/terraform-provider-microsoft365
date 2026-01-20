# Example: Get detailed conditional access template configuration
# This example demonstrates how to retrieve a template and access its detailed configuration
# including conditions, grant controls, and session controls

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "detailed" {
  name = "Block legacy authentication"

  timeouts = {
    read = "1m"
  }
}

# Output detailed template configuration
output "template_details" {
  value = {
    template_id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.template_id
    name        = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.name
    description = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.description
    scenarios   = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.scenarios
    details     = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details
  }
  description = "Detailed conditional access template configuration"
  sensitive   = false
}

# Example: Use template data to understand policy structure
# You can reference specific parts of the template details
output "template_grant_controls" {
  value = {
    operator          = try(data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details.grant_controls.operator, null)
    built_in_controls = try(data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details.grant_controls.built_in_controls, null)
  }
  description = "Grant controls from the template"
}

output "template_conditions" {
  value = {
    client_app_types = try(data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details.conditions.client_app_types, null)
  }
  description = "Conditions from the template"
}
