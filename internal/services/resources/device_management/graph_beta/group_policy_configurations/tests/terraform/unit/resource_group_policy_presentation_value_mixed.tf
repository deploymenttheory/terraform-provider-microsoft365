resource "microsoft365_graph_beta_device_management_group_policy_configuration" "group_policy_presentation_value_mixed" {
  display_name       = "unit-test-group-policy-presentation-value-mixed"
  description        = "unit test for group policy with multiple boolean presentation values, text presentation value, and decimal presentation value"
  role_scope_tag_ids = ["0"]

  # Definition 1: 5cdf0ce8-f338-4bb3-b2c2-d00e5a255514 with 3 boolean presentation values
  definition_values {
    enabled       = true
    definition_id = "5cdf0ce8-f338-4bb3-b2c2-d00e5a255514"
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueBoolean"
      presentation_id = "95e90bbb-f30f-4d15-a3a6-00bcebfa02e9"
      boolean_value   = true
    }
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueBoolean"
      presentation_id = "ee6784ea-dca6-4a86-8058-c556add7f60f"
      boolean_value   = true
    }
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueBoolean"
      presentation_id = "1626672c-d7cc-400f-bd95-fdedcb90bf10"
      boolean_value   = true
    }
  }

  # Definition 2: afa83413-18ab-44d7-b829-8d707a738816 with no presentation values
  definition_values {
    enabled       = true
    definition_id = "afa83413-18ab-44d7-b829-8d707a738816"
  }

  # Definition 3: 83991532-1ff7-45e2-b583-3159c069649f with no presentation values
  definition_values {
    enabled       = true
    definition_id = "83991532-1ff7-45e2-b583-3159c069649f"
  }

  # Definition 4: d5b711bc-57ff-45f2-9a81-7da7d607547e with no presentation values
  definition_values {
    enabled       = true
    definition_id = "d5b711bc-57ff-45f2-9a81-7da7d607547e"
  }

  # Definition 5: 109b4590-814f-4e94-853c-e96c420e7514 with 1 text presentation value
  definition_values {
    enabled       = true
    definition_id = "109b4590-814f-4e94-853c-e96c420e7514"
    
    presentation_values {
      odata_type      = "#microsoft.graph.groupPolicyPresentationValueText"
      presentation_id = "609703ae-1fcd-4994-a1e3-5b5ec543c388"
      text_value      = "0"
    }
  }

  # Definition 6: d1eab7d5-c54f-4bfc-a084-e873fb33c8e6 with no presentation values
  definition_values {
    enabled       = true
    definition_id = "d1eab7d5-c54f-4bfc-a084-e873fb33c8e6"
  }

  # Definition 7: 238bd4f4-1b30-4b10-b141-dc223d7d969b with no presentation values
  definition_values {
    enabled       = true
    definition_id = "238bd4f4-1b30-4b10-b141-dc223d7d969b"
  }

  # Definition 8: 5ac47f6a-44b5-492c-a293-c39894a921be with no presentation values
  definition_values {
    enabled       = true
    definition_id = "5ac47f6a-44b5-492c-a293-c39894a921be"
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}