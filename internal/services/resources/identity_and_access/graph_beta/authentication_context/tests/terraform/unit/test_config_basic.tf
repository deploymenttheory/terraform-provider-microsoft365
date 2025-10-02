resource "microsoft365_graph_beta_identity_and_access_authentication_context" "test" {
  id           = "c90"
  display_name = "Test Authentication Context"
  description  = "Test authentication context for unit testing"
  is_available = true

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}