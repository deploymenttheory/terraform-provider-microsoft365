resource "microsoft365_graph_beta_identity_and_access_directory_settings" "group_unified" {
  template_type               = "Group.Unified"
  overwrite_existing_settings = true

  group_unified {
    # Naming Policy Settings
    prefix_suffix_naming_requirement = "GRP_[GroupName]_[Department]"
    custom_blocked_words_list        = "CEO,President,Admin,Executive,Confidential"
    enable_ms_standard_blocked_words = true

    # Group Creation Settings
    enable_group_creation           = false
    group_creation_allowed_group_id = "12345678-1234-1234-1234-123456789012"

    # Guest Access Settings
    allow_guests_to_access_groups  = false
    allow_guests_to_be_group_owner = false
    allow_to_add_guests            = false
    guest_usage_guidelines_url     = "https://contoso.com/guest-guidelines"

    # Classification Settings
    classification_list         = "Low,Medium,High,Confidential"
    default_classification      = "Medium"
    classification_descriptions = "[{\"Value\":\"Low\",\"Description\":\"Low business impact\"},{\"Value\":\"Medium\",\"Description\":\"Medium business impact\"},{\"Value\":\"High\",\"Description\":\"High business impact\"},{\"Value\":\"Confidential\",\"Description\":\"Confidential information\"}]"

    # Other Settings
    usage_guidelines_url                = "https://contoso.com/group-guidelines"
    enable_mip_labels                   = true
    new_unified_group_writeback_default = false
  }

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

