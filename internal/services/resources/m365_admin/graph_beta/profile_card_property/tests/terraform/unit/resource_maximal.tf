# Maximal profile card property: hidden, with a custom label and two localized labels.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "maximal" {
  directory_property_name = "CustomAttribute2"
  is_visible              = false

  annotations = [
    {
      display_name = "Cost Center"
      localizations = [
        {
          language_tag = "de-DE"
          display_name = "Kostenstelle"
        },
        {
          language_tag = "nb-NO"
          display_name = "Kostnadssenter"
        }
      ]
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
