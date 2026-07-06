resource "microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy" "example" {
  name           = "Web Content Filtering Policy"
  description    = "Global Secure Access web content filtering policy managed by Terraform"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy_rule" "example" {
  web_content_filtering_policy_id = microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy.example.id

  name        = "Example Web Content Filtering Rule"
  description = "Allow matching web traffic"
  priority    = 100
  action      = "allow"
  status      = "enabled"

  urls_or_fqdns  = ["*.example.com"]
  web_categories = ["AlcoholAndTobacco"]
  http_methods   = ["get"]
  session_types  = ["user", "agent"]
}

resource "microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy_rule" "with_headers" {
  web_content_filtering_policy_id = microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy.example.id

  name        = "Allow With Custom Headers"
  description = "Allow matching web traffic and add custom headers"
  priority    = 200
  action      = "allow"
  status      = "enabled"

  urls_or_fqdns = ["headers.example.com"]

  custom_headers = [
    {
      header_name  = "X-Managed-By"
      header_value = "Terraform"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy_rule" "category_only" {
  web_content_filtering_policy_id = microsoft365_graph_beta_identity_and_access_network_web_content_filtering_policy.example.id

  name        = "Block AI Agents Category"
  description = "Block traffic that matches a selected web category"
  priority    = 300
  action      = "block"
  status      = "enabled"

  web_categories = ["AIAgents"]
  session_types  = ["user", "agent"]
}
