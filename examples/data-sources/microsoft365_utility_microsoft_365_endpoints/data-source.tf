# Example 1: Get all worldwide Microsoft 365 endpoints
data "microsoft365_utility_microsoft_365_endpoints" "all_worldwide" {
  instance = "worldwide"
}

# Output all endpoints
output "all_endpoints_count" {
  value = length(data.microsoft365_utility_microsoft_365_endpoints.all_worldwide.endpoints)
}

# Example 2: Get only required Exchange Online endpoints
data "microsoft365_utility_microsoft_365_endpoints" "exchange_required" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  required_only = true
}

# Example 3: Get Optimize category endpoints (highest priority)
data "microsoft365_utility_microsoft_365_endpoints" "optimize_only" {
  instance   = "worldwide"
  categories = ["Optimize"]
}

# Example 4: Get ExpressRoute-enabled endpoints
data "microsoft365_utility_microsoft_365_endpoints" "expressroute" {
  instance      = "worldwide"
  express_route = true
}

# Example 5: Get Teams endpoints for US Government DoD
data "microsoft365_utility_microsoft_365_endpoints" "teams_dod" {
  instance      = "usgov-dod"
  service_areas = ["Skype"] # Teams is under "Skype" service area
}

# Example 6: Get SharePoint and OneDrive endpoints with multiple filters
data "microsoft365_utility_microsoft_365_endpoints" "sharepoint_allow" {
  instance      = "worldwide"
  service_areas = ["SharePoint"]
  categories    = ["Allow", "Optimize"]
  required_only = true
}

# Example 7: China cloud endpoints
data "microsoft365_utility_microsoft_365_endpoints" "china_all" {
  instance = "china"
}

# Example 8: US Government GCC High endpoints
data "microsoft365_utility_microsoft_365_endpoints" "gcchigh_all" {
  instance = "usgov-gcchigh"
}

# Output examples - Extract specific data for use in firewall rules, etc.
output "exchange_optimize_urls" {
  description = "Exchange Optimize category URLs"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize_only.endpoints :
    endpoint.urls
    if endpoint.service_area == "Exchange" && endpoint.urls != null
  ]
}

output "exchange_optimize_ips" {
  description = "Exchange Optimize category IP ranges"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize_only.endpoints :
    endpoint.ips
    if endpoint.service_area == "Exchange" && endpoint.ips != null
  ]
}

output "teams_udp_ports" {
  description = "Teams UDP ports for media traffic"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.teams_dod.endpoints :
    endpoint.udp_ports
    if endpoint.udp_ports != null && endpoint.udp_ports != ""
  ]
}

# Example 9: Create locals for firewall automation
locals {
  # Extract all unique TCP ports from Optimize endpoints
  optimize_tcp_ports = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize_only.endpoints :
    split(",", endpoint.tcp_ports)
    if endpoint.tcp_ports != null && endpoint.tcp_ports != ""
  ]))

  # Extract all IP ranges for required endpoints
  required_ip_ranges = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.exchange_required.endpoints :
    endpoint.ips
    if endpoint.ips != null
  ]))

  # Create a map of service areas to their URLs
  service_area_urls = {
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.all_worldwide.endpoints :
    "${endpoint.service_area}_${endpoint.id}" => {
      service_area = endpoint.service_area_display_name
      category     = endpoint.category
      urls         = endpoint.urls
      required     = endpoint.required
      notes        = endpoint.notes
    }
    if endpoint.urls != null
  }
}

output "optimize_tcp_ports" {
  description = "Unique TCP ports for Optimize category endpoints"
  value       = local.optimize_tcp_ports
}

output "required_ip_ranges" {
  description = "All IP ranges for required Exchange endpoints"
  value       = local.required_ip_ranges
}

