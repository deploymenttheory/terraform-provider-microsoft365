# Replace notification_url with an HTTPS endpoint that responds to Microsoft Graph validation.
# See: https://learn.microsoft.com/en-us/graph/webhooks#notification-endpoint-validation
resource "microsoft365_graph_beta_change_notifications_subscription" "example" {
  change_type      = "updated"
  notification_url = "https://your-host.example/graph-notify"
  resource         = "users"
  # UTC; must be within the maximum subscription lifetime for the monitored resource.
  expiration_date_time         = "2030-01-01T12:00:00Z"
  client_state                 = "opaque-shared-secret"
  latest_supported_tls_version = "v1_2"
}
