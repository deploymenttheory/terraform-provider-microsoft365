resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_005" {
  group_policy_configuration_id = "config-005"
  policy_name                   = "Test Policy MultiText"
  class_type                    = "machine"
  category_path                 = "\\Test\\MultiText"
  enabled                       = true

  values = [
    {
      label = "Filter List"
      value = "FilterDriver1\nFilterDriver2\nFilterDriver3"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
