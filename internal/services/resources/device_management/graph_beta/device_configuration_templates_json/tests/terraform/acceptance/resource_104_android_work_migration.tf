resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Keep this profile unassigned. Assignment starts an irreversible device migration.
resource "microsoft365_graph_beta_device_management_device_configuration_templates_json" "test" {
  display_name       = "acc-test-android-work-migration-${random_string.suffix.result}"
  description        = "Device configuration template test"
  role_scope_tag_ids = ["0"]
  settings = jsonencode({
    "@odata.type"                                 = "#microsoft.graph.androidWorkProfileMigrationConfiguration"
    "deviceManagementApplicabilityRuleOsEdition"  = null
    "deviceManagementApplicabilityRuleOsVersion"  = null
    "deviceManagementApplicabilityRuleDeviceMode" = null
    "disableMigration"                            = false
  })
}
