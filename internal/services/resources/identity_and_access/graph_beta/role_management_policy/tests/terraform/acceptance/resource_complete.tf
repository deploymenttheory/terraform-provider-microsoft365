# Manages every supported rule of the Usage Summary Reports Reader PIM policy. Every rule keeps the
# Microsoft Entra default except Expiration_Admin_Eligibility, whose maximum duration drops from P365D to P180D.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "75934031-6c7e-415a-99d7-48dbd49e875e"
  scope_id           = "/"
  scope_type         = "DirectoryRole"

  expiration_rules = {
    Expiration_Admin_Eligibility = {
      is_expiration_required = false
      maximum_duration       = "P180D"
    }
    Expiration_Admin_Assignment = {
      is_expiration_required = false
      maximum_duration       = "P180D"
    }
    Expiration_EndUser_Assignment = {
      is_expiration_required = true
      maximum_duration       = "PT8H"
    }
  }

  enablement_rules = {
    Enablement_Admin_Eligibility = {
      enabled_rules = []
    }
    Enablement_Admin_Assignment = {
      enabled_rules = ["Justification"]
    }
    Enablement_EndUser_Assignment = {
      enabled_rules = ["Justification"]
    }
  }

  notification_rules = {
    Notification_Admin_Admin_Eligibility = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Requestor_Admin_Eligibility = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Approver_Admin_Eligibility = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Admin_Admin_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Requestor_Admin_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Approver_Admin_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Admin_EndUser_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Requestor_EndUser_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
    Notification_Approver_EndUser_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
  }

  authentication_context_rules = {
    AuthenticationContext_EndUser_Assignment = {
      is_enabled = false
    }
  }
}
