# Minimal configuration: only the device quota is managed.
# Every other setting keeps its current tenant value.
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  user_device_quota = 49
}
