# Example 1: Add a user as an owner to a security group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "user_to_security_group" {
  group_id          = "1132b215-826f-42a9-8cfe-1643d19d17fd"  # Security group UUID
  owner_id          = "2243c326-937g-53f0-c9df-2e68f106b901"  # User UUID
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: Add a user as an owner to a Microsoft 365 group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "user_to_m365_group" {
  group_id          = "3354d437-048h-64g1-d0ef-3f79g217c012"  # Microsoft 365 group UUID
  owner_id          = "4465e548-159i-75h2-e1fg-4g80h328d123"  # User UUID
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Add a service principal as an owner to a security group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "service_principal_to_security_group" {
  group_id          = "5576f659-260j-86i3-f2gh-5i91i439e234"  # Security group UUID
  owner_id          = "6687g760-371k-97j4-g3hi-6j02j540f345"  # Service principal UUID
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 4: Add a service principal as an owner to a Microsoft 365 group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "service_principal_to_m365_group" {
  group_id          = "7798h871-482l-08k5-h4ij-7k13k651g456"  # Microsoft 365 group UUID
  owner_id          = "8809i982-593m-19l6-i5jk-8l24l762h567"  # Service principal UUID
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 5: Using data sources to get IDs dynamically
data "microsoft365_graph_beta_groups_group" "target_group" {
  display_name = "Sales Team Security Group"
}

data "microsoft365_graph_beta_user" "target_user" {
  user_principal_name = "john.doe@contoso.com"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "dynamic_assignment" {
  group_id          = data.microsoft365_graph_beta_groups_group.target_group.id
  owner_id          = data.microsoft365_graph_beta_user.target_user.id
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 6: Multiple owner assignments to the same group
resource "microsoft365_graph_beta_groups_group_owner_assignment" "multiple_user_owners" {
  for_each = toset([
    "3354d437-048h-64g1-d0ef-3f79g217c012",  # User 1
    "4465e548-159i-75h2-e1fg-4g80h328d123",  # User 2
    "5576f659-260j-86i3-f2gh-5i91i439e234"   # User 3
  ])

  group_id          = "7798h871-482l-08k5-h4ij-7k13k651g456"  # Target security group
  owner_id          = each.value
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 7: Mixed owner types (users and service principals) for the same group
locals {
  owners = [
    {
      id   = "9910j093-604n-20m7-j6kl-9m35m873i678"
      type = "User"
    },
    {
      id   = "0021k104-715o-31n8-k7lm-0n46n984j789"
      type = "User"
    },
    {
      id   = "1132l215-826p-42o9-l8mn-1d57o095k890"
      type = "ServicePrincipal"
    }
  ]
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "mixed_owner_types" {
  for_each = { for idx, owner in local.owners : "${owner.type}_${idx}" => owner }

  group_id          = "2243m326-937q-53p0-m9no-2e68p106l901"  # Target group
  owner_id          = each.value.id
  owner_object_type = each.value.type

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 8: Creating a group and immediately adding owners
resource "microsoft365_graph_beta_groups_group" "example_group" {
  display_name     = "Example Project Team"
  mail_nickname    = "example-project-team"
  description      = "Security group for Example Project Team"
  security_enabled = true
  mail_enabled     = false
  group_types      = []
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "project_team_owners" {
  depends_on = [microsoft365_graph_beta_groups_group.example_group]

  for_each = toset([
    "1132b215-826f-42a9-8cfe-1643d19d17fd",  # Project Lead
    "2243c326-937g-53f0-c9df-2e68f106b901",  # Team Manager
  ])

  group_id          = microsoft365_graph_beta_groups_group.example_group.id
  owner_id          = each.value
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 9: Service principal from managed identity as group owner
data "azuread_service_principal" "managed_identity" {
  display_name = "my-app-managed-identity"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "managed_identity_owner" {
  group_id          = "3354d437-048h-64g1-d0ef-3f79g217c012"  # Target group
  owner_id          = data.azuread_service_principal.managed_identity.object_id
  owner_object_type = "ServicePrincipal"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 10: Conditional owner assignment based on group type
data "microsoft365_graph_beta_groups_group" "conditional_group" {
  display_name = "Conditional Target Group"
}

resource "microsoft365_graph_beta_groups_group_owner_assignment" "conditional_owner" {
  # Only add owner if the group is a security group
  count = contains(data.microsoft365_graph_beta_groups_group.conditional_group.group_types, "Unified") ? 0 : 1

  group_id          = data.microsoft365_graph_beta_groups_group.conditional_group.id
  owner_id          = "4465e548-159i-75h2-e1fg-4g80h328d123"  # User UUID
  owner_object_type = "User"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
} 