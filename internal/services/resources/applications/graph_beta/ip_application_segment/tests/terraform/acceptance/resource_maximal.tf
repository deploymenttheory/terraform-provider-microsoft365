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


# ==============================================================================
# IP Application Segment - Maximal Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_maximal" {
  application_object_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_app.id
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
