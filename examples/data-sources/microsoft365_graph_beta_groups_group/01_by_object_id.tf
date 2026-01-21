# Example 1: Look up group by object_id
# This example shows all possible output attributes

data "microsoft365_graph_beta_groups_group" "by_object_id" {
  object_id = "12345678-1234-1234-1234-123456789012"
}

# All available outputs
output "group_id" {
  description = "The unique identifier for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.id
}

output "object_id" {
  description = "The object ID of the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.object_id
}

output "display_name" {
  description = "The display name for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.display_name
}

output "description" {
  description = "The optional description of the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.description
}

output "classification" {
  description = "A classification for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.classification
}

output "mail_nickname" {
  description = "The mail alias for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.mail_nickname
}

output "mail_enabled" {
  description = "Whether the group is mail-enabled"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.mail_enabled
}

output "security_enabled" {
  description = "Whether the group is a security group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.security_enabled
}

output "group_types" {
  description = "List of group types (e.g., DynamicMembership, Unified)"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.group_types
}

output "visibility" {
  description = "Group join policy and content visibility"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.visibility
}

output "assignable_to_role" {
  description = "Whether group can be assigned to Azure AD role"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.assignable_to_role
}

output "membership_rule" {
  description = "The rule for dynamic membership"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.membership_rule
}

output "membership_rule_processing_state" {
  description = "Dynamic membership processing state (On/Paused)"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.membership_rule_processing_state
}

output "created_date_time" {
  description = "When the group was created"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.created_date_time
}

output "mail" {
  description = "The SMTP address for the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.mail
}

output "proxy_addresses" {
  description = "Email addresses that direct to the same mailbox"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.proxy_addresses
}

output "assigned_licenses" {
  description = "Licenses assigned to the group"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.assigned_licenses
}

output "has_members_with_license_errors" {
  description = "Whether members have license errors"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.has_members_with_license_errors
}

output "hide_from_address_lists" {
  description = "Whether hidden from Outlook address lists"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.hide_from_address_lists
}

output "hide_from_outlook_clients" {
  description = "Whether hidden from Outlook clients"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.hide_from_outlook_clients
}

output "onpremises_sync_enabled" {
  description = "Whether synced from on-premises directory"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_sync_enabled
}

output "onpremises_last_sync_date_time" {
  description = "Last sync time from on-premises"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_last_sync_date_time
}

output "onpremises_sam_account_name" {
  description = "On-premises SAM account name"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_sam_account_name
}

output "onpremises_domain_name" {
  description = "On-premises FQDN"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_domain_name
}

output "onpremises_netbios_name" {
  description = "On-premises NetBIOS name"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_netbios_name
}

output "onpremises_security_identifier" {
  description = "On-premises security identifier (SID)"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.onpremises_security_identifier
}

output "members" {
  description = "List of member object IDs"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.members
}

output "owners" {
  description = "List of owner object IDs"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.owners
}

output "dynamic_membership_enabled" {
  description = "Whether dynamic membership is enabled"
  value       = data.microsoft365_graph_beta_groups_group.by_object_id.dynamic_membership_enabled
}
