# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# AU004: Update Test - Step 4 (Pause Dynamic Membership)
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au004_update" {
  display_name                     = "acc-test-au004-update-${random_string.suffix.result}"
  description                      = "Paused dynamic membership"
  visibility                       = "HiddenMembership"
  membership_type                  = "Dynamic"
  membership_rule                  = "(user.country -eq \"United States\")"
  membership_rule_processing_state = "Paused"
  hard_delete                      = true
}
