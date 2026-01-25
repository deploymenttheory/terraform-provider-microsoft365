# Test 10: Removed Policy Configuration
# Purpose: Use removed block to control destruction order for cleanup

# Dependencies
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Removed block for policy to control destruction order
removed {
  from = microsoft365_graph_beta_device_management_app_control_for_business_policy.maximal

  lifecycle {
    destroy = true
  }
}
