resource "random_string" "suffix_ud_maximal" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Device Security Group
# ==============================================================================

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client_ud_maximal" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_maximal_device_group" {
  display_name     = "acc-test-autopilot-dpp-ud-maximal-device-group-${random_string.suffix_ud_maximal.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-maximal-dg-${random_string.suffix_ud_maximal.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Device security group for autopilot device preparation policy user-driven maximal acceptance test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "time_sleep" "wait_for_device_group_ud_maximal" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_ud_maximal_device_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "acc_test_ud_maximal_device_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.acc_test_ud_maximal_device_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client_ud_maximal.id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_device_group_ud_maximal]
}

# ==============================================================================
# WinGet App Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "acc_test_ud_maximal_app" {
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

resource "microsoft365_graph_beta_device_management_windows_platform_script" "acc_test_ud_maximal_script" {
  display_name       = "acc-test-autopilot-dpp-ud-maximal-script-${random_string.suffix_ud_maximal.result}"
  description        = "Test script for autopilot device preparation policy user-driven maximal acceptance test"
  role_scope_tag_ids = ["0"]

  script_content = <<EOT
    Write-Host "Autopilot Device Preparation - User-Driven Maximal Test Script"
    Write-Host "Device setup in progress..."
    exit 0
  EOT

  run_as_account          = "system"
  enforce_signature_check = false
  file_name               = "autopilot_ud_maximal_test.ps1"
  run_as_32_bit           = false

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

# Allow Microsoft Graph replication to settle before the policy is created.
resource "time_sleep" "wait_for_dependencies_ud_maximal" {
  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.acc_test_ud_maximal_device_group_owner,
    microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_maximal_app,
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_ud_maximal_script,
  ]

  create_duration = "60s"
}

# ==============================================================================
# Policy Under Test - User-Driven Mode Maximal (No Assignments)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_maximal" {
  name                  = "acc-test-autopilot-dpp-ud-maximal-${random_string.suffix_ud_maximal.result}"
  description           = "User-driven mode maximal acceptance test with enhanced mode features"
  role_scope_tag_ids    = ["0"]
  device_security_group = microsoft365_graph_beta_groups_group.acc_test_ud_maximal_device_group.id

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1"
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0"
    account_type    = "enrollment_autopilot_dpp_accountype_1"
  }

  oobe_settings = {
    timeout_in_minutes   = 120
    custom_error_message = "Please contact your IT administrator for assistance with device setup."
    allow_skip           = true
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_maximal_app.id
      app_type = "winGetApp"
    }
  ]

  allowed_scripts = [
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_ud_maximal_script.id
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_dependencies_ud_maximal]
}
