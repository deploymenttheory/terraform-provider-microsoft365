resource "microsoft365_graph_beta_identity_and_access_attribute_set" "test" {
  id                     = "Engineering"
  description            = "Attributes for engineering team"
  max_attributes_per_set = 25

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
