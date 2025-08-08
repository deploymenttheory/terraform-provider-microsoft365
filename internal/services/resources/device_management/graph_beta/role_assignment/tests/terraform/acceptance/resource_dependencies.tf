# Groups used for Role Assignment testing
# These groups serve as dependencies for role assignment tests

# Test Group 1 - Policy Managers
resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "Acceptance Test Dependency - Policy Managers"
  description      = "Test group for policy managers used in role assignments"
  mail_nickname    = "policy-mgrs-test"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 2 - Device Administrators
resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "Acceptance Test Dependency - Device Administrators"
  description      = "Test group for device administrators used in role assignments"
  mail_nickname    = "device-admins-test"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 3 - Application Managers
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "Acceptance Test Dependency - Application Managers"
  description      = "Test group for application managers used in role assignments"
  mail_nickname    = "app-mgrs-test"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Test Group 4 - Security Operators
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "Acceptance Test Dependency - Security Operators"
  description      = "Test group for security operators used in role assignments"
  mail_nickname    = "sec-ops-test"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Outputs for easy reference in tests
output "test_group_ids" {
  description = "Group IDs for use in role assignment tests"
  value = {
    policy_managers       = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    device_administrators = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    application_managers  = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    security_operators    = microsoft365_graph_beta_groups_group.acc_test_group_4.id
  }
}

# Convenience outputs for individual group IDs
output "policy_managers_id" {
  description = "Policy Managers group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_1.id
}

output "device_administrators_id" {
  description = "Device Administrators group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_2.id
}

output "application_managers_id" {
  description = "Application Managers group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_3.id
}

output "security_operators_id" {
  description = "Security Operators group ID"
  value       = microsoft365_graph_beta_groups_group.acc_test_group_4.id
}