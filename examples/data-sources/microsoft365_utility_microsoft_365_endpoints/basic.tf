# Basic Example: Get Microsoft 365 Endpoints

# Get all worldwide endpoints
data "microsoft365_utility_microsoft_365_endpoints" "all" {
  instance = "worldwide"
}

# Get only required Exchange Online endpoints
data "microsoft365_utility_microsoft_365_endpoints" "exchange_required" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  required_only = true
}

# Get Optimize category endpoints (direct routing, no proxy)
data "microsoft365_utility_microsoft_365_endpoints" "optimize" {
  instance   = "worldwide"
  categories = ["Optimize"]
}

# Get Teams/Skype endpoints
data "microsoft365_utility_microsoft_365_endpoints" "teams" {
  instance      = "worldwide"
  service_areas = ["Skype"]
}

# Get ExpressRoute-capable endpoints
data "microsoft365_utility_microsoft_365_endpoints" "expressroute" {
  instance      = "worldwide"
  express_route = true
}

# Outputs
output "total_endpoints" {
  description = "Total number of Microsoft 365 endpoints"
  value       = length(data.microsoft365_utility_microsoft_365_endpoints.all.endpoints)
}

output "exchange_urls" {
  description = "Exchange Online URLs"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.exchange_required.endpoints :
    endpoint.urls if endpoint.urls != null
  ]
}

output "optimize_ip_ranges" {
  description = "IP ranges for Optimize category (direct routing)"
  value = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    endpoint.ips if endpoint.ips != null
  ]))
}

output "teams_udp_ports" {
  description = "UDP ports for Teams media traffic"
  value = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.teams.endpoints :
    split(",", endpoint.udp_ports)
    if endpoint.udp_ports != null && endpoint.udp_ports != ""
  ]))
}

# Example: Extract specific endpoint details
output "optimize_endpoint_details" {
  description = "Detailed information about Optimize endpoints"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints : {
      id            = endpoint.id
      service       = endpoint.service_area_display_name
      category      = endpoint.category
      required      = endpoint.required
      urls          = endpoint.urls
      ips           = endpoint.ips
      tcp_ports     = endpoint.tcp_ports
      udp_ports     = endpoint.udp_ports
      express_route = endpoint.express_route
      notes         = endpoint.notes
    }
  ]
}

