resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "exclusion_assignment" {
  name         = "Test Exclusion Assignment Settings Catalog Policy - Unit"
  description  = ""
  platforms    = "macOS"
  technologies = ["mdm", "appleRemoteManagement"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    name     = "Test Exclusion Assignment Settings Catalog Policy Configuration - Unit"
    settings = []
  }

  assignments = [
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "77777777-7777-7777-7777-777777777777"
      filter_type = "include"
      filter_id   = "88888888-8888-8888-8888-888888888888"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}