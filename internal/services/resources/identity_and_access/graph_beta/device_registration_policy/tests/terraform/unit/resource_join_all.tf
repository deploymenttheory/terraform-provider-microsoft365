resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  azure_ad_join = {
    allowed_to_join = {
      type = "all"
    }
  }
}
