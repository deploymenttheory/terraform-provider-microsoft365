resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_007" {
  group_policy_configuration_id = "config-007"
  policy_name                   = "Test Policy MultiText"
  class_type                    = "machine"
  category_path                 = "\\Test\\MultiText"
  enabled                       = true

  values = [
    { label = "Filter List", value = "Driver1\nDriver2" }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
