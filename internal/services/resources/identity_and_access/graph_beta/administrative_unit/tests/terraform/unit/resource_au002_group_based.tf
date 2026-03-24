# ==============================================================================
# AU002: Group-Based Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au002_group_based" {
  display_name = "AU002: Group-Based Administrative Unit"
  description  = "Administrative unit for group-based testing"
}
