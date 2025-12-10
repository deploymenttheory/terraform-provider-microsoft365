# Maximal Agent User configuration - all fields
resource "microsoft365_graph_beta_agents_agent_user" "test_maximal" {
  # Required fields
  display_name        = "Unit Test Agent User Maximal"
  agent_identity_id   = "11111111-1111-1111-1111-111111111111"
  account_enabled     = true
  user_principal_name = "testagentusermaximal@contoso.onmicrosoft.com"
  mail_nickname       = "testagentusermaximal"
  sponsor_ids         = ["22222222-2222-2222-2222-222222222222", "33333333-3333-3333-3333-333333333333"]
  hard_delete         = true

  # Optional name fields
  given_name = "Test"
  surname    = "AgentUser"

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
}

