# Example: Assign an app role to a service principal
# This grants the "User.Read.All" permission from Microsoft Graph to a client service principal

resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "example" {
  # The Object ID of the service principal that exposes the app role (e.g., Microsoft Graph)
  resource_object_id = "00000003-0000-0000-c000-000000000000" # Microsoft Graph service principal Object ID

  # The app role ID to assign (e.g., User.Read.All = df021288-bdef-4463-88db-98f22de89214)
  app_role_id = "df021288-bdef-4463-88db-98f22de89214"

  # The Object ID of the service principal being granted the app role
  target_service_principal_object_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Your application's service principal Object ID

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

# Example: Assign app role to a group
# This grants application permissions to all members of the security group

resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "group_assignment" {
  # The Object ID of the service principal that exposes the app role
  resource_object_id = var.resource_service_principal_object_id

  # The app role ID from the resource's appRoles collection
  app_role_id = var.app_role_id

  # The Object ID of the security group being granted the app role
  target_service_principal_object_id = var.security_group_id
}

# Example: Default app role assignment (no specific role)
# Use this when the application doesn't define specific app roles

resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "default_role" {
  resource_object_id = var.resource_service_principal_object_id

  # Default app role ID when no specific roles are defined
  app_role_id = "00000000-0000-0000-0000-000000000000"

  target_service_principal_object_id = var.client_service_principal_object_id
}
