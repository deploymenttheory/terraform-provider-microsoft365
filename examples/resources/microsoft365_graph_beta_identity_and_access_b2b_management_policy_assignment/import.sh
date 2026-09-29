#!/bin/bash
# Import using the composite ID format: {b2b_management_policy_id}/{directory_object_id}
# The Microsoft Graph API does not return an assignment-specific ID.

# {b2b_management_policy_id} - GUID of the B2B management policy
# {directory_object_id} - Object ID of the application or service principal
terraform import microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment.application 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002
