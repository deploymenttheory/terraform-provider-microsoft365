# ==============================================================================
# Test 01: Users - Round-Robin Strategy (No Seed)
#
# Purpose: Verify round-robin distribution without seed uses API order
# and distributes evenly
#
# Use Case: Quick one-time equal split where reproducibility isn't needed
#
# Expected Behavior:
# - Exactly equal shard sizes (within ±1)
# - Uses API order (may change between Terraform runs)
# - Fast processing (no shuffle overhead)
# ==============================================================================

# ==============================================================================
# Test Data Setup - Create 9 Test Users
# Note: For each loops are not supported by 
# ==============================================================================

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

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

resource "microsoft365_graph_beta_users_user" "test_user_6" {
  display_name        = "acc-test-sharder-user-6-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-6-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser6${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_7" {
  display_name        = "acc-test-sharder-user-7-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-7-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser7${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_8" {
  display_name        = "acc-test-sharder-user-8-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-8-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser8${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

# ==============================================================================
# Wait for User Provisioning (30 seconds)
# ==============================================================================

resource "time_sleep" "wait_for_users" {
  depends_on = [
    microsoft365_graph_beta_users_user.test_user_0,
    microsoft365_graph_beta_users_user.test_user_1,
    microsoft365_graph_beta_users_user.test_user_2,
    microsoft365_graph_beta_users_user.test_user_3,
    microsoft365_graph_beta_users_user.test_user_4,
    microsoft365_graph_beta_users_user.test_user_5,
    microsoft365_graph_beta_users_user.test_user_6,
    microsoft365_graph_beta_users_user.test_user_7,
    microsoft365_graph_beta_users_user.test_user_8,
  ]

  create_duration = "30s"
}

# ==============================================================================
# Test Data Source - GUID List Sharder
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "test" {
  depends_on = [time_sleep.wait_for_users]

  resource_type = "users"
  odata_filter  = "startswith(displayName,'acc-test-sharder-user-') and endswith(displayName,'${random_string.test_id.result}')"
  shard_count   = 3
  strategy      = "round-robin"
  # No seed - uses API order (non-deterministic)
  # Filter - only queries our 9 test users
}

output "shard_0_count" {
  description = "Users in shard 0"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in shard 1"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in shard 2"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "total_users_distributed" {
  description = "Total users distributed across all shards"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "sharder_id" {
  description = "Sharder datasource ID (deterministic hash)"
  value       = data.microsoft365_utility_guid_list_sharder.test.id
}

output "shard_0_first_guid" {
  description = "First GUID in shard 0 (to verify GUID format)"
  value       = tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"])[0]
}

output "shard_1_first_guid" {
  description = "First GUID in shard 1 (to verify GUID format)"
  value       = tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"])[0]
}

output "shard_2_first_guid" {
  description = "First GUID in shard 2 (to verify GUID format)"
  value       = tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])[0]
}

output "shard_0_all_guids" {
  description = "All GUIDs in shard 0 as a set"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]
}

output "shard_1_all_guids" {
  description = "All GUIDs in shard 1 as a set"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]
}

output "shard_2_all_guids" {
  description = "All GUIDs in shard 2 as a set"
  value       = data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"]
}

output "all_guids_valid" {
  description = "Check if all GUIDs match the GUID pattern"
  value = alltrue([
    for guid in tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
    ]) && alltrue([
    for guid in tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
    ]) && alltrue([
    for guid in tolist(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"]) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}