resource "random_string" "custom_sec_att_user_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_users_user" "with_custom_security_attributes" {
  display_name        = "acc-test-user-custom-sec-att-${random_string.custom_sec_att_user_id.result}"
  user_principal_name = "acc-test-user-custom-sec-att-${random_string.custom_sec_att_user_id.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-user-custom-sec-att-${random_string.custom_sec_att_user_id.result}"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }

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
