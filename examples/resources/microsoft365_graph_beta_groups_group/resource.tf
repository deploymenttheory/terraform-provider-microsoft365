# Security Group Example
resource "microsoft365_graph_beta_group" "security_group" {
  display_name     = "Security Team"
  description      = "Security team members with access to security resources"
  mail_nickname    = "security-team"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Microsoft 365 Group Example
resource "microsoft365_graph_beta_group" "m365_group" {
  display_name     = "Marketing Team"
  description      = "Marketing team collaboration group"
  mail_nickname    = "marketing-team"
  mail_enabled     = true
  security_enabled = true
  group_types      = ["Unified"]
  visibility       = "Public"

  preferred_language      = "en-US"
  preferred_data_location = "US"
  theme                   = "Blue"
  classification          = "General"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Dynamic Security Group Example
resource "microsoft365_graph_beta_group" "dynamic_group" {
  display_name     = "Dynamic IT Department"
  description      = "All users in IT department (dynamic membership)"
  mail_nickname    = "dynamic-it"
  mail_enabled     = false
  security_enabled = true
  group_types      = ["DynamicMembership"]
  visibility       = "Private"

  membership_rule                  = "(user.department -eq \"IT\")"
  membership_rule_processing_state = "On"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Role Assignable Group Example
resource "microsoft365_graph_beta_group" "role_assignable_group" {
  display_name          = "Azure AD Administrators"
  description           = "Group that can be assigned to Azure AD roles"
  mail_nickname         = "aad-admins"
  mail_enabled          = false
  security_enabled      = true
  visibility            = "Private"
  is_assignable_to_role = true

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Distribution Group Example
resource "microsoft365_graph_beta_group" "distribution_group" {
  display_name     = "Company Announcements"
  description      = "Distribution list for company-wide announcements"
  mail_nickname    = "company-announce"
  mail_enabled     = true
  security_enabled = false
  visibility       = "Public"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
} 