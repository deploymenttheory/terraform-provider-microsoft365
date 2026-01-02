# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

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

resource "microsoft365_graph_beta_groups_group" "cal005_exclude" {
  display_name     = "EID_UA_CAL005_EXCLUDE"
  mail_nickname    = "eid-ua-cal005-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAL005_EXCLUDE"
}

# ==============================================================================
# Named Location Dependencies
# ==============================================================================

# Semi-trusted partner networks
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_partner_networks" {
  display_name = "CAL005 Semi-Trusted Partner Networks - ${random_string.suffix.result}"

  ipv4_ranges = [
    "198.18.0.0/24", # Example: Partner A network
    "198.18.1.0/24", # Example: Partner B network
  ]
}

# Semi-trusted public spaces
resource "microsoft365_graph_beta_identity_and_access_named_location" "semi_trusted_public_spaces" {
  display_name = "CAL005 Semi-Trusted Public Spaces - ${random_string.suffix.result}"

  ipv4_ranges = [
    "198.51.100.0/24", # Example: Public WiFi location 1
    "203.0.113.0/24",  # Example: Public WiFi location 2
  ]
}

# ==============================================================================
# Propagation Delay for Named Locations
# ==============================================================================

# Allow time for named locations to propagate in Microsoft Entra ID
resource "time_sleep" "wait_for_named_locations" {
  depends_on = [
    microsoft365_graph_beta_identity_and_access_named_location.semi_trusted_partner_networks,
    microsoft365_graph_beta_identity_and_access_named_location.semi_trusted_public_spaces
  ]

  create_duration = "30s"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAL005: Less-Trusted Locations Require Compliance
# Requires compliant device when accessing from less-trusted locations.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cal005_less_trusted_locations_compliant" {
  display_name = "acc-test-cal005-selected: Grant access for All users on less-trusted locations when Browser and Modern Auth Clients and Compliant ${random_string.suffix.result}"
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
        microsoft365_graph_beta_groups_group.cal005_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = ["Office365"]
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    locations = {
      # Note: Add specific less-trusted location IDs
      include_locations = [
        microsoft365_graph_beta_identity_and_access_named_location.semi_trusted_partner_networks.id,
        microsoft365_graph_beta_identity_and_access_named_location.semi_trusted_public_spaces.id
      ]
      exclude_locations = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }
}

