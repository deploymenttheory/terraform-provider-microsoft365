# ==============================================================================
# Acceptance Test 06: Users - Size Strategy (With Seed)
#
# Purpose: Verify absolute size-based distribution with seed creates exact-count
# shards AND reproducible membership across Terraform runs
#
# Use Case: "Same pilot group size every time, but use seed to vary which users
# are in each ring across different rollouts (e.g. MFA vs Windows Updates)"
#
# Expected Behavior:
# - shard_0: exactly 3 users (absolute size — seed only changes which users)
# - shard_1: exactly 4 users (absolute size — seed only changes which users)
# - shard_2: exactly 2 users (remaining — -1 sentinel, 9 - 3 - 4)
# - Total: 9
# - Same seed + same input set = same user-to-shard assignment on every run
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

resource "microsoft365_utility_guid_list_sharder" "test" {
  depends_on              = [time_sleep.wait_for_users]
  resource_type           = "users"
  odata_filter            = "startswith(displayName,'acc-test-sharder-user-') and endswith(displayName,'${random_string.test_id.result}')"
  shard_sizes             = [3, 4, -1]
  strategy                = "size"
  recalculate_on_next_run = false
  seed                    = "mfa-rollout-2024"
}

output "total_distributed" {
  description = "Total users distributed across all shards (should be 9)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_0_count" {
  description = "Users in pilot shard — should be exactly 3 (absolute size)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in broader pilot shard — should be exactly 4 (absolute size)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in full rollout shard — should be exactly 2 (9 - 3 - 4, via -1 sentinel)"
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
