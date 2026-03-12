# Tenant restrictions — prevents users on managed devices from accessing
# external tenants. Blocks all users from accessing all external applications.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  tenant_restrictions = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
        }
      ]
    }
    applications = {
      access_type = "blocked"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }
}
