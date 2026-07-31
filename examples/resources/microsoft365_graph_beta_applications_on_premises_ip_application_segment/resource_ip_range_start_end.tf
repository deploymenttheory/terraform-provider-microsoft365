# IP Application Segment with IP Range (start..end notation)
# This example demonstrates how to configure an application segment for a
# contiguous range of IP addresses using the start..end notation required by
# the Graph API for destination_type = "ipRange".

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "ip_range_start_end" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "192.168.1.1..192.168.1.10"
  destination_type      = "ipRange"
  ports                 = ["80-80"]
  protocol              = ["tcp"]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
