---
page_title: "microsoft365_graph_beta_identity_and_access_named_location Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft 365 Named Locations using the /identity/conditionalAccess/namedLocations endpoint. Named Locations define network locations that can be used in Conditional Access policies. Supports both IP-based and country-based named locations.
---

# microsoft365_graph_beta_identity_and_access_named_location (Resource)

Manages Microsoft 365 Named Locations using the `/identity/conditionalAccess/namedLocations` endpoint. Named Locations define network locations that can be used in Conditional Access policies. Supports both IP-based and country-based named locations.

## Microsoft Documentation

- [namedLocation resource type](https://learn.microsoft.com/en-us/graph/api/resources/namedlocation?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Policy.ReadWrite.ConditionalAccess`, `Policy.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.28.0-alpha | Experimental | Initial release |
| v0.38.0-alpha | Preview | Refactored to use Kiota SDK |


## Example Usage

### Country-Based Named Locations

#### High Risk Countries (Client IP)
```terraform
resource "microsoft365_graph_beta_identity_and_access_named_location" "high_risk_countries_blocked_by_client_ip" {
  display_name                          = "High Risk Countries Blocked by Client IP"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false
  hard_delete                           = true

  countries_and_regions = [
    "AF", # Afghanistan
    "BY", # Belarus
    "CN", # China
    "CU", # Cuba
    "ER", # Eritrea
    "IR", # Iran
    "IQ", # Iraq
    "KP", # North Korea
    "LY", # Libya
    "MM", # Myanmar (Burma)
    "RU", # Russia
    "SO", # Somalia
    "SD", # Sudan
    "SS", # South Sudan
    "SY", # Syria
    "VE", # Venezuela
    "YE", # Yemen
  ]
}
```

#### High Risk Countries (Authenticator GPS)
```terraform
resource "microsoft365_graph_beta_identity_and_access_named_location" "high_risk_countries_blocked_by_authenticator_gps" {
  display_name                          = "High Risk Countries Blocked by Authenticator GPS"
  country_lookup_method                 = "authenticatorAppGps"
  include_unknown_countries_and_regions = false
  hard_delete                           = true

  countries_and_regions = [
    "AF", # Afghanistan
    "BY", # Belarus
    "CN", # China
    "CU", # Cuba
    "ER", # Eritrea
    "IR", # Iran
    "IQ", # Iraq
    "KP", # North Korea
    "LY", # Libya
    "MM", # Myanmar (Burma)
    "RU", # Russia
    "SO", # Somalia
    "SD", # Sudan
    "SS", # South Sudan
    "SY", # Syria
    "VE", # Venezuela
    "YE", # Yemen
  ]
}
```

#### Trusted Countries and Regions
```terraform
# Trusted Countries/Regions
# Note: Country-based locations cannot be marked as is_trusted. 
# To use in CA policies, reference this location by ID explicitly.
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_countries" {
  display_name                          = "Trusted - Countries and Regions"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false
  hard_delete                           = true

  countries_and_regions = [
    "US", # United States
    "CA", # Canada
    "GB", # United Kingdom
    "IE", # Ireland
    "AU", # Australia
    "NZ", # New Zealand
    "DE", # Germany
    "FR", # France
    "NL", # Netherlands
    "BE", # Belgium
    "SE", # Sweden
    "NO", # Norway
    "DK", # Denmark
    "FI", # Finland
    "CH", # Switzerland
    "AT", # Austria
    "ES", # Spain
    "IT", # Italy
    "PT", # Portugal
    "PL", # Poland
    "CZ", # Czech Republic
    "SG", # Singapore
    "JP", # Japan
    "KR", # South Korea
    "IL", # Israel
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

#### Semi-Trusted Countries (Requiring Compliance)
```terraform
# Less-Trusted Countries (require device compliance)
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_countries" {
  display_name                          = "Semi-Trusted - Countries Requiring Compliance"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false
  hard_delete                           = true

  countries_and_regions = [
    "BR", # Brazil
    "MX", # Mexico
    "AR", # Argentina
    "CO", # Colombia
    "PE", # Peru
    "CL", # Chile
    "IN", # India
    "PH", # Philippines
    "TH", # Thailand
    "MY", # Malaysia
    "ID", # Indonesia
    "VN", # Vietnam
    "ZA", # South Africa
    "EG", # Egypt
    "MA", # Morocco
    "KE", # Kenya
    "NG", # Nigeria
    "PK", # Pakistan
    "BD", # Bangladesh
    "TR", # Turkey
    "UA", # Ukraine
    "RO", # Romania
    "BG", # Bulgaria
    "GR", # Greece
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

### IP-Based Named Locations

#### Trusted Corporate Headquarters
```terraform
# ==============================================================================
# IP-based locations marked as is_trusted=true will be automatically included
# in the "AllTrusted" built-in location reference in Conditional Access policies.
# Note: Country-based locations do NOT support is_trusted attribute and must be
# referenced explicitly by ID in CA policies if needed.
# ==============================================================================

# Corporate Headquarters Network
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_corporate_hq" {
  display_name = "Trusted - Corporate Headquarters"
  is_trusted   = true

  ipv4_ranges = [
    "203.0.113.0/24", # Example: HQ public IP range
    "203.0.114.0/24", # Example: HQ additional subnet
  ]

  ipv6_ranges = [
    "2001:db8:1234::/48", # Example: HQ IPv6 range
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

#### Trusted Corporate VPN
```terraform
# Corporate VPN Endpoints
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_corporate_vpn" {
  display_name = "Trusted - Corporate VPN"
  is_trusted   = true

  ipv4_ranges = [
    "198.51.100.0/24", # Example: VPN endpoint pool
  ]

  ipv6_ranges = [
    "2001:db8:5678::/48", # Example: VPN IPv6 pool
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

#### Trusted Regional Offices
```terraform
# Regional Office Networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_regional_offices" {
  display_name = "Trusted - Regional Offices"
  is_trusted   = true

  ipv4_ranges = [
    "192.0.2.0/24",    # Example: EMEA Office
    "198.51.100.0/24", # Example: APAC Office
    "203.0.113.0/24",  # Example: Americas Office
  ]

  ipv6_ranges = [
    "2001:db8:abcd::/48", # Example: Regional IPv6
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

#### Semi-Trusted Partner Networks
```terraform
# ==============================================================================
# These locations require additional security controls like device compliance
# but are not outright blocked. NOT marked as trusted.
# ==============================================================================

# Partner/Vendor Networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_partner_networks" {
  display_name = "Semi-Trusted - Partner Networks"

  ipv4_ranges = [
    "198.18.0.0/24", # Example: Partner A network
    "198.18.1.0/24", # Example: Partner B network
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

#### Allowed Service Account Sources
```terraform
# ==============================================================================
# Restrictive locations for specific accounts (service accounts, regional users)
# These define the ONLY locations where certain accounts can sign in from.
# ==============================================================================

# Service Account Source IPs
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_service_account_sources" {
  display_name = "Allowed - Service Account Sources"
  is_trusted   = true
  hard_delete  = true

  ipv4_ranges = [
    "10.100.0.10/32", # Example: Build server
    "10.100.0.11/32", # Example: Automation server
    "10.100.0.12/32", # Example: Monitoring server
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

#### Allowed Regional Offices (Region-Locked Accounts)
```terraform
# Specific Regional Offices (for region-locked accounts)

# EMEA Office Only
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_emea_office_only" {
  display_name = "Allowed - EMEA Office Only"
  is_trusted   = true
  hard_delete  = true

  ipv4_ranges = [
    "203.0.113.0/24", # Example: EMEA office IP range
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# APAC Office Only
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_apac_office_only" {
  display_name = "Allowed - APAC Office Only"
  is_trusted   = true
  hard_delete  = true

  ipv4_ranges = [
    "198.51.100.0/24", # Example: APAC office IP range
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Americas Office Only
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_americas_office_only" {
  display_name = "Allowed - Americas Office Only"
  is_trusted   = true
  hard_delete  = true

  ipv4_ranges = [
    "192.0.2.0/24", # Example: Americas office IP range
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the Named Location.

### Optional

- `countries_and_regions` (Set of String) Set of countries and/or regions in two-letter format specified by ISO 3166-2 (e.g., 'US', 'GB', 'CA'). Used for country named locations only.
- `country_lookup_method` (String) Provides the method used to decide which country the user is located in. Possible values are `clientIpAddress` and `authenticatorAppGps`. Used for country named locations only.
- `hard_delete` (Boolean) When `true`, the named location will be permanently deleted (hard delete) during destroy. When `false` (default), the named location will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. Note: This field defaults to `false` on import since the API does not return this value.
- `include_unknown_countries_and_regions` (Boolean) True if IP addresses that don't map to a country or region should be included in the named location. Used for country named locations only.
- `ipv4_ranges` (Set of String) Set of IPv4 CIDR ranges that define this IP named location. Each range should be specified in CIDR notation (e.g., '192.168.1.0/24'). Used for IP named locations only.
- `ipv6_ranges` (Set of String) Set of IPv6 CIDR ranges that define this IP named location. Each range should be specified in CIDR notation (e.g., '2001:db8::/32'). Used for IP named locations only.
- `is_trusted` (Boolean) Indicates whether the IP named location is trusted. Only applies to IP named locations.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The creation date and time of the named location.
- `id` (String) String (identifier)
- `modified_date_time` (String) The last modified date and time of the named location.

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

Represents a Microsoft Entra ID named location defined by countries and regions. Named locations are custom rules that define network locations which can then be used in a Conditional Access policy.
Represents a Microsoft Entra ID named location defined by IP ranges. Named locations are custom rules that define network locations that can then be used in a Conditional Access policy.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using composite ID format: {id}
terraform import microsoft365_graph_beta_identity_and_access_named_location.example 00000000-0000-0000-0000-000000000000
``` 