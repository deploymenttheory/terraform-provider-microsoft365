resource "microsoft365_graph_beta_applications_on_premises_ip_application_segment" "ip_segment_range_start_end" {
  application_object_id = "12345678-1234-1234-1234-123456789012"
  destination_host      = "192.168.1.1..192.168.1.10"
  destination_type      = "ipRange"
  ports                 = ["80-80"]
  protocol              = ["tcp"]
}
