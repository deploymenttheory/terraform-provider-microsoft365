resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "member_group_1" {
  display_name     = "acc-test-audience-member-1-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-member-1-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "member_group_2" {
  display_name     = "acc-test-audience-member-2-${random_string.suffix.result}"
  mail_enabled     = false
  mail_nickname    = "acc-test-member-2-${random_string.suffix.result}"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.member_group_1,
    microsoft365_graph_beta_groups_group.member_group_2
  ]

  create_duration = "30s"
}

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {
  depends_on = [time_sleep.wait_for_groups]

  member_type = "updatableAssetGroup"

  members = [
    microsoft365_graph_beta_groups_group.member_group_1.id,
    microsoft365_graph_beta_groups_group.member_group_2.id
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "10m"
  }
}
