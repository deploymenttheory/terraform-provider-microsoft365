#!/bin/bash
# Import format: {user_id}
# The user_id can be either the user's object ID (UUID) or User Principal Name (UPN)

# Example 1: Import using user object ID (UUID)
terraform import microsoft365_graph_beta_users_user_mailbox_settings.example "12345678-1234-1234-1234-123456789012"

# Example 2: Import using User Principal Name (UPN)
terraform import microsoft365_graph_beta_users_user_mailbox_settings.example "john.doe@example.com"

