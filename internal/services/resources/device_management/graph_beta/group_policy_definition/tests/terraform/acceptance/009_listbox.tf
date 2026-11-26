resource "random_id" "test_009" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_009" {
  display_name = "acc-test-gpd-listbox-4905-${random_id.test_009.hex}"
  description  = "Owned ListBox acceptance fixture for issue 4905"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_009" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_009.id
  policy_name                   = "Site to Zone Assignment List"
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\Internet Explorer\\Internet Control Panel\\Security Page"
  enabled                       = true
  values = [{
    label = "Enter the zone assignments here."
    value = jsonencode([{ name = "https://example.test", value = "2" }, { name = "https://other.test", value = "1" }])
  }]
  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
