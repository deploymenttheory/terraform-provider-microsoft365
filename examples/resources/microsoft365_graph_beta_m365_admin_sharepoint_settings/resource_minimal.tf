# Keep a deleted user's OneDrive for one year and hide the OneDrive sync button.
# Every other setting keeps its current tenant value.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "example" {
  deleted_user_personal_site_retention_period_in_days = 365
  is_sync_button_hidden_on_personal_site              = true
}
