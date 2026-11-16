# Custom directory role with every optional attribute set.
# template_id becomes the role definition ID and cannot be changed after creation.
resource "microsoft365_graph_beta_identity_and_access_role_definition" "device_support" {
  display_name = "Device Support"
  description  = "Read devices, groups and applications, and read BitLocker key metadata"
  is_enabled   = true
  template_id  = "d5b6a8c4-3f2e-4a1b-9c7d-2e8f6a4b1c3d"
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
