resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_minimal" {
  application_id   = "12345678-1234-1234-1234-123456789012"
  destination_host = "192.168.1.100"
  destination_type = "ipAddress"
  ports            = ["80-80"]
  protocol         = "tcp"
}

