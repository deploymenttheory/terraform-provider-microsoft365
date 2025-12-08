# Import an existing app role assignment
# The import ID is the app role assignment ID returned by Microsoft Graph

# {app_role_assignment_id} - The unique identifier of the app role assignment
terraform import microsoft365_graph_beta_applications_service_principal_app_role_assigned_to.example {app_role_assignment_id}

