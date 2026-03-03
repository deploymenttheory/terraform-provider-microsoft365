resource "random_string" "suffix_ud_minimal" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Device Security Group
# ==============================================================================

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client_ud_minimal" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_ud_minimal_device_group" {
  display_name     = "acc-test-autopilot-dpp-ud-minimal-device-group-${random_string.suffix_ud_minimal.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ud-minimal-dg-${random_string.suffix_ud_minimal.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Device security group for autopilot device preparation policy user-driven minimal acceptance test"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

resource "time_sleep" "wait_for_device_group_ud_minimal" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_ud_minimal_device_group,
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "acc_test_ud_minimal_device_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.acc_test_ud_minimal_device_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client_ud_minimal.id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_device_group_ud_minimal]
}

# ==============================================================================
# WinGet App Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "acc_test_ud_minimal_app" {
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
resource "time_sleep" "wait_for_dependencies_ud_minimal" {
  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.acc_test_ud_minimal_device_group_owner,
    microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_minimal_app,
  ]

  create_duration = "60s"
}

# ==============================================================================
# Policy Under Test - User-Driven Mode Minimal (No Assignments)
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_minimal" {
  name                  = "acc-test-autopilot-dpp-ud-minimal-${random_string.suffix_ud_minimal.result}"
  description           = "User-driven mode minimal acceptance test"
  role_scope_tag_ids    = ["0"]
  device_security_group = microsoft365_graph_beta_groups_group.acc_test_ud_minimal_device_group.id

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
      app_id   = microsoft365_graph_beta_device_and_app_management_win_get_app.acc_test_ud_minimal_app.id
      app_type = "winGetApp"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [time_sleep.wait_for_dependencies_ud_minimal]
}
