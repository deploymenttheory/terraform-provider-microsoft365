# Lookup device by device ID (Azure Device Registration Service ID)
# This is the deviceId property from the device object
data "microsoft365_graph_beta_identity_and_access_device" "test" {
  device_id = "06771871-1375-494e-97f9-ab87ba64edeb"
}
