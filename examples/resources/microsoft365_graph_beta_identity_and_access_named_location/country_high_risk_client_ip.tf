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

