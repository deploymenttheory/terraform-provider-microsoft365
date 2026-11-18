resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-4874-group-${random_string.suffix.result}"
  mail_nickname    = "acctest4874${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Group import default preservation"
  hard_delete      = true
}
