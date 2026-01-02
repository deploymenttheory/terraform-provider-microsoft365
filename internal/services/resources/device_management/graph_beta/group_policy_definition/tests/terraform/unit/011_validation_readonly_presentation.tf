resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_011" {
  group_policy_configuration_id = "config-010"
  policy_name                   = "Test Policy With Read-Only Text"
  class_type                    = "machine"
  category_path                 = "\\Test\\ReadOnly"
  enabled                       = true

  values = [
    {
      label = "Read-Only Text Label"
      value = "cannot set value on text presentation"
    }
  ]
}
