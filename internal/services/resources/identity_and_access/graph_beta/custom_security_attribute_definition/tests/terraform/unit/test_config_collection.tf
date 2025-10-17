resource "microsoft365_graph_beta_identity_and_access_custom_security_attribute_definition" "test" {
  attribute_set               = "HumanResources"
  name                        = "Skills"
  description                 = "Skills and competencies of the employee"
  type                        = "String"
  status                      = "Available"
  is_collection               = true
  is_searchable               = true
  use_pre_defined_values_only = false

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "10m"
  }
}

