# Automatically install a prerequisite before the application that needs it.
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "runtime" {
  source_id       = "00000000-0000-0000-0000-000000000001" # Application that requires the runtime
  target_id       = "00000000-0000-0000-0000-000000000002" # Prerequisite runtime app
  dependency_type = "autoInstall"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Require an existing prerequisite without installing it automatically.
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "existing_prerequisite" {
  source_id       = "00000000-0000-0000-0000-000000000003"
  target_id       = "00000000-0000-0000-0000-000000000004"
  dependency_type = "detect"
}
