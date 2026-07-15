resource "microsoft365_graph_beta_identity_and_access_network_content_policy_rule" "test" {
  content_policy_id  = "00000000-0000-0000-0000-000000000301"
  name               = "unit-test-content-policy-rule-updated"
  description        = "updated by Terraform"
  action             = "block"
  priority           = 102
  status             = "disabled"
  activities         = ["upload"]
  content_types      = ["application/pdf"]
  text_content_types = ["html"]
  destinations = [
    {
      type   = "fqdn"
      values = ["updated.example.com"]
    }
  ]
  session_types = ["user"]
}
