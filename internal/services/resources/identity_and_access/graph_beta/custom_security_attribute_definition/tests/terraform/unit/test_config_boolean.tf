resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "Security"
  name                        = "HasClearance"
  description                 = "Indicates if user has security clearance"
  type                        = "Boolean"
  status                      = "Available"
  is_collection               = false
  is_searchable               = true
  use_pre_defined_values_only = false

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

