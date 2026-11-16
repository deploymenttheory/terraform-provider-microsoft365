# Maximal custom directory role: every optional attribute, a custom template ID and two permission blocks.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "maximal" {
  display_name = "unit-test-role-definition-maximal"
  description  = "Unit test custom role with every optional field set"
  is_enabled   = false
  template_id  = "22222222-2222-2222-2222-222222222222"
  version      = "1"

  role_permissions = [
    {
      allowed_resource_actions = [
        "microsoft.directory/applications/standard/read",
        "microsoft.directory/groups/standard/read"
      ]
    },
    {
      allowed_resource_actions = [
        "microsoft.directory/devices/standard/read",
        "microsoft.directory/bitlockerKeys/metadata/read"
      ]
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
