resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_014" {
  group_policy_configuration_id = "config-014"
  policy_name                   = "Test Policy Without Presentations"
  class_type                    = "machine"
  category_path                 = "\\Test\\NoPresentations"
  enabled                       = true

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
