---
page_title: "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages the default configuration for cross-tenant access policy in Microsoft Entra ID using the /policies/crossTenantAccessPolicy/default endpoint.
  This is a singleton resource — one default configuration exists per tenant and cannot be created or deleted via the Microsoft Graph API. The create operation uses an UPDATE (PATCH) request to configure the default settings. On destroy, the resource can optionally restore the default configuration to system defaults by setting restore_defaults_on_destroy = true (using the resetToSystemDefault API), or simply remove it from Terraform state while leaving the configuration in place (the default behaviour).
  See the Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta for details.
---

# microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings (Resource)

Manages the default configuration for cross-tenant access policy in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy/default` endpoint.

This is a **singleton resource** — one default configuration exists per tenant and cannot be created or deleted via the Microsoft Graph API. The `create` operation uses an UPDATE (PATCH) request to configure the default settings. On `destroy`, the resource can optionally restore the default configuration to system defaults by setting `restore_defaults_on_destroy = true` (using the resetToSystemDefault API), or simply remove it from Terraform state while leaving the configuration in place (the default behaviour).

See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta) for details.

## Microsoft Documentation

- [crossTenantAccessPolicyConfigurationDefault resource type](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta)
- [Get crossTenantAccessPolicyConfigurationDefault](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-get?view=graph-rest-beta&tabs=http)
- [Update crossTenantAccessPolicyConfigurationDefault](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-update?view=graph-rest-beta&tabs=http)
- [resetToSystemDefault](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-resettosystemdefault?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Directory.Read.All`
- `Policy.ReadWrite.CrossTenantAccess`
- `User.Read.All`
- `User.ReadBasic.All`

Find out more about the permissions required for managing cross-tenant access policies at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.49.0 | Experimental | Initial release |

## Important Notes

- **Singleton Resource**: This is a singleton resource — one default cross-tenant access policy configuration exists per tenant and cannot be created or deleted via the Microsoft Graph API. It is automatically provisioned by Microsoft Entra ID.
- **Create Behavior**: The `create` operation uses a PATCH request (`/policies/crossTenantAccessPolicy/default`) to apply the desired configuration.
- **Destroy Behavior**: On `destroy`, setting `restore_defaults_on_destroy = true` will call `resetToSystemDefault` to restore the policy to service defaults and verify `is_service_default = true`. Setting it to `false` (the default) removes the resource from Terraform state only, leaving the configuration in Entra ID unchanged.

## Example Usage

### Minimal

```terraform
# Minimal example — only outbound B2B collaboration configured.
# All other settings remain at their tenant defaults.
# On destroy, the configuration is reset to system defaults.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

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
```

### B2B Collaboration

```terraform
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
```

### B2B Direct Connect with Specific User and Group Targets

```terraform
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
```

### Inbound Trust

```terraform
# Inbound trust — configures which claims from external tenants are trusted
# when evaluating Conditional Access policies for inbound B2B users.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }
}
```

### Invitation Redemption Identity Provider Configuration

```terraform
# Invitation redemption — controls the order in which identity providers are
# tried when a B2B guest redeems an invitation, and the fallback provider if
# none of the primary options succeed.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  invitation_redemption_identity_provider_configuration = {
    primary_identity_provider_precedence_order = [
      "azureActiveDirectory",
      "externalFederation",
      "socialIdentityProviders"
    ]
    fallback_identity_provider = "emailOneTimePasscode"
  }
}
```

### Tenant Restrictions

```terraform
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
```

### Automatic User Consent Settings

```terraform
# Automatic user consent — controls whether users in this tenant can
# automatically consent to cross-tenant access without admin approval.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
```

### Complete

```terraform
# Complete example — all available blocks and settings configured together.
#
# This resource manages the singleton default cross-tenant access policy
# configuration (crossTenantAccessPolicyConfigurationDefault) that applies
# to all external tenants unless overridden by a partner-specific policy.
#
# Because no POST/DELETE endpoints exist, the provider issues a PATCH on create
# and update, and optionally calls resetToSystemDefault on destroy when
# restore_defaults_on_destroy = true.
#
# Dependencies:
#   b2b_direct_connect_outbound.users_and_groups.targets references specific
#   user and group IDs. A time_sleep ensures those objects are fully propagated
#   in Entra ID before the policy PATCH is issued.

# ==============================================================================
# User Dependency
# ==============================================================================

resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "example-cta-user"
  user_principal_name = "example-cta-user@contoso.com"
  mail_nickname       = "example-cta-user"
  account_enabled     = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "example_group_1" {
  display_name     = "example-cta-group-1"
  mail_nickname    = "example-cta-group-1"
  mail_enabled     = false
  security_enabled = true
  description      = "Group 1 blocked from outbound direct connect"
}

resource "microsoft365_graph_beta_groups_group" "example_group_2" {
  display_name     = "example-cta-group-2"
  mail_nickname    = "example-cta-group-2"
  mail_enabled     = false
  security_enabled = true
  description      = "Group 2 blocked from outbound direct connect"
}

# ==============================================================================
# Propagation Wait
# ==============================================================================

resource "time_sleep" "wait_30_seconds" {
  depends_on = [
    microsoft365_graph_beta_users_user.example_user,
    microsoft365_graph_beta_groups_group.example_group_1,
    microsoft365_graph_beta_groups_group.example_group_2,
  ]
  create_duration = "30s"
}

# ==============================================================================
# Cross-Tenant Access Default Settings
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  # When true, destroys the resource by calling resetToSystemDefault and
  # verifying is_service_default = true. When false (default), Terraform
  # removes the resource from state only and leaves the configuration in place.
  restore_defaults_on_destroy = true

  depends_on = [time_sleep.wait_30_seconds]

  # --------------------------------------------------------------------------
  # B2B Collaboration Inbound
  # Controls inbound B2B guest access from external tenants into this tenant.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # B2B Collaboration Outbound
  # Controls which users in this tenant can be invited as guests to external
  # tenants via B2B collaboration.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # B2B Direct Connect Inbound
  # Controls inbound Teams Connect shared channels from external tenants.
  # Note: users_and_groups only supports "AllUsers" — individual user or
  # group GUIDs are not valid for this direction.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # B2B Direct Connect Outbound
  # Controls which users in this tenant can use Teams Connect shared channels
  # with external tenants. Supports specific user and group GUIDs as targets.
  # Note: "AllUsers" cannot be mixed with specific user or group GUIDs in the
  # same targets set.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # Inbound Trust
  # Determines which claims from external tenants are honoured when evaluating
  # Conditional Access policies for inbound B2B users.
  # --------------------------------------------------------------------------
  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }

  # --------------------------------------------------------------------------
  # Invitation Redemption Identity Provider Configuration
  # Sets the ordered list of identity providers tried when a B2B guest redeems
  # an invitation. The API normalises this list to the canonical enum order
  # regardless of input, so the order here is informational only.
  # --------------------------------------------------------------------------
  invitation_redemption_identity_provider_configuration = {
    primary_identity_provider_precedence_order = [
      "azureActiveDirectory",
      "externalFederation",
      "socialIdentityProviders"
    ]
    fallback_identity_provider = "emailOneTimePasscode"
  }

  # --------------------------------------------------------------------------
  # Tenant Restrictions
  # Prevents users on managed devices from signing in to other tenants.
  # --------------------------------------------------------------------------
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

  # --------------------------------------------------------------------------
  # Automatic User Consent Settings
  # Controls whether users can automatically consent to cross-tenant
  # applications without requiring admin approval.
  # --------------------------------------------------------------------------
  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `automatic_user_consent_settings` (Attributes) Determines the default configuration for automatic user consent settings. The `inbound_allowed` and `outbound_allowed` properties are always `false` and can't be updated in the default configuration. Read-only. (see [below for nested schema](#nestedatt--automatic_user_consent_settings))
- `b2b_collaboration_inbound` (Attributes) B2B collaboration inbound access settings lets you collaborate with people outside of your organization by allowing them to sign in using their own identities. These users become guests in your Microsoft Entra tenant. You can invite external users directly or you can set up self-service sign-up so they can request access to your resources. (see [below for nested schema](#nestedatt--b2b_collaboration_inbound))
- `b2b_collaboration_outbound` (Attributes) B2B collaboration outbound access settings determine whether your users can be invited to external Microsoft Entra tenants for B2B collaboration and added to their directories as guests. These default settings apply to all external Microsoft Entra tenants except those with organization-specific settings. Below, specify whether your users and groups can be invited for B2B collaboration and the external applications they can access. (see [below for nested schema](#nestedatt--b2b_collaboration_outbound))
- `b2b_direct_connect_inbound` (Attributes) B2B direct connect inbound access settings determine whether users from external Microsoft Entra tenants can access your resources without being added to your tenant as guests. By selecting 'Allow access' below, you're permitting users and groups from other organizations to connect with you. To establish a connection, an admin from the other organization must also enable B2B direct connect. (see [below for nested schema](#nestedatt--b2b_direct_connect_inbound))
- `b2b_direct_connect_outbound` (Attributes) Outbound access settings determine how your users and groups can interact with apps and resources in external organizations. The default settings apply to all your cross-tenant scenarios unless you configure organizational settings to override them for a specific organization. Default settings can be modified but not deleted. (see [below for nested schema](#nestedatt--b2b_direct_connect_outbound))
- `inbound_trust` (Attributes) Configure whether your Conditional Access policies will accept claims from other Microsoft Entra tenants when external users access your resources. The default settings apply to all external Microsoft Entra tenants except those with organization-specific settings. You'll first need to configure Conditional Access for guest users on all cloud apps if you want to require multifactor authentication or require a device to be compliant or Microsoft Entra hybrid joined. (see [below for nested schema](#nestedatt--inbound_trust))
- `invitation_redemption_identity_provider_configuration` (Attributes) Defines the priority order based on which an identity provider is selected during invitation redemption for a guest user. (see [below for nested schema](#nestedatt--invitation_redemption_identity_provider_configuration))
- `restore_defaults_on_destroy` (Boolean) Controls behaviour when this resource is destroyed. When `true`, Terraform will issue a POST request to `/policies/crossTenantAccessPolicy/default/resetToSystemDefault` to reset the default configuration to system defaults, then verify that `is_service_default` is `true` before removing from state. When `false` (the default), Terraform removes the resource from state only — the existing default configuration is left unchanged in Microsoft Entra ID.
- `tenant_restrictions` (Attributes) Tenant restrictions lets you control whether your users can access external applications from your network or devices using external accounts, including accounts issued to them by external organizations and accounts they've created in unknown tenants. Below, select which external applications to allow or block. These default settings apply to all external Microsoft Entra tenants except those with organization-specific settings. (see [below for nested schema](#nestedatt--tenant_restrictions))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for the cross-tenant access default settings. This is a singleton resource; the value is always `crossTenantAccessDefaultSettings`.
- `is_service_default` (Boolean) If `true`, the default configuration is set to the system default configuration. If `false`, the default settings are customized. This is a read-only computed value.

<a id="nestedatt--automatic_user_consent_settings"></a>
### Nested Schema for `automatic_user_consent_settings`

Optional:

- `inbound_allowed` (Boolean) Specifies whether inbound automatic user consent is allowed. This is always `false` in the default configuration.
- `outbound_allowed` (Boolean) Specifies whether outbound automatic user consent is allowed. This is always `false` in the default configuration.


<a id="nestedatt--b2b_collaboration_inbound"></a>
### Nested Schema for `b2b_collaboration_inbound`

Optional:

- `applications` (Attributes) Specifies whether to allow or block access for applications. (see [below for nested schema](#nestedatt--b2b_collaboration_inbound--applications))
- `users_and_groups` (Attributes) Specifies whether to allow or block access for users and groups. (see [below for nested schema](#nestedatt--b2b_collaboration_inbound--users_and_groups))

<a id="nestedatt--b2b_collaboration_inbound--applications"></a>
### Nested Schema for `b2b_collaboration_inbound.applications`

Required:

- `targets` (Attributes Set) The set of application targets to allow or block. (see [below for nested schema](#nestedatt--b2b_collaboration_inbound--applications--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_collaboration_inbound--applications--targets"></a>
### Nested Schema for `b2b_collaboration_inbound.applications.targets`

Required:

- `target` (String) The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.
- `target_type` (String) The type of target. Must be `application`.



<a id="nestedatt--b2b_collaboration_inbound--users_and_groups"></a>
### Nested Schema for `b2b_collaboration_inbound.users_and_groups`

Required:

- `targets` (Attributes Set) The set of user and group targets to allow or block. (see [below for nested schema](#nestedatt--b2b_collaboration_inbound--users_and_groups--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_collaboration_inbound--users_and_groups--targets"></a>
### Nested Schema for `b2b_collaboration_inbound.users_and_groups.targets`

Required:

- `target` (String) The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.
- `target_type` (String) The type of target. Possible values: `user`, `group`.




<a id="nestedatt--b2b_collaboration_outbound"></a>
### Nested Schema for `b2b_collaboration_outbound`

Optional:

- `applications` (Attributes) Specifies whether to allow or block access for applications. (see [below for nested schema](#nestedatt--b2b_collaboration_outbound--applications))
- `users_and_groups` (Attributes) Specifies whether to allow or block access for users and groups. (see [below for nested schema](#nestedatt--b2b_collaboration_outbound--users_and_groups))

<a id="nestedatt--b2b_collaboration_outbound--applications"></a>
### Nested Schema for `b2b_collaboration_outbound.applications`

Required:

- `targets` (Attributes Set) The set of application targets to allow or block. (see [below for nested schema](#nestedatt--b2b_collaboration_outbound--applications--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_collaboration_outbound--applications--targets"></a>
### Nested Schema for `b2b_collaboration_outbound.applications.targets`

Required:

- `target` (String) The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.
- `target_type` (String) The type of target. Must be `application`.



<a id="nestedatt--b2b_collaboration_outbound--users_and_groups"></a>
### Nested Schema for `b2b_collaboration_outbound.users_and_groups`

Required:

- `targets` (Attributes Set) The set of user and group targets to allow or block. (see [below for nested schema](#nestedatt--b2b_collaboration_outbound--users_and_groups--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_collaboration_outbound--users_and_groups--targets"></a>
### Nested Schema for `b2b_collaboration_outbound.users_and_groups.targets`

Required:

- `target` (String) The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.
- `target_type` (String) The type of target. Possible values: `user`, `group`.




<a id="nestedatt--b2b_direct_connect_inbound"></a>
### Nested Schema for `b2b_direct_connect_inbound`

Optional:

- `applications` (Attributes) Specifies whether to allow or block access for applications. (see [below for nested schema](#nestedatt--b2b_direct_connect_inbound--applications))
- `users_and_groups` (Attributes) Specifies whether to allow or block access for users and groups. (see [below for nested schema](#nestedatt--b2b_direct_connect_inbound--users_and_groups))

<a id="nestedatt--b2b_direct_connect_inbound--applications"></a>
### Nested Schema for `b2b_direct_connect_inbound.applications`

Required:

- `targets` (Attributes Set) The set of application targets to allow or block. (see [below for nested schema](#nestedatt--b2b_direct_connect_inbound--applications--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_direct_connect_inbound--applications--targets"></a>
### Nested Schema for `b2b_direct_connect_inbound.applications.targets`

Required:

- `target` (String) The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.
- `target_type` (String) The type of target. Must be `application`.



<a id="nestedatt--b2b_direct_connect_inbound--users_and_groups"></a>
### Nested Schema for `b2b_direct_connect_inbound.users_and_groups`

Required:

- `targets` (Attributes Set) The set of user and group targets to allow or block. (see [below for nested schema](#nestedatt--b2b_direct_connect_inbound--users_and_groups--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_direct_connect_inbound--users_and_groups--targets"></a>
### Nested Schema for `b2b_direct_connect_inbound.users_and_groups.targets`

Required:

- `target` (String) The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.
- `target_type` (String) The type of target. Possible values: `user`, `group`.




<a id="nestedatt--b2b_direct_connect_outbound"></a>
### Nested Schema for `b2b_direct_connect_outbound`

Optional:

- `applications` (Attributes) Specifies whether to allow or block access for applications. (see [below for nested schema](#nestedatt--b2b_direct_connect_outbound--applications))
- `users_and_groups` (Attributes) Specifies whether to allow or block access for users and groups. (see [below for nested schema](#nestedatt--b2b_direct_connect_outbound--users_and_groups))

<a id="nestedatt--b2b_direct_connect_outbound--applications"></a>
### Nested Schema for `b2b_direct_connect_outbound.applications`

Required:

- `targets` (Attributes Set) The set of application targets to allow or block. (see [below for nested schema](#nestedatt--b2b_direct_connect_outbound--applications--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_direct_connect_outbound--applications--targets"></a>
### Nested Schema for `b2b_direct_connect_outbound.applications.targets`

Required:

- `target` (String) The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.
- `target_type` (String) The type of target. Must be `application`.



<a id="nestedatt--b2b_direct_connect_outbound--users_and_groups"></a>
### Nested Schema for `b2b_direct_connect_outbound.users_and_groups`

Required:

- `targets` (Attributes Set) The set of user and group targets to allow or block. (see [below for nested schema](#nestedatt--b2b_direct_connect_outbound--users_and_groups--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--b2b_direct_connect_outbound--users_and_groups--targets"></a>
### Nested Schema for `b2b_direct_connect_outbound.users_and_groups.targets`

Required:

- `target` (String) The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.
- `target_type` (String) The type of target. Possible values: `user`, `group`.




<a id="nestedatt--inbound_trust"></a>
### Nested Schema for `inbound_trust`

Optional:

- `is_compliant_device_accepted` (Boolean) Specifies whether to trust compliant device claims from external Microsoft Entra organizations.
- `is_hybrid_azure_ad_joined_device_accepted` (Boolean) Specifies whether to trust hybrid Azure AD joined device claims from external Microsoft Entra organizations.
- `is_mfa_accepted` (Boolean) Specifies whether to trust MFA claims from external Microsoft Entra organizations.


<a id="nestedatt--invitation_redemption_identity_provider_configuration"></a>
### Nested Schema for `invitation_redemption_identity_provider_configuration`

Required:

- `fallback_identity_provider` (String) Fallback identity providers are used when none of the primary identity providers are applicable. You must always have at least one fallback provider set to prevent users from being blocked while redeeming an invitation. Possible values: `defaultConfiguredIdp`, `emailOneTimePasscode`.
- `primary_identity_provider_precedence_order` (List of String) Users will redeem their invitations using the default order set by Microsoft. You can enable and specify the order of identity providers that your guest users can sign in with when they redeem their invitation. Possible values are: `externalFederation`, `azureActiveDirectory`, `socialIdentityProviders`. By not specifying an invitation redemption identity provider type it will set set to disabled.


<a id="nestedatt--tenant_restrictions"></a>
### Nested Schema for `tenant_restrictions`

Required:

- `applications` (Attributes) Specifies whether to allow or block access for applications. (see [below for nested schema](#nestedatt--tenant_restrictions--applications))
- `users_and_groups` (Attributes) Specifies whether to allow or block access for users and groups. (see [below for nested schema](#nestedatt--tenant_restrictions--users_and_groups))

<a id="nestedatt--tenant_restrictions--applications"></a>
### Nested Schema for `tenant_restrictions.applications`

Required:

- `targets` (Attributes Set) The set of application targets to allow or block. (see [below for nested schema](#nestedatt--tenant_restrictions--applications--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--tenant_restrictions--applications--targets"></a>
### Nested Schema for `tenant_restrictions.applications.targets`

Required:

- `target` (String) The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.
- `target_type` (String) The type of target. Must be `application`.



<a id="nestedatt--tenant_restrictions--users_and_groups"></a>
### Nested Schema for `tenant_restrictions.users_and_groups`

Required:

- `targets` (Attributes Set) The set of user and group targets to allow or block. (see [below for nested schema](#nestedatt--tenant_restrictions--users_and_groups--targets))

Optional:

- `access_type` (String) The access type. Possible values: `allowed`, `blocked`.

<a id="nestedatt--tenant_restrictions--users_and_groups--targets"></a>
### Nested Schema for `tenant_restrictions.users_and_groups.targets`

Required:

- `target` (String) The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.
- `target_type` (String) The type of target. Possible values: `user`, `group`.




<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

```shell
# Import without resetting defaults on destroy (state-only removal, the safe default)
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings.example crossTenantAccessDefaultSettings

# Import and reset to system defaults on destroy (calls resetToSystemDefault)
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings.example "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true"

# Import explicitly keeping configuration on destroy
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings.example "crossTenantAccessDefaultSettings:restore_defaults_on_destroy=false"
```
