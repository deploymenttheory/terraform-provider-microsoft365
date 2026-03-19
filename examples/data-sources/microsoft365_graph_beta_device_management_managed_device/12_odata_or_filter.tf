# Example: Get managed devices using OData filter with OR logic
data "microsoft365_graph_beta_device_management_managed_device" "ios_or_android" {
  odata_query = "operatingSystem eq 'iOS' or operatingSystem eq 'Android'"
}

# Output: iOS and Android devices
output "mobile_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items :
    {
      id               = device.id
      device_name      = device.device_name
      operating_system = device.operating_system
      os_version       = device.os_version
      model            = device.model
      manufacturer     = device.manufacturer
    }
  ]
  description = "List of iOS and Android mobile devices"
}

# Output: Device count by OS
output "mobile_device_summary" {
  value = {
    total = length(data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items)
    ios = length([
      for device in data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items :
      device if device.operating_system == "iOS"
    ])
    android = length([
      for device in data.microsoft365_graph_beta_device_management_managed_device.ios_or_android.items :
      device if device.operating_system == "Android"
    ])
  }
  description = "Summary of mobile devices by operating system"
}
