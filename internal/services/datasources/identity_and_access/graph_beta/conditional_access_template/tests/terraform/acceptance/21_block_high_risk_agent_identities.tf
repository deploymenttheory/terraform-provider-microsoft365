# Block high risk agent identities from accessing resources
# Template discovery is independent of policy creation and its tenant prerequisites.
data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_high_risk_agents" {
  template_id = "6bdcbbb7-ebd1-4f8f-a02f-652db0e3665d"
}
