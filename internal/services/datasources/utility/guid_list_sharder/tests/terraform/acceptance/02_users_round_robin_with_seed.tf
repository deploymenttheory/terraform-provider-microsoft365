# ==============================================================================
# Acceptance Test 02: Users - Round-Robin Strategy (With Seed)
#
# Purpose: Verify round-robin distribution with seed produces exactly equal
# shard sizes AND reproducible results
#
# Use Case: A/B testing, capacity planning, or when you need exact equal
# distribution that you can recreate
#
# Expected Behavior:
# - Exactly equal shard sizes (within ±1)
# - Deterministic shuffle before round-robin
# - Same seed = same distribution every time
# ==============================================================================

# Generate unique test ID for resource names
resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

# ============================================================================
# Test Resources: Create 6 test users for seeded round-robin distribution
# ============================================================================

resource "microsoft365_graph_beta_users_user" "test_user_0" {
  display_name        = "acc-test-sharder-user-0-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-0-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser0${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_1" {
  display_name        = "acc-test-sharder-user-1-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-1-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser1${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_2" {
  display_name        = "acc-test-sharder-user-2-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-2-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser2${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_3" {
  display_name        = "acc-test-sharder-user-3-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-3-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser3${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_4" {
  display_name        = "acc-test-sharder-user-4-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-4-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser4${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_5" {
  display_name        = "acc-test-sharder-user-5-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-5-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser5${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

# Wait for user provisioning
resource "time_sleep" "wait_for_users" {
  depends_on = [
    microsoft365_graph_beta_users_user.test_user_0,
    microsoft365_graph_beta_users_user.test_user_1,
    microsoft365_graph_beta_users_user.test_user_2,
    microsoft365_graph_beta_users_user.test_user_3,
    microsoft365_graph_beta_users_user.test_user_4,
    microsoft365_graph_beta_users_user.test_user_5,
  ]
  create_duration = "30s"
}

# ============================================================================
# Datasource Under Test: GUID List Sharder with Seed
# ============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  depends_on    = [time_sleep.wait_for_users]
  resource_type = "users"
  odata_query   = "startswith(displayName,'acc-test-sharder-user-') and endswith(displayName,'${random_string.test_id.result}')"
  shard_count   = 2
  strategy      = "round-robin"
  seed          = "ab-test-2024" # Makes distribution reproducible
}

# ============================================================================
# Test Outputs: Validate seeded round-robin distribution
# ============================================================================

output "total_users_distributed" {
  description = "Total number of users distributed across all shards (should be 6)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "group_a_count" {
  description = "Users in Group A (should be exactly 3)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "group_b_count" {
  description = "Users in Group B (should be exactly 3)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "is_balanced" {
  description = "Confirms equal split (difference should be 0 or 1)"
  value       = abs(length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) - length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]))
}

output "shard_0_first_guid" {
  description = "First GUID from shard_0 (for validation)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) > 0 ? data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"][0] : ""
}

output "shard_1_first_guid" {
  description = "First GUID from shard_1 (for validation)"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) > 0 ? data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"][0] : ""
}

# Comprehensive GUID validation for all shards
output "shard_0_all_guids" {
  description = "All GUIDs in shard_0 are valid"
  value = alltrue([
    for guid in data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"] :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}

output "shard_1_all_guids" {
  description = "All GUIDs in shard_1 are valid"
  value = alltrue([
    for guid in data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"] :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}

output "all_guids_valid" {
  description = "Comprehensive validation: ALL GUIDs across ALL shards are valid GUID format"
  value = alltrue([
    for guid in concat(
      tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]),
      tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
    ) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}
