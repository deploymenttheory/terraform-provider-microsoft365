# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
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
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cal001_exclude" {
  display_name     = "EID_UA_CAL001_EXCLUDE"
  mail_nickname    = "eid-ua-cal001-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAL001_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Named Location Dependencies
# ==============================================================================

# High risk countries blocked by client IP
resource "microsoft365_graph_beta_identity_and_access_named_location" "high_risk_countries_blocked_by_client_ip" {
  display_name                          = "CAL001 High Risk Countries (Client IP) - ${random_string.suffix.result}"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

  countries_and_regions = [
    "KP", # North Korea
    "IR", # Iran
  ]
}

# High risk countries blocked by authenticator GPS
resource "microsoft365_graph_beta_identity_and_access_named_location" "high_risk_countries_blocked_by_authenticator_gps" {
  display_name                          = "CAL001 High Risk Countries (GPS) - ${random_string.suffix.result}"
  country_lookup_method                 = "authenticatorAppGps"
  include_unknown_countries_and_regions = false

  countries_and_regions = [
    "KP", # North Korea
    "IR", # Iran
  ]
}

# ==============================================================================
# Propagation Delay for Named Locations
# ==============================================================================

# Allow time for named locations to propagate in Microsoft Entra ID
resource "time_sleep" "wait_for_named_locations" {
  depends_on = [
    microsoft365_graph_beta_identity_and_access_named_location.high_risk_countries_blocked_by_client_ip,
    microsoft365_graph_beta_identity_and_access_named_location.high_risk_countries_blocked_by_authenticator_gps
  ]

  create_duration = "30s"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAL001: Block Specified Locations
# Blocks access from specified untrusted locations for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal001_block_locations" {
  display_name = "acc-test-cal001-all: Block specified locations for All users when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  depends_on = [
    time_sleep.wait_for_named_locations
  ]

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cal001_exclude.id
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
      # Note: Add specific blocked location IDs
      include_locations = [
        microsoft365_graph_beta_identity_and_access_named_location.high_risk_countries_blocked_by_client_ip.id,
        microsoft365_graph_beta_identity_and_access_named_location.high_risk_countries_blocked_by_authenticator_gps.id
      ]
      exclude_locations = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

