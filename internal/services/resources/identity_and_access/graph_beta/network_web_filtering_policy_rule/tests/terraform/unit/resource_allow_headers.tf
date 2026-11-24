resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule" "test" {
  web_filtering_policy_id = "11111111-2222-3333-4444-555555555555"
  name                    = "unit-test-allow-rule"
  description             = "Allow rule with custom headers"
  priority                = 100
  action                  = "allow"
  status                  = "enabled"
  urls_or_fqdns           = ["www.example.com", "*.contoso.com"]
  web_categories          = ["AlcoholAndTobacco"]
  http_methods            = ["get", "post"]
  session_types           = ["user", "agent"]

  custom_headers = [
    {
      name  = "Restrict-Access-To-Tenants"
      value = "contoso.com"
    },
    {
      name  = "Restrict-Access-Context"
      value = "00000000-0000-0000-0000-000000000000"
    },
  ]
}
