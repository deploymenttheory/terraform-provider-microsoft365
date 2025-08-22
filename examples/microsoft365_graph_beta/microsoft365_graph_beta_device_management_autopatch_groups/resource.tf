resource "microsoft365_graph_beta_device_management_autopatch_groups" "auto_patch_group" {
  name        = "auto-patch-group"
  description = ""

  global_user_managed_aad_groups = []

  # Deployment Groups
  deployment_groups = [
    {
      name = "auto-patch-group - Test"
      user_managed_aad_groups = [
        {
          id   = "00000000-0000-0000-0000-000000000000"
          name = "group-name-01"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "auto-patch-group - Test"
        is_update_settings_modified = false
        device_configuration_setting = {
          policy_id            = "000"
          update_behavior      = "AutoInstallAndRestart"
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
      name = "auto-patch-group - Ring1"
      user_managed_aad_groups = [
        {
          id   = "00000000-0000-0000-0000-000000000000"
          name = "group-name-02"
          type = 0
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "auto-patch-group - Ring1"
        is_update_settings_modified = false
        device_configuration_setting = {
          policy_id            = "000"
          update_behavior      = "AutoInstallAndRestart"
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
      name = "auto-patch-group - Last"
      user_managed_aad_groups = [
        {
          id   = "00000000-0000-0000-0000-000000000000"
          name = "group-name-03"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "auto-patch-group - Last"
        is_update_settings_modified = false
        device_configuration_setting = {
          policy_id            = "000"
          update_behavior      = "AutoInstallAndRestart"
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

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}