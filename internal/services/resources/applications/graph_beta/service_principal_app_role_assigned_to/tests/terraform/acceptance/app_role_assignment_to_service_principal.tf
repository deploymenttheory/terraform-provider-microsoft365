# Service Principal app role assignment test
# This test assigns the User.Read.All permission from Microsoft Graph to a test service principal
#
# Test Resource Creation Strategy:
# 1. Create a test application
# 2. Create a service principal from that application
# 3. Assign Microsoft Graph app role to the test service principal
# 4. Clean up with hard_delete to ensure no residual resources

# ==============================================================================
# Random Suffix for Unique Names
# ==============================================================================

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Application and Service Principal
# ==============================================================================

# Create a test application
resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "acc-test-approle-sp-${random_string.test_id.result}"
  description  = "Test application for app role assignment acceptance test"
  hard_delete  = true
}

# Wait for application to be fully created
resource "time_sleep" "wait_for_app" {
  depends_on      = [microsoft365_graph_beta_applications_application.test]
  create_duration = "15s"
}

# Create service principal from the test application
# This is the target that will receive the app role assignment
resource "microsoft365_graph_beta_applications_service_principal" "test_target" {
  app_id = microsoft365_graph_beta_applications_application.test.app_id

  depends_on = [time_sleep.wait_for_app]
}

# Wait for service principal to be fully created
resource "time_sleep" "wait_for_sp" {
  depends_on      = [microsoft365_graph_beta_applications_service_principal.test_target]
  create_duration = "15s"
}

# ==============================================================================
# Microsoft Graph Service Principal (Resource with App Roles)
# ==============================================================================

# Data source to get Microsoft Graph service principal (the resource that has app roles)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  app_id = "00000003-0000-0000-c000-000000000000"
}

# ==============================================================================
# App Role Assignment
# ==============================================================================

# Create the app role assignment - User.Read.All
# Assigns Microsoft Graph's User.Read.All permission to our test service principal
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "user_read_all" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.id
  app_role_id                        = "df021288-bdef-4463-88db-98f22de89214" # User.Read.All
  target_service_principal_object_id = microsoft365_graph_beta_applications_service_principal.test_target.id

  depends_on = [time_sleep.wait_for_sp]
}
