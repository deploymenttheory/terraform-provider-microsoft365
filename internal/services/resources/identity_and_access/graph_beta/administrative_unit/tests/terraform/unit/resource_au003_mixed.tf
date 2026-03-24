# ==============================================================================
# AU003: Mixed User and Group Administrative Unit
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au003_mixed" {
  display_name = "AU003: Mixed User and Group Administrative Unit"
  description  = "Administrative unit for mixed user and group testing"
  visibility   = "HiddenMembership"
}
