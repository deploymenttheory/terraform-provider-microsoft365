---
page_title: "microsoft365_utility_microsoft_365_endpoint_reference Data Source - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  Retrieves Microsoft 365 network endpoints from the official IP Address and URL Web Service using the https://endpoints.office.com API. This data source is used to query current IP addresses, URLs, and ports for firewall and proxy configuration.
---

# microsoft365_utility_microsoft_365_endpoint_reference

Retrieves Microsoft 365 network endpoints from the official [Microsoft 365 IP Address and URL Web Service](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-ip-web-service).

This datasource queries `https://endpoints.office.com` to get current IP addresses, URLs, and ports for Microsoft 365 services. 
It's designed for automating firewall rules, proxy configurations, SD-WAN policies, and PAC files with always up-to-date Microsoft 365 network requirements.

## Background

Microsoft publishes network endpoints for Microsoft 365 services through a REST API that is updated regularly. These endpoints include:

- **FQDNs/URLs**: Destination domains and URL patterns (may include wildcards like `*.office.com`)
- **IP Ranges**: IPv4 and IPv6 address ranges in CIDR notation
- **Ports**: TCP and UDP port numbers required for each service
- **Categories**: Network optimization categories (`Optimize`, `Allow`, `Default`)
- **Service Areas**: Exchange Online, SharePoint Online, Microsoft Teams, and common services

Network administrators use this data to:
- Configure firewalls and security appliances
- Define proxy bypass lists (PAC files)
- Implement SD-WAN and QoS policies
- Plan Azure ExpressRoute for Microsoft 365
- Optimize network paths for latency-sensitive traffic

## Network Categories

Microsoft defines three network categories based on optimization requirements:

### Optimize (Highest Priority)
- **Description**: Latency-sensitive services with highest traffic volume
- **Examples**: Exchange Online mailbox access, Teams real-time media (audio/video)
- **Recommended Treatment**:
  - Direct internet routing (bypass proxy and packet inspection)
  - Lowest latency network path
  - QoS marking: DSCP EF (46) for real-time media
  - Never proxy or perform TLS inspection
- **Typical Impact**: 70% of Microsoft 365 bandwidth, most sensitive to network conditions

### Allow (Medium Priority)
- **Description**: Required endpoints for core functionality with lower latency sensitivity
- **Examples**: Exchange mail flow (SMTP), SharePoint/OneDrive file operations
- **Recommended Treatment**:
  - Direct routing preferred but proxy acceptable
  - Standard QoS policies
  - TLS inspection acceptable with caution
- **Typical Impact**: 20% of Microsoft 365 bandwidth

### Default (Standard Priority)
- **Description**: Optional services or low-priority traffic
- **Examples**: Office CDNs, telemetry, updates, optional features
- **Recommended Treatment**:
  - Can route through standard proxy
  - Normal firewall inspection
  - Standard QoS
- **Typical Impact**: 10% of Microsoft 365 bandwidth

See [Microsoft 365 Network Connectivity Principles](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-network-connectivity-principles) for detailed guidance.

## Supported Cloud Instances

| Instance Value | Description | API Endpoint |
|----------------|-------------|--------------|
| `worldwide` | Worldwide commercial cloud (includes US GCC) | `https://endpoints.office.com/endpoints/worldwide` |
| `usgov-dod` | US Government DoD | `https://endpoints.office.com/endpoints/USGOVDoD` |
| `usgov-gcchigh` | US Government GCC High | `https://endpoints.office.com/endpoints/USGOVGCCHigh` |
| `china` | Microsoft 365 operated by 21Vianet (China) | `https://endpoints.office.com/endpoints/China` |

## Service Areas

- **MEM**: Microsoft Endpoint Manager (Intune, Autopilot, Windows Updates, Remote Help)
- **Exchange**: Exchange Online, Outlook, Exchange Online Protection (EOP)
- **SharePoint**: SharePoint Online and OneDrive for Business
- **Skype**: Microsoft Teams and Skype for Business Online
- **Common**: Microsoft Entra ID, Office 365 portal, shared services

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.36.0-alpha | Experimental | Initial release |

## Example Usage

### Basic Usage

```terraform
# Basic Example: Get Microsoft 365 Endpoints

# Get all worldwide endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "all" {
  instance = "worldwide"
}

# Get only required Exchange Online endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "exchange_required" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  required_only = true
}

# Get Optimize category endpoints (direct routing, no proxy)
data "microsoft365_utility_microsoft_365_endpoint_reference" "optimize" {
  instance   = "worldwide"
  categories = ["Optimize"]
}

# Get Teams/Skype endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "teams" {
  instance      = "worldwide"
  service_areas = ["Skype"]
}

# Get ExpressRoute-capable endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "expressroute" {
  instance      = "worldwide"
  express_route = true
}

# Outputs
output "total_endpoints" {
  description = "Total number of Microsoft 365 endpoints"
  value       = length(data.microsoft365_utility_microsoft_365_endpoint_reference.all.endpoints)
}

output "exchange_urls" {
  description = "Exchange Online URLs"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.exchange_required.endpoints :
    endpoint.urls if endpoint.urls != null
  ]
}

output "optimize_ip_ranges" {
  description = "IP ranges for Optimize category (direct routing)"
  value = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
    endpoint.ips if endpoint.ips != null
  ]))
}

output "teams_udp_ports" {
  description = "UDP ports for Teams media traffic"
  value = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.teams.endpoints :
    split(",", endpoint.udp_ports)
    if endpoint.udp_ports != null && endpoint.udp_ports != ""
  ]))
}

# Example: Extract specific endpoint details
output "optimize_endpoint_details" {
  description = "Detailed information about Optimize endpoints"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints : {
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
```

### Comprehensive Filtering and Outputs

```terraform
# Example 1: Get all worldwide Microsoft 365 endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "all_worldwide" {
  instance = "worldwide"
}

# Output all endpoints
output "all_endpoints_count" {
  value = length(data.microsoft365_utility_microsoft_365_endpoint_reference.all_worldwide.endpoints)
}

# Example 2: Get only required Exchange Online endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "exchange_required" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  required_only = true
}

# Example 3: Get Optimize category endpoints (highest priority)
data "microsoft365_utility_microsoft_365_endpoint_reference" "optimize_only" {
  instance   = "worldwide"
  categories = ["Optimize"]
}

# Example 4: Get ExpressRoute-enabled endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "expressroute" {
  instance      = "worldwide"
  express_route = true
}

# Example 5: Get Teams endpoints for US Government DoD
data "microsoft365_utility_microsoft_365_endpoint_reference" "teams_dod" {
  instance      = "usgov-dod"
  service_areas = ["Skype"] # Teams is under "Skype" service area
}

# Example 6: Get SharePoint and OneDrive endpoints with multiple filters
data "microsoft365_utility_microsoft_365_endpoint_reference" "sharepoint_allow" {
  instance      = "worldwide"
  service_areas = ["SharePoint"]
  categories    = ["Allow", "Optimize"]
  required_only = true
}

# Example 7: China cloud endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "china_all" {
  instance = "china"
}

# Example 8: US Government GCC High endpoints
data "microsoft365_utility_microsoft_365_endpoint_reference" "gcchigh_all" {
  instance = "usgov-gcchigh"
}

# Output examples - Extract specific data for use in firewall rules, etc.
output "exchange_optimize_urls" {
  description = "Exchange Optimize category URLs"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize_only.endpoints :
    endpoint.urls
    if endpoint.service_area == "Exchange" && endpoint.urls != null
  ]
}

output "exchange_optimize_ips" {
  description = "Exchange Optimize category IP ranges"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize_only.endpoints :
    endpoint.ips
    if endpoint.service_area == "Exchange" && endpoint.ips != null
  ]
}

output "teams_udp_ports" {
  description = "Teams UDP ports for media traffic"
  value = [
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.teams_dod.endpoints :
    endpoint.udp_ports
    if endpoint.udp_ports != null && endpoint.udp_ports != ""
  ]
}

# Example 9: Create locals for firewall automation
locals {
  # Extract all unique TCP ports from Optimize endpoints
  optimize_tcp_ports = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize_only.endpoints :
    split(",", endpoint.tcp_ports)
    if endpoint.tcp_ports != null && endpoint.tcp_ports != ""
  ]))

  # Extract all IP ranges for required endpoints
  required_ip_ranges = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.exchange_required.endpoints :
    endpoint.ips
    if endpoint.ips != null
  ]))

  # Create a map of service areas to their URLs
  service_area_urls = {
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.all_worldwide.endpoints :
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
```

### Advanced Firewall Automation

```terraform
# Practical Example: Firewall Rule Automation with Microsoft 365 Endpoints
# This example demonstrates how to use the endpoints datasource to automate
# firewall rule creation for Microsoft 365 services.

# ============================================================================
# 1. Fetch Microsoft 365 Endpoints with Different Priority Levels
# ============================================================================

# Critical (Optimize) - Direct routing recommended, latency sensitive
data "microsoft365_utility_microsoft_365_endpoint_reference" "optimize" {
  instance   = "worldwide"
  categories = ["Optimize"]
}

# Important (Allow) - Direct routing recommended, medium priority
data "microsoft365_utility_microsoft_365_endpoint_reference" "allow" {
  instance   = "worldwide"
  categories = ["Allow"]
}

# Standard (Default) - Can route through proxy
data "microsoft365_utility_microsoft_365_endpoint_reference" "default" {
  instance   = "worldwide"
  categories = ["Default"]
}

# ============================================================================
# 2. Process Endpoints for Firewall Rules
# ============================================================================

locals {
  # Extract Optimize category endpoints (highest priority)
  optimize_endpoints = {
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
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
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.allow.endpoints :
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
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
    endpoint.ips
    if endpoint.ips != null
  ]))

  allow_ip_ranges = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.allow.endpoints :
    endpoint.ips
    if endpoint.ips != null
  ]))

  # Extract all FQDNs by category
  optimize_fqdns = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
    endpoint.urls
    if endpoint.urls != null
  ]))

  allow_fqdns = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.allow.endpoints :
    endpoint.urls
    if endpoint.urls != null
  ]))

  # Extract Teams media UDP ports (critical for audio/video)
  teams_udp_ports = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
    split(",", endpoint.udp_ports)
    if endpoint.service_area == "Skype" && endpoint.udp_ports != null && endpoint.udp_ports != ""
  ]))

  # Extract Teams media IP ranges
  teams_media_ips = distinct(flatten([
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
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
data "microsoft365_utility_microsoft_365_endpoint_reference" "exchange" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  required_only = true
}

output "exchange_endpoints" {
  description = "Exchange Online required endpoints"
  value = {
    service = "Exchange Online"
    urls = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.exchange.endpoints :
      endpoint.urls
      if endpoint.urls != null
    ]))
    ips = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.exchange.endpoints :
      endpoint.ips
      if endpoint.ips != null
    ]))
    ports = {
      tcp = "25,80,143,443,587,993,995"
    }
  }
}

# SharePoint Online and OneDrive
data "microsoft365_utility_microsoft_365_endpoint_reference" "sharepoint" {
  instance      = "worldwide"
  service_areas = ["SharePoint"]
  required_only = true
}

output "sharepoint_endpoints" {
  description = "SharePoint Online and OneDrive required endpoints"
  value = {
    service = "SharePoint Online / OneDrive"
    urls = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.sharepoint.endpoints :
      endpoint.urls
      if endpoint.urls != null
    ]))
    ips = distinct(flatten([
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.sharepoint.endpoints :
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
    for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints :
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

data "microsoft365_utility_microsoft_365_endpoint_reference" "expressroute" {
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
      for endpoint in data.microsoft365_utility_microsoft_365_endpoint_reference.expressroute.endpoints :
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
      length(data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints) +
      length(data.microsoft365_utility_microsoft_365_endpoint_reference.allow.endpoints) +
      length(data.microsoft365_utility_microsoft_365_endpoint_reference.default.endpoints)
    )

    optimize = {
      count     = length(data.microsoft365_utility_microsoft_365_endpoint_reference.optimize.endpoints)
      treatment = "Direct routing, no proxy, no inspection"
      qos       = "High priority, low latency"
    }

    allow = {
      count     = length(data.microsoft365_utility_microsoft_365_endpoint_reference.allow.endpoints)
      treatment = "Direct routing preferred, proxy acceptable"
      qos       = "Medium priority"
    }

    default = {
      count     = length(data.microsoft365_utility_microsoft_365_endpoint_reference.default.endpoints)
      treatment = "Can route through proxy and inspection"
      qos       = "Standard priority"
    }

    expressroute_enabled = length(data.microsoft365_utility_microsoft_365_endpoint_reference.expressroute.endpoints)
  }
}
```

### Microsoft Endpoint Manager (Intune) Endpoints

```terraform
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
```

## Argument Reference

* `instance` - (Required) The Microsoft 365 cloud instance to query. Valid values: `worldwide`, `usgov-dod`, `usgov-gcchigh`, `china`.

* `service_areas` - (Optional) Filter endpoints by service area. Valid values: `MEM`, `Exchange`, `SharePoint`, `Skype`, `Common`. If omitted, returns all service areas.

* `categories` - (Optional) Filter endpoints by network optimization category. Valid values: `Optimize`, `Allow`, `Default`. If omitted, returns all categories.

* `required_only` - (Optional) If `true`, only returns endpoints marked as required by Microsoft. Optional endpoints provide enhanced functionality but are not necessary for core service operation. Defaults to `false`.

* `express_route` - (Optional) If `true`, only returns endpoints that support Azure ExpressRoute for Microsoft 365. Useful for organizations using ExpressRoute for optimized connectivity. Defaults to `false`.

* `timeouts` - (Optional) Timeout configuration block. See [Timeouts](#timeouts) below.

## Additional Resources

- [Microsoft 365 Network Connectivity Principles](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-network-connectivity-principles)
- [Managing Microsoft 365 Endpoints](https://learn.microsoft.com/en-us/microsoft-365/enterprise/managing-office-365-endpoints)
- [Microsoft 365 IP Address and URL Web Service](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-ip-web-service)
- [Microsoft 365 Endpoints Documentation](https://learn.microsoft.com/en-us/microsoft-365/enterprise/microsoft-365-endpoints)
- [Unified cloud.microsoft Domain](https://learn.microsoft.com/en-us/microsoft-365/enterprise/cloud-microsoft-domain)
- [Office 365 URLs and IP address ranges](https://learn.microsoft.com/en-us/microsoft-365/enterprise/urls-and-ip-address-ranges)

