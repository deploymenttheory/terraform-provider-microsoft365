# AU006: Dynamic Department-Based Administrative Unit
# Creates an administrative unit that automatically includes users from a specific department
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au006_sales" {
  display_name                     = "Sales Department"
  description                      = "Administrative unit for all Sales department users"
  membership_type                  = "Dynamic"
  membership_rule                  = "(user.department -eq \"Sales\")"
  membership_rule_processing_state = "On"
  visibility                       = "Public"
}
