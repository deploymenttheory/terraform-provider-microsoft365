#!/bin/bash

# Import scripts for Microsoft 365 Group Lifecycle Policy
# Replace {policy_id} with the actual policy ID from your Microsoft 365 tenant

# Import the default policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_expiration_policy.default_policy {default_policy_id}

# Import the project groups policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_expiration_policy.project_groups_policy {project_policy_id}

# Import the department groups policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_expiration_policy.department_groups_policy {department_policy_id}

# Import the disabled policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_expiration_policy.disabled_policy {disabled_policy_id}

# Import the compliance policy
terraform import microsoft365_graph_beta_groups_group_lifecycle_expiration_policy.compliance_policy {compliance_policy_id}

# Note: You can import individual policies as needed. Not all policies need to be imported.
# To find policy IDs, you can use Microsoft Graph API or PowerShell:
# Get-MgGroupLifecyclePolicy 