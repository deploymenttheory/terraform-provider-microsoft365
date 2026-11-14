# Invalid: users and groups are only allowed for enumerated memberships.
resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  azure_ad_join = {
    allowed_to_join = {
      type   = "all"
      groups = ["55555555-5555-5555-5555-555555555555"]
    }
  }
}
