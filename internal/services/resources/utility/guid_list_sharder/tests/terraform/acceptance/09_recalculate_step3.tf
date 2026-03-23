# ==============================================================================
# Acceptance Test 09 — Step 3 of 3: recalculate_on_next_run = true (unlock)
#
# Purpose: Prove that switching recalculate_on_next_run to true causes the
# Update path to re-query the Graph API and reshard from the current tenant
# membership. The 3 users added in Step 2 must now appear in the shard state.
#
# 9 test users exist in the tenant (unchanged from Step 2).
# Only the recalculate_on_next_run flag changes — all other config is identical.
#
# Expected: total_distributed = 9 (all current tenant members resharded)
#           recalculate_flag  = true
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
  shard_count             = 2
  strategy                = "round-robin"
  seed                    = "lifecycle-test-seed"
  recalculate_on_next_run = true
}

output "total_distributed" {
  description = "Total users in state — must be 9 (reshard picked up all current tenant members)"
  value       = length(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
}

output "recalculate_flag" {
  description = "Current value of recalculate_on_next_run"
  value       = microsoft365_utility_guid_list_sharder.test.recalculate_on_next_run
}

output "all_guids_valid" {
  description = "All GUIDs match the GUID format after reshard"
  value = alltrue([
    for guid in concat(
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_0"]),
      tolist(microsoft365_utility_guid_list_sharder.test.shards["shard_1"])
    ) :
    can(regex("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$", guid))
  ])
}
