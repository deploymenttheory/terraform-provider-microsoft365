resource "microsoft365_graph_beta_device_and_app_management_macos_vpp_app" "example" {
  display_name                = "Example MacOS VPP App"
  description                 = "Example MacOS VPP application managed by Terraform"
  publisher                   = "Example Publisher"
  bundle_id                   = "com.example.macosvppapp"
  vpp_token_id                = "00000000-0000-0000-0000-000000000000" # Replace with actual VPP token ID
  vpp_token_organization_name = "Example Organization"
  vpp_token_account_type      = "business" # Possible values: business, education
  vpp_token_apple_id          = "example@organization.com"

  # Optional fields
  information_url         = "https://example.com/app-info"
  privacy_information_url = "https://example.com/privacy"
  owner                   = "IT Department"
  developer               = "Example Developer"
  notes                   = "Managed by Terraform"
  is_featured             = false

  # Role scope tags (optional)
  role_scope_tag_ids = ["0"]

  # App icon (optional)
  # app_icon {
  #   icon_file_path_source = "/path/to/icon.png"
  #   # OR
  #   # icon_url_source     = "https://example.com/icon.png"
  # }

  # Licensing type (optional)
  licensing_type {
    support_user_licensing    = true
    support_device_licensing  = true
    supports_user_licensing   = true
    supports_device_licensing = true
  }

  # Categories (optional)
  # categories = ["00000000-0000-0000-0000-000000000000"] # Replace with actual category IDs

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
} 