# Example 2: Maximal mailbox settings configuration
# This example shows all available configuration options for user mailbox settings
resource "microsoft365_graph_beta_users_user_mailbox_settings" "maximal" {
  user_id                                   = "jane.smith@example.com"
  time_zone                                 = "Pacific Standard Time"
  date_format                               = "dd/MM/yyyy"
  time_format                               = "HH:mm"
  delegate_meeting_message_delivery_options = "sendToDelegateOnly"

  # Configure automatic replies (Out of Office)
  automatic_replies_setting = {
    status            = "scheduled"
    external_audience = "all"

    scheduled_start_date_time = {
      date_time = "2024-12-20T00:00:00"
      time_zone = "Pacific Standard Time"
    }

    scheduled_end_date_time = {
      date_time = "2024-12-30T00:00:00"
      time_zone = "Pacific Standard Time"
    }

    internal_reply_message = "<html><body><p>I'm out of office and will respond when I return.</p></body></html>"
    external_reply_message = "<html><body><p>I'm currently out of office. For urgent matters, please contact support@example.com.</p></body></html>"
  }

  # Configure language/locale settings
  language = {
    locale = "en-GB"
  }

  # Configure working hours
  working_hours = {
    days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_time   = "09:00:00"
    end_time     = "17:00:00"

    time_zone = {
      name = "Pacific Standard Time"
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

