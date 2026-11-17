# Invalid: the retention period must be between 30 and 3650 days.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  deleted_user_personal_site_retention_period_in_days = 29
}
