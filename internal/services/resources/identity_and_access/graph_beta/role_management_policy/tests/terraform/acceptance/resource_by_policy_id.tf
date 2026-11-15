# Adopts the PIM policy of the built-in Usage Summary Reports Reader role by its policy ID.
variable "policy_id" {
  type = string
}

resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  policy_id = var.policy_id

  expiration_rules = {
    Expiration_Admin_Eligibility = {
      is_expiration_required = false
      maximum_duration       = "P180D"
    }
  }
}
