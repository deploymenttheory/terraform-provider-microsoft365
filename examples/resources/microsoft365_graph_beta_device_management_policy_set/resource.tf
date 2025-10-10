# Policy Set with multiple policy types
resource "microsoft365_graph_beta_device_management_policy_set" "example" {
  display_name = "example intune policy set"
  description  = "maximal example"

  items = [
    # Windows MSI line-of-business app example
    {
      type       = "app"
      payload_id = "00000000-0000-0000-0000-000000000000"
      intent     = "required"
    },
    # iOS Store App example
    {
      type       = "app" 
      payload_id = "00000000-0000-0000-0000-000000000000"
      intent     = "required"
      settings = {
        odata_type                  = "#microsoft.graph.iosStoreAppAssignmentSettings"
        vpn_configuration_id        = "00000000-0000-0000-0000-000000000000"
        uninstall_on_device_removal = true
        is_removable                = true
      }
    },
    # Targeted Managed App Configuration Policy Set Items
    {
      type       = "app_configuration_policy"
      payload_id = "A_00000000-0000-0000-0000-000000000000"
    },
    {
      type       = "app_configuration_policy"
      payload_id = "A_00000000-0000-0000-0000-000000000000"
    },
    # Managed App Protection Policy Set Items
    {
      type       = "app_protection_policy"
      payload_id = "T_00000000-0000-0000-0000-000000000000"
    },
    {
      type       = "app_protection_policy"
      payload_id = "T_00000000-0000-0000-0000-000000000000"
    },
    # Device Configuration Policy Set Items
    {
      type       = "device_configuration_profile"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },  
    {
      type       = "device_configuration_profile"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
    # Device Management Configuration Policy Set Items
    {
      type       = "device_management_configuration_policy"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type       = "device_management_configuration_policy"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
    # Device Compliance Policy Set Items
    {
      type       = "device_compliance_policy"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type       = "device_compliance_policy"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
    # Windows Autopilot Deployment Profile Policy Set Items
    {
      type       = "windows_autopilot_deployment_profile"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type       = "windows_autopilot_deployment_profile"
      payload_id = "00000000-0000-0000-0000-000000000000"
    },
  ]

  # Assignments for the policy set
  assignments = [
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000000"
    },
  ]


  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}