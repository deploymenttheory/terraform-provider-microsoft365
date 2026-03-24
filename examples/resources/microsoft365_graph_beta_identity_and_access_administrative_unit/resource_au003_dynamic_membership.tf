# AU003: Dynamic Membership Administrative Unit
# Creates an administrative unit with dynamic membership based on user attributes
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au003_dynamic" {
  display_name                     = "US-Based Users"
  description                      = "Administrative unit for all users located in the United States"
  membership_type                  = "Dynamic"
  membership_rule                  = "(user.country -eq \"United States\")"
  membership_rule_processing_state = "On"
}
