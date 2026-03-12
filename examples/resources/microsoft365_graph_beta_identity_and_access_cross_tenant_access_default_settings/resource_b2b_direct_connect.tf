# B2B direct connect — controls Teams Connect shared channels.
#
# Inbound: only supports "AllUsers" for users_and_groups. Applications may
# reference specific application IDs such as "Office365".
#
# Outbound: supports specific user and group GUIDs in addition to "AllUsers".
# When specific targets are used, "AllUsers" cannot be mixed with them.

resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "example-direct-connect-user"
  user_principal_name = "example-direct-connect-user@contoso.com"
  mail_nickname       = "example-direct-connect-user"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

resource "microsoft365_graph_beta_groups_group" "example_group_1" {
  display_name     = "example-direct-connect-group-1"
  mail_nickname    = "example-direct-connect-group-1"
  mail_enabled     = false
  security_enabled = true
}

resource "microsoft365_graph_beta_groups_group" "example_group_2" {
  display_name     = "example-direct-connect-group-2"
  mail_nickname    = "example-direct-connect-group-2"
  mail_enabled     = false
  security_enabled = true
}

# Allow time for the user and groups to fully propagate before the policy
# references their IDs.
resource "time_sleep" "wait_30_seconds" {
  depends_on = [
    microsoft365_graph_beta_users_user.example_user,
    microsoft365_graph_beta_groups_group.example_group_1,
    microsoft365_graph_beta_groups_group.example_group_2,
  ]
  create_duration = "30s"
}

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  depends_on = [time_sleep.wait_30_seconds]

  # Inbound direct connect: only AllUsers is valid for users_and_groups.
  b2b_direct_connect_inbound = {
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
          target      = "Office365"
          target_type = "application"
        }
      ]
    }
  }

  # Outbound direct connect: specific users and groups are blocked.
  b2b_direct_connect_outbound = {
    users_and_groups = {
      access_type = "blocked"
      targets = [
        {
          target      = microsoft365_graph_beta_users_user.example_user.id
          target_type = "user"
        },
        {
          target      = microsoft365_graph_beta_groups_group.example_group_1.id
          target_type = "group"
        },
        {
          target      = microsoft365_graph_beta_groups_group.example_group_2.id
          target_type = "group"
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
