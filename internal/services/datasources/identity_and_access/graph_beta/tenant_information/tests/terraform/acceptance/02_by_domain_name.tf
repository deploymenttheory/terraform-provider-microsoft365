data "microsoft365_graph_beta_identity_and_access_tenant_information" "by_domain_name" {
  filter_type  = "domain_name"
  filter_value = "deploymenttheory.com"
}

