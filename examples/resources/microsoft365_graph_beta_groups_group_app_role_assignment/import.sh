#!/bin/bash

# Import scripts for Microsoft 365 Group App Role Assignment
# {group_id}/{assignment_id}

# Import a group app role assignment
terraform import microsoft365_graph_beta_groups_group_app_role_assignment.example {group_id}/{assignment_id}

# Import a read-only app role assignment
terraform import microsoft365_graph_beta_groups_group_app_role_assignment.read_only {group_id}/{assignment_id}

# Import a full-access app role assignment
terraform import microsoft365_graph_beta_groups_group_app_role_assignment.full_access {group_id}/{assignment_id}

# Note: You can import individual app role assignments as needed.
# To find assignment IDs, you can use Microsoft Graph API or PowerShell:
# Get-MgGroupAppRoleAssignment -GroupId {group_id}

