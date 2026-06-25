# Example: Assign a token lifetime policy to a service principal
# Only one token lifetime policy can be assigned to a service principal at a time.
# The policy must be created separately using the token_lifetime_policy resource.

resource "microsoft365_graph_beta_applications_token_lifetime_policy" "example" {
  display_name = "example-token-lifetime-policy"

  definition = [
    "{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"
  ]
}

resource "microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment" "example" {
  service_principal_id     = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Object ID of the service principal
  token_lifetime_policy_id = microsoft365_graph_beta_applications_token_lifetime_policy.example.id
}
