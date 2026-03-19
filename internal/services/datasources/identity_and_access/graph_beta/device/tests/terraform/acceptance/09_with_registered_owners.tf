# Lookup device and its registered owners
data "microsoft365_graph_beta_identity_and_access_device" "test" {
  object_id              = "23ace577-ee29-416f-8566-11c948310bff"
  list_registered_owners = true
}
