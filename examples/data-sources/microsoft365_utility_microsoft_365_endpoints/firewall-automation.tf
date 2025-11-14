# Practical Example: Firewall Rule Automation with Microsoft 365 Endpoints
# This example demonstrates how to use the endpoints datasource to automate
# firewall rule creation for Microsoft 365 services.

# ============================================================================
# 1. Fetch Microsoft 365 Endpoints with Different Priority Levels
# ============================================================================

# Critical (Optimize) - Direct routing recommended, latency sensitive
data "microsoft365_utility_microsoft_365_endpoints" "optimize" {
  instance   = "worldwide"
  categories = ["Optimize"]
}

# Important (Allow) - Direct routing recommended, medium priority
data "microsoft365_utility_microsoft_365_endpoints" "allow" {
  instance   = "worldwide"
  categories = ["Allow"]
}

# Standard (Default) - Can route through proxy
data "microsoft365_utility_microsoft_365_endpoints" "default" {
  instance   = "worldwide"
  categories = ["Default"]
}

# ============================================================================
# 2. Process Endpoints for Firewall Rules
# ============================================================================

locals {
  # Extract Optimize category endpoints (highest priority)
  optimize_endpoints = {
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    "optimize_${endpoint.id}" => {
      id            = endpoint.id
      service_area  = endpoint.service_area_display_name
      urls          = endpoint.urls != null ? endpoint.urls : []
      ips           = endpoint.ips != null ? endpoint.ips : []
      tcp_ports     = endpoint.tcp_ports != null ? endpoint.tcp_ports : ""
      udp_ports     = endpoint.udp_ports != null ? endpoint.udp_ports : ""
      category      = endpoint.category
      required      = endpoint.required
      express_route = endpoint.express_route
    }
  }

  # Extract Allow category endpoints
  allow_endpoints = {
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.allow.endpoints :
    "allow_${endpoint.id}" => {
      id            = endpoint.id
      service_area  = endpoint.service_area_display_name
      urls          = endpoint.urls != null ? endpoint.urls : []
      ips           = endpoint.ips != null ? endpoint.ips : []
      tcp_ports     = endpoint.tcp_ports != null ? endpoint.tcp_ports : ""
      udp_ports     = endpoint.udp_ports != null ? endpoint.udp_ports : ""
      category      = endpoint.category
      required      = endpoint.required
      express_route = endpoint.express_route
    }
  }

  # Extract all IP ranges by category for easy reference
  optimize_ip_ranges = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    endpoint.ips
    if endpoint.ips != null
  ]))

  allow_ip_ranges = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.allow.endpoints :
    endpoint.ips
    if endpoint.ips != null
  ]))

  # Extract all FQDNs by category
  optimize_fqdns = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    endpoint.urls
    if endpoint.urls != null
  ]))

  allow_fqdns = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.allow.endpoints :
    endpoint.urls
    if endpoint.urls != null
  ]))

  # Extract Teams media UDP ports (critical for audio/video)
  teams_udp_ports = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    split(",", endpoint.udp_ports)
    if endpoint.service_area == "Skype" && endpoint.udp_ports != null && endpoint.udp_ports != ""
  ]))

  # Extract Teams media IP ranges
  teams_media_ips = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    endpoint.ips
    if endpoint.service_area == "Skype" && endpoint.ips != null
  ]))
}

# ============================================================================
# 3. Output Structured Data for Firewall Configuration
# ============================================================================

output "firewall_rules_optimize" {
  description = "Optimize category endpoints requiring direct routing (bypass proxy/inspection)"
  value = {
    category    = "Optimize"
    description = "Latency-sensitive Microsoft 365 traffic - route directly, do not inspect"
    priority    = 100
    action      = "ALLOW"
    ip_ranges   = local.optimize_ip_ranges
    fqdns       = local.optimize_fqdns
    tcp_ports   = "80,443"
    udp_ports   = "443,3478,3479,3480,3481"
  }
}

output "firewall_rules_allow" {
  description = "Allow category endpoints (can be proxied but direct preferred)"
  value = {
    category    = "Allow"
    description = "Required Microsoft 365 traffic - proxy acceptable but not recommended"
    priority    = 200
    action      = "ALLOW"
    ip_ranges   = local.allow_ip_ranges
    fqdns       = local.allow_fqdns
    tcp_ports   = "25,80,143,443,587,993,995"
    udp_ports   = ""
  }
}

output "teams_media_config" {
  description = "Teams real-time media configuration (critical for call quality)"
  value = {
    service     = "Microsoft Teams Media"
    description = "UDP traffic for Teams audio/video - MUST NOT be proxied or inspected"
    priority    = 50
    action      = "ALLOW"
    ip_ranges   = local.teams_media_ips
    udp_ports   = join(",", local.teams_udp_ports)
    qos_marking = "DSCP EF (46)" # Expedited Forwarding for real-time traffic
  }
}

# ============================================================================
# 4. Service-Specific Endpoints
# ============================================================================

# Exchange Online
data "microsoft365_utility_microsoft_365_endpoints" "exchange" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  required_only = true
}

output "exchange_endpoints" {
  description = "Exchange Online required endpoints"
  value = {
    service = "Exchange Online"
    urls = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoints.exchange.endpoints :
      endpoint.urls
      if endpoint.urls != null
    ]))
    ips = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoints.exchange.endpoints :
      endpoint.ips
      if endpoint.ips != null
    ]))
    ports = {
      tcp = "25,80,143,443,587,993,995"
    }
  }
}

# SharePoint Online and OneDrive
data "microsoft365_utility_microsoft_365_endpoints" "sharepoint" {
  instance      = "worldwide"
  service_areas = ["SharePoint"]
  required_only = true
}

output "sharepoint_endpoints" {
  description = "SharePoint Online and OneDrive required endpoints"
  value = {
    service = "SharePoint Online / OneDrive"
    urls = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoints.sharepoint.endpoints :
      endpoint.urls
      if endpoint.urls != null
    ]))
    ips = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoints.sharepoint.endpoints :
      endpoint.ips
      if endpoint.ips != null
    ]))
    ports = {
      tcp = "80,443"
    }
  }
}

# ============================================================================
# 5. PAC File / Proxy Bypass List Generation
# ============================================================================

locals {
  # Generate PAC file bypass list for Optimize category
  pac_bypass_domains = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints :
    endpoint.urls
    if endpoint.urls != null
  ]

  # Flatten and format for PAC file
  pac_bypass_list = distinct(flatten(local.pac_bypass_domains))
}

output "pac_file_bypass_list" {
  description = "Domains to bypass proxy (Optimize category) - use in PAC files"
  value       = local.pac_bypass_list
}

# ============================================================================
# 6. SD-WAN / QoS Configuration
# ============================================================================

output "sdwan_policy_optimize" {
  description = "SD-WAN policy for Optimize traffic - use direct internet breakout"
  value = {
    name        = "M365-Optimize-Direct"
    action      = "DIRECT_INTERNET"
    priority    = 10
    bandwidth   = "HIGH"
    latency     = "LOW"
    packet_loss = "LOW"
    jitter      = "LOW"
    ip_ranges   = local.optimize_ip_ranges
    fqdns       = local.optimize_fqdns
  }
}

output "sdwan_policy_teams_media" {
  description = "SD-WAN policy for Teams media - highest QoS priority"
  value = {
    name          = "Teams-Media-RTP"
    action        = "DIRECT_INTERNET"
    priority      = 1
    qos_dscp      = "EF" # Expedited Forwarding (46)
    min_bandwidth = "5Mbps"
    ip_ranges     = local.teams_media_ips
    udp_ports     = join(",", local.teams_udp_ports)
  }
}

# ============================================================================
# 7. ExpressRoute for Microsoft 365 Configuration
# ============================================================================

data "microsoft365_utility_microsoft_365_endpoints" "expressroute" {
  instance      = "worldwide"
  express_route = true
}

output "expressroute_prefixes" {
  description = "IP prefixes advertised over ExpressRoute for Microsoft 365"
  value = {
    enabled      = true
    community    = "12076:5010" # Microsoft 365 BGP community
    route_filter = "M365"
    ip_prefixes = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoints.expressroute.endpoints :
      endpoint.ips
      if endpoint.ips != null
    ]))
  }
}

# ============================================================================
# 8. Network Summary
# ============================================================================

output "network_summary" {
  description = "Summary of Microsoft 365 network requirements"
  value = {
    total_endpoints = (
      length(data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints) +
      length(data.microsoft365_utility_microsoft_365_endpoints.allow.endpoints) +
      length(data.microsoft365_utility_microsoft_365_endpoints.default.endpoints)
    )

    optimize = {
      count     = length(data.microsoft365_utility_microsoft_365_endpoints.optimize.endpoints)
      treatment = "Direct routing, no proxy, no inspection"
      qos       = "High priority, low latency"
    }

    allow = {
      count     = length(data.microsoft365_utility_microsoft_365_endpoints.allow.endpoints)
      treatment = "Direct routing preferred, proxy acceptable"
      qos       = "Medium priority"
    }

    default = {
      count     = length(data.microsoft365_utility_microsoft_365_endpoints.default.endpoints)
      treatment = "Can route through proxy and inspection"
      qos       = "Standard priority"
    }

    expressroute_enabled = length(data.microsoft365_utility_microsoft_365_endpoints.expressroute.endpoints)
  }
}

