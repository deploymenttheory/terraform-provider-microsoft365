resource "microsoft365_graph_beta_identity_and_access_named_location" "high_risk_countries_blocked_by_client_ip" {
  display_name                          = "High Risk Countries Blocked by Client IP"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

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

resource "microsoft365_graph_beta_identity_and_access_named_location" "high_risk_countries_blocked_by_authenticator_gps" {
  display_name                          = "High Risk Countries Blocked by Authenticator GPS"
  country_lookup_method                 = "authenticatorAppGps"
  include_unknown_countries_and_regions = false

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
    # "203.0.113.0/24",      # Example: HQ public IP range
    # "203.0.114.0/24",      # Example: HQ additional subnet
  ]

  ipv6_ranges = [
    # "2001:db8:1234::/48",  # Example: HQ IPv6 range
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Corporate VPN Endpoints
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_corporate_vpn" {
  display_name = "Trusted - Corporate VPN"
  is_trusted   = true

  ipv4_ranges = [
    # "198.51.100.0/24",     # Example: VPN endpoint pool
  ]

  ipv6_ranges = [
    # "2001:db8:5678::/48",  # Example: VPN IPv6 pool
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Regional Office Networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_regional_offices" {
  display_name = "Trusted - Regional Offices"
  is_trusted   = true

  ipv4_ranges = [
    # "192.0.2.0/24",        # Example: EMEA Office
    # "198.51.100.0/24",     # Example: APAC Office
    # "203.0.113.0/24",      # Example: Americas Office
  ]

  ipv6_ranges = [
    # "2001:db8:abcd::/48",  # Example: Regional IPv6
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Cloud Infrastructure (Azure, AWS, etc.)
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_cloud_infrastructure" {
  display_name = "Trusted - Cloud Infrastructure"
  is_trusted   = true

  ipv4_ranges = [
    # "10.0.0.0/8",          # Example: Azure VNet CIDR
    # "172.16.0.0/12",       # Example: AWS VPC CIDR
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Trusted Countries/Regions
# Note: Country-based locations cannot be marked as is_trusted. 
# To use in CA policies, reference this location by ID explicitly.
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_countries" {
  display_name                          = "Trusted - Countries and Regions"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

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

# Authorized Home Office IPs (for remote executives/privileged users)
resource "microsoft365_graph_beta_identity_and_access_named_location" "trusted_authorized_home_offices" {
  display_name = "Trusted - Authorized Home Offices"
  is_trusted   = true

  ipv4_ranges = [
    # "203.0.113.50/32",     # Example: Executive home office
    # "198.51.100.75/32",    # Example: CTO home office
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# ==============================================================================
# These locations require additional security controls like device compliance
# but are not outright blocked. NOT marked as trusted.
# ==============================================================================

# Partner/Vendor Networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_partner_networks" {
  display_name = "Semi-Trusted - Partner Networks"

  ipv4_ranges = [
    # "198.18.0.0/24",       # Example: Partner A network
    # "198.18.1.0/24",       # Example: Partner B network
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Less-Trusted Countries (require device compliance)
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_countries" {
  display_name                          = "Semi-Trusted - Countries Requiring Compliance"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

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

# Public WiFi / Co-working Spaces (if you can identify these)
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_public_spaces" {
  display_name = "Semi-Trusted - Public WiFi and Co-working Spaces"

  ipv4_ranges = [
    # "192.0.2.0/24",        # Example: Known co-working space
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}


# ==============================================================================
# Restrictive locations for specific accounts (service accounts, regional users)
# These define the ONLY locations where certain accounts can sign in from.
# ==============================================================================

# Service Account Source IPs
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_service_account_sources" {
  display_name = "Allowed - Service Account Sources"
  is_trusted   = true

  ipv4_ranges = [
    # "10.100.0.10/32",      # Example: Build server
    # "10.100.0.11/32",      # Example: Automation server
    # "10.100.0.12/32",      # Example: Monitoring server
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Specific Regional Office (for region-locked accounts)
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_emea_office_only" {
  display_name = "Allowed - EMEA Office Only"
  is_trusted   = true

  ipv4_ranges = [
    # "203.0.113.0/24",      # Example: EMEA office IP range
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Specific Regional Office (for region-locked accounts)
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_apac_office_only" {
  display_name = "Allowed - APAC Office Only"
  is_trusted   = true

  ipv4_ranges = [
    # "198.51.100.0/24",     # Example: APAC office IP range
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Specific Regional Office (for region-locked accounts)
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_americas_office_only" {
  display_name = "Allowed - Americas Office Only"
  is_trusted   = true

  ipv4_ranges = [
    # "192.0.2.0/24",        # Example: Americas office IP range
  ]

  ipv6_ranges = []

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}