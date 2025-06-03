# Basic Terms and Conditions Assignment to all licensed users
resource "microsoft365_graph_beta_device_management_terms_and_conditions_assignment" "all_users" {
  terms_and_conditions_id = microsoft365_graph_beta_device_management_terms_and_conditions.company_terms.id

  target = {
    target_type = "allLicensedUsers"
  }
}

# Assignment to a specific Azure AD group
resource "microsoft365_graph_beta_device_management_terms_and_conditions_assignment" "specific_group" {
  terms_and_conditions_id = microsoft365_graph_beta_device_management_terms_and_conditions.company_terms.id

  target = {
    target_type = "groupAssignment"
    group_id    = "12345678-1234-1234-1234-123456789012" # IT Department group
  }
}

# Assignment to SCCM collection
resource "microsoft365_graph_beta_device_management_terms_and_conditions_assignment" "sccm_collection" {
  terms_and_conditions_id = microsoft365_graph_beta_device_management_terms_and_conditions.company_terms.id

  target = {
    target_type   = "configurationManagerCollection"
    collection_id = "MEM00012345" # Custom SCCM collection
  }
}