data "microsoft365_graph_beta_device_management_managed_device" "odata_advanced" {
  filter_type   = "odata"
  odata_filter  = "operatingSystem eq 'Windows'"
  odata_orderby = "deviceName"
  odata_select  = "id,deviceName,operatingSystem,complianceState"
}

