# Lookup device by Entra Device Object ID
data "microsoft365_graph_beta_identity_and_access_device" "test" {
  object_id = "23ace577-ee29-416f-8566-11c948310bff"
}
