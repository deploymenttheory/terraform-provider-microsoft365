data "microsoft365_graph_beta_identity_and_access_tenant_information" "by_tenant_id" {
  filter_type  = "tenant_id"
  filter_value = "6babcaad-604b-40ac-a9d7-9fd97c0b779f"
}

