
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_windows_updates_autopatch_device_registration" "test_004" {
  update_category = "feature"
  entra_device_object_ids = [
    "0243c10a-fb67-4262-b253-fc510717d1dc"
  ]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
