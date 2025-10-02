resource "microsoft365_graph_beta_identity_and_access_authentication_context" "test" {
  id           = "c91"
  display_name = "Initial Context"
  description  = "Initial description"
  is_available = true

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}