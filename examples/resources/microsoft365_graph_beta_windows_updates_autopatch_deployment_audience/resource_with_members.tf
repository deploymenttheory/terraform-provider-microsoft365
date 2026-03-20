# Deployment audience with members — adds updatable asset groups as members.
# A 30-second time_sleep ensures the groups have fully propagated before assignment.

resource "microsoft365_graph_beta_groups_group" "member_group_1" {
  display_name     = "autopatch-audience-member-1"
  mail_enabled     = false
  mail_nickname    = "autopatch-member-1"
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "member_group_2" {
  display_name     = "autopatch-audience-member-2"
  mail_enabled     = false
  mail_nickname    = "autopatch-member-2"
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

resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "example" {
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
