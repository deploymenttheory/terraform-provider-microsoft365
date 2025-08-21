resource "microsoft365_graph_beta_device_management_autopatch_groups" "auto_patch_group" {
  name        = "auto-patch-group"
  description = ""

  # Global User Managed AAD Groups (empty in the example)
  global_user_managed_aad_groups = []

  # Deployment Groups
  deployment_groups = [
    {
      aad_id = "00000000-0000-0000-0000-000000000000"
      name   = "auto-patch-group - Test"
      user_managed_aad_groups = [
        {
          id   = "410a28bd-9c9f-403f-b1b2-4a0bd04e98d9"
          name = "[Azure]-[ConditonalAccess]-[Prod]-[CAD003-PolicyExclude]-[UG]"
          type = 0
        }
      ]
      failed_prerequisite_check_count = 0
      deployment_group_policy_settings = {
        aad_group_name               = "auto-patch-group - Test"
        is_update_settings_modified  = false
        device_configuration_setting = {
          policy_id           = "000"
          update_behavior     = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 1
            deferral     = 0
            grace_period = 0
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
      }
    },
    {
      aad_id = "00000000-0000-0000-0000-000000000000"
      name   = "auto-patch-group - Ring1"
      user_managed_aad_groups = [
        {
          id   = "35d09841-af73-43e6-a59f-024fef1b6b95"
          name = "[Azure]-[ConditonalAccess]-[Prod]-[CAD002-PolicyExclude]-[UG]"
          type = 0
        }
      ]
      failed_prerequisite_check_count = 0
      deployment_group_policy_settings = {
        aad_group_name               = "auto-patch-group - Ring1"
        is_update_settings_modified  = false
        device_configuration_setting = {
          policy_id           = "000"
          update_behavior     = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 2
            deferral     = 1
            grace_period = 2
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
      }
    },
    {
      aad_id = "00000000-0000-0000-0000-000000000000"
      name   = "auto-patch-group - Last"
      user_managed_aad_groups = [
        {
          id   = "48fe6d79-f045-448a-bd74-716db27f0783"
          name = "[Azure]-[ConditonalAccess]-[Prod]-[CAD005-PolicyExclude]-[UG]"
          type = 0
        }
      ]
      failed_prerequisite_check_count = 0
      deployment_group_policy_settings = {
        aad_group_name               = "auto-patch-group - Last"
        is_update_settings_modified  = false
        device_configuration_setting = {
          policy_id           = "000"
          update_behavior     = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 3
            deferral     = 5
            grace_period = 2
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
      }
    }
  ]

  # Driver update settings
  enable_driver_update = true
  
  # Scope tags
  scope_tags = [0]
  
  # Enabled content types
  enabled_content_types = 31

  timeouts {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}