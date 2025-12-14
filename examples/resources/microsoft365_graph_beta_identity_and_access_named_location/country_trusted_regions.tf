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

