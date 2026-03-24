# ==============================================================================
# AU001: User-Based Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au001_user_based" {
  display_name = "AU001: User-Based Administrative Unit"
  description  = "Administrative unit for user-based testing"
}
