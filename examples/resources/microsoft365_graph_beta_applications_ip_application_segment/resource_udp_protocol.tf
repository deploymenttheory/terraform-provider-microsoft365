# IP Application Segment with UDP Protocol
# This example demonstrates how to configure an application segment using UDP protocol
# instead of TCP, useful for applications like VoIP or video conferencing.

resource "microsoft365_graph_beta_applications_ip_application_segment" "udp_app" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "voip.contoso.com"
  destination_type      = "fqdn"
  ports                 = ["5060-5061", "10000-20000"]
  protocol              = "udp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
