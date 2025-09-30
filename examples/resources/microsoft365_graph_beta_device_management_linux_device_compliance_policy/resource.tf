resource "microsoft365_graph_beta_device_management_linux_device_compliance_policy" "json_config" {
  name        = "linux device compliance policy"
  description = "example"

  # Distribution restrictions matching the JSON
  distribution_allowed_distros = [
    {
      type            = "ubuntu"
      minimum_version = "11"
      maximum_version = "10"
    },
    {
      type            = "rhel"
      minimum_version = "10"
      maximum_version = "9"
    }
  ]

  # Custom compliance configuration matching the JSON
  custom_compliance_required         = true
  custom_compliance_discovery_script = "996bf3d2-958b-478d-b3d4-4c5017a5650e"
  custom_compliance_rules = jsonencode({
    "Rules" = [
      {
        "SettingName" = "BiosVersion"
        "Operator"    = "GreaterEquals"
        "DataType"    = "Version"
        "Operand"     = "2.3"
        "MoreInfoUrl" = "https://your-website.com"
        "RemediationStrings" = [
          {
            "Language"    = "en_US"
            "Title"       = "BIOS Version needs to be upgraded to at least 2.3. Value discovered was {ActualValue}."
            "Description" = "BIOS must be updated. Please refer to the link above"
          }
        ]
      },
      {
        "SettingName" = "TPMChipPresent"
        "Operator"    = "IsEquals"
        "DataType"    = "Boolean"
        "Operand"     = true
        "MoreInfoUrl" = "https://bing.com"
        "RemediationStrings" = [
          {
            "Language"    = "en_US"
            "Title"       = "TPM chip must be enabled."
            "Description" = "TPM chip must be enabled. Please refer to the link above"
          }
        ]
      }
    ]
  })

  # Device encryption required
  device_encryption_required = true

  # Password policy settings - all set to minimum value 1
  password_policy_minimum_digits    = 1
  password_policy_minimum_length    = 1
  password_policy_minimum_lowercase = 1
  password_policy_minimum_symbols   = 1
  password_policy_minimum_uppercase = 1

  # Optional scheduled action
  scheduled_actions = [
    {
      rule_name = "PasswordRequired"
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 0
        }
      ]
    }
  ]
  # Optional Assignments
  assignments = [
    # Assignment targeting all devices with an include filter
    {
      type = "allDevicesAssignmentTarget"
    },
    # Assignment targeting all licensed users with an exclude filter
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    # Assignment targeting a specific group with include filter
    {
      type     = "groupAssignmentTarget"
      group_id = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
    },
    # Assignment targeting a specific group with exclude filter
    {
      type     = "groupAssignmentTarget"
      group_id = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
    },
    # Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "b8c661c2-fa9a-4351-af86-adc1729c343f"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "f6ebd6ff-501e-4b3d-a00b-a2e102c3fa0f"
    }
  ]
}