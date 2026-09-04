# Creates a device preparation policy and looks it up by name.

resource "random_string" "suffix_by_name" {
  length  = 8
  special = false
  upper   = false
}

data "microsoft365_graph_beta_applications_service_principal" "intune_provisioning_client_by_name" {
  app_id = "f1346770-5b25-470b-88bd-d5744ab7952c"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_by_name_device_group" {
  display_name     = "acc-test-autopilot-dpp-ds-name-device-group-${random_string.suffix_by_name.result}"
  mail_nickname    = "acc-test-autopilot-dpp-ds-name-dg-${random_string.suffix_by_name.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Device security group for the device preparation policy data source acceptance test"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "acc_test_by_name_device_group_owner" {
  group_id          = microsoft365_graph_beta_groups_group.acc_test_by_name_device_group.id
  owner_id          = data.microsoft365_graph_beta_applications_service_principal.intune_provisioning_client_by_name.id
  owner_object_type = "ServicePrincipal"
}

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  name               = "acc-test-autopilot-dpp-ds-name-${random_string.suffix_by_name.result}"
  description        = "Device preparation policy data source acceptance test - lookup by name"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1"
  }

  device_security_group = microsoft365_graph_beta_groups_group.acc_test_by_name_device_group.id

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }

  depends_on = [
    microsoft365_graph_beta_groups_group_owner_assignment.acc_test_by_name_device_group_owner,
  ]
}

data "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "test" {
  name = microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy.test.name
}
