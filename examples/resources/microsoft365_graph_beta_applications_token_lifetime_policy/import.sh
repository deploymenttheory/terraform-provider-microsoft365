#!/bin/bash
# Import using the token lifetime policy ID (GUID) from Microsoft Graph

# {id} - The unique identifier (GUID) of the token lifetime policy
terraform import microsoft365_graph_beta_applications_token_lifetime_policy.example 00000000-0000-0000-0000-000000000000
