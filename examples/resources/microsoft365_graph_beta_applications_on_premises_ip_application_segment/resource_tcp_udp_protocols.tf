# IP Application Segment with TCP and UDP Protocols
# A set models the protocols without depending on input order.

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "tcp_udp_app" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "192.168.1.100"
  destination_type      = "ipAddress"
  ports                 = ["443-443"]
  protocol              = ["tcp", "udp"]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
