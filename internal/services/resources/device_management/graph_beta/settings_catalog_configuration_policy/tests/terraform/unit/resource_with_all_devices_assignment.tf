resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "all_devices_assignment" {
  name         = "Test All Devices Assignment Settings Catalog Policy - Unit"
  platforms    = "macOS"
  technologies = ["mdm", "appleRemoteManagement"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    settings = []
  }

  assignments = [
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "include"
      filter_id   = "6666666-6666-6666-6666-666666666666"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}