# Maximal Agent User Example
# This example shows all available fields for creating an agent user

resource "microsoft365_graph_beta_agents_agent_user" "example" {
  # Required fields
  display_name        = "Example Agent User"
  agent_identity_id   = "00000000-0000-0000-0000-000000000000" # ID of parent agent identity
  account_enabled     = true
  user_principal_name = "agent-user@contoso.com" # Must match your tenant's verified domain
  mail_nickname       = "agent-user"
  sponsor_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]
  hard_delete = true

  # Optional name fields
  given_name = "Agent"
  surname    = "User"

  # Optional organizational fields
  job_title       = "AI Agent"
  department      = "Engineering"
  company_name    = "Contoso"
  office_location = "Building A"

  # Optional address fields
  city           = "Seattle"
  state          = "WA"
  country        = "US"
  postal_code    = "98101"
  street_address = "123 Main Street"

  # Optional locale fields
  usage_location     = "US"
  preferred_language = "en-US"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

