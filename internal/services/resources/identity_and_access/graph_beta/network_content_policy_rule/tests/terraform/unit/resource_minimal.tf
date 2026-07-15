resource "microsoft365_graph_beta_identity_and_access_network_content_policy_rule" "test" {
  content_policy_id  = "00000000-0000-0000-0000-000000000301"
  name               = "unit-test-content-policy-rule"
  description        = "managed by Terraform"
  action             = "scanPurview"
  priority           = 101
  status             = "enabled"
  activities         = ["download", "upload"]
  content_types      = ["text/csv", "application/pdf"]
  text_content_types = ["json", "plain"]
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
