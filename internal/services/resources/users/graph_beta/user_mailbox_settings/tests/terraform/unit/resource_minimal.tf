resource "microsoft365_graph_beta_users_user_mailbox_settings" "minimal" {
  user_id     = "00000000-0000-0000-0000-000000000001"
  time_zone   = "UTC"
  date_format = "MM/dd/yyyy"
  time_format = "hh:mm tt"
}

