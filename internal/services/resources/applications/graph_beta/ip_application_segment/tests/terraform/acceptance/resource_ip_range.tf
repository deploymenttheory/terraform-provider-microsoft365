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
# IP Application Segment - IP Range Configuration
# ==============================================================================

resource "microsoft365_graph_beta_applications_ip_application_segment" "ip_segment_range" {
  application_object_id = microsoft365_graph_beta_agents_agent_identity_blueprint.test_app.id
  destination_host      = "192.168.1.0/24"
  destination_type      = "ipRangeCidr"
  ports                 = ["443-443"]
  protocol              = "tcp"
}
