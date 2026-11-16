# ==============================================================================
# Mobility Management Policy - adopt without changes
# ==============================================================================

# Adopts the Microsoft Intune MDM user scope without configuring any setting, so nothing is written.
# Never set applies_to here: the MDM scope controls automatic Windows enrollment for real users.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "minimal" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"
}
