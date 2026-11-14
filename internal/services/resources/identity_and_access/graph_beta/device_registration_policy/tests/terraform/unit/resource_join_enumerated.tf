resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  azure_ad_join = {
    allowed_to_join = {
      type   = "enumerated"
      groups = ["55555555-5555-5555-5555-555555555555"]
    }
  }
}
