# Example: macOS PKG App Resource

resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "example_app" {
  display_name    = "Example macOS PKG App"
  description     = "This is an example macOS PKG app managed by Terraform"
  publisher       = "Example Publisher"
  file_name       = "example_app.pkg"
  primary_bundle_id = "com.example.app"
  primary_bundle_version = "1.0.0"

  large_icon = {
    type  = "image/png"
    value = filebase64("path/to/icon.png")
  }

  is_featured = true
  privacy_information_url = "https://example.com/privacy"
  information_url = "https://example.com/info"
  owner = "IT Department"
  developer = "Example Developer"
  notes = "Example notes for the app"

  role_scope_tag_ids = ["tag1", "tag2"]

  ignore_version_detection = false

  included_apps = [
    {
      bundle_id = "com.example.includedapp1"
      bundle_version = "1.0.0"
    },
    {
      bundle_id = "com.example.includedapp2"
      bundle_version = "2.0.0"
    }
  ]

  minimum_supported_operating_system = {
    v10_14 = true
    v10_15 = true
    v11_0 = true
    v12_0 = true
    v13_0 = true
    v14_0 = true
  }

  pre_install_script = {
    script_content = file("path/to/pre_install_script.sh")
  }

  post_install_script = {
    script_content = file("path/to/post_install_script.sh")
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}