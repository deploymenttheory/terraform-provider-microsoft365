# ==============================================================================
# Profile Card Property
# ==============================================================================

# Minimal profile card property for the update test: CustomAttribute15 with no custom label.
# CustomAttribute15 must not already be configured in the test tenant; the test checks this first.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "update" {
  directory_property_name = "CustomAttribute15"
}
