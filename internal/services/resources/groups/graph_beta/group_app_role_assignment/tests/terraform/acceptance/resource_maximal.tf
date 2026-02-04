# Group app role assignment test - maximal with custom timeouts
# Creates isolated test resources to avoid dependencies on production resources

# ==============================================================================
# Random Suffix for Unique Names
# ==============================================================================

resource "random_string" "maximal_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Group (Target for App Role Assignment)
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "maximal" {
  display_name     = "acc-test-group-approle-maximal-${random_string.maximal_suffix.result}"
  mail_nickname    = "approletestmax${random_string.maximal_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group for app role assignment with custom timeouts"
  hard_delete      = true
}

# ==============================================================================
# MileIQ Admin Center Service Principal (has app roles)
# ==============================================================================

# Get MileIQ Admin Center service principal (has app roles available for groups)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  display_name = "MileIQ Admin Center"
}

# Wait for group to be fully created
resource "time_sleep" "wait_for_group" {
  depends_on      = [microsoft365_graph_beta_groups_group.maximal]
  create_duration = "15s"
}

# ==============================================================================
# App Role Assignment
# ==============================================================================

# Assign Microsoft Graph's app role to the group with custom timeouts
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "maximal" {
  target_group_id    = microsoft365_graph_beta_groups_group.maximal.id
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.msgraph.id
  app_role_id        = "ea358ccf-c4a8-48ac-8b94-2558ae2f7a5c" # mdladmincenterrole.admin

  timeouts = {
    create = "300s"
    read   = "300s"
    delete = "300s"
  }

  depends_on = [time_sleep.wait_for_group]
}
