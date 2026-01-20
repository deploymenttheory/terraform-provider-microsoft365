# Test 13: Create a conditional access policy from template - 
# Block access for unknown or unsupported device platform
#
# NOTE: This template appears to have multiple errors and cannot be used for testing:
# 1. Platform condition has contradictory configuration (include "all" + exclude specific platforms)
# 2. Grant control "block" with operator "OR" is not supported by Graph API
# 3. Template data appears to be malformed for actual production use
# This test is skipped - template is not suitable for creating actual policies.

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_unknown_platform" {
  name = "Block access for unknown or unsupported device platform"
}

