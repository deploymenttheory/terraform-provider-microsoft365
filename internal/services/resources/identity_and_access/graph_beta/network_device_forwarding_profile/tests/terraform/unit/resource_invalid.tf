resource "microsoft365_graph_beta_identity_and_access_network_device_forwarding_profile" "test" {
  name     = "unit-invalid"
  priority = 250
  state    = "disabled"

  assignment_rules = {
    membership = "enumerated"
    user_ids   = ["aaaaaaaa-1111-4111-8111-111111111111"]
    platforms  = ["windows"]
  }
}
