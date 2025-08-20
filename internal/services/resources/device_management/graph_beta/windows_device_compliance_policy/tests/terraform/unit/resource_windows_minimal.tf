resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "minimal" {
  display_name                               = "Windows Compliance Policy Minimal"
  description                                = "Minimal Windows compliance policy for testing"
  password_required                          = true
  password_block_simple                      = true
  password_minimum_length                    = 8
  password_required_type                     = "alphanumeric"
  password_previous_password_block_count     = 5
  password_expiration_days                   = 90
  password_minutes_of_inactivity_before_lock = 15
  password_required_to_unlock_from_idle      = true
  require_healthy_device_report              = true
  os_minimum_version                         = "10.0.19041.0"
  os_maximum_version                         = "10.0.22631.3155"
}