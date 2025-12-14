# CAAU001: AI/Agentic Identities Conditional Access Policy
# Block agent identity access for All agentic identities to all agent resources when Risk is medium or higher
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "caau001_agent_risk_block" {
  display_name = "CAAU001-ALL: Block agent identity access for All agentic identities to all agent resources when Risk is medium or higher v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["None"]
      exclude_users  = []
      include_groups = []
      exclude_groups = []
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["AllAgentIdResources"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    client_applications = {
      include_service_principals          = []
      exclude_service_principals          = []
      include_agent_id_service_principals = ["All"]
      exclude_agent_id_service_principals = []
    }

    sign_in_risk_levels           = []
    service_principal_risk_levels = []
    agent_id_risk_levels          = ["high", "medium"]
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}
