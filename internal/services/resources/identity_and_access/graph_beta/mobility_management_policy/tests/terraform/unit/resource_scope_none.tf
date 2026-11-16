# Scope none: removes every included group and changes a URL in the same apply.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "maximal" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"

  compliance_url   = "https://portal.manage.microsoft.com/?portalAction=Compliance&unit-test=none"
  discovery_url    = "https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"
  terms_of_use_url = "https://portal.manage.microsoft.com/TermsofUse.aspx?unit-test=maximal"

  is_mdm_enrollment_during_registration_disabled = false

  applies_to = "none"
}
