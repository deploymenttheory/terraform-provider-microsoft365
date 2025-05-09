# Get all browser site lists
data "microsoft365_graph_beta_m365_admin_browser_site_list" "all" {
  filter_type = "all"
}

# Output the first item's details from our "all" query
output "first_site_list" {
  value = length(data.microsoft365_graph_beta_m365_admin_browser_site_list.all.items) > 0 ? {
    id           = data.microsoft365_graph_beta_m365_admin_browser_site_list.all.items[0].id
    display_name = data.microsoft365_graph_beta_m365_admin_browser_site_list.all.items[0].display_name
  } : null
}

# Output all site lists as a collection
output "all_site_lists" {
  value = [
    for item in data.microsoft365_graph_beta_m365_admin_browser_site_list.all.items : {
      id           = item.id
      display_name = item.display_name
    }
  ]
}

# Output just the names of all site lists
output "site_list_names" {
  value = [
    for item in data.microsoft365_graph_beta_m365_admin_browser_site_list.all.items : 
    item.display_name
  ]
}

# Get a specific browser site list by ID
data "microsoft365_graph_beta_m365_admin_browser_site_list" "by_id" {
  filter_type  = "id"
  filter_value = "dec86081-2085-4917-8fe2-8580bdbd74d2"
}

# Get browser site lists by display name (partial match)
data "microsoft365_graph_beta_m365_admin_browser_site_list" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Example Browser Site List"
}

# Output for the site list queried by ID
output "site_list_by_id" {
  value = length(data.microsoft365_graph_beta_m365_admin_browser_site_list.by_id.items) > 0 ? {
    id           = data.microsoft365_graph_beta_m365_admin_browser_site_list.by_id.items[0].id
    display_name = data.microsoft365_graph_beta_m365_admin_browser_site_list.by_id.items[0].display_name
  } : null
  description = "Details of the browser site list with ID 'dec86081-2085-4917-8fe2-8580bdbd74d2'"
}

# Output for the site lists filtered by display name
output "site_lists_by_display_name" {
  value = [
    for item in data.microsoft365_graph_beta_m365_admin_browser_site_list.by_display_name.items : {
      id           = item.id
      display_name = item.display_name
    }
  ]
  description = "All browser site lists containing 'Example Browser Site List' in their display name"
}

# Count of results from display name search
output "display_name_match_count" {
  value = length(data.microsoft365_graph_beta_m365_admin_browser_site_list.by_display_name.items)
  description = "Number of browser site lists that match the display name filter"
}

# Check if ID lookup returned a result
output "id_lookup_found" {
  value = length(data.microsoft365_graph_beta_m365_admin_browser_site_list.by_id.items) > 0
  description = "Whether the browser site list with the specified ID was found"
}

# Combined output showing all results from both queries
output "all_queried_site_lists" {
  value = distinct(concat(
    data.microsoft365_graph_beta_m365_admin_browser_site_list.by_id.items,
    data.microsoft365_graph_beta_m365_admin_browser_site_list.by_display_name.items
  ))
  description = "Combined results from both ID and display name queries, with duplicates removed"
}