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

# Maximal profile card property: hidden, with a custom label and two localized labels.
# CustomAttribute14 must not already be configured in the test tenant; the test checks this first.
resource "microsoft365_graph_beta_m365_admin_profile_card_property" "maximal" {
  directory_property_name = "CustomAttribute14"
  is_visible              = false

  annotations = [
    {
      display_name = "acc-test-profile-card-maximal-${random_string.suffix.result}"
      localizations = [
        {
          language_tag = "de-DE"
          display_name = "acc-test-profile-card-de-${random_string.suffix.result}"
        },
        {
          language_tag = "nb-NO"
          display_name = "acc-test-profile-card-nb-${random_string.suffix.result}"
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
