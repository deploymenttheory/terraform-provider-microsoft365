# Minimal configuration: only the deleted user OneDrive retention period is managed.
# Every other setting keeps its current tenant value.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  deleted_user_personal_site_retention_period_in_days = 60
}
