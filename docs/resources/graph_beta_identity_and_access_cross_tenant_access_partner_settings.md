---
page_title: "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages partner-specific cross-tenant access settings in Microsoft Entra ID using the /policies/crossTenantAccessPolicy/partners endpoint. This resource is used to configure B2B collaboration, B2B direct connect, inbound trust, and tenant restrictions for a specific partner organization.
---

# microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings (Resource)

Manages partner-specific cross-tenant access settings in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy/partners` endpoint. This resource is used to configure B2B collaboration, B2B direct connect, inbound trust, and tenant restrictions for a specific partner organization.

## Microsoft Documentation

- [crossTenantAccessPolicyConfigurationPartner resource type](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationpartner?view=graph-rest-beta)
- [Create crossTenantAccessPolicyConfigurationPartner](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicy-post-partners?view=graph-rest-beta&tabs=http)
- [Get crossTenantAccessPolicyConfigurationPartner](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-get?view=graph-rest-beta&tabs=http)
- [Update crossTenantAccessPolicyConfigurationPartner](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-update?view=graph-rest-beta&tabs=http)
- [Delete crossTenantAccessPolicyConfigurationPartner](https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-delete?view=graph-rest-beta&tabs=http)
- [Hard Delete (policydeletableitem-delete)](https://learn.microsoft.com/en-us/graph/api/policydeletableitem-delete?view=graph-rest-beta&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `Directory.Read.All`
- `Policy.Read.All`
- `Policy.ReadWrite.CrossTenantAccess`
- `User.Read.All`
- `User.ReadBasic.All`

Find out more about the permissions required for managing cross-tenant access policies at Microsoft Learn [here](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationpartner?view=graph-rest-beta).

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.49.0 | Experimental | Initial release |
| v0.50.0 | Experimental | Numerous big fixes and full test harness with additional examples added |

## Important Notes

- **Partner-Specific Configuration**: This resource configures cross-tenant access settings for a specific partner tenant identified by `tenant_id`. Each partner can have unique B2B collaboration, B2B direct connect, inbound trust, tenant restrictions, and automatic user consent settings.
- **Create Behavior**: The `create` operation uses a POST request (`/policies/crossTenantAccessPolicy/partners`) to create a new partner configuration.
- **Update Behavior**: The `update` operation uses a PATCH request to modify the partner configuration.
- **Delete Behavior**: The `delete` operation supports both soft delete (default) and hard delete:
  - **Soft Delete** (`hard_delete = false`, default): Moves the partner configuration to deleted items (can be restored within 30 days).
  - **Hard Delete** (`hard_delete = true`): Permanently removes the partner configuration from deleted items using `/directory/deletedItems/{tenantId}`.
- **Import with Hard Delete**: When importing, use the format `tenant_id:hard_delete=true` to enable hard delete on destroy. Example: `terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings.example "12345678-1234-1234-1234-123456789012:hard_delete=true"`

## Example Usage

### Minimal

```terraform
# Minimal example — only outbound B2B collaboration configured.
# All other settings remain at their partner defaults.
# On destroy, the partner configuration is soft deleted (can be restored within 30 days).

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = false

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
# B2B collaboration — controls inbound and outbound B2B guest access for a specific partner.
# Both directions are configured to allow all users and all applications.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

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

### B2B Direct Connect

```terraform
# B2B direct connect — blocks direct connect for both inbound and outbound directions.
# This prevents Teams Connect shared channels with the specified partner tenant.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

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
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }

  b2b_direct_connect_outbound = {
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

### Inbound Trust

```terraform
# Inbound trust — accepts MFA, compliant devices, and hybrid Azure AD joined devices
# from the partner tenant. This allows users from the partner to satisfy conditional
# access policies using their home tenant's device and authentication state.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

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

  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }
}
```

### Tenant Restrictions

```terraform
# Tenant restrictions — blocks all users and applications from the partner tenant
# when accessing resources in your tenant. This provides an additional layer of
# control beyond B2B collaboration settings.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

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
# Automatic user consent settings — disables automatic consent for both inbound
# and outbound collaboration. Users will need to explicitly consent to share data
# with the partner tenant.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

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

  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
```

### Complete

```terraform
# Complete example — all available blocks configured for a specific partner tenant.
# This demonstrates the full range of cross-tenant access controls including B2B
# collaboration, B2B direct connect, inbound trust, tenant restrictions, and
# automatic user consent settings.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id   = "12345678-1234-1234-1234-123456789012"
  hard_delete = true

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
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }

  b2b_direct_connect_outbound = {
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

  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }

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

  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `tenant_id` (String) The tenant ID of the partner Microsoft Entra organization. This is a GUID that uniquely identifies the partner tenant.

### Optional

- `automatic_user_consent_settings` (Attributes) Automatic user consent settings for the partner organization. (see [below for nested schema](#nestedatt--automatic_user_consent_settings))
- `b2b_collaboration_inbound` (Attributes) B2B collaboration inbound access settings for the partner organization. (see [below for nested schema](#nestedatt--b2b_collaboration_inbound))
- `b2b_collaboration_outbound` (Attributes) B2B collaboration outbound access settings for the partner organization. (see [below for nested schema](#nestedatt--b2b_collaboration_outbound))
- `b2b_direct_connect_inbound` (Attributes) B2B direct connect inbound access settings for the partner organization. (see [below for nested schema](#nestedatt--b2b_direct_connect_inbound))
- `b2b_direct_connect_outbound` (Attributes) B2B direct connect outbound access settings for the partner organization. (see [below for nested schema](#nestedatt--b2b_direct_connect_outbound))
- `hard_delete` (Boolean) When `true`, the partner configuration will be permanently deleted (hard delete) during destroy. When `false` (default), the partner configuration will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. **Note**: Hard delete permanently removes the configuration and cannot be undone.
- `inbound_trust` (Attributes) Inbound trust settings for accepting claims from the partner organization. (see [below for nested schema](#nestedatt--inbound_trust))
- `tenant_restrictions` (Attributes) Tenant restrictions settings for the partner organization. (see [below for nested schema](#nestedatt--tenant_restrictions))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for the partner configuration. This is the same as the `tenant_id`.
- `is_in_multi_tenant_organization` (Boolean) Identifies whether the partner organization is part of a multi-tenant organization with the local tenant.
- `is_service_provider` (Boolean) Identifies whether the partner-specific configuration is a cloud service provider for your organization. **Important**: This field can only be set when using delegated (user) authentication. When using application (client credentials) authentication, this field must be omitted entirely - the API will reject requests with 403's that explicitly set this field to either `true` or `false`. This is a read-only computed field when using service principal authentication.

<a id="nestedatt--automatic_user_consent_settings"></a>
### Nested Schema for `automatic_user_consent_settings`

Optional:

- `inbound_allowed` (Boolean) Specifies whether automatic user consent is allowed for inbound flows.
- `outbound_allowed` (Boolean) Specifies whether automatic user consent is allowed for outbound flows.


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

- `is_compliant_device_accepted` (Boolean) Specifies whether to accept compliant device claims from the partner organization.
- `is_hybrid_azure_ad_joined_device_accepted` (Boolean) Specifies whether to accept hybrid Azure AD joined device claims from the partner organization.
- `is_mfa_accepted` (Boolean) Specifies whether to accept MFA claims from the partner organization.


<a id="nestedatt--tenant_restrictions"></a>
### Nested Schema for `tenant_restrictions`

Optional:

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
# Import with soft delete (default)
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings.example "12345678-1234-1234-1234-123456789012"

# Import with hard delete enabled
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
```
