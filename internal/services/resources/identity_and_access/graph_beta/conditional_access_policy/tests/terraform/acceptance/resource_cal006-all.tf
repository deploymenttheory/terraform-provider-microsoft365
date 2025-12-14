# ==============================================================================
# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Named Location Dependencies
# ==============================================================================

# APAC Office Location
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_apac_office_only" {
  display_name = "Allowed - APAC Office Only ${random_string.suffix.result}"
  is_trusted   = true

  ipv4_ranges = [
    "198.51.100.0/24", # Example: APAC office IP range
  ]

  ipv6_ranges = [
    "2001:db8:1234::/48", # Example: APAC office IPv6 range
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# EMEA Office Location
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_emea_office_only" {
  display_name = "Allowed - EMEA Office Only ${random_string.suffix.result}"
  is_trusted   = true

  ipv4_ranges = [
    "203.0.113.0/24", # Example: EMEA office IP range
  ]

  ipv6_ranges = [
    "2001:db8:1234::/48", # Example: EMEA office IPv6 range
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# Hazelwood Office Location
resource "microsoft365_graph_beta_identity_and_access_named_location" "allowed_hazelwood_office_only" {
  display_name = "Allowed - Hazelwood Office Only ${random_string.suffix.result}"
  is_trusted   = true

  ipv4_ranges = [
    "82.44.54.0/24", # Example: Hazelwood office IP range
  ]

  ipv6_ranges = [
    "2001:db8:5678::/48", # Example: Hazelwood office IPv6 range
  ]

  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

# Break Glass Emergency Access Accounts
resource "microsoft365_graph_beta_groups_group" "breakglass" {
  display_name     = "EID_UA_ConAcc-Breakglass"
  mail_nickname    = "eid-ua-conacc-breakglass"
  mail_enabled     = false
  security_enabled = true
  description      = "Group containing Break Glass Accounts"
}

resource "microsoft365_graph_beta_groups_group" "cal006_exclude" {
  display_name     = "EID_UA_CAL006_EXCLUDE"
  mail_nickname    = "eid-ua-cal006-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAL006_EXCLUDE"
}

resource "microsoft365_graph_beta_groups_group" "cal006_include" {
  display_name     = "EID_UA_CAL006_INCLUDE"
  mail_nickname    = "eid-ua-cal006-include"
  mail_enabled     = false
  security_enabled = true
  description      = "uincludeion group for CA policy CAL006_INCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAL006: Allow Access Only from Specified Locations
# Restricts access to only specified trusted locations for specific accounts.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal006_allow_only_specified_locations" {
  display_name = "acc-test-cal006-all: Only Allow Access from specified locations for specific accounts when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cal006_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cal006_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    locations = {
      include_locations = ["All"]
      # Note: Add specific allowed location IDs to exclude_locations
      exclude_locations = [
        microsoft365_graph_beta_identity_and_access_named_location.allowed_apac_office_only.id,
        microsoft365_graph_beta_identity_and_access_named_location.allowed_emea_office_only.id,
        microsoft365_graph_beta_identity_and_access_named_location.allowed_hazelwood_office_only.id,
      ]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

