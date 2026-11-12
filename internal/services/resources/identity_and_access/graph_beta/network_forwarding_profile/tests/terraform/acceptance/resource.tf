resource "random_pet" "suffix" { length = 2 }

resource "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "test" {
  name                    = "tfacc-fp-${random_pet.suffix.id}"
  description             = "acceptance test"
  traffic_forwarding_type = "private"
  priority                = var.priority
  state                   = "disabled"

  local_network_settings = {
    allow_local_subnet_bypass = true
  }

  forwarding_source_matching = {
    session_type = "any"
  }

  assignment_rules = {
    membership = "all"
    platforms  = ["windows", "macOS"]
  }
}
