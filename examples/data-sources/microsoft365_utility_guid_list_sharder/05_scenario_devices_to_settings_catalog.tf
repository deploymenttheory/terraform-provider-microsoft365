# Scenario 1: Devices → Groups → Settings Catalog Deployment Rings
# Use case: Roll out Windows Update policies in phases across managed devices

# Distribute Windows devices into 4 deployment rings (5%, 15%, 30%, 50%)
data "microsoft365_utility_guid_list_sharder" "windows_devices" {
  resource_type     = "devices"
  odata_filter      = "operatingSystem eq 'Windows' and trustType eq 'AzureAd'"
  shard_percentages = [5, 15, 30, 50]
  strategy          = "percentage"
  seed              = "windows-updates-2024"
}

# Create groups for each deployment ring
resource "microsoft365_graph_beta_group" "ring_0_validation" {
  display_name     = "Windows Updates - Ring 0 (5% Validation)"
  mail_nickname    = "win-updates-ring-0"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_0"]
}

resource "microsoft365_graph_beta_group" "ring_1_pilot" {
  display_name     = "Windows Updates - Ring 1 (15% Pilot)"
  mail_nickname    = "win-updates-ring-1"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_1"]
}

resource "microsoft365_graph_beta_group" "ring_2_broad" {
  display_name     = "Windows Updates - Ring 2 (30% Broad)"
  mail_nickname    = "win-updates-ring-2"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_2"]
}

resource "microsoft365_graph_beta_group" "ring_3_production" {
  display_name     = "Windows Updates - Ring 3 (50% Production)"
  mail_nickname    = "win-updates-ring-3"
  security_enabled = true
  mail_enabled     = false

  members = data.microsoft365_utility_guid_list_sharder.windows_devices.shards["shard_3"]
}

# Settings Catalog Policy for Ring 0 (immediate deployment)
resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy_json" "ring_0_validation" {
  name               = "Windows Updates - Ring 0 (5% Validation)"
  description        = "Immediate deployment for validation devices"
  platforms          = "windows10"
  technologies       = ["mdm"]
  role_scope_tag_ids = ["0"]

  settings = jsonencode({
    settings = [
      {
        id = "0"
        settingInstance = {
          "@odata.type"       = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
          settingDefinitionId = "device_vendor_msft_policy_config_update_deferqualityupdatesperiodindays"
          choiceSettingValue = {
            value    = "device_vendor_msft_policy_config_update_deferqualityupdatesperiodindays_0"
            children = []
          }
        }
      }
    ]
  })

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_0_validation.id
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_1_pilot.id
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_2_broad.id
      filter_type = "none"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_group.ring_3_production.id
      filter_type = "none"
    }
  ]
}