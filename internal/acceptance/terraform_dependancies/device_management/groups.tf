# Groups used for acceptance testing
# These groups serve as dependencies.

resource "random_string" "group_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Test Group 1
resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "Acceptance Test Dependency - ${random_string.group_suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s" # Increased for cleanup
  }
}

# Test Group 2 - Device Management Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "Acceptance Test Dependency - ${random_string.group_suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s" # Increased for cleanup
  }
}

# Test Group 3 - Helpdesk Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "Acceptance Test Dependency - ${random_string.group_suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s" # Increased for cleanup
  }
}

# Test Group 4 - Security Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "Acceptance Test Dependency - ${random_string.group_suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s" # Increased for cleanup
  }
}

# Test Group 5 - Microsoft 365 Group - mail-enabled
resource "microsoft365_graph_beta_groups_group" "acc_test_group_5_mail_enabled" {
  display_name     = "acc-test-group-5-mail-enabled-${random_string.group_suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-${random_string.group_suffix.result}"
  mail_enabled     = true
  security_enabled = false
  group_types      = ["Unified"]
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Wait for groups to propagate in Azure AD
resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2,
    microsoft365_graph_beta_groups_group.acc_test_group_3,
    microsoft365_graph_beta_groups_group.acc_test_group_4,
    microsoft365_graph_beta_groups_group.acc_test_group_5_mail_enabled
  ]

  create_duration = "15s"
}

# Outputs for easy reference in tests
output "test_group_ids" {
  description = "Group IDs for use in terms and conditions assignment tests"
  value = {
    group_1 = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    group_2 = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    group_3 = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    group_4 = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    group_5 = microsoft365_graph_beta_groups_group.acc_test_group_5_mail_enabled.id
  }
}