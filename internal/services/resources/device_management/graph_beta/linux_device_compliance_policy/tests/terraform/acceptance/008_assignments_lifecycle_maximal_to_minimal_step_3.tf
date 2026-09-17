resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# All assignment targets are newly created empty groups.
resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_1" {
  display_name     = "acc-test-linux-compliance-008-include-${random_string.test_suffix.result}"
  mail_nickname    = "acc-linux-008-include-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_2" {
  display_name     = "acc-test-linux-compliance-008-second-${random_string.test_suffix.result}"
  mail_nickname    = "acc-linux-008-second-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_008_3" {
  display_name     = "acc-test-linux-compliance-008-exclude-${random_string.test_suffix.result}"
  mail_nickname    = "acc-linux-008-exclude-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "test_008" {
  name                       = "acc-test-linux-compliance-008-${random_string.test_suffix.result}"
  device_encryption_required = false
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "5m"
  }
}
