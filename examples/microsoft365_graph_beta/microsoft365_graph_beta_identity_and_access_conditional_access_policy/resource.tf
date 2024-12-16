resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "example_policy" {
  display_name = "Example Conditional Access Policy"
  state        = "enabled"

  conditions = {
    applications ={
      include_applications = ["All"]
      exclude_applications = ["MicrosoftAdminPortals"]
      application_filter ={
        mode = "include"
        rule = "(appId -eq '11111111-1111-1111-1111-111111111111')"
      }
      include_user_actions = ["urn:user:registersecurityinfo"]
    }

    authentication_flows = {
      transfer_methods = "deviceCodeFlow"
    }

    users = {
      include_users  = ["All"]
      exclude_users  = ["11111111-1111-1111-1111-111111111111"]
      include_groups = ["22222222-2222-2222-2222-222222222222"]
    }

    client_applications ={
      include_service_principals = ["ServicePrincipalsInMyTenant"]
      exclude_service_principals = ["33333333-3333-3333-3333-333333333333"]
      service_principal_filter ={
        mode = "include"
        rule = "(servicePrincipalId -eq '44444444-4444-4444-4444-444444444444')"
      }
    }

    client_app_types = ["all"]

    locations ={
      include_locations = ["All"]
      exclude_locations = ["55555555-5555-5555-5555-555555555555"]
    }

    platforms ={
      include_platforms = ["android", "iOS"]
      exclude_platforms = ["windows"]
    }

    device_states ={
      include_states = ["All"]
      exclude_states = ["Compliant"]
    }

    devices ={
      include_devices = ["All"]
      exclude_devices = ["DomainJoined"]
      device_filter ={
        mode = "exclude"
        rule = "(device.deviceId -eq '66666666-6666-6666-6666-666666666666')"
      }
    }

    sign_in_risk_levels = ["medium", "high"]
    user_risk_levels    = ["low", "medium"]
  }

  grant_controls ={
    operator                      = "OR"
    built_in_controls             = ["mfa", "compliantDevice"]
    custom_authentication_factors = ["77777777-7777-7777-7777-777777777777"]
    terms_of_use                  = ["88888888-8888-8888-8888-888888888888"]

    authentication_strength ={
      id                     = "99999999-9999-9999-9999-999999999999"
      display_name           = "Example Authentication Strength"
      description            = "A description for the authentication strength."
      policy_type            = "required"
      requirements_satisfied = "mfa"
      allowed_combinations   = ["password", "biometric"]
    }
  }

  session_controls ={
    application_enforced_restrictions ={
      is_enabled = true
    }

    cloud_app_security ={
      is_enabled              = true
      cloud_app_security_type = "monitorOnly"
    }

    continuous_access_evaluation ={
      mode = "strictEnforcement"
    }

    persistent_browser ={
      is_enabled = true
      mode       = "always"
    }

    sign_in_frequency ={
      is_enabled          = true
      type                = "days"
      value               = 1
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
    }

    disable_resilience_defaults = false

    secure_sign_in_session ={
      is_enabled = true
    }
  }
  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}
