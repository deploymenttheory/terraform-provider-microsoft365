# Deployment audience with members and exclusions — adds updatable asset groups
# as members and specifies exclusion groups. Exclusions take precedence over members.

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

resource "microsoft365_graph_beta_groups_group" "exclusion_group" {
  display_name     = "autopatch-audience-exclusion"
  mail_enabled     = false
  mail_nickname    = "autopatch-exclusion"
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.member_group_1,
    microsoft365_graph_beta_groups_group.member_group_2,
    microsoft365_graph_beta_groups_group.exclusion_group
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

  exclusions = [
    microsoft365_graph_beta_groups_group.exclusion_group.id
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "10m"
  }
}
