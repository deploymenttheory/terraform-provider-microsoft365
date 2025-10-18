data "microsoft365_graph_beta_device_management_managed_device" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'compliant'"
}

