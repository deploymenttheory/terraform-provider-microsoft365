# Lookup device by display name
data "microsoft365_graph_beta_identity_and_access_device" "test" {
  display_name = "DT-000481110457"
}
