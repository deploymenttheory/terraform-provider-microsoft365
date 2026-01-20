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
# IP Application Segment - FQDN Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_fqdn" {
  application_object_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_app.id
  destination_host      = "app.example.com"
  destination_type      = "fqdn"
  ports                 = ["443-443", "8443-8443"]
  protocol              = "tcp"
}
