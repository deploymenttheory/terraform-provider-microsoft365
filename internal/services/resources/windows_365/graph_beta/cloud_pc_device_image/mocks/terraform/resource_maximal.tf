resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "maximal" {
  display_name              = "Test Maximal Cloud PC Device Image - Unique"
  version                   = "2.1.5"
  source_image_resource_id  = "/subscriptions/87654321-4321-4321-4321-cba987654321/resourceGroups/test-maximal-rg/providers/Microsoft.Compute/images/test-maximal-image"
  
  timeouts = {
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "60s"
  }
}