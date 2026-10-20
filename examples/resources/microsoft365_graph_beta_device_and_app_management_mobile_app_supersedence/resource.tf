# This example demonstrates how to create a supersedence relationship between two mobile apps in Intune
# The supersedence relationship allows you to upgrade or replace an older app with a newer version

# Example 1: Create an "update" supersedence relationship between two apps
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "update_example" {
  source_id         = "00000000-0000-0000-0000-000000000001" # ID of the newer, superseding app
  target_id         = "00000000-0000-0000-0000-000000000002" # ID of the older, superseded app
  supersedence_type = "update"                               # The source app updates the target app

  # Optional timeouts
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 2: Create a "replace" supersedence relationship between two apps
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "replace_example" {
  source_id         = "00000000-0000-0000-0000-000000000003" # ID of the replacement app
  target_id         = "00000000-0000-0000-0000-000000000004" # ID of the app being replaced
  supersedence_type = "replace"                              # Uninstall the target before installing the source
}
