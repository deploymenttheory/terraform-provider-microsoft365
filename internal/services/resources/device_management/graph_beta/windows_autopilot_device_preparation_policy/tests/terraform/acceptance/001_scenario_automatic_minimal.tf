resource "random_string" "suffix_auto_minimal" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# WinGet App Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "acc_test_auto_minimal_app" {
  package_identifier              = "9NZVDKPMR9RD"
  automatically_generate_metadata = true

  install_experience = {
    run_as_account = "user"
  }

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

# Allow Microsoft Graph replication to settle before the policy is created.
resource "time_sleep" "wait_for_app_auto_minimal" {
  depends_on = [
    microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_auto_minimal_app,
  ]

  create_duration = "60s"
}

# ==============================================================================
# Policy Under Test - Automatic Mode Minimal
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "auto_minimal" {
  name               = "acc-test-autopilot-dpp-auto-minimal-${random_string.suffix_auto_minimal.result}"
  description        = "Windows Autopilot device preparation policy - automatic mode minimal acceptance test"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1"
  }

  allowed_apps = [
    {
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_auto_minimal_app.id
      app_type = "winGetApp"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_app_auto_minimal]
}
