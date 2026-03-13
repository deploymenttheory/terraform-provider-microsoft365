resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "example" {
  tenant_id           = "12345678-1234-1234-1234-123456789012"
  is_service_provider = false
  hard_delete         = false

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
    is_hybrid_azure_ad_joined_device_accepted = false
  }

  automatic_user_consent_settings = {
    inbound_allowed  = true
    outbound_allowed = true
  }
}
