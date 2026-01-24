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

################################################################################
# Test Data Setup - Create 100 Test Users
################################################################################

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

locals {
  test_users = { for i in range(100) : "user_${i}" => {
    display_name        = "acc-test-sharder-user-${i}-${random_string.test_id.result}"
    user_principal_name = "acc-test-sharder-user-${i}-${random_string.test_id.result}@deploymenttheory.com"
    mail_nickname       = "acctestsharderuser${i}${random_string.test_id.result}"
  } }
}

resource "microsoft365_graph_beta_users_user" "test_users" {
  for_each = local.test_users

  display_name        = each.value.display_name
  user_principal_name = each.value.user_principal_name
  mail_nickname       = each.value.mail_nickname
  account_enabled     = true
  hard_delete         = true

  password_profile = {
    force_change_password_next_sign_in = false
    password                           = "TempPass123!"
  }
}

################################################################################
# Wait for User Provisioning (60 seconds)
################################################################################

resource "time_sleep" "wait_for_users" {
  depends_on = [microsoft365_graph_beta_users_user.test_users]

  create_duration = "60s"
}

################################################################################
# Test Data Source - GUID List Sharder
################################################################################

data "microsoft365_utility_guid_list_sharder" "test" {
  depends_on = [time_sleep.wait_for_users]

  resource_type = "users"
  odata_query   = "$filter=accountEnabled eq true"
  shard_count   = 3
  strategy      = "round-robin"
  # No seed - uses API order (non-deterministic)
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

output "total_users" {
  description = "Total users distributed"
  value       = length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.test.shards["shard_2"])
}
