# ==============================================================================
# Profile Card Property
# ==============================================================================

# Minimal profile card property: surface CustomAttribute13 on the profile card with no custom label.
# CustomAttribute13 must not already be configured in the test tenant; the test checks this first.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "minimal" {
  directory_property_name = "CustomAttribute13"
}
