resource "microsoft365_graph_beta_identity_and_access_directory_settings" "custom_policy_settings" {
  template_type               = "Custom Policy Settings"
  overwrite_existing_settings = true

  custom_policy_settings {
    custom_conditional_access_policy_url = "https://contoso.com/custom-ca-policy"
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

