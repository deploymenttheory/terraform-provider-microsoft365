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
  display_name = "acc-test-app-minimal-${random_string.suffix.result}"
  description  = "Minimal acceptance test application for IP segment"

  prevent_duplicate_names = false
  hard_delete             = true
}

# ==============================================================================
# IP Application Segment - Minimal Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_minimal" {
  application_object_id = microsoft365_graph_beta_applications_application.test_minimal.id
  destination_host      = "192.168.1.100"
  destination_type      = "ipAddress"
  ports                 = ["80-80"]
  protocol              = "tcp"
}
