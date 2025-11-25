// Resource must only exist after the lifecycle policy has been created.
// group_id must be a valid Microsoft 365 group.
resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy_assignment" "test" {
  group_id   = microsoft365_graph_beta_groups_group.m365_group.id
  depends_on = [microsoft365_graph_beta_groups_group_lifecycle_expiration_policy.selected]
}

// Dependencies

resource "microsoft365_graph_beta_groups_group_lifecycle_expiration_policy" "selected" {
  group_lifetime_in_days        = 365
  managed_group_types           = "Selected"
  alternate_notification_emails = "admin@deploymenttheory.com;notifications@deploymenttheory.com"
  //overwrite_existing_policy     = true
}

resource "microsoft365_graph_beta_groups_group" "m365_group" {
  display_name          = "acc-m365-group-assigned-${random_string.group_suffix.result}"
  mail_enabled          = true
  security_enabled      = true
  group_types           = ["Unified"]
  description           = "Acceptance test - M365 group with assigned membership"
  mail_nickname         = "accm365g6${random_string.group_suffix.result}"
  is_assignable_to_role = true
  visibility            = "Private"
}

resource "random_string" "group_suffix" {
  length  = 6
  special = false
  upper   = false
}
