resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-flexible-fic-app"
  description  = "Application using a flexible federated identity credential"
}

# Flexible FIC: match the token with a claims matching expression instead of a
# fixed subject. One wildcard is allowed, so a single credential can cover, for
# example, every environment in a repository. claims_matching_expression and
# subject are mutually exclusive, so subject is omitted here.
resource "microsoft365_graph_beta_applications_application_federated_identity_credential" "claims_matching" {
  application_id             = microsoft365_graph_beta_applications_application.example.id
  name                       = "github-actions-any-environment"
  description                = "GitHub Actions for any environment in the repository"
  issuer                     = "https://token.actions.githubusercontent.com"
  claims_matching_expression = "claims['sub'] matches 'repo:myorg/myrepo:environment:*'"
  audiences                  = ["api://AzureADTokenExchange"]
}
