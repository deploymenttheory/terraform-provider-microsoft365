resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_allowed_value" "test" {
  custom_security_attribute_definition_id = "Engineering_Department"
  id                                      = "Human Resources"
  is_active                               = true
}

