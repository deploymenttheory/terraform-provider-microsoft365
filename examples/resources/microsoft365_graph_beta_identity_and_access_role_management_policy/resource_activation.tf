# Harden activation of the Security Administrator role: activations last at most 4 hours, require MFA,
# a justification and a ticket, and require the "c1" Conditional Access authentication context.
# Approvers receive critical notifications only, sent to an additional mailbox.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "security_administrator" {
  role_definition_id = "194ae4cb-b126-40b2-bd5b-6091b380977d"
  scope_id           = "/"
  scope_type         = "DirectoryRole"

  expiration_rules = {
    Expiration_EndUser_Assignment = {
      is_expiration_required = true
      maximum_duration       = "PT4H"
    }
  }

  enablement_rules = {
    Enablement_EndUser_Assignment = {
      enabled_rules = ["MultiFactorAuthentication", "Justification", "Ticketing"]
    }
  }

  authentication_context_rules = {
    AuthenticationContext_EndUser_Assignment = {
      is_enabled  = true
      claim_value = "c1"
    }
  }

  notification_rules = {
    Notification_Approver_EndUser_Assignment = {
      notification_level            = "Critical"
      is_default_recipients_enabled = true
      notification_recipients       = ["pim-approvers@contoso.com"]
    }
  }
}
