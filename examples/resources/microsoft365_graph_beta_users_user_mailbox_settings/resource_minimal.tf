# Example 1: Minimal mailbox settings configuration
# This example shows the minimum required configuration for user mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "minimal" {
  user_id   = "john.doe@example.com" # Can be user ID (UUID) or UPN
  time_zone = "UTC"
}

