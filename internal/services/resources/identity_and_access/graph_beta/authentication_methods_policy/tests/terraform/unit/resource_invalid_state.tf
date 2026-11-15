resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  email = {
    state = "on"
  }
}
