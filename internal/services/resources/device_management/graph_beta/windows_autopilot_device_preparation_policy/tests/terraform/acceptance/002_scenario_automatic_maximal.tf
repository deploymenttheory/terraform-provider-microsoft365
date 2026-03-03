resource "random_string" "suffix_auto_maximal" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# WinGet App Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "acc_test_auto_maximal_app" {
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

# ==============================================================================
# Windows Platform Script Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_platform_script" "acc_test_auto_maximal_script" {
  display_name       = "acc-test-autopilot-dpp-auto-maximal-script-${random_string.suffix_auto_maximal.result}"
  description        = "Test script for autopilot device preparation policy automatic maximal acceptance test"
  role_scope_tag_ids = ["0"]

  script_content = <<EOT
    Write-Host "Autopilot Device Preparation - Automatic Maximal Test Script"
    Write-Host "Device setup in progress..."
    exit 0
  EOT

  run_as_account          = "system"
  enforce_signature_check = false
  file_name               = "autopilot_auto_maximal_test.ps1"
  run_as_32_bit           = false

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

# Allow Microsoft Graph replication to settle before the policy is created.
resource "time_sleep" "wait_for_dependencies_auto_maximal" {
  depends_on = [
    microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_auto_maximal_app,
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_auto_maximal_script,
  ]

  create_duration = "60s"
}

# ==============================================================================
# Policy Under Test - Automatic Mode Maximal
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "auto_maximal" {
  name               = "acc-test-autopilot-dpp-auto-maximal-${random_string.suffix_auto_maximal.result}"
  description        = "Automatic mode maximal acceptance test with apps and scripts"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1"
  }

  allowed_apps = [
    {
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_auto_maximal_app.id
      app_type = "winGetApp"
    }
  ]

  allowed_scripts = [
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_auto_maximal_script.id
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_dependencies_auto_maximal]
}
