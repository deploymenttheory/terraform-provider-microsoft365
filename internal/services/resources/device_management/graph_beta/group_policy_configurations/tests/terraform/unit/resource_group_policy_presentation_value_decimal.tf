resource "microsoft365_graph_beta_device_management_group_policy_configuration" "group_policy_presentation_value_decimal" {
  display_name       = "unit-test-group-policy-presentation-value-decimal"
  description        = "Comprehensive test description for group policy configuration with all fields"
  role_scope_tag_ids = ["0", "1", "2"]

  definition_values {
    enabled       = true
    definition_id = "98d69f26-2201-4aed-8927-d20c29b24ed5"
    
    presentation_values {
      odata_type     = "#microsoft.graph.groupPolicyPresentationValueDecimal"
      presentation_id = "5c8d10ad-8a28-4fe2-bd16-170d88dceb82"
      decimal_value  = 1024
    }
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
