# Profile card property with a custom label and no localized labels.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "label_only" {
  directory_property_name = "CustomAttribute3"

  annotations = [
    {
      display_name = "Employee Number"
    }
  ]
}
