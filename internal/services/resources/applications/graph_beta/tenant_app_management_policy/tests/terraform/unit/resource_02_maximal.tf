resource "microsoft365_graph_beta_applications_tenant_app_management_policy" "maximal" {
  is_enabled   = true
  display_name = "Custom Tenant App Management Policy"
  description  = "Enforces comprehensive app management restrictions"

  application_restrictions = {
    password_credentials = [
      {
        restriction_type                           = "passwordAddition"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        
        exclude_actors = {
          custom_security_attributes = [
            {
              id       = "PolicyExemptions_AppManagementExemption"
              operator = "equals"
              value    = "ExemptFromPasswordRestriction"
            }
          ]
        }
      },
      {
        restriction_type                           = "passwordLifetime"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        max_lifetime                               = "P90D"
      },
      {
        restriction_type                           = "symmetricKeyLifetime"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        max_lifetime                               = "P30D"
      }
    ]

    key_credentials = [
      {
        restriction_type                           = "asymmetricKeyLifetime"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        max_lifetime                               = "P365D"
      },
      {
        restriction_type                           = "trustedCertificateAuthority"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        certificate_based_application_configuration_ids = [
          "eec5ba11-2fc0-4113-83a2-ed986ed13743",
          "bb8e164b-f9ed-4b98-bc45-65eddc14f4c1"
        ]
      }
    ]

    identifier_uris = {
      non_default_uri_addition = {
        restrict_for_apps_created_after_date_time = "2024-01-01T00:00:00Z"
        exclude_apps_receiving_v2_tokens          = true
        exclude_saml                              = true
        
        exclude_actors = {
          custom_security_attributes = [
            {
              id       = "PolicyExemptions_AppManagementExemption"
              operator = "equals"
              value    = "ExemptFromIdentifierUriRestriction"
            }
          ]
        }
      }
    }
  }

  service_principal_restrictions = {
    password_credentials = [
      {
        restriction_type                           = "passwordLifetime"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        max_lifetime                               = "P90D"
      }
    ]

    key_credentials = [
      {
        restriction_type                           = "asymmetricKeyLifetime"
        restrict_for_apps_created_after_date_time  = "2024-01-01T00:00:00Z"
        max_lifetime                               = "P365D"
      }
    ]
  }

  restore_to_default_upon_delete = false

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
