
resource "microsoft365_graph_beta_device_management_group_policy_text_value" "site_url" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example_with_assignments.id

  policy_name   = "AD attribute containing Personal Site URL"
  class_type    = "user"
  enabled       = true
  value         = "wwwHomePage2"
  category_path = "\\Microsoft Office 2016\\Server Settings"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_text_value" "vhd_sddl_fslogix_profile_containers" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example_with_assignments.id

  policy_name   = "Attached VHD SDDL"
  class_type    = "machine"
  enabled       = true
  value         = "D:P(A;;FA;;;BA)(A;;FRFW;;;AU)" // gives Full access for admins, read/write for authenticated users
  category_path = "\\FSLogix\\Profile Containers"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_text_value" "vhd_sddl_fslogix_odfc_containers" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example_with_assignments.id

  policy_name   = "Attached VHD SDDL"
  class_type    = "machine"
  enabled       = true
  value         = "D:P:(A;;GA;;;WD" // gives everyone full access
  category_path = "\\FSLogix\\ODFC Containers"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_text_value" "vhd_access_mode_fslogix_odfc_containers" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.example_with_assignments.id

  policy_name   = "VHD Access Mode"
  class_type    = "machine"
  enabled       = true
  value         = "1" // Values:  0 = Direct Access, 1 = DiffDisk on Network, 2 = Local DiffDisk, 3 = Unique Disk per Session
  category_path = "\\FSLogix\\ODFC Containers"

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}