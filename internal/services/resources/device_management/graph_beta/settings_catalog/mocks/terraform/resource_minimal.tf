resource "microsoft365_graph_beta_device_management_settings_catalog" "minimal" {
  name        = "Minimal Settings Catalog"
  description = "Minimal settings catalog policy"
  platforms   = "windows10"
} 