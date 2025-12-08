# Example: Assign Microsoft Entra Agent ID permissions to a service principal
# These permissions are required for managing AI agent identities in Microsoft Entra
# Reference: https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities

# Get the Microsoft Graph service principal (resource that defines the permissions)
data "microsoft365_graph_beta_applications_service_principal" "msgraph" {
  filter_type  = "display_name"
  filter_value = "Microsoft Graph"
}

# Example: Agent Identity Read permissions
# AgentIdentity.Read.All - Read all agent identities
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_identity_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "b2b8f011-2898-4234-9092-5059f6c1ebfa" # AgentIdentity.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Identity ReadWrite permissions
# AgentIdentity.ReadWrite.All - Read and write all agent identities
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_identity_readwrite" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "dcf7150a-88d4-4fe6-9be1-c2744c455397" # AgentIdentity.ReadWrite.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Identity Delete/Restore permissions
# AgentIdentity.DeleteRestore.All - Delete and restore agent identities
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_identity_delete_restore" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "5b016f9b-18eb-41d4-869a-66931914d1c8" # AgentIdentity.DeleteRestore.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Registry permissions
# AgentCardManifest.Read.All - Read agent card manifests
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_card_manifest_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "3ee18438-e6e5-4858-8f1c-d7b723b45213" # AgentCardManifest.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Collection permissions
# AgentCollection.Read.All - Read all agent collections
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_collection_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "e65ee1da-d1d5-467b-bdd0-3e9bb94e6e0c" # AgentCollection.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Instance permissions
# AgentInstance.Read.All - Read all agent instances
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_instance_read" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "799a4732-85b8-4c67-b048-75f0e88a232b" # AgentInstance.Read.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Example: Agent Instance ReadWrite permissions
# AgentInstance.ReadWrite.All - Read and write all agent instances
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "agent_instance_readwrite" {
  resource_object_id                 = data.microsoft365_graph_beta_applications_service_principal.msgraph.items[0].id
  app_role_id                        = "07abdd95-78dc-4353-bd32-09f880ea43d0" # AgentInstance.ReadWrite.All (App-only)
  target_service_principal_object_id = var.client_service_principal_object_id
}

# Variable for the target service principal
variable "client_service_principal_object_id" {
  description = "The Object ID of the service principal to grant agent identity permissions to"
  type        = string
}

