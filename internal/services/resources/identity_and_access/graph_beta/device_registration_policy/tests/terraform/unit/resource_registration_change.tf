# Invalid for a tenant where Intune manages registration (isAdminConfigurable = false).
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  azure_ad_registration = {
    allowed_to_register = {
      type = "none"
    }
  }
}
