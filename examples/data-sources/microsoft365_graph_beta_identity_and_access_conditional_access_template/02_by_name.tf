# Example: Get conditional access template by name
# This example demonstrates how to retrieve a conditional access template using its name
# Useful when you know the template name and need to find its ID or configuration details

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "by_name" {
  name = "Require multifactor authentication for admins"

  timeouts = {
    read = "1m"
  }
}

# Output template details
output "template_by_name" {
  value = {
    template_id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.template_id
    name        = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.name
    description = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.description
    scenarios   = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.scenarios
  }
  description = "Conditional access template retrieved by name"
}
