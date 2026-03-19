# Example: Get devices using an OData query

# Filter for Windows devices that are enabled
data "microsoft365_graph_beta_identity_and_access_device" "windows_enabled" {
  odata_query = "operatingSystem eq 'Windows' and accountEnabled eq true"
}

output "windows_device_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.windows_enabled.items)
}

# Filter for compliant devices
data "microsoft365_graph_beta_identity_and_access_device" "compliant" {
  odata_query = "isCompliant eq true"
}

output "compliant_device_count" {
  value = length(data.microsoft365_graph_beta_identity_and_access_device.compliant.items)
}
