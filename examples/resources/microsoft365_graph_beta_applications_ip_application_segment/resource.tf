resource "microsoft365_graph_beta_applications_ip_application_segment" "example" {
  application_id   = "00000000-0000-0000-0000-000000000000"
  destination_host = "internal.contoso.com"
  destination_type = "fqdn"
  ports            = ["80-80", "443-443"]
  protocol         = "tcp"

  timeouts {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}

