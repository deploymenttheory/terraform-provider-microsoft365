# IP Application Segment with IP Range (CIDR notation)
# This example demonstrates how to configure an application segment for an entire
# IP subnet using CIDR notation.

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_range" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "192.168.1.0/24"
  destination_type      = "ipRangeCidr"
  ports                 = ["443-443"]
  protocol              = "tcp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
