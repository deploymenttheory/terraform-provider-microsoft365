# Microsoft Defender for Endpoint connector.
# The connector ID is the fixed partner ID of Microsoft Defender for Endpoint.
# If the tenant already has this connector (the "Microsoft Intune connection" is
# turned on in the Microsoft Defender portal), import it first so Terraform
# manages the existing settings instead of replacing them on create.
resource "microsoft365_graph_beta_device_management_mobile_threat_defense_connector" "microsoft_defender" {
  id = "fc780465-2017-40d4-a0c5-307022471b92"

  # Use partner data in compliance evaluations
  android_enabled = true
  ios_enabled     = true
  mac_enabled     = false
  windows_enabled = true

  # Use partner data in app protection (MAM) evaluations
  android_mobile_application_management_enabled = true
  ios_mobile_application_management_enabled     = true
  windows_mobile_application_management_enabled = false

  # Compliance behaviour
  partner_unsupported_os_version_blocked     = true
  partner_unresponsiveness_threshold_in_days = 7

  # iOS/iPadOS metadata collection
  allow_partner_to_collect_ios_application_metadata          = true
  allow_partner_to_collect_ios_personal_application_metadata = false
  allow_partner_to_collect_ios_certificate_metadata          = false
  allow_partner_to_collect_ios_personal_certificate_metadata = false

  # Defender for Endpoint only: configuration profile management through Defender
  microsoft_defender_for_endpoint_attach_enabled = true

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
