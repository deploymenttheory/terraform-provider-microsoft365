# This example demonstrates how to create a supersedence relationship between two mobile apps in Intune
# The supersedence relationship allows you to upgrade or replace an older app with a newer version

# Example 1: Create an constants.TfOperationUpdate supersedence relationship between two apps
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "update_example" {
  source_id         = "00000000-0000-0000-0000-000000000001" # ID of the older app version
  target_id         = "00000000-0000-0000-0000-000000000002" # ID of the newer app version
  supersedence_type = constants.TfOperationUpdate                               # Indicates this is an update to the existing app

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
  source_id         = "00000000-0000-0000-0000-000000000003" # ID of the app being replaced
  target_id         = "00000000-0000-0000-0000-000000000004" # ID of the replacement app
  supersedence_type = "replace"                              # Indicates this app replaces the source app
}