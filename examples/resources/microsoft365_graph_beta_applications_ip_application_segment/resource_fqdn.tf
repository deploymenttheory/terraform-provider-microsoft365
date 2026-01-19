# IP Application Segment with Fully Qualified Domain Name (FQDN)
# This example demonstrates how to configure an application segment using a specific
# hostname with multiple ports.

resource "microsoft365_graph_beta_applications_ip_application_segment" "fqdn" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "app.contoso.com"
  destination_type      = "fqdn"
  ports                 = ["443-443", "8443-8443"]
  protocol              = "tcp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
