resource "microsoft365_graph_beta_applications_tenant_app_management_policy" "minimal" {
  is_enabled = true

  application_restrictions = {
    password_credentials = [
      {
        restriction_type                           = "passwordLifetime"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        max_lifetime                               = "P90D"
      }
    ]
  }
}
