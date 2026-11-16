# Not found: a policy ID with no mobility management application behind it.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "missing" {
  policy_type = "device"
  policy_id   = "99999999-9999-9999-9999-999999999999"
}
