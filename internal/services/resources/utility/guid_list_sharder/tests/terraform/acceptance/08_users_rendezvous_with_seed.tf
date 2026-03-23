# ==============================================================================
# Acceptance Test 08: Users - Rendezvous (HRW) Strategy (With Seed)
#
# Purpose: Verify Highest Random Weight hashing with an explicit seed distributes
# all users across shards, producing a reproducible assignment that differs from
# the no-seed variant
#
# Use Case: Multiple independent rollout streams where users should be spread
# differently across rings — e.g. user X is in ring_0 for MFA but ring_2 for
# Windows Updates. Different seeds produce independent distributions.
#
# Expected Behavior:
# - All 12 users are distributed (no user is lost)
# - 4 shards are created (shards.% = 4)
# - Distribution is probabilistic — individual shard counts are NOT asserted.
#   With 12 users and 4 shards the expected count per shard is ~3, but
#   deviation is normal and does not indicate a defect.
# - Always deterministic: same GUIDs + same seed → same assignment every run
# - Different seed from test 07 → different user-to-shard assignments
# - All output GUIDs match the GUID format
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

resource "microsoft365_graph_beta_users_user" "test_user_10" {
  display_name        = "acc-test-sharder-user-10-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-10-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser10${random_string.test_id.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

resource "microsoft365_graph_beta_users_user" "test_user_11" {
  display_name        = "acc-test-sharder-user-11-${random_string.test_id.result}"
  user_principal_name = "acc-test-sharder-user-11-${random_string.test_id.result}@deploymenttheory.com"
  mail_nickname       = "acctestsharderuser11${random_string.test_id.result}"
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
    microsoft365_graph_beta_users_user.test_user_10,
    microsoft365_graph_beta_users_user.test_user_11,
  ]
  create_duration = "30s"
}

resource "microsoft365_utility_guid_list_sharder" "test" {
  depends_on              = [time_sleep.wait_for_users]
  resource_type           = "users"
  odata_filter            = "startswith(displayName,'acc-test-sharder-user-') and endswith(displayName,'${random_string.test_id.result}')"
  shard_count             = 4
  strategy                = "rendezvous"
  recalculate_on_next_run = false
  seed                    = "deployment-ring-2024"
}

output "total_distributed" {
  description = "Total users distributed across all shards (should be 12)"
  value = (
    length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) +
    length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) +
    length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"]) +
    length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
  )
}

output "shard_0_count" {
  description = "Users in shard_0 (probabilistic — not asserted as exact value)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"])
}

output "shard_1_count" {
  description = "Users in shard_1 (probabilistic — not asserted as exact value)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "shard_2_count" {
  description = "Users in shard_2 (probabilistic — not asserted as exact value)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}

output "shard_3_count" {
  description = "Users in shard_3 (probabilistic — not asserted as exact value)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
}

output "all_guids_valid" {
  description = "All GUIDs across all shards match the GUID format"
  value = alltrue([
    for guid in concat(
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]),
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_1"]),
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_2"]),
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_3"])
    ) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}
