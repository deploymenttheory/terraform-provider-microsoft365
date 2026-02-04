# To find service principals in your tenant, use:
# Get-MgServicePrincipal -Filter "appId eq '00000003-0000-0000-c000-000000000000'" | Select-Object Id, DisplayName, AppId
# The "Id" property is what you need for resource_id

# Get the Microsoft Graph service principal (resource that defines the permissions)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  display_name = "Microsoft Graph"
}

# Example 1: Assign default access to Microsoft Graph
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "graph_default" {
  target_group_id    = "12345678-1234-1234-1234-123456789012"                                 # UUID of your group
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.msgraph.id # Object ID of Microsoft Graph service principal
  app_role_id        = "00000000-0000-0000-0000-000000000000"                                 # Default role (basic access)

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Example 2: Assign specific app role to a group
# Common Microsoft Graph App Roles:
# - Directory.Read.All: df021288-bdef-4463-88db-98f22de89214
# - User.Read.All: a154be20-db9c-4678-8ab7-66f6cc099a59
# - Group.Read.All: 5b567255-7703-4780-807c-7be8301ae99b
resource "microsoft365_graph_beta_groups_group_app_role_assignment" "graph_directory_read" {
  target_group_id    = "12345678-1234-1234-1234-123456789012"                                 # UUID of your group
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.msgraph.id # Object ID of Microsoft Graph service principal
  app_role_id        = "df021288-bdef-4463-88db-98f22de89214"                                 # Directory.Read.All
}

# Example 3: Assign role to SharePoint Online
# Get the SharePoint Online service principal
# The App ID for SharePoint Online is always: 00000003-0000-0ff1-ce00-000000000000
data "microsoft365_graph_beta_applications_service_principal" "sharepoint" {
  display_name = "SharePoint Online"
}

resource "microsoft365_graph_beta_groups_group_app_role_assignment" "sharepoint" {
  target_group_id    = "12345678-1234-1234-1234-123456789012"                                    # UUID of your group
  resource_object_id = data.microsoft365_graph_beta_applications_service_principal.sharepoint.id # Object ID of SharePoint service principal
  app_role_id        = "678536fe-1083-478a-9c59-b99265e6b0d3"                                    # Example SharePoint app role
}

