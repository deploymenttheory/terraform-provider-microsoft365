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

resource "microsoft365_graph_beta_groups_group" "cau002_exclude" {
  display_name     = "EID_UA_CAU002_EXCLUDE"
  mail_nickname    = "eid-ua-cau002-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAU002_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU002: Require MFA for All Users
# Requires MFA for all users (with admin role exclusions).
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau002_all_users_mfa" {
  display_name = "acc-test-cau002-all: Grant Require MFA for All users when Browser and Modern Auth Clients ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau002_exclude.id
      ]
      include_roles = []
      exclude_roles = [
        "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3",
        "cf1c38e5-3621-4004-a7cb-879624dced7c",
        "c4e39bd9-1100-46d3-8c65-fb160da0071f",
        "25a516ed-2fa0-40ea-a2d0-12923a21473a",
        "aaf43236-0c0d-4d5f-883a-6955382ac081",
        "7698a772-787b-4ac8-901f-60d6b08affd2",
        "158c047a-c907-4556-b7ef-446551a6b5f7",
        "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9",
        "9360feb5-f418-4baa-8175-e2a00bac4301",
        "62e90394-69f5-4237-9190-012177145e10",
        "f2ef992c-3afb-46b9-b7cf-a126ee74c451",
        "729827e3-9c14-49f7-bb1b-9608f156bbb8",
        "8ac3fc64-6eca-42ea-9e69-59f4c7b60eb2",
        "744ec460-397e-42ad-a462-8b3f9747a02c",
        "966707d0-3269-4727-9be2-8c3a10f19b9d",
        "7be44c8a-adaf-4e2a-84d6-ab2649e08a13",
        "e8611ab8-c189-46e8-94e1-60213ab1f814",
        "194ae4cb-b126-40b2-bd5b-6091b380977d",
        "5f2222b1-57c3-48ba-8ad5-d4759f1fde6f",
        "5d6b6bb7-de71-4623-b4af-96380a352509",
        "fe930be7-5e62-47db-91af-98c3a49a38b1",
        "b0f54661-2d74-4c50-afa3-1ec803f12efe",
        "29232cdf-9323-42fd-ade2-1d097af3e4de"
      ]

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []


  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000002" # Maps to "multifactor_authentication"
    }
  }
}

