# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Application Dependency
# ==============================================================================

resource "microsoft365_graph_beta_applications_application" "test_minimal" {
  display_name = "acc-test-app-iprange-${random_string.suffix.result}"
  description  = "IP Range acceptance test application for IP segment"

  prevent_duplicate_names = false
  hard_delete             = true
}


# ==============================================================================
# IP Application Segment - IP Range Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_range" {
  application_object_id = microsoft365_graph_beta_applications_application.test_minimal.id
  destination_host      = "192.168.1.0/24"
  destination_type      = "ipRangeCidr"
  ports                 = ["443-443"]
  protocol              = "tcp"
}
