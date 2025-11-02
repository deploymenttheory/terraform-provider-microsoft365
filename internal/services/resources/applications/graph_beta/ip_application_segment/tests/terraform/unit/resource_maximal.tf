resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_maximal" {
  application_id   = "12345678-1234-1234-1234-123456789012"
  destination_host = "*.example.com"
  destination_type = "dnsSuffix"
  ports = [
    "80-80",
    "443-443",
    "8080-8080",
    "8443-8443"
  ]
  protocol = "tcp"
}

