# Lookup device by Azure AD Device ID
data "microsoft365_graph_beta_identity_and_access_device" "test" {
  device_id = "06771871-1375-494e-97f9-ab87ba64edeb"
}
