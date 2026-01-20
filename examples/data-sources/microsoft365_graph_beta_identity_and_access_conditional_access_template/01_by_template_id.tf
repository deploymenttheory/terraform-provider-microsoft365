# Example: Get conditional access template by template ID
# This example demonstrates how to retrieve a conditional access template using its ID (GUID)
# Useful when you know the template ID and need to retrieve its configuration details

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "by_id" {
  template_id = "c7503427-338e-4c5e-902d-abe252abfb43" # Require multifactor authentication for admins

  timeouts = {
    read = "1m"
  }
}

# Output template details
output "template_by_id" {
  value = {
    template_id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.template_id
    name        = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.name
    description = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.description
    scenarios   = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.scenarios
  }
  description = "Conditional access template retrieved by template ID"
}
