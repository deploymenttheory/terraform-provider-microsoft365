resource "microsoft365_graph_beta_identity_and_access_directory_settings" "consent_policy_settings" {
  template_type               = "Consent Policy Settings"
  overwrite_existing_settings = true

  consent_policy_settings {
    enable_group_specific_consent                               = true
    block_user_consent_for_risky_apps                           = true
    enable_admin_consent_requests                               = true
    constrain_group_specific_consent_to_members_of_group_id = "87654321-4321-4321-4321-210987654321"
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

