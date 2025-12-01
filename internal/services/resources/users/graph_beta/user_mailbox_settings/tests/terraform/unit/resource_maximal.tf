resource "microsoft365_graph_beta_users_user_mailbox_settings" "maximal" {
  user_id                                   = "00000000-0000-0000-0000-000000000002"
  time_zone                                 = "UTC"
  date_format                               = "MM/dd/yyyy"
  time_format                               = "hh:mm tt"
  delegate_meeting_message_delivery_options = "sendToDelegateOnly"

  automatic_replies_setting = {
    status            = "scheduled"
    external_audience = "all"

    scheduled_start_date_time = {
      date_time = "2016-03-14T07:00:00"
      time_zone = "UTC"
    }

    scheduled_end_date_time = {
      date_time = "2016-03-28T07:00:00"
      time_zone = "UTC"
    }

    internal_reply_message = "<html>\n<body>\n<p>I'm at our company's worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n"
    external_reply_message = "<html>\n<body>\n<p>I'm at the Contoso worldwide reunion and will respond to your message as soon as I return.<br>\n</p></body>\n</html>\n"
  }

  language = {
    locale = "en-US"
  }

  working_hours = {
    days_of_week = ["monday", "tuesday", "wednesday", "thursday", "friday"]
    start_time   = "08:00:00"
    end_time     = "17:00:00"

    time_zone = {
      name = "Pacific Standard Time"
    }
  }
}

