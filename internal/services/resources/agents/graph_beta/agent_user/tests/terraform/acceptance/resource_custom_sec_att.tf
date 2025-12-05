resource "random_string" "custom_sec_att_user_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_agent_user" "with_custom_security_attributes" {
  display_name        = "acc-test-agent-user-custom-sec-att-${random_string.custom_sec_att_user_id.result}"
  user_principal_name = "acc-test-agent-user-custom-sec-att-${random_string.custom_sec_att_user_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-agent-user-custom-sec-att-${random_string.custom_sec_att_user_id.result}"
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
