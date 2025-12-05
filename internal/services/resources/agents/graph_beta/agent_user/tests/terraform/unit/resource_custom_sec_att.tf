resource "microsoft365_graph_beta_users_agent_user" "with_custom_security_attributes" {
  display_name        = "unit-test-agent-user-custom-sec-att"
  user_principal_name = "unit-test-agent-user-custom-sec-att@deploymenttheory.com"
  mail_nickname       = "unit-test-agent-user-custom-sec-att"
  account_enabled     = true
  identity_parent_id  = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

  custom_security_attributes = [
    {
      attribute_set = "Engineering"
      attributes = [
        {
          name          = "Project"
          string_values = ["Baker", "Cascade"]
        },
        {
          name         = "LastTrainingDate"
          string_value = "2024-10-15"
        },
      ]
    },
    {
      attribute_set = "Marketing"
      attributes = [
        {
          name       = "IsContractor"
          bool_value = false
        }
      ]
    }
  ]
}
