# Example: Get all MEM (Microsoft Endpoint Manager / Intune) endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "mem_all" {
  instance      = "worldwide"
  service_areas = ["MEM"]
}

# Output: Total count of MEM endpoints
output "mem_endpoints_count" {
  description = "Total number of MEM endpoints (Expected: 15 for Worldwide)"
  value       = length(data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints)
}

# Output: All MEM URLs/FQDNs (flattened list)
output "mem_urls" {
  description = "All MEM/Intune URLs that need to be allowed"
  value = flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints :
    coalesce(endpoint.urls, [])
  ])
}

# Output: All MEM IP ranges (flattened list)
output "mem_ip_ranges" {
  description = "All MEM/Intune IP ranges in CIDR notation"
  value = flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints :
    coalesce(endpoint.ips, [])
  ])
}

# Output: Detailed breakdown of each MEM endpoint
output "mem_endpoint_details" {
  description = "Detailed information for each MEM endpoint"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints : {
      id            = endpoint.id
      display_name  = endpoint.service_area_display_name
      category      = endpoint.category
      required      = endpoint.required
      express_route = endpoint.express_route
      tcp_ports     = endpoint.tcp_ports
      udp_ports     = endpoint.udp_ports
      url_count     = length(coalesce(endpoint.urls, []))
      ip_count      = length(coalesce(endpoint.ips, []))
      urls          = coalesce(endpoint.urls, [])
      ips           = coalesce(endpoint.ips, [])
      notes         = endpoint.notes
    }
  ]
}

# Example: Get only REQUIRED MEM endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "mem_required" {
  instance      = "worldwide"
  service_areas = ["MEM"]
  required_only = true
}

output "mem_required_urls" {
  description = "Required MEM URLs only"
  value = flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_required.endpoints :
    coalesce(endpoint.urls, [])
  ])
}

# Example: Generate firewall rules for MEM endpoints
locals {
  mem_firewall_rules = flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints : [
      for url in coalesce(endpoint.urls, []) : {
        destination   = url
        tcp_ports     = endpoint.tcp_ports != "" ? endpoint.tcp_ports : null
        udp_ports     = endpoint.udp_ports != "" ? endpoint.udp_ports : null
        category      = endpoint.category
        required      = endpoint.required
        express_route = endpoint.express_route
      }
    ]
  ])

  # Group MEM endpoints by category
  mem_by_category = {
    for category in distinct([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints :
      endpoint.category
    ]) :
    category => [
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints :
      endpoint if endpoint.category == category
    ]
  }
}

output "mem_firewall_rules" {
  description = "Firewall rules for MEM endpoints"
  value       = local.mem_firewall_rules
}

output "mem_by_category" {
  description = "MEM endpoints grouped by category (Allow/Default)"
  value = {
    for category, endpoints in local.mem_by_category :
    category => {
      count = length(endpoints)
      urls  = distinct(flatten([for ep in endpoints : coalesce(ep.urls, [])]))
      ips   = distinct(flatten([for ep in endpoints : coalesce(ep.ips, [])]))
    }
  }
}

# Example: Extract specific Intune services
output "intune_core_management" {
  description = "Core Intune management endpoints (ID 163)"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints :
    {
      urls = coalesce(endpoint.urls, [])
      ips  = coalesce(endpoint.ips, [])
    }
    if endpoint.id == 163
  ]
}

output "windows_update_urls" {
  description = "Windows Update delivery endpoints for Intune"
  value = flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.mem_all.endpoints :
    coalesce(endpoint.urls, [])
    if(endpoint.notes != null && (
      strcontains(lower(endpoint.notes), "windows") ||
      strcontains(lower(endpoint.notes), constants.TfOperationUpdate)
    )) ||
    length([
      for url in coalesce(endpoint.urls, []) :
      url if strcontains(lower(url), constants.TfOperationUpdate) || strcontains(lower(url), "windowsupdate")
    ]) > 0
  ])
}

