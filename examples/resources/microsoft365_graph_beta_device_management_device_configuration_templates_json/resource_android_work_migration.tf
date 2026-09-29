# Android Enterprise — Move to Android Management API.
# This example disables migration. Set disableMigration to false to enable it.
# Keep the profile unassigned while evaluating it; assigning an enabled policy migrates devices.

# Keep this profile unassigned. Assignment starts an irreversible device migration.
resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "android_work_migration" {
  display_name       = "Android Work Migration"
  description        = "Template configuration example"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileMigrationConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "disableMigration"                            = true
  })
}
