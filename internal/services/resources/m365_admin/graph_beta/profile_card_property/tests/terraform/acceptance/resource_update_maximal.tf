# ==============================================================================
# Random Suffix for Unique Labels
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Profile Card Property
# ==============================================================================

# CustomAttribute15 updated in place: hidden, with a custom label and one localized label.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "update" {
  directory_property_name = "CustomAttribute15"
  is_visible              = false

  annotations = [
    {
      display_name = "acc-test-profile-card-updated-${random_string.suffix.result}"
      localizations = [
        {
          language_tag = "de-DE"
          display_name = "acc-test-profile-card-de-${random_string.suffix.result}"
        }
      ]
    }
  ]
}
