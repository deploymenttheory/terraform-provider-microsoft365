resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_category" "description" {
  display_name = "Test Description Device Category - ${random_string.test_suffix.result}"
  description  = "This is a test device category with description"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
