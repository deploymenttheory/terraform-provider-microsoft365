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
  display_name = "acc-test-app-maximal-${random_string.suffix.result}"
  description  = "Maximal acceptance test application for IP segment"

  prevent_duplicate_names = false
  hard_delete             = true
}


# ==============================================================================
# IP Application Segment - Maximal Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_maximal" {
  application_object_id = microsoft365_graph_beta_applications_application.test_minimal.id
  destination_host      = "*.example.com"
  destination_type      = "dnsSuffix"
  ports = [
    "80-80",
    "443-443",
    "8080-8080",
    "8443-8443"
  ]
  protocol = "tcp"
}
