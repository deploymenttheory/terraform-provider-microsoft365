resource "microsoft365_graph_beta_identity_and_access_external_identities_policy" "test" {
  allow_external_identities_to_leave    = true
  allow_deleted_identities_data_removal = true
}
