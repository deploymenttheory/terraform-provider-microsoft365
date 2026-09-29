resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "include" {
  display_name     = "acc-test-json-include-${random_string.suffix.result}"
  mail_nickname    = "acc-test-json-include-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}
resource "microsoft365_graph_beta_groups_group" "exclude" {
  display_name     = "acc-test-json-exclude-${random_string.suffix.result}"
  mail_nickname    = "acc-test-json-exclude-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}
resource "time_sleep" "directory_consistency" {
  depends_on      = [microsoft365_graph_beta_groups_group.include, microsoft365_graph_beta_groups_group.exclude]
  create_duration = "30s"
}
