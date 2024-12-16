resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "whatsapp" {
  package_identifier = "9NKSQGP7F2NH" # The unique identifier for the app obtained from msft app store

  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }
}