resource "random_id" "listbox" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "listbox" {
  display_name = "listbox-${random_id.listbox.hex}"
  description  = "Site to zone assignments"
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "listbox" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.listbox.id
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

