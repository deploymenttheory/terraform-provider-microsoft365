resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "Engineering"
  name                        = "MinimalAttribute"
  type                        = "String"
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

