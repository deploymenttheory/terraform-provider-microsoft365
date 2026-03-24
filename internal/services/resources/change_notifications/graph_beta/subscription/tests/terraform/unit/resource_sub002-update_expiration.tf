# SUB002: same subscription as SUB001 with updated expiration (mocked Graph)
resource "microsoft365_graph_beta_change_notifications_subscription" "sub001_minimal" {
  change_type                  = "updated"
  notification_url             = "https://example.com/webhook"
  resource                     = "users"
  expiration_date_time         = "2030-06-01T12:00:00Z"
  client_state                 = "unit-test-secret"
  latest_supported_tls_version = "v1_2"
}
