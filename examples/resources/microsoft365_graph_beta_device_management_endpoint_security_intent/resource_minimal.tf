# BitLocker disk encryption intent with the template defaults
resource "microsoft365_graph_beta_device_management_endpoint_security_intent" "bitlocker_minimal" {
  template_id  = "d1174162-1dd2-4976-affc-6667049ab0ae"
  display_name = "BitLocker - template defaults"
}
