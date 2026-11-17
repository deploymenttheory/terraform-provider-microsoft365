# Invalid: sharing_capability must be one of the documented values.
resource "microsoft365_graph_beta_m365_admin_sharepoint_settings" "test" {
  sharing_capability = "everyone"
}
