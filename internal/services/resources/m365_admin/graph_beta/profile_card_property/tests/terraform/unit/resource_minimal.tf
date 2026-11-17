# Minimal profile card property: surface CustomAttribute1 on the profile card with no custom label.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "minimal" {
  directory_property_name = "CustomAttribute1"
}
