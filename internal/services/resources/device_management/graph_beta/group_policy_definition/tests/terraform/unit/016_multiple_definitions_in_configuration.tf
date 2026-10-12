resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_016_textbox" {
  group_policy_configuration_id = "config-016"
  policy_name                   = "Test Policy TextBox"
  class_type                    = "machine"
  category_path                 = "\\Test\\TextBox"
  enabled                       = true

  values = [
    {
      label = "Text Setting"
      value = "test-value-string"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_016_decimal" {
  group_policy_configuration_id = "config-016"
  policy_name                   = "Test Policy Decimal"
  class_type                    = "machine"
  category_path                 = "\\Test\\Decimal"
  enabled                       = true

  values = [
    {
      label = "Timeout Setting"
      value = "7200"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
