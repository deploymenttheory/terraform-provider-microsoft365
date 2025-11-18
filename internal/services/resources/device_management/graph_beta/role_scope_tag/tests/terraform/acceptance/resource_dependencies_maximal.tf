# Groups used for Role Scope Tag maximal configuration testing
# Only creates groups 3 and 4 which are assigned in resource_maximal.tf

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

resource "random_string" "group_suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

# Test Group 3 - Helpdesk Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "acc-test-rst-hd-${random_string.group_suffix.result}"
  description      = "Test group for helpdesk staff used in role scope tag assignments"
  mail_nickname    = "rst-hd-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 4 - Security Team
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "acc-test-rst-sec-${random_string.group_suffix.result}"
  description      = "Test group for security staff used in role scope tag assignments"
  mail_nickname    = "rst-sec-${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

