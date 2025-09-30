# Example 1: Add a user to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "user_to_security_group" {
  group_id           = "1132b215-826f-42a9-8cfe-1643d19d17fd" # Security group UUID
  member_id          = "2243c326-937g-53f0-c9df-2e68f106b901" # User UUID
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: Add a user to a Microsoft 365 group
resource "microsoft365_graph_beta_groups_group_member_assignment" "user_to_m365_group" {
  group_id           = "3354d437-048h-64g1-d0ef-3f79g217c012" # Microsoft 365 group UUID
  member_id          = "4465e548-159i-75h2-e1fg-4g80h328d123" # User UUID
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Add a security group to another security group (nested groups)
resource "microsoft365_graph_beta_groups_group_member_assignment" "group_to_security_group" {
  group_id           = "5576f659-260j-86i3-f2gh-5i91i439e234" # Parent security group UUID
  member_id          = "6687g760-371k-97j4-g3hi-6j02j540f345" # Member security group UUID
  member_object_type = "Group"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 4: Add a device to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "device_to_security_group" {
  group_id           = "7798h871-482l-08k5-h4ij-7k13k651g456" # Security group UUID
  member_id          = "8809i982-593m-19l6-i5jk-8l24l762h567" # Device UUID
  member_object_type = "Device"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 5: Add a service principal to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "service_principal_to_security_group" {
  group_id           = "9910j093-604n-20m7-j6kl-9m35m873i678" # Security group UUID
  member_id          = "0021k104-715o-31n8-k7lm-0n46n984j789" # Service principal UUID
  member_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 6: Add an organizational contact to a security group
resource "microsoft365_graph_beta_groups_group_member_assignment" "contact_to_security_group" {
  group_id           = "1132l215-826p-42o9-l8mn-1d57o095k890" # Security group UUID
  member_id          = "2243m326-937q-53p0-m9no-2e68p106l901" # Organizational contact UUID
  member_object_type = "OrganizationalContact"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 7: Using data sources to get IDs dynamically
data "microsoft365_graph_beta_groups_group" "target_group" {
  display_name = "Sales Team Security Group"
}

data "microsoft365_graph_beta_user" "target_user" {
  user_principal_name = "john.doe@contoso.com"
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "dynamic_assignment" {
  group_id           = data.microsoft365_graph_beta_groups_group.target_group.id
  member_id          = data.microsoft365_graph_beta_user.target_user.id
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 8: Multiple member assignments to the same group
resource "microsoft365_graph_beta_groups_group_member_assignment" "multiple_users" {
  for_each = toset([
    "3354d437-048h-64g1-d0ef-3f79g217c012", # User 1
    "4465e548-159i-75h2-e1fg-4g80h328d123", # User 2
    "5576f659-260j-86i3-f2gh-5i91i439e234"  # User 3
  ])

  group_id           = "7798h871-482l-08k5-h4ij-7k13k651g456" # Target security group
  member_id          = each.value
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 9: Creating a group and immediately adding members
resource "microsoft365_graph_beta_groups_group" "example_group" {
  display_name     = "Example Project Team"
  mail_nickname    = "example-project-team"
  description      = "Security group for Example Project Team members"
  security_enabled = true
  mail_enabled     = false
  group_types      = []
}

resource "microsoft365_graph_beta_groups_group_member_assignment" "project_team_members" {
  depends_on = [microsoft365_graph_beta_groups_group.example_group]

  for_each = toset([
    "1132b215-826f-42a9-8cfe-1643d19d17fd", # Project Manager
    "2243c326-937g-53f0-c9df-2e68f106b901", # Developer 1
    "3354d437-048h-64g1-d0ef-3f79g217c012"  # Developer 2
  ])

  group_id           = microsoft365_graph_beta_groups_group.example_group.id
  member_id          = each.value
  member_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
} 