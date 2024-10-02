resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "example" {
  package_identifier = "9nzvdkpmr9rd"  # The unique identifier for the app

  # Install experience settings
  install_experience {
    run_as_account = "user"  # Can be 'system' or 'user'
  }

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Role scope tags
  role_scope_tag_ids = ["tag-id-1", "tag-id-2"]

  # Assignments (example of group assignments)
  assignments = {
    target_group_id = "group-id-1"
    install_intent  = "available"
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }
}
