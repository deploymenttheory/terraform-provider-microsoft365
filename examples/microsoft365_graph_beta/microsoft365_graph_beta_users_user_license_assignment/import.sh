#!/bin/bash
# Import using user object ID
terraform import microsoft365_graph_beta_user_license_assignment.example "12345678-1234-1234-1234-123456789012"

# Import using user principal name (UPN)
terraform import microsoft365_graph_beta_user_license_assignment.example "john.doe@example.com" 