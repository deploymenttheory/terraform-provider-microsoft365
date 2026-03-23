# Filter apps by category name (case-insensitive partial match)
# Note: This uses local filtering as categories are expanded relationships
data "microsoft365_graph_beta_device_and_app_management_mobile_app" "by_category" {
  category = "Productivity" # Finds all apps in Productivity category

  timeouts = {
    read = "30s" # Category filtering may take longer as it fetches categories for each app
  }
}

output "productivity_apps" {
  value = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.by_category.items : {
      id           = app.id
      display_name = app.display_name
      publisher    = app.publisher
      categories   = app.categories
      is_assigned  = app.is_assigned
    }
  ]
  description = "All apps in the Productivity category"
}
