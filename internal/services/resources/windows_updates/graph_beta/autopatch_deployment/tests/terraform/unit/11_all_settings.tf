resource "microsoft365_graph_beta_windows_updates_autopatch_deployment" "test" {
  content = {
    catalog_entry_id   = "d0c03fbb-43b9-4dff-840b-974ef227384d"
    catalog_entry_type = "qualityUpdate"
  }

  settings = {
    schedule = {
      start_date_time = "2024-01-15T10:00:00Z"
      gradual_rollout = {
        duration_between_offers = "P1D"
        devices_per_offer       = 50
      }
    }
    monitoring = {
      monitoring_rules = [
        {
          signal    = "rollback"
          threshold = 10
          action    = "pauseDeployment"
        }
      ]
    }
    user_experience = {
      days_until_forced_reboot = 14
      offer_as_optional        = true
    }
    expedite = {
      is_expedited      = true
      is_readiness_test = false
    }
    content_applicability = {
      safeguard = {
        disabled_safeguard_profiles = [
          {
            category = "likelyIssues"
          }
        ]
      }
    }
  }
}
