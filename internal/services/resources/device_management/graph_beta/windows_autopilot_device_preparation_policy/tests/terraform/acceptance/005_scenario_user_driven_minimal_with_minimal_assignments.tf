resource "random_string" "suffix_ud_min_assign" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Device Security Group
# ==============================================================================

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client_ud_min_assign" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_min_assign_device_group" {
  display_name     = "acc-test-autopilot-dpp-ud-min-assign-device-group-${random_string.suffix_ud_min_assign.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-min-assign-dg-${random_string.suffix_ud_min_assign.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Device security group for autopilot device preparation policy user-driven minimal with assignments test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "time_sleep" "wait_for_device_group_ud_min_assign" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_ud_min_assign_device_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "acc_test_ud_min_assign_device_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.acc_test_ud_min_assign_device_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client_ud_min_assign.id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_device_group_ud_min_assign]
}

# ==============================================================================
# Assignment Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_min_assign_group_1" {
  display_name     = "acc-test-autopilot-dpp-ud-min-assign-group-1-${random_string.suffix_ud_min_assign.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-min-assign-g1-${random_string.suffix_ud_min_assign.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test assignment group for autopilot device preparation policy user-driven minimal with assignments test"
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

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "acc_test_ud_min_assign_app" {
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
resource "time_sleep" "wait_for_dependencies_ud_min_assign" {
  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.acc_test_ud_min_assign_device_group_owner,
    microsoft365_graph_beta_groups_group.acc_test_ud_min_assign_group_1,
    microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_min_assign_app,
  ]

  create_duration = "60s"
}

# ==============================================================================
# Policy Under Test - User-Driven Mode Minimal with Minimal Assignments
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_min_assign" {
  name                  = "acc-test-autopilot-dpp-ud-min-assign-${random_string.suffix_ud_min_assign.result}"
  description           = "User-driven mode minimal with minimal assignments acceptance test"
  role_scope_tag_ids    = ["0"]
  device_security_group = microsoft365_graph_beta_groups_group.acc_test_ud_min_assign_device_group.id

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_0"
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0"
    account_type    = "enrollment_autopilot_dpp_accountype_0"
  }

  oobe_settings = {
    timeout_in_minutes   = 60
    custom_error_message = "Contact your organization's support person for help."
    allow_skip           = false
    allow_diagnostics    = false
  }

  allowed_apps = [
    {
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_min_assign_app.id
      app_type = "winGetApp"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_ud_min_assign_group_1.id
    },
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_dependencies_ud_min_assign]
}
