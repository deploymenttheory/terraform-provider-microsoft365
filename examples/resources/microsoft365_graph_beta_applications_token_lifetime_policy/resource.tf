# Example: Token lifetime policy with custom access token lifetime
# Configures a token lifetime policy that sets access token lifetime to 1 hour.
# For the JSON structure of the definition, see:
# https://learn.microsoft.com/en-us/entra/identity-platform/configure-token-lifetimes

resource "microsoft365_graph_beta_applications_token_lifetime_policy" "example" {
  display_name = "example-token-lifetime-policy"

  definition = [
    "{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"
  ]

  is_organization_default = false
}

# Example: Organization default token lifetime policy
# Only one organization default policy can exist at a time.

resource "microsoft365_graph_beta_applications_token_lifetime_policy" "org_default" {
  display_name = "organization-default-token-lifetime-policy"

  definition = [
    "{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"
  ]

  is_organization_default = true
}
