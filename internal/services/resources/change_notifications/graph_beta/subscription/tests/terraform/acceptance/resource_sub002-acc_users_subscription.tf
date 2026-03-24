# ==============================================================================
# SUB002: acceptance update — longer expiration (same subscription as SUB001)
# ==============================================================================

locals {
  expiration = timeadd(plantimestamp(), "72h")
}

resource "microsoft365_graph_beta_change_notifications_subscription" "sub001_mail" {
  change_type                  = "created"
  notification_url             = "https://webhook.azurewebsites.net/api/send/myNotifyClient"
  resource                     = "me/mailFolders('Inbox')/messages"
  expiration_date_time         = local.expiration
  client_state                 = "secretClientValue"
  latest_supported_tls_version = "v1_2"
}
