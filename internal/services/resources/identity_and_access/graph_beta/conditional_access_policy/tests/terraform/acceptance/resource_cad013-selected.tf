# # # CAD013: Selected Apps - Compliant Device Requirement
# # # Requires compliant device for access to selected applications.
# # resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad013_selected_apps_compliant" {
# #   display_name = "CAD013-Selected: Grant access for All users when Browser and Modern Auth Clients and Compliant-v1.0"
# #   state        = "enabledForReportingButNotEnforced"

# #   conditions = {
# #     client_app_types = ["browser", "mobileAppsAndDesktopClients"]

# #     users = {
# #       include_users  = ["All"]
# #       exclude_users  = []
# #       include_groups = []
# #       exclude_groups = [
# #         microsoft365_graph_beta_groups_group.breakglass.id,
# #         microsoft365_graph_beta_groups_group.cad013_exclude.id
# #       ]
# #       include_roles = []
# #       exclude_roles = []
# #     }

# #     applications = {
# #       include_applications = [
# #         "a4f2693f-129c-4b96-982b-2c364b8314d7", # Specific application 1
# #         "499b84ac-1321-427f-aa17-267ca6975798", # Specific application 2
# #         "996def3d-b36c-4153-8607-a6fd3c01b89f", # Specific application 3
# #         "797f4846-ba00-4fd7-ba43-dac1f8f63013"  # Specific application 4
# #       ]
# #       exclude_applications                             = []
# #       include_user_actions                             = []
# #       include_authentication_context_class_references = []
# #     }

# #     platforms = {
# #       include_platforms = ["all"]
# #       exclude_platforms = []
# #     }

# #     sign_in_risk_levels = []
# #   }

# #   grant_controls = {
# #     operator          = "OR"
# #     built_in_controls = ["compliantDevice", "domainJoinedDevice"]
# #     custom_authentication_factors = []
# #   }

# #   timeouts = {
# #     create = "150s"
# #     read   = "150s"
# #     update = "150s"
# #     delete = "150s"
# #   }
# # }

