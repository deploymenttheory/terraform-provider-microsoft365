# Use a heredoc when maintaining the definition as JSON rather than HCL.
# JSON formatting and whitespace do not cause a recurring Terraform diff.
resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "heredoc" {
  display_name            = "B2B invitations - JSON definition"
  is_organization_default = false

  definition = [<<-JSON
    {
      "B2BManagementPolicy": {
        "InvitationsAllowedAndBlockedDomainsPolicy": {
          "AllowedDomains": [
            "example.com",
            "example.net"
          ]
        }
      }
    }
  JSON
  ]
}
