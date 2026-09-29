resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "test" {
  display_name = "acc-test-b2b-${random_uuid.test.result}"
  definition = [<<-JSON
{
  "B2BManagementPolicy": {
    "InvitationsAllowedAndBlockedDomainsPolicy": {
      "AllowedDomains": [
        "contoso.example",
        "fabrikam.example"
      ]
    }
  }
}
JSON
  ]
}
