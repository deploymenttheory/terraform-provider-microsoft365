resource "random_string" "suffix_ud_max_assign" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Device Security Group
# ==============================================================================

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client_ud_max_assign" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_max_assign_device_group" {
  display_name     = "acc-test-autopilot-dpp-ud-max-assign-device-group-${random_string.suffix_ud_max_assign.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-max-assign-dg-${random_string.suffix_ud_max_assign.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Device security group for autopilot device preparation policy user-driven maximal assignments test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "time_sleep" "wait_for_device_group_ud_max_assign" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_device_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "acc_test_ud_max_assign_device_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_device_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client_ud_max_assign.id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_device_group_ud_max_assign]
}

# ==============================================================================
# Assignment Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_max_assign_group_1" {
  display_name     = "acc-test-autopilot-dpp-ud-max-assign-group-1-${random_string.suffix_ud_max_assign.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-max-assign-g1-${random_string.suffix_ud_max_assign.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test assignment group 1 for autopilot device preparation policy user-driven maximal assignments test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_max_assign_group_2" {
  display_name     = "acc-test-autopilot-dpp-ud-max-assign-group-2-${random_string.suffix_ud_max_assign.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-max-assign-g2-${random_string.suffix_ud_max_assign.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test assignment group 2 for autopilot device preparation policy user-driven maximal assignments test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_max_assign_group_3" {
  display_name     = "acc-test-autopilot-dpp-ud-max-assign-group-3-${random_string.suffix_ud_max_assign.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-max-assign-g3-${random_string.suffix_ud_max_assign.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test assignment group 3 for autopilot device preparation policy user-driven maximal assignments test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# ==============================================================================
# WinGet App Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "acc_test_ud_max_assign_app" {
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

resource "microsoft365_graph_beta_device_management_windows_platform_script" "acc_test_ud_max_assign_script" {
  display_name       = "acc-test-autopilot-dpp-ud-max-assign-script-${random_string.suffix_ud_max_assign.result}"
  description        = "Test script for autopilot device preparation policy user-driven maximal assignments test"
  role_scope_tag_ids = ["0"]

  script_content = <<EOT
    Write-Host "Autopilot Device Preparation - User-Driven Maximal Assignments Test Script"
    Write-Host "Device setup in progress..."
    exit 0
  EOT

  run_as_account          = "system"
  enforce_signature_check = false
  file_name               = "autopilot_ud_max_assign_test.ps1"
  run_as_32_bit           = false

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "60s"
    delete = "60s"
  }
}

# Allow Microsoft Graph replication to settle before the policy is created.
resource "time_sleep" "wait_for_dependencies_ud_max_assign" {
  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.acc_test_ud_max_assign_device_group_owner,
    microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_group_1,
    microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_group_2,
    microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_group_3,
    microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_max_assign_app,
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_ud_max_assign_script,
  ]

  create_duration = "60s"
}

# ==============================================================================
# Policy Under Test - User-Driven Mode with Maximal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_max_assign" {
  name                  = "acc-test-autopilot-dpp-ud-max-assign-${random_string.suffix_ud_max_assign.result}"
  description           = "User-driven mode with maximal assignments and allLicensedUsers acceptance test"
  role_scope_tag_ids    = ["0"]
  device_security_group = microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_device_group.id

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1"
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0"
    account_type    = "enrollment_autopilot_dpp_accountype_1"
  }

  oobe_settings = {
    timeout_in_minutes   = 90
    custom_error_message = "Please contact your IT administrator for assistance."
    allow_skip           = true
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_max_assign_app.id
      app_type = "winGetApp"
    }
  ]

  allowed_scripts = [
    microsoft365_graph_beta_device_management_windows_platform_script.acc_test_ud_max_assign_script.id
  ]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_group_2.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_ud_max_assign_group_3.id
    },
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_dependencies_ud_max_assign]
}
