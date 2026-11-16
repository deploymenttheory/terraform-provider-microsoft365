# Groups changed: one group kept, one removed and one added.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "maximal" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  compliance_url   = "https://portal.manage.microsoft.com/?portalAction=Compliance&unit-test=maximal"
  discovery_url    = "https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"
  terms_of_use_url = "https://portal.manage.microsoft.com/TermsofUse.aspx?unit-test=maximal"

  is_mdm_enrollment_during_registration_disabled = true

  applies_to = "selected"
  included_groups = [
    "22222222-2222-2222-2222-222222222222",
    "33333333-3333-3333-3333-333333333333",
  ]
}
