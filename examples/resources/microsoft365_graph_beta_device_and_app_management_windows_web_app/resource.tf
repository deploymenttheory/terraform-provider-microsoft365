resource "microsoft365_graph_beta_device_and_app_management_windows_web_app" "example" {
  display_name = "Company Portal"
  description  = "Company Portal Web App"
  publisher    = "Contoso IT"
  app_url      = "https://portal.contoso.com"

  # Optional fields
  developer               = "Contoso Development Team"
  notes                   = "Use this web app to access the company portal"
  owner                   = "IT Department"
  privacy_information_url = "https://privacy.contoso.com"
  information_url         = "https://help.contoso.com/portal"
  is_featured             = true

  # Categories (requires existing category or use inbuilt categories)
  # categories = ["Business", "Productivity"]

  # Role scope tags (requires existing scope tags)
  # role_scope_tag_ids = ["0", "1"]

  # App icon (uncomment to use)
  # app_icon {
  #   icon_file_path_source = "./path/to/icon.png"
  #   # OR
  #   # icon_url_source = "https://example.com/icon.png"
  # }

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
} 