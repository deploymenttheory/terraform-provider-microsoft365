# Third-party Mobile Threat Defense partner connector (Lookout for Work).
# Each partner has a fixed connector ID that is the same in every tenant.
# The partner_state attribute stays "notSetUp" until the partner console
# completes the connection to Intune.
resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "lookout" {
  id = "2c7790de-8b02-4814-85cf-e0c59380dee8"

  android_enabled = true
  ios_enabled     = true

  partner_unsupported_os_version_blocked     = false
  partner_unresponsiveness_threshold_in_days = 14

  allow_partner_to_collect_ios_application_metadata = true
}
