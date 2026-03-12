# B2B collaboration — controls inbound and outbound B2B guest access.
# Both directions are configured to allow all users and all applications.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  b2b_collaboration_inbound = {
    users_and_groups = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
        }
      ]
    }
    applications = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }

  b2b_collaboration_outbound = {
    users_and_groups = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
        }
      ]
    }
    applications = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }
}
