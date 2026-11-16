# App policy: the Microsoft Intune MAM user scope, scoped to one group.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "app" {
  policy_type = "app"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  discovery_url = "https://wip.mam.manage.microsoft.com/Enroll"

  applies_to = "selected"
  included_groups = [
    "44444444-4444-4444-4444-444444444444",
  ]
}
