resource "microsoft365_graph_beta_identity_and_access_attribute_set" "example" {
  id                     = "Engineering"
  description            = "Attributes for engineering team"
  max_attributes_per_set = 25 // max is 500

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}