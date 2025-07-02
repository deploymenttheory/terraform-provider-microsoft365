resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "maximal" {
  display_name = "Comprehensive Security Policy - Maximal"
  state        = "enabled"

  conditions = {
    client_app_types              = ["all"]
    user_risk_levels              = ["high"]
    sign_in_risk_levels           = ["high", "medium"]
    service_principal_risk_levels = ["high", "medium"]

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = ["00000002-0000-0ff1-ce00-000000000000"] # Office 365 Exchange Online
      include_user_actions                            = ["urn:user:registersecurityinfo"]
      include_authentication_context_class_references = ["c00000000-0000-0000-0000-000000000001"]
      application_filter = {
        mode = "exclude"
        rule = "device.deviceOwnership -eq \"Company\""
      }
    }

    users = {
      include_users  = ["All"]
      exclude_users  = ["GuestsOrExternalUsers"]
      include_groups = ["a1b2c3d4-e5f6-7890-abcd-ef1234567890"]
      exclude_groups = ["f1e2d3c4-b5a6-9870-fedc-ba0987654321"]
      include_roles  = ["62e90394-69f5-4237-9190-012177145e10"] # Global Administrator
      exclude_roles  = ["e3973bdf-4987-49ae-837a-ba8e231c7286"] # Security Reader

      include_guests_or_external_users = {
        guest_or_external_user_types = "internalGuest,b2bCollaborationGuest"
        external_tenants = {
          membership_kind = "enumerated"
          members         = ["12345678-1234-1234-1234-123456789012"]
        }
      }
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = ["iOS", "android"]
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = ["AllTrusted", "11111111-1111-1111-1111-111111111111"]
    }

    devices = {
      include_devices       = ["All"]
      exclude_devices       = ["22222222-2222-2222-2222-222222222222"]
      include_device_states = ["domainJoined"]
      exclude_device_states = ["compliant"]
      device_filter = {
        mode = "include"
        rule = "device.isCompliant -eq True -and device.deviceOwnership -eq \"Company\""
      }
    }
  }

  grant_controls = {
    operator                      = "AND"
    built_in_controls             = ["mfa", "compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = ["33333333-3333-3333-3333-333333333333"]
    terms_of_use                  = ["44444444-4444-4444-4444-444444444444"]
    authentication_strength = {
      id                     = "00000000-0000-0000-0000-000000000004"
      display_name           = "Multifactor authentication"
      description            = "Combinations of methods that satisfy strong authentication, such as a password + SMS"
      policy_type            = "builtIn"
      requirements_satisfied = "mfa"
      allowed_combinations = [
        "password,sms",
        "password,voice",
        "password,hardwareOath",
        "password,softwareOath",
        "password,microsoftAuthenticatorPush",
        "windowsHelloForBusiness",
        "fido2"
      ]
    }
  }

  session_controls = {
    disable_resilience_defaults = false

    application_enforced_restrictions = {
      is_enabled = true
    }

    cloud_app_security = {
      is_enabled              = true
      cloud_app_security_type = "mcasConfigured"
    }

    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 4
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
    }

    persistent_browser = {
      is_enabled = true
      mode       = "always"
    }

    continuous_access_evaluation = {
      mode = "strict"
    }

    secure_sign_in_session = {
      is_enabled = true
    }
  }
}