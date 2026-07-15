resource "microsoft365_graph_beta_identity_and_access_network_content_policy" "example" {
  name           = "Content Policy"
  description    = "Global Secure Access content policy managed by Terraform"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_content_policy_rule" "example" {
  content_policy_id  = microsoft365_graph_beta_identity_and_access_network_content_policy.example.id
  name               = "Inspect uploaded and downloaded files"
  description        = "Scan matching files with Microsoft Purview"
  action             = "scanPurview"
  priority           = 101
  status             = "enabled"
  activities         = ["download", "upload"]
  content_types      = ["text/csv", "application/pdf"]
  text_content_types = ["json", "plain", "html", "xml"]

  destinations = [
    {
      type   = "web_category"
      values = ["AlcoholAndTobacco"]
    },
    {
      type   = "fqdn"
      values = ["example.com", "*.example.com"]
    },
    {
      type   = "url"
      values = ["https://example.com/path"]
    }
  ]

  session_types = ["user", "agent"]
}
