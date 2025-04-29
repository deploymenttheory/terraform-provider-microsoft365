# Return all catalog items
data "microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item" "all" {
  filter_type = "all"
}

output "all_windows_updates" {
  description = "All Windows Update Catalog Items retrieved from Intune"
  value = [
    for item in data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.all.items : {
      id                 = item.id
      display_name       = item.display_name
      release_date_time  = item.release_date_time
      end_of_support_date = item.end_of_support_date
    }
  ]
}

# Filter by ID
data "microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item" "by_id" {
  filter_type  = "id"
  filter_value = "de352f933dc512e9f8af0fedbbbdcda1b6ef79c0741b81554bbd6743af8b4d89"
}

output "specific_update" {
  description = "Details of the specific Windows Update Catalog Item retrieved by ID"
  value = length(data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_id.items) > 0 ? {
    id                  = data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_id.items[0].id
    display_name        = data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_id.items[0].display_name
    release_date_time   = data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_id.items[0].release_date_time
    end_of_support_date = data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_id.items[0].end_of_support_date
  } : null
}

# Filter by display name (will return all items containing "Security Update" in their name)
data "microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item" "by_name" {
  filter_type  = "display_name"
  filter_value = "02/25/2025 - 2025.02 D Update for Windows 10 and later"
}

output "updates_by_name" {
  description = "Windows Update Catalog Items matching the display name '02/25/2025 - 2025.02 D Update for Windows 10 and later'"
  value = [
    for item in data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_name.items : {
      id                  = item.id
      display_name        = item.display_name
      release_date_time   = item.release_date_time
      end_of_support_date = item.end_of_support_date
    }
  ]
}

# Filter by release date
data "microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item" "by_release_date" {
  filter_type  = "release_date_time"
  filter_value = "2025-01-14T00:00:00Z"
}

output "updates_by_release_date" {
  description = "Windows Update Catalog Items released on January 14, 2025"
  value = [
    for item in data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_release_date.items : {
      id                  = item.id
      display_name        = item.display_name
      release_date_time   = item.release_date_time
      end_of_support_date = item.end_of_support_date
    }
  ]
}

# Filter by end of support date
data "microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item" "by_end_of_support" {
  filter_type  = "end_of_support_date"
  filter_value = "2027-10-12T00:00:00Z"
}

output "updates_by_end_of_support" {
  description = "Windows Update Catalog Items with end of support date on October 12, 2027"
  value = [
    for item in data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_end_of_support.items : {
      id                  = item.id
      display_name        = item.display_name
      release_date_time   = item.release_date_time
      end_of_support_date = item.end_of_support_date
    }
  ]
}

output "end_of_support_count" {
  description = "Number of updates with end of support date on October 12, 2027"
  value = length(data.microsoft365_graph_beta_device_and_app_management_windows_update_catalog_item.by_end_of_support.items)
}