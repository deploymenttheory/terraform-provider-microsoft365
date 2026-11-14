# Manage only the device quota. Every other setting keeps its current tenant value.
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "example" {
  user_device_quota = 50
}
