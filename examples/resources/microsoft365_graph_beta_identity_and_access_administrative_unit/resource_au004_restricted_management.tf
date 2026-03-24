# AU004: Restricted Management Administrative Unit
# Creates an administrative unit with restricted member management
# Only administrators with specific permissions can manage members
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au004_restricted" {
  display_name                    = "Managed Devices"
  description                     = "Administrative unit for managed devices with restricted management"
  is_member_management_restricted = true
}
