# AU001: Basic Administrative Unit
# Creates a simple administrative unit with assigned membership
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au001_basic" {
  display_name = "Finance Department"
  description  = "Administrative unit for Finance department users and resources"
}
