action "microsoft365_graph_beta_device_management_managed_device_move_devices_to_ou" "minimal" {
  config {
    organizational_unit_path = "OU=Workstations,DC=contoso,DC=com"
    managed_device_ids       = ["00000000-0000-0000-0000-000000000001"]
  }
}

