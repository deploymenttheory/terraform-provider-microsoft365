resource "microsoft365_graph_beta_device_and_app_management_ios_ipados_web_clip" "example" {
  display_name  = "Company Portal"
  description   = "Company Portal Web Clip"
  publisher     = "Contoso IT"
  app_url       = "https://portal.contoso.com"
  
  full_screen_enabled               = true
  ignore_manifest_scope             = true
  pre_composed_icon_enabled         = true
  use_managed_browser               = false
  target_application_bundle_identifier = "com.apple.mobilesafari"
  
  # Optional fields
  developer               = "Contoso Development Team"
  notes                   = "Use this web clip to access the company portal"
  owner                   = "IT Department"
  privacy_information_url = "https://privacy.contoso.com"
  information_url         = "https://help.contoso.com/portal"
  is_featured             = true
  
  # Categories (requires existing category or use inbuilt categories)
  # categories = ["3263b037-b2f7-4383-87b8-7515d30e5a76", "85e202d5-e967-4cad-9ae9-603149b5d258"]
  
  # Role scope tags (requires existing scope tags)
  # role_scope_tag_ids = ["0", "1"]
  
  # App icon (uncomment to use)
  # app_icon {
  #   icon_file_path_source = "./path/to/icon.png"
  #   # OR
  #   # icon_url_source = "https://example.com/icon.png"
  # }
  
  timeouts =  {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
} 