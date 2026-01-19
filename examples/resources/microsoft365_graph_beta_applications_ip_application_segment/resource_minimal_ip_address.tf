# Basic IP Application Segment with single IP address
# This example demonstrates the minimal configuration for an IP application segment
# targeting a single IP address.

resource "microsoft365_graph_beta_applications_ip_application_segment" "minimal_ip" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "192.168.1.100"
  destination_type      = "ipAddress"
  ports                 = ["80-80"]
  protocol              = "tcp"
}
