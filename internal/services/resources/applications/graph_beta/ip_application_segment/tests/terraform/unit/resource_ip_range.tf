resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_range" {
  application_object_id = "12345678-1234-1234-1234-123456789012"
  destination_host      = "192.168.1.0/24"
  destination_type      = "ipRangeCidr"
  ports                 = ["443-443"]
  protocol              = "tcp"
}

