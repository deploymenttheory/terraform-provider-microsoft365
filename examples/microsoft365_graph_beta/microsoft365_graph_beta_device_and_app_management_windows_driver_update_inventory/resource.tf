resource "microsoft365_graph_beta_device_and_app_management_windows_driver_update_inventory" "example" {
  # Required attributes
  name                             = "Intel HD Graphics Driver"
  version                          = "27.20.100.8681"
  manufacturer                     = "Intel Corporation"
  approval_status                  = "approved"                             # Possible values: "needsReview", "declined", "approved", "suspended"
  category                         = "recommended"                          # Possible values: "recommended", "previouslyApproved", "other"
  windows_driver_update_profile_id = "12345678-1234-1234-1234-123456789012" # ID of the Windows Driver Update Profile

  # Optional attributes
  release_date_time = "2024-12-15T00:00:00Z"
  driver_class      = "Display"
  deploy_date_time  = "2025-01-15T00:00:00Z" # Only needed if approval_status is "approved"

  # Optional timeouts
  timeouts = {
    create = "3m"
    update = "3m"
    read   = "3m"
    delete = "3m"
  }
}