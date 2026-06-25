#!/bin/bash
# Import using the composite ID format: {service_principal_id}/{token_lifetime_policy_id}
# The Microsoft Graph API does not return an assignment-specific ID.
# Use the Object ID of the service principal and the GUID of the token lifetime policy.

# {service_principal_id} - Object ID of the service principal
# {token_lifetime_policy_id} - GUID of the token lifetime policy
terraform import microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment.example 00000000-0000-0000-0000-000000000001/00000000-0000-0000-0000-000000000002
