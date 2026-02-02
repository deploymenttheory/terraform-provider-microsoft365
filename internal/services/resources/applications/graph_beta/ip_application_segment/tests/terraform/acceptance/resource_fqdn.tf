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
  display_name = "acc-test-app-fqdn-${random_string.suffix.result}"
  description  = "FQDN acceptance test application for IP segment"

  prevent_duplicate_names = false
  hard_delete             = true
}


# ==============================================================================
# IP Application Segment - FQDN Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_fqdn" {
  application_object_id = microsoft365_graph_beta_applications_application.test_minimal.id
  destination_host      = "app.example.com"
  destination_type      = "fqdn"
  ports                 = ["443-443", "8443-8443"]
  protocol              = "tcp"
}
