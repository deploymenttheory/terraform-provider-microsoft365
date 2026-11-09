resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "concurrent" {
  count = 3

  application_object_id = "12345678-1234-1234-1234-123456789012"
  destination_host      = "app${count.index}.example.com"
  destination_type      = "fqdn"
  ports                 = ["443-443"]
  protocol              = ["tcp"]

  timeouts = {
    create = "8s"
  }
}
