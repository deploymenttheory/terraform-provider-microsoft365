# ==============================================================================
# SUB001: acceptance create — aligned with Microsoft Graph subscription example payload
# (changeType, notificationUrl, resource, clientState, latestSupportedTlsVersion).
# Expiration uses plan time + offset; the doc sample date is not used.
# See https://learn.microsoft.com/en-us/graph/api/subscription-post-subscriptions
# ==============================================================================

resource "microsoft365_graph_beta_change_notifications_subscription" "sub001_mail" {
  change_type                  = "created"
  notification_url             = "https://webhook.azurewebsites.net/api/send/myNotifyClient"
  resource                     = "me/mailFolders('Inbox')/messages"
  expiration_date_time         = "48h"
  client_state                 = "secretClientValue"
  latest_supported_tls_version = "v1_2"
}
