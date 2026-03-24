# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# AU004: Update Test - Step 1 (Initial Configuration)
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_administrative_unit" "au004_update" {
  display_name = "acc-test-au004-update-${random_string.suffix.result}"
  description  = "Initial description for update testing"
  hard_delete  = true
}
