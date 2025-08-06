resource "microsoft365_graph_beta_device_management_role_assignment" "maximal" {
  display_name       = "Test Maximal Role Assignment - Unique"
  description        = "Comprehensive role assignment for testing with all features"
  role_definition_id = "9e0cc482-82df-4ab2-a24c-0c23a3f52e1e" # Help Desk Operator
  
  members = [
    "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
    "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2",
    "35d09841-af73-43e6-a59f-024fef1b6b95"
  ]
  
  scope_configuration {
    type = "ResourceScopes"
    resource_scopes = [
      "ea8e2fb8-e909-44e6-bae7-56757cf6f347",
      "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2",
      "35d09841-af73-43e6-a59f-024fef1b6b95"
    ]
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}