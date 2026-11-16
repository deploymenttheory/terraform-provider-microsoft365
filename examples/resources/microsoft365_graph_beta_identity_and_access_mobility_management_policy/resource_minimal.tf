# Adopt the Microsoft Intune MDM user scope without changing any setting.
# Settings that are not configured keep their current tenant value.
resource "microsoft365_graph_beta_identity_and_access_mobility_management_policy" "mdm" {
  policy_type = "device"
  policy_id   = "0000000a-0000-0000-c000-000000000000"
}
