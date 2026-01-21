---
page_title: "microsoft365_graph_beta_identity_and_access_conditional_access_template Data Source - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Retrieves conditional access policy templates from Microsoft Entra ID using the /identity/conditionalAccess/templates endpoint. This data source is used to discover pre-configured security policy templates with conditions, grant controls, and session controls for common scenarios.
---

# microsoft365_graph_beta_identity_and_access_conditional_access_template (Data Source)

Retrieves conditional access policy templates from Microsoft Entra ID using the `/identity/conditionalAccess/templates` endpoint. This data source is used to discover pre-configured security policy templates with conditions, grant controls, and session controls for common scenarios.

## Microsoft Documentation

- [Conditional access template resource type](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesstemplate?view=graph-rest-beta)
- [List conditional access templates](https://learn.microsoft.com/en-us/graph/api/conditionalaccessroot-list-templates?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `Policy.Read.All`

**Optional:**
- `None` `[N/A]`

## Filtering

This data source supports filtering using the following attributes:

- `template_id` (Optional): The unique identifier (GUID) of the conditional access template. Exactly one of `template_id` or `name` must be specified.
- `name` (Optional): The name of the conditional access template. Exactly one of `template_id` or `name` must be specified.

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.41.0-alpha | Experimental | Initial release |

## Example Usage

### Example 1: Get conditional access template by template ID

```terraform
# Example: Get conditional access template by template ID
# This example demonstrates how to retrieve a conditional access template using its ID (GUID)
# Useful when you know the template ID and need to retrieve its configuration details

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "by_id" {
  template_id = "c7503427-338e-4c5e-902d-abe252abfb43" # Require multifactor authentication for admins

  timeouts = {
    read = "1m"
  }
}

# Output template details
output "template_by_id" {
  value = {
    template_id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.template_id
    name        = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.name
    description = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.description
    scenarios   = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_id.scenarios
  }
  description = "Conditional access template retrieved by template ID"
}
```

### Example 2: Get conditional access template by name

```terraform
# Example: Get conditional access template by name
# This example demonstrates how to retrieve a conditional access template using its name
# Useful when you know the template name and need to find its ID or configuration details

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "by_name" {
  name = "Require multifactor authentication for admins"

  timeouts = {
    read = "1m"
  }
}

# Output template details
output "template_by_name" {
  value = {
    template_id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.template_id
    name        = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.name
    description = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.description
    scenarios   = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.by_name.scenarios
  }
  description = "Conditional access template retrieved by name"
}
```

### Example 3: Get detailed template configuration

```terraform
# Example: Get detailed conditional access template configuration
# This example demonstrates how to retrieve a template and access its detailed configuration
# including conditions, grant controls, and session controls

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "detailed" {
  name = "Block legacy authentication"

  timeouts = {
    read = "1m"
  }
}

# Output detailed template configuration
output "template_details" {
  value = {
    template_id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.template_id
    name        = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.name
    description = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.description
    scenarios   = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.scenarios
    details     = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details
  }
  description = "Detailed conditional access template configuration"
  sensitive   = false
}

# Example: Use template data to understand policy structure
# You can reference specific parts of the template details
output "template_grant_controls" {
  value = {
    operator          = try(data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details.grant_controls.operator, null)
    built_in_controls = try(data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details.grant_controls.built_in_controls, null)
  }
  description = "Grant controls from the template"
}

output "template_conditions" {
  value = {
    client_app_types = try(data.microsoft365_graph_beta_identity_and_access_conditional_access_template.detailed.details.conditions.client_app_types, null)
  }
  description = "Conditions from the template"
}
```

### Example 4: Create policy from template - Require MFA for admins

```terraform
# Test 04: Create a conditional access policy from template - 
# Require multifactor authentication for admins

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_admins" {
  name = "Require multifactor authentication for admins"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_admins" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-admins-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.users.exclude_roles
      include_users  = []
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.conditions.users.exclude_groups
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admins.details.grant_controls.terms_of_use
  }
}
```

### Example 5: Create policy from template - Block legacy authentication

```terraform
# Test 05: Create a conditional access policy from template - 
# Block legacy authentication

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_legacy_auth" {
  name = "Block legacy authentication"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_block_legacy_auth" {
  display_name = "acc-test-ca-policy-template-block-legacy-authentication-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.terms_of_use
  }
}
```

### Example 6: Create policy from template - Securing security info registration

```terraform
# Test 06: Create a conditional access policy from template - 
# Securing security info registration

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "securing_security_info" {
  name = "Securing security info registration"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_securing_security_info" {
  display_name = "acc-test-ca-policy-template-securing-security-info-registration-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.users.exclude_roles
    }

    locations = {
      include_locations = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.locations.include_locations
      exclude_locations = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.conditions.locations.exclude_locations
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.securing_security_info.details.grant_controls.terms_of_use
  }
}
```

### Example 7: Create policy from template - Require MFA for all users

```terraform
# Test 07: Create a conditional access policy from template - 
# Require multifactor authentication for all users

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_all_users" {
  name = "Require multifactor authentication for all users"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_all_users" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-all-users-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_all_users.details.grant_controls.terms_of_use
  }
}
```

### Example 8: Create policy from template - Require MFA for guest access

```terraform
# Test 08: Create a conditional access policy from template - 
# Require multifactor authentication for guest access

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_guest_access" {
  name = "Require multifactor authentication for guest access"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_guest_access" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-guest-access-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_guest_access.details.grant_controls.terms_of_use
  }
}
```

### Example 9: Create policy from template - Require MFA for Azure management

```terraform
# Test 09: Create a conditional access policy from template - 
# Require multifactor authentication for Azure management

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_azure_management" {
  name = "Require multifactor authentication for Azure management"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_azure_management" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-azure-management-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_azure_management.details.grant_controls.terms_of_use
  }
}
```

### Example 10: Create policy from template - Require MFA for risky sign-ins

```terraform
# Test 10: Create a conditional access policy from template - 
# Require multifactor authentication for risky sign-ins

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_risky_signins" {
  name = "Require multifactor authentication for risky sign-ins"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_risky_signins" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-risky-signins-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.terms_of_use
  }
}
```

### Example 11: Create policy from template - Require password change for high-risk users

```terraform
# Test 11: Create a conditional access policy from template - 
# Require password change for high-risk users

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "password_change_high_risk" {
  name = "Require password change for high-risk users"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_password_change_high_risk" {
  display_name = "acc-test-ca-policy-template-require-password-change-for-high-risk-users-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.password_change_high_risk.details.grant_controls.terms_of_use
  }
}
```

### Example 12: Create policy from template - Require compliant or hybrid device for admins

```terraform
# Test 12: Create a conditional access policy from template - 
# Require compliant or hybrid Azure AD joined device for admins

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "compliant_device_admins" {
  name = "Require compliant or hybrid Azure AD joined device for admins"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_compliant_device_admins" {
  display_name = "acc-test-ca-policy-template-require-compliant-or-hybrid-device-for-admins-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = [] # Overriding template's "None" value which conflicts with role assignments
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_admins.details.grant_controls.terms_of_use
  }
}
```

### Example 14: Create policy from template - Require compliant device or MFA for all users

```terraform
# Test 14: Create a conditional access policy from template - 
# Require compliant or hybrid Azure AD joined device or multifactor authentication for all users

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "compliant_device_or_mfa_all_users" {
  name = "Require compliant or hybrid Azure AD joined device or multifactor authentication for all users"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_compliant_device_or_mfa_all_users" {
  display_name = "acc-test-ca-policy-template-require-compliant-device-or-mfa-for-all-users-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.compliant_device_or_mfa_all_users.details.grant_controls.terms_of_use
  }
}
```

### Example 15: Create policy from template - Use application enforced restrictions for O365

```terraform
# Test 15: Create a conditional access policy from template - 
# Use application enforced restrictions for O365 apps

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "app_enforced_restrictions_o365" {
  name = "Use application enforced restrictions for O365 apps"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_app_enforced_restrictions_o365" {
  display_name = "acc-test-ca-policy-template-use-application-enforced-restrictions-for-o365-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.exclude_roles
    }
  }

  # Template has null grant_controls, but we need to provide an empty structure to avoid inconsistency
  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    terms_of_use                  = []
  }

  session_controls = {
    application_enforced_restrictions = {
      is_enabled = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.session_controls.application_enforced_restrictions.is_enabled
    }
  }
}
```

### Example 16: Create policy from template - Require phishing-resistant MFA for admins

```terraform
# Test 16: Create a conditional access policy from template - 
# Require phishing-resistant multifactor authentication for admins

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "phishing_resistant_mfa_admins" {
  name = "Require phishing-resistant multifactor authentication for admins"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_phishing_resistant_mfa_admins" {
  display_name = "acc-test-ca-policy-template-require-phishing-resistant-mfa-for-admins-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = "OR" # Template specifies AND, but API normalizes to OR when only authentication_strength is used
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.grant_controls.authentication_strength.id
    }
  }
}
```

### Example 17: Create policy from template - Require MFA for Microsoft admin portals

```terraform
# Test 17: Create a conditional access policy from template - 
# Require multifactor authentication for Microsoft admin portals

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_admin_portals" {
  name = "Require multifactor authentication for Microsoft admin portals"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_admin_portals" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-admin-portals-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.grant_controls.terms_of_use
    authentication_strength = {
      id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_admin_portals.details.grant_controls.authentication_strength.id
    }
  }
}
```

### Example 18: Create policy from template - Block access to Office365 apps for users with insider risk

```terraform
# Test 18: Create a conditional access policy from template - 
# Block access to Office365 apps for users with insider risk

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_o365_insider_risk" {
  name = "Block access to Office365 apps for users with insider risk"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_block_o365_insider_risk" {
  display_name = "acc-test-ca-policy-template-block-access-o365-insider-risk-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.service_principal_risk_levels
    insider_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.insider_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.exclude_roles

      exclude_guests_or_external_users = {
        guest_or_external_user_types = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.exclude_guests_or_external_users.guest_or_external_user_types
        external_tenants = {
          membership_kind = "all" # Template has externalTenants: null, but external_tenants is required in the resource
        }
      }
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.terms_of_use
  }
}
```

### Example 19: Create policy from template - Require MDM-enrolled and compliant device

```terraform
# Test 19: Create a conditional access policy from template - 
# Require MDM-enrolled and compliant device to access cloud apps for all users (Preview)

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mdm_compliant_device" {
  name = "Require MDM-enrolled and compliant device to access cloud apps for all users (Preview)"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mdm_compliant_device" {
  display_name = "acc-test-ca-policy-template-require-mdm-enrolled-compliant-device-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.service_principal_risk_levels
    agent_id_risk_levels          = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.agent_id_risk_levels
    insider_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.insider_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mdm_compliant_device.details.grant_controls.terms_of_use
  }
}
```

### Example 21: Create policy from template - Block high-risk agent identities

```terraform
# Test 21: Create a conditional access policy from template - 
# Block high risk agent identities from accessing resources

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_high_risk_agents" {
  name = "Block high risk agent identities from accessing resources"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_block_high_risk_agents" {
  display_name = "acc-test-ca-policy-template-block-high-risk-agent-identities-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.service_principal_risk_levels
    agent_id_risk_levels          = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.agent_id_risk_levels
    insider_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.insider_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.users.include_users
      exclude_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.users.exclude_users
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.users.exclude_roles
    }

    client_applications = {
      include_service_principals          = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.client_applications.include_service_principals
      include_agent_id_service_principals = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.client_applications.include_agent_id_service_principals
      exclude_service_principals          = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.client_applications.exclude_service_principals
      exclude_agent_id_service_principals = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.conditions.client_applications.exclude_agent_id_service_principals
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_high_risk_agents.details.grant_controls.terms_of_use
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the conditional access template.
- `template_id` (String) The unique identifier (GUID) of the conditional access template.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `description` (String) Description of what the conditional access template does.
- `details` (Attributes) The policy configuration details including conditions, grant controls, and session controls. (see [below for nested schema](#nestedatt--details))
- `id` (String) The unique identifier for this data source operation.
- `scenarios` (Set of String) Set of scenarios this template applies to (e.g., secureFoundation, zeroTrust, remoteWork, protectAdmins).

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--details"></a>
### Nested Schema for `details`

Read-Only:

- `conditions` (Attributes) The conditions that must be met for the policy to apply. (see [below for nested schema](#nestedatt--details--conditions))
- `grant_controls` (Attributes) The grant controls applied by the policy. (see [below for nested schema](#nestedatt--details--grant_controls))
- `session_controls` (Attributes) The session controls applied by the policy. (see [below for nested schema](#nestedatt--details--session_controls))

<a id="nestedatt--details--conditions"></a>
### Nested Schema for `details.conditions`

Read-Only:

- `agent_id_risk_levels` (Set of String) Agent identity risk levels included in the policy.
- `applications` (Attributes) Application conditions. (see [below for nested schema](#nestedatt--details--conditions--applications))
- `client_app_types` (List of String) Client application types included in the policy.
- `client_applications` (Attributes) Client application conditions. (see [below for nested schema](#nestedatt--details--conditions--client_applications))
- `devices` (Attributes) Device conditions. (see [below for nested schema](#nestedatt--details--conditions--devices))
- `insider_risk_levels` (Set of String) Insider risk levels included in the policy.
- `locations` (Attributes) Location conditions. (see [below for nested schema](#nestedatt--details--conditions--locations))
- `platforms` (Attributes) Platform conditions. (see [below for nested schema](#nestedatt--details--conditions--platforms))
- `service_principal_risk_levels` (List of String) Service principal risk levels included in the policy.
- `sign_in_risk_levels` (List of String) Sign-in risk levels included in the policy.
- `user_risk_levels` (List of String) User risk levels included in the policy.
- `users` (Attributes) User and group conditions. (see [below for nested schema](#nestedatt--details--conditions--users))

<a id="nestedatt--details--conditions--applications"></a>
### Nested Schema for `details.conditions.applications`

Read-Only:

- `exclude_applications` (List of String) Applications excluded from the policy.
- `include_applications` (List of String) Applications included in the policy.
- `include_authentication_context_class_references` (List of String) Authentication context class references included in the policy.
- `include_user_actions` (List of String) User actions included in the policy.


<a id="nestedatt--details--conditions--client_applications"></a>
### Nested Schema for `details.conditions.client_applications`

Read-Only:

- `exclude_agent_id_service_principals` (List of String) Agent identity service principals excluded from the policy.
- `exclude_service_principals` (List of String) Service principals excluded from the policy.
- `include_agent_id_service_principals` (List of String) Agent identity service principals included in the policy.
- `include_service_principals` (List of String) Service principals included in the policy.


<a id="nestedatt--details--conditions--devices"></a>
### Nested Schema for `details.conditions.devices`

Read-Only:

- `device_filter` (Attributes) Device filter configuration. (see [below for nested schema](#nestedatt--details--conditions--devices--device_filter))
- `exclude_device_states` (List of String) Device states excluded from the policy.
- `exclude_devices` (List of String) Devices excluded from the policy.
- `include_device_states` (List of String) Device states included in the policy.
- `include_devices` (List of String) Devices included in the policy.

<a id="nestedatt--details--conditions--devices--device_filter"></a>
### Nested Schema for `details.conditions.devices.device_filter`

Read-Only:

- `mode` (String) Filter mode (include or exclude).
- `rule` (String) Filter rule expression.



<a id="nestedatt--details--conditions--locations"></a>
### Nested Schema for `details.conditions.locations`

Read-Only:

- `exclude_locations` (List of String) Locations excluded from the policy.
- `include_locations` (List of String) Locations included in the policy.


<a id="nestedatt--details--conditions--platforms"></a>
### Nested Schema for `details.conditions.platforms`

Read-Only:

- `exclude_platforms` (List of String) Platforms excluded from the policy.
- `include_platforms` (List of String) Platforms included in the policy.


<a id="nestedatt--details--conditions--users"></a>
### Nested Schema for `details.conditions.users`

Read-Only:

- `exclude_groups` (List of String) Groups excluded from the policy.
- `exclude_guests_or_external_users` (Attributes) Guest or external user exclusion conditions. (see [below for nested schema](#nestedatt--details--conditions--users--exclude_guests_or_external_users))
- `exclude_roles` (List of String) Roles excluded from the policy.
- `exclude_users` (List of String) Users excluded from the policy.
- `include_groups` (List of String) Groups included in the policy.
- `include_guests_or_external_users` (Attributes) Guest or external user inclusion conditions. (see [below for nested schema](#nestedatt--details--conditions--users--include_guests_or_external_users))
- `include_roles` (List of String) Roles included in the policy.
- `include_users` (List of String) Users included in the policy.

<a id="nestedatt--details--conditions--users--exclude_guests_or_external_users"></a>
### Nested Schema for `details.conditions.users.exclude_guests_or_external_users`

Read-Only:

- `external_tenants` (Attributes) External tenant configuration. (see [below for nested schema](#nestedatt--details--conditions--users--exclude_guests_or_external_users--external_tenants))
- `guest_or_external_user_types` (Set of String) Types of guest or external users.

<a id="nestedatt--details--conditions--users--exclude_guests_or_external_users--external_tenants"></a>
### Nested Schema for `details.conditions.users.exclude_guests_or_external_users.external_tenants`

Read-Only:

- `membership_kind` (String) Membership kind (e.g., all).



<a id="nestedatt--details--conditions--users--include_guests_or_external_users"></a>
### Nested Schema for `details.conditions.users.include_guests_or_external_users`

Read-Only:

- `external_tenants` (Attributes) External tenant configuration. (see [below for nested schema](#nestedatt--details--conditions--users--include_guests_or_external_users--external_tenants))
- `guest_or_external_user_types` (Set of String) Types of guest or external users.

<a id="nestedatt--details--conditions--users--include_guests_or_external_users--external_tenants"></a>
### Nested Schema for `details.conditions.users.include_guests_or_external_users.external_tenants`

Read-Only:

- `membership_kind` (String) Membership kind (e.g., all).





<a id="nestedatt--details--grant_controls"></a>
### Nested Schema for `details.grant_controls`

Read-Only:

- `authentication_strength` (Attributes) Authentication strength requirements. (see [below for nested schema](#nestedatt--details--grant_controls--authentication_strength))
- `built_in_controls` (List of String) Built-in grant controls (e.g., mfa, compliantDevice, domainJoinedDevice).
- `custom_authentication_factors` (List of String) Custom authentication factors.
- `operator` (String) Logical operator for grant controls (AND or OR).
- `terms_of_use` (List of String) Terms of use.

<a id="nestedatt--details--grant_controls--authentication_strength"></a>
### Nested Schema for `details.grant_controls.authentication_strength`

Read-Only:

- `allowed_combinations` (List of String) The allowed authentication method combinations.
- `created_date_time` (String) The date and time when the authentication strength policy was created.
- `description` (String) The description of the authentication strength policy.
- `display_name` (String) The display name of the authentication strength policy.
- `id` (String) The unique identifier of the authentication strength policy.
- `modified_date_time` (String) The date and time when the authentication strength policy was last modified.
- `policy_type` (String) The type of the authentication strength policy (e.g., builtIn).
- `requirements_satisfied` (String) The requirements satisfied by this authentication strength (e.g., mfa).



<a id="nestedatt--details--session_controls"></a>
### Nested Schema for `details.session_controls`

Read-Only:

- `application_enforced_restrictions` (Attributes) Application enforced restrictions. (see [below for nested schema](#nestedatt--details--session_controls--application_enforced_restrictions))
- `cloud_app_security` (Attributes) Session control to apply cloud app security. (see [below for nested schema](#nestedatt--details--session_controls--cloud_app_security))
- `continuous_access_evaluation` (Attributes) Session control for continuous access evaluation settings. (see [below for nested schema](#nestedatt--details--session_controls--continuous_access_evaluation))
- `disable_resilience_defaults` (Boolean) Session control that determines whether it's acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not.
- `global_secure_access_filtering_profile` (Attributes) Session control to link to Global Secure Access security profiles or filtering profiles. (see [below for nested schema](#nestedatt--details--session_controls--global_secure_access_filtering_profile))
- `persistent_browser` (Attributes) Persistent browser session settings. (see [below for nested schema](#nestedatt--details--session_controls--persistent_browser))
- `secure_sign_in_session` (Attributes) Session control to require sign in sessions to be bound to a device. (see [below for nested schema](#nestedatt--details--session_controls--secure_sign_in_session))
- `sign_in_frequency` (Attributes) Sign-in frequency settings. (see [below for nested schema](#nestedatt--details--session_controls--sign_in_frequency))

<a id="nestedatt--details--session_controls--application_enforced_restrictions"></a>
### Nested Schema for `details.session_controls.application_enforced_restrictions`

Read-Only:

- `is_enabled` (Boolean) Whether application enforced restrictions are enabled.


<a id="nestedatt--details--session_controls--cloud_app_security"></a>
### Nested Schema for `details.session_controls.cloud_app_security`

Read-Only:

- `cloud_app_security_type` (String) The possible values are: mcasConfigured, monitorOnly, blockDownloads.
- `is_enabled` (Boolean) Specifies whether the session control is enabled.


<a id="nestedatt--details--session_controls--continuous_access_evaluation"></a>
### Nested Schema for `details.session_controls.continuous_access_evaluation`

Read-Only:

- `mode` (String) Specifies continuous access evaluation settings. The possible values are: strictEnforcement, disabled, unknownFutureValue, strictLocation.


<a id="nestedatt--details--session_controls--global_secure_access_filtering_profile"></a>
### Nested Schema for `details.session_controls.global_secure_access_filtering_profile`

Read-Only:

- `is_enabled` (Boolean) Specifies whether the session control is enabled.
- `profile_id` (String) Specifies the distinct identifier that is assigned to the security profile or filtering profile.


<a id="nestedatt--details--session_controls--persistent_browser"></a>
### Nested Schema for `details.session_controls.persistent_browser`

Read-Only:

- `is_enabled` (Boolean) Whether persistent browser session is enabled.
- `mode` (String) Persistent browser mode (e.g., always, never).


<a id="nestedatt--details--session_controls--secure_sign_in_session"></a>
### Nested Schema for `details.session_controls.secure_sign_in_session`

Read-Only:

- `is_enabled` (Boolean) Specifies whether the session control is enabled.


<a id="nestedatt--details--session_controls--sign_in_frequency"></a>
### Nested Schema for `details.session_controls.sign_in_frequency`

Read-Only:

- `authentication_type` (String) The authentication type this frequency applies to.
- `frequency_interval` (String) The frequency interval (e.g., timeBased, everyTime).
- `is_enabled` (Boolean) Whether sign-in frequency is enabled.
- `type` (String) The frequency type (e.g., hours, days).
- `value` (Number) The frequency value.
