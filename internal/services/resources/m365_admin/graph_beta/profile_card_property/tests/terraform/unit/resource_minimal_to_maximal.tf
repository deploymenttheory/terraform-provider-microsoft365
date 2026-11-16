# Minimal profile card property updated in place: hidden, with a custom label and one localized label.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "minimal" {
  directory_property_name = "CustomAttribute1"
  is_visible              = false

  annotations = [
    {
      display_name = "Cost Center"
      localizations = [
        {
          language_tag = "de-DE"
          display_name = "Kostenstelle"
        }
      ]
    }
  ]
}
