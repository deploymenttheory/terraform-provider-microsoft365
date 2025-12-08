# Service Principal app role assignment test
# This test assigns the User.Read.All permission from Microsoft Graph to a service principal
#
# Dependency Strategy:
# - Uses microsoft365 provider data sources to get service principals (this provider first)
# - The test assigns a permission to an existing service principal found via data source
# - Uses random provider for test uniqueness (required by test framework)

########################################################################################
# Dependencies
########################################################################################

# Random UUID for test uniqueness
resource "random_uuid" "test" {}

# Data source to get Microsoft Graph service principal (the resource that has app roles)
# IMPORTANT: Use app_id filter to get the ACTUAL Microsoft Graph SP (00000003-0000-0000-c000-000000000000)
# display_name filter returns multiple results (PowerShell, Change Tracking, Connectors Core, etc.)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  filter_type  = "app_id"
  filter_value = "00000003-0000-0000-c000-000000000000"
}

# Get the test service principal (SP-CPSS-GLBL-AGENTS-C-01) that will receive the app role
# This is the agents SP which we know exists and can receive app roles
data "microsoft365_graph_beta_applications_service_principal" "target" {
  filter_type  = "app_id"
  filter_value = "5b64cedc-ccc7-4896-8b1c-24a1bf84b101"
}

########################################################################################
# Test Resource
########################################################################################

# Create the app role assignment - User.Read.All
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "user_read_all" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "df021288-bdef-4463-88db-98f22de89214"
  target_service_principal_object_id = data.microsoft365_graph_beta_applications_service_principal.target.items[0].id
}
