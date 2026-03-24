# ==============================================================================
# AU004: Dynamic Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au004_dynamic" {
  display_name                     = "AU004: Dynamic Administrative Unit"
  description                      = "Administrative unit with dynamic membership"
  membership_type                  = "Dynamic"
  membership_rule                  = "(user.country -eq \"United States\")"
  membership_rule_processing_state = "On"
}
