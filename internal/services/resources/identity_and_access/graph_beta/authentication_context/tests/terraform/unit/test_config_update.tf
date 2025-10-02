resource "microsoft365_graph_beta_identity_and_access_authentication_context" "test" {
  id           = "c91"
  display_name = "Updated Test Authentication Context"
  description  = "Updated test authentication context for unit testing"
  is_available = false

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}