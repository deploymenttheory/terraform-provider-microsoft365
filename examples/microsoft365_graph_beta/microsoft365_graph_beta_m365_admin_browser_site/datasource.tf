# Example 1: Query all browser sites in a specific site list
data "microsoft365_graph_beta_m365_admin_browser_site" "all_sites" {
  filter_type                   = "all"
  browser_site_list_assignment_id = "00000000-0000-0000-0000-000000000000"
}

# Output all site URLs
output "all_site_urls" {
  value = [for site in data.microsoft365_graph_beta_m365_admin_browser_site.all_sites.items : site.web_url]
}

# Example 2: Find a browser site by its ID
data "microsoft365_graph_beta_m365_admin_browser_site" "specific_site" {
  filter_type                   = "id"
  filter_value                  = "11111111-1111-1111-1111-111111111111"
  browser_site_list_assignment_id = "00000000-0000-0000-0000-000000000000"
}

# Access the specific site
output "specific_site_url" {
  value = length(data.microsoft365_graph_beta_m365_admin_browser_site.specific_site.items) > 0 ? data.microsoft365_graph_beta_m365_admin_browser_site.specific_site.items[0].web_url : null
}

# Example 3: Find browser sites by URL pattern
data "microsoft365_graph_beta_m365_admin_browser_site" "contoso_sites" {
  filter_type                   = "web_url"
  filter_value                  = "contoso.com"
  browser_site_list_assignment_id = "00000000-0000-0000-0000-000000000000"
}

# Output matching sites
output "contoso_sites" {
  value = data.microsoft365_graph_beta_m365_admin_browser_site.contoso_sites.items
}

# Example 4: Using the data source with resources
# First, get the existing browser site list
data "microsoft365_graph_beta_m365_admin_browser_site_list" "existing_list" {
  filter_type  = "display_name"
  filter_value = "My Browser Site List"
}

# Then query all sites in that list
data "microsoft365_graph_beta_m365_admin_browser_site" "sites_in_list" {
  filter_type                   = "all"
  browser_site_list_assignment_id = data.microsoft365_graph_beta_m365_admin_browser_site_list.existing_list.items[0].id
}

# Create a new browser site in the same list
resource "microsoft365_graph_beta_m365_admin_browser_site" "new_site" {
  browser_site_list_assignment_id = data.microsoft365_graph_beta_m365_admin_browser_site_list.existing_list.items[0].id
  web_url                         = "https://newapp.contoso.com"
  allow_redirect                  = true
  compatibility_mode              = "internetExplorer11"
  merge_type                      = "default"
  target_environment              = "internetExplorerMode"
  comment                         = "Added via Terraform"
}

# Example 5: Output count of browser sites
output "site_count" {
  value = length(data.microsoft365_graph_beta_m365_admin_browser_site.sites_in_list.items)
  description = "Total number of browser sites in the list"
}

# Example 6: Create a map of site IDs to URLs
output "site_map" {
  value = { for site in data.microsoft365_graph_beta_m365_admin_browser_site.sites_in_list.items : site.id => site.web_url }
  description = "Map of browser site IDs to their URLs"
}