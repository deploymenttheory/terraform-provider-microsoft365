resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_fqdn" {
  application_id   = "12345678-1234-1234-1234-123456789012"
  destination_host = "app.example.com"
  destination_type = "fqdn"
  ports            = ["443-443", "8443-8443"]
  protocol         = "tcp"
}

