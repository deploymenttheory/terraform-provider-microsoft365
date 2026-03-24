# AU005: Administrative Unit with Hard Delete
# Creates an administrative unit that will be permanently deleted when destroyed
# instead of being moved to the deleted items container
resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au005_hard_delete" {
  display_name = "Temporary Project Team"
  description  = "Administrative unit for temporary project that will be permanently deleted"
  hard_delete  = true
}
