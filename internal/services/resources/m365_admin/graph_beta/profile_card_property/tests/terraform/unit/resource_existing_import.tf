# Pre-existing profile card property stored with lower-case casing, imported by its exact key.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "existing" {
  directory_property_name = "customAttribute9"

  annotations = [
    {
      display_name = "Department Code"
    }
  ]
}
