resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_015" {
  group_policy_configuration_id = "config-015"
  policy_name                   = "Test Policy Presentations Error"
  class_type                    = "machine"
  category_path                 = "\\Test\\PresentationsError"
  enabled                       = true

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
