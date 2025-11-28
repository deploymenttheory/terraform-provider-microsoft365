resource "microsoft365_graph_beta_device_management_autopatch_groups" "unit-test_autopatch_group" {
  name        = "unit-test-autopatch-group"
  description = "unit-test"

  global_user_managed_aad_groups = [
    {
      id   = "550dba96-abd7-4ef0-9cf1-25be705f676c"
      type = "None"
    },
    {
      id   = "6a08c3a0-1693-4089-80cb-f9c2f8063a3b"
      type = "None"
    }
  ]

  deployment_groups = [
    {
      name = "unit-test-autopatch-group - test"
      user_managed_aad_groups = [
        {
          id   = "5a0832d3-19b9-4f78-9b50-906774ac4d49"
          type = "Device"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "unit-test-autopatch-group - test"
        is_update_settings_modified = false
        device_configuration_setting = {
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
        dnf_update_cloud_setting = {
          approval_type               = "Automatic"
          deployment_deferral_in_days = 0
        }
        office_dcv2_setting = {
          deadline                  = 1
          deferral                  = 0
          hide_update_notifications = false
          target_channel            = "MonthlyEnterprise"
        }
        edge_dcv2_setting = {
          target_channel = "Beta"
        }
        feature_update_anchor_cloud_setting = {
          target_os_version                                         = "Windows 11, version 25H2"
          install_latest_windows10_on_windows11_ineligible_device = true
        }
      }
    },
    {
      name         = "unit-test - Ring1"
      distribution = 75
      user_managed_aad_groups = [
        {
          id   = "0fe3d2cb-62ae-4fa4-858f-a122061ada62"
          type = "Device"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "unit-test - Ring1"
        is_update_settings_modified = false
        device_configuration_setting = {
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
        dnf_update_cloud_setting = {
          approval_type               = "Automatic"
          deployment_deferral_in_days = 1
        }
        office_dcv2_setting = {
          deadline                  = 2
          deferral                  = 1
          hide_update_notifications = false
          target_channel            = "MonthlyEnterprise"
        }
        edge_dcv2_setting = {
          target_channel = "Stable"
        }
        feature_update_anchor_cloud_setting = {
          target_os_version                                         = "Windows 11, version 25H2"
          install_latest_windows10_on_windows11_ineligible_device = true
        }
      }
    },
    {
      name         = "unit-test-autopatch-group - Ring2"
      distribution = 25
      user_managed_aad_groups = [
        {
          id   = "20f4b274-0ca2-406a-ae20-12cf99730d62"
          type = "Device"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "unit-test-autopatch-group - Ring2"
        is_update_settings_modified = false
        device_configuration_setting = {
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
        dnf_update_cloud_setting = {
          approval_type = "Manual"
        }
        office_dcv2_setting = {
          deadline                  = 3
          deferral                  = 5
          hide_update_notifications = false
          target_channel            = "MonthlyEnterprise"
        }
        edge_dcv2_setting = {
          target_channel = "Stable"
        }
        feature_update_anchor_cloud_setting = {
          target_os_version                                         = "Windows 11, version 25H2"
          install_latest_windows10_on_windows11_ineligible_device = true
        }
      }
    },
    {
      name = "unit-test-autopatch-group - Last"
      user_managed_aad_groups = [
        {
          id   = "6c57971c-4369-4569-9d4b-29c9351665e1"
          type = "Device"
        }
      ]
      deployment_group_policy_settings = {
        aad_group_name              = "unit-test-autopatch-group  - Last"
        is_update_settings_modified = false
        device_configuration_setting = {
          update_behavior      = "AutoInstallAndRestart"
          notification_setting = "DefaultNotifications"
          quality_deployment_settings = {
            deadline     = 5
            deferral     = 9
            grace_period = 2
          }
          feature_deployment_settings = {
            deadline = 5
            deferral = 0
          }
        }
        dnf_update_cloud_setting = {
          approval_type = "Manual"
        }
        office_dcv2_setting = {
          deadline                  = 5
          deferral                  = 9
          hide_update_notifications = false
          target_channel            = "MonthlyEnterprise"
        }
        edge_dcv2_setting = {
          target_channel = "Stable"
        }
        feature_update_anchor_cloud_setting = {
          target_os_version                                         = "Windows 11, version 25H2"
          install_latest_windows10_on_windows11_ineligible_device = true
        }
      }
    }
  ]

  scope_tags = ["0", "1232", "1234"]

  enable_driver_update  = true
  enabled_content_types = 31
}

