resource "microsoft365_graph_beta_windows_365_cloud_pc_device_image" "minimal" {
  display_name              = "Test Minimal Cloud PC Device Image - Unique"
  version                   = "1.0.0"
  source_image_resource_id  = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/test-rg/providers/Microsoft.Compute/images/test-minimal-image"
  
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}