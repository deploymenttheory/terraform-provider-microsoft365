resource "microsoft365_graph_beta_identity_and_access_authentication_methods_policy" "test" {
  temporary_access_pass = {
    state          = "enabled"
    default_length = 4
  }
}
