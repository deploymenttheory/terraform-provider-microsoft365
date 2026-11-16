# Invalid: is_mdm_enrollment_during_registration_disabled only exists on device policies.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "invalid" {
  policy_type = "app"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  is_mdm_enrollment_during_registration_disabled = true
}
