resource "microsoft365_graph_beta_groups_group" "intune_enrollment" {
  display_name     = "Intune automatic enrollment"
  mail_nickname    = "intune-automatic-enrollment"
  mail_enabled     = false
  security_enabled = true
}

# Limit Microsoft Intune MDM automatic enrollment to one group and set every MDM setting.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "mdm" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  compliance_url   = "https://portal.manage.microsoft.com/?portalAction=Compliance"
  discovery_url    = "https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"
  terms_of_use_url = "https://portal.manage.microsoft.com/TermsofUse.aspx"

  is_mdm_enrollment_during_registration_disabled = false

  applies_to      = "selected"
  included_groups = [microsoft365_graph_beta_groups_group.intune_enrollment.id]
}

# Turn off the Microsoft Intune MAM user scope.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "mam" {
  policy_type = "app"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  applies_to = "none"
}
