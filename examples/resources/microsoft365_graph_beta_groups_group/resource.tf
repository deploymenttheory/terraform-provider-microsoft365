# Example 1: Basic Security Group with Assigned Membership
# Creates a standard security group where members are manually assigned.
# This is the most common type of security group used for access control.
resource "microsoft365_graph_beta_groups_group" "security_basic" {
  display_name     = "Engineering Team"
  mail_nickname    = "engineering-team"
  mail_enabled     = false
  security_enabled = true
  description      = "Security group for engineering team members"
}

# Example 2: Security Group with Dynamic User Membership
# Creates a security group that automatically adds/removes users based on a membership rule.
# Useful for automatically managing group membership based on user attributes.
resource "microsoft365_graph_beta_groups_group" "security_dynamic_users" {
  display_name                     = "Active Employees"
  mail_nickname                    = "active-employees"
  mail_enabled                     = false
  security_enabled                 = true
  description                      = "Security group containing all active employees"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}

# Example 3: Security Group with Dynamic Device Membership
# Creates a security group that automatically includes devices based on a membership rule.
# Ideal for device management scenarios like Conditional Access or Intune policies.
resource "microsoft365_graph_beta_groups_group" "security_dynamic_devices" {
  display_name                     = "Corporate Managed Devices"
  mail_nickname                    = "corporate-devices"
  mail_enabled                     = false
  security_enabled                 = true
  description                      = "Security group containing all corporate managed devices"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(device.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}

# Example 4: Role-Assignable Security Group
# Creates a security group that can be assigned to Entra ID roles.
# Note: Requires elevated permissions and visibility must be "Private".
# Once created, is_assignable_to_role cannot be changed.
resource "microsoft365_graph_beta_groups_group" "security_role_assignable" {
  display_name          = "Privileged Access Administrators"
  mail_nickname         = "privileged-admins"
  mail_enabled          = false
  security_enabled      = true
  description           = "Security group for privileged access administration"
  is_assignable_to_role = true
  visibility            = "Private"
}

# Example 5: Microsoft 365 Group with Dynamic User Membership
# Creates a Microsoft 365 group (formerly Office 365 group) with automatic membership.
# Includes Teams, SharePoint, Outlook, and other Microsoft 365 services.
resource "microsoft365_graph_beta_groups_group" "m365_dynamic_users" {
  display_name                     = "Marketing Department"
  mail_nickname                    = "marketing-dept"
  mail_enabled                     = true
  security_enabled                 = true
  group_types                      = ["Unified", "DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
  visibility                       = "Private"
}

# Example 6: Microsoft 365 Group with Role Assignment
# Creates a Microsoft 365 group that can be assigned to Entra ID roles.
# Combines collaboration features with privileged access management.
# Note: Requires elevated permissions and visibility must be "Private".
resource "microsoft365_graph_beta_groups_group" "m365_role_assignable" {
  display_name          = "Executive Leadership Team"
  mail_nickname         = "executive-team"
  mail_enabled          = true
  security_enabled      = true
  group_types           = ["Unified"]
  description           = "Microsoft 365 group for executive leadership"
  is_assignable_to_role = true
  visibility            = "Private"
}