resource "microsoft365_graph_beta_identity_and_access_directory_settings" "group_unified_guest" {
  template_type               = "Group.Unified.Guest"
  overwrite_existing_settings = true

  group_unified_guest {
    allow_to_add_guests = true
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

