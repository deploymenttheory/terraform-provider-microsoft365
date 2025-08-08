resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "all_users_assignment" {
  name         = "Test All Users Assignment Settings Catalog Policy - Unit"
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
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "include"
      filter_id   = "5555555-5555-5555-5555-555555555555"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}