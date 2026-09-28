# Example: Apply a B2B management policy to a service principal and an application
# B2B management policies can only be applied to applications and service principals.
# The provider detects which of the two the directory object is.

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy" "example" {
  display_name = "example-b2b-management-policy"

  definition = [
    jsonencode({
      B2BManagementPolicy = {
        InvitationsAllowedAndBlockedDomainsPolicy = {
          BlockedDomains = ["example.com"]
        }
      }
    })
  ]
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "service_principal" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.example.id
  directory_object_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Object ID of the service principal
}

resource "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment" "application" {
  b2b_management_policy_id = microsoft365_graph_beta_identity_and_access_b2b_management_policy.example.id
  directory_object_id      = "yyyyyyyy-yyyy-yyyy-yyyy-yyyyyyyyyyyy" # Object ID of the application
}
