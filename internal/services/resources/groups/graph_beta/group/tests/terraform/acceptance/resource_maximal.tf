resource "random_string" "group_suffix" {
  length  = 8
  special = false
}

resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-group-updated-${random_string.group_suffix.result}"
  description      = "Updated description for acceptance testing"
  mail_nickname    = "acctestgroup${random_string.group_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
}
