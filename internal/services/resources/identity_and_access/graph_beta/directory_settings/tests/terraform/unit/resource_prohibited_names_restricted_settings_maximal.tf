resource "microsoft365_graph_beta_identity_and_access_directory_settings" "prohibited_names_restricted_settings" {
  template_type               = "Prohibited Names Restricted Settings"
  overwrite_existing_settings = true

  prohibited_names_restricted_settings {
    custom_allowed_sub_strings_list   = "contoso,fabrikam,northwind"
    custom_allowed_whole_words_list   = "ContosoApp,FabrikamSolution,NorthwindTraders"
    do_not_validate_against_trademark = true
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

