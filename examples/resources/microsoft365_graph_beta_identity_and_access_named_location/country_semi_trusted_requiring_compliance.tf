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

