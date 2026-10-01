resource "microsoft365_graph_beta_identity_and_access_external_identities_policy" "example" {
  allow_external_identities_to_leave    = true
  allow_deleted_identities_data_removal = false

  timeouts = {
    create = "5m"
    read   = "1m"
    update = "5m"
  }
}
