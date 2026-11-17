# Minimal profile card property: show the Alias attribute on users' profile cards.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "alias" {
  directory_property_name = "Alias"
}
