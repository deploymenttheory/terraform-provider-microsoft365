resource "microsoft365_graph_beta_identity_and_access_network_content_policy_rule" "test" {
  content_policy_id = "00000000-0000-0000-0000-000000000301"
  name              = "invalid-action"
  action            = "deny"
  priority          = 101
  status            = "enabled"
  activities        = ["download"]
  content_types     = ["application/pdf"]
  destinations = [{
    type   = "fqdn"
    values = ["example.com"]
  }]
  session_types = ["user"]
}
