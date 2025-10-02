resource "microsoft365_graph_beta_identity_and_access_authentication_context" "test" {
  id           = "c92"
  display_name = "Import Test Context"
  description  = "Context for import testing"
  is_available = true
}