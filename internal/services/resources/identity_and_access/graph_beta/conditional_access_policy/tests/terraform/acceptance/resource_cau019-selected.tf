# ==============================================================================
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

resource "microsoft365_graph_beta_groups_group" "cau019_exclude" {
  display_name     = "EID_UA_CAU019_EXCLUDE"
  mail_nickname    = "eid-ua-cau019-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAU019_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU019: Allow Only Approved Apps for Guests
# Blocks access to all applications for guests except approved apps. This is the
# inverse of CAU003 - allows specific approved apps for guest users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau019_allow_only_approved_apps_guests" {
  display_name = "acc-test-cau019-selected: Only allow approved apps for guests when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau019_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      include_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications = ["All"]
      exclude_applications = [
        "00000002-0000-0000-c000-000000000000",
        "65d91a3d-ab74-42e6-8a2f-0add61688c74",
        "1b912ec3-a9dd-4c4d-a53e-76aa7adb28d7",
        "ea890292-c8c8-4433-b5ea-b09d0668e1a6",
        "0000000c-0000-0000-c000-000000000000",
        "19db86c3-b2b9-44cc-b339-36da233a3be2",
        "4660504c-45b3-4674-a709-71951a6b0763",
        "8c59ead7-d703-4a27-9e55-c96a0054c8d2",
        "2793995e-0a7d-40d7-bd35-6968ba142197",
        "Office365"
      ]
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

