resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_007" {
  group_policy_configuration_id = "config-007"
  policy_name                   = "Test Policy TextBox"
  class_type                    = "machine"
  category_path                 = "\\Test\\TextBox"
  enabled                       = true

  values = [
    { label = "Text Setting", value = "updated-text-value" }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
