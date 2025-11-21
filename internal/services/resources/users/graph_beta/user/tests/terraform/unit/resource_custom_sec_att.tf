resource "microsoft365_graph_beta_users_user" "with_custom_security_attributes" {
  display_name        = "unit-test-user-custom-sec-att"
  user_principal_name = "unit-test-user-custom-sec-att@deploymenttheory.com"
  mail_nickname       = "unit-test-user-custom-sec-att"
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
