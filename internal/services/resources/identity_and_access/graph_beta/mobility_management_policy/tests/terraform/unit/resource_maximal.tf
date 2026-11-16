# Maximal: every writable MDM setting, scoped to two groups.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "maximal" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  compliance_url   = "https://portal.manage.microsoft.com/?portalAction=Compliance&unit-test=maximal"
  discovery_url    = "https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"
  terms_of_use_url = "https://portal.manage.microsoft.com/TermsofUse.aspx?unit-test=maximal"

  is_mdm_enrollment_during_registration_disabled = true

  applies_to = "selected"
  included_groups = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
  ]
}
