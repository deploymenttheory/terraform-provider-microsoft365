resource "microsoft365_graph_beta_identity_and_access_directory_settings" "prohibited_names_settings" {
  template_type               = "Prohibited Names Settings"
  overwrite_existing_settings = true

  prohibited_names_settings {
    custom_blocked_sub_strings_list = "microsoft,windows,azure,office"
    custom_blocked_whole_words_list = "Microsoft,Windows,Azure,Office365,Outlook"
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

