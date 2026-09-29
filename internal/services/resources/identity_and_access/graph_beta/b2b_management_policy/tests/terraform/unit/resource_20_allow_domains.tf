resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "test" {
  display_name            = "unit-test-b2b"
  is_organization_default = false
  definition = [jsonencode({
    "B2BManagementPolicy" : {
      "InvitationsAllowedAndBlockedDomainsPolicy" : {
        "AllowedDomains" : [
          "contoso.example",
          "fabrikam.example"
        ]
      }
    }
  })]
}
