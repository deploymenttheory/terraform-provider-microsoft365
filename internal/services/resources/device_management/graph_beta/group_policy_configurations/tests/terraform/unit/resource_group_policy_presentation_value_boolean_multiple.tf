resource "microsoft365_graph_beta_device_management_group_policy_configuration" "group_policy_presentation_value_boolean_multiple" {
  display_name       = "unit-test-group-policy-presentation-value-boolean-multiple"
  description        = "unit test for group policy with multiple boolean presentation values"
  role_scope_tag_ids = ["0"]

  definition_values {
    enabled       = true
    definition_id = "5cdf0ce8-f338-4bb3-b2c2-d00e5a255514"
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueBoolean"
      presentation_id = "95e90bbb-f30f-4d15-a3a6-00bcebfa02e9"
      value          = "true"
    }
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueBoolean"
      presentation_id = "ee6784ea-dca6-4a86-8058-c556add7f60f"
      value          = "true"
    }
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueBoolean"
      presentation_id = "1626672c-d7cc-400f-bd95-fdedcb90bf10"
      value          = "true"
    }
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}