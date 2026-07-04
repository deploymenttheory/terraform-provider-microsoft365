# IP Application Segment with wildcard FQDN
# The application-scoped Graph endpoint accepts wildcard hosts when
# destination_type is fqdn. dnsSuffix is reserved for Quick Access configuration.

resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "wildcard_fqdn" {
  application_object_id = "00000000-0000-0000-0000-000000000000"
  destination_host      = "*.internal.contoso.com"
  destination_type      = "fqdn"
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
