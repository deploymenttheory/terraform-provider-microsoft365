resource "microsoft365_graph_beta_identity_and_access_attribute_set" "test" {
  id                     = "Engineering"
  description            = "Updated attributes for engineering team"
  max_attributes_per_set = 50

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
