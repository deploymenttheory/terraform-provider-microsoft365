data "microsoft365_graph_beta_device_and_app_management_mobile_app" "win32_apps" {
  filter_type     = "all"
  app_type_filter = "win32LobApp"
}

