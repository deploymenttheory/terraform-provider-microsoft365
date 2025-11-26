resource "microsoft365_graph_beta_identity_and_access_directory_settings" "application" {
  template_type               = "Application"
  overwrite_existing_settings = true

  application {
    enable_access_check_for_privileged_application_updates = true
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

