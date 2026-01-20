# IP Application Segment with DNS Suffix (Wildcard Domain)
# This example demonstrates how to configure an application segment using a wildcard
# domain to match all subdomains.

resource "microsoft365_graph_beta_applications_ip_application_segment" "dns_suffix" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "*.internal.contoso.com"
  destination_type      = "dnsSuffix"
  ports = [
    "80-80",
    "443-443",
    "8080-8080",
    "8443-8443"
  ]
  protocol = "tcp"

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
