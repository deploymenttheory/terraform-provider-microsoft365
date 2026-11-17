# ==============================================================================
# Profile Card Property
# ==============================================================================

# Invalid directory property name: CustomAttribute16 is not a profile card attribute.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "invalid" {
  directory_property_name = "CustomAttribute16"
}
