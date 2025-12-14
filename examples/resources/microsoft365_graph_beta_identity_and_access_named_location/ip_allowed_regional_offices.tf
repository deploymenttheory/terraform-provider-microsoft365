# Specific Regional Offices (for region-locked accounts)

# EMEA Office Only
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_emea_office_only" {
  display_name = "Allowed - EMEA Office Only"
  is_trusted   = true

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

