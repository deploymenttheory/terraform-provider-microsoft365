# ==============================================================================
# Acceptance Test 04: Users - Percentage Strategy (With Seed)
#
# Purpose: Verify percentage-based distribution with seed produces custom-sized
# shards AND reproducible membership across Terraform runs
#
# Use Case: Structured phased rollout (10% → 30% → 60%) where the SAME users
# must appear in the same phase on every run — e.g. a support team needs to
# know which pilot users to expect
#
# Expected Behavior:
# - shard_0: 1 user  (10% of 10)
# - shard_1: 3 users (30% of 10)
# - shard_2: 6 users (60% of 10 — last shard absorbs all remaining)
# - Shard sizes are identical to no-seed variant; seed only affects WHO is assigned
# - Same seed + same input set = same assignment on every run
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

resource "microsoft365_graph_beta_users_user" "test_user_9" {
  display_name        = "acc-test-sharder-user-9-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-9-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser9${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

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
    microsoft365_graph_beta_users_user.test_user_9,
  ]
  create_duration = "30s"
}

resource "microsoft365_utility_guid_list_sharder" "test" {
  depends_on              = [time_sleep.wait_for_users]
  resource_type           = "users"
  odata_filter            = "startswith(displayName,'acc-test-sharder-user-') and endswith(displayName,'${random_string.test_id.result}')"
  shard_percentages       = [10, 30, 60]
  strategy                = "percentage"
  recalculate_on_next_run = false
  seed                    = "mfa-phased-2024"
}

output "total_distributed" {
  description = "Total users distributed across all shards (should be 10)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_0_count" {
  description = "Users in pilot shard — should be exactly 1 (10% of 10)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in broader pilot shard — should be exactly 3 (30% of 10)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in full rollout shard — should be exactly 6 (remaining 60% of 10)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_0_first_guid" {
  description = "First GUID in shard_0 (validates GUID format)"
  value       = tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_0"])[0]
}

output "shard_1_first_guid" {
  description = "First GUID in shard_1 (validates GUID format)"
  value       = tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])[0]
}

output "shard_2_first_guid" {
  description = "First GUID in shard_2 (validates GUID format)"
  value       = tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])[0]
}

output "all_guids_valid" {
  description = "All GUIDs across all shards match the GUID format"
  value = alltrue([
    for guid in concat(
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]),
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]),
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
    ) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}
