# ==============================================================================
# Test 20: Integration - Group Member Assignment
#
# Purpose: Demonstrate how shards integrate directly with group resources
# for splitting large groups into manageable subgroups
#
# Use Case: Split a large department group into pilot subgroups
#
# Note: This is a demonstration. Actual resource creation would happen in
# acceptance tests, not unit tests
# ==============================================================================

data "microsoft365_utility_guid_list_sharder" "split_department" {
  resource_type = "group_members"
  group_id      = "12345678-1234-1234-1234-123456789abc" # Sales Department
  odata_query   = "accountEnabled eq true"
  shard_count   = 3
  strategy      = "round-robin"
  seed          = "sales-pilot-groups-2024"
}

# Example: Create pilot groups from shards
# resource "microsoft365_graph_beta_groups_group" "sales_pilot_a" {
#   display_name     = "Sales Pilot Group A"
#   mail_nickname    = "sales-pilot-a"
#   security_enabled = true
#   
#   # Shard output is a set - can be used directly
#   group_members = data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_0"]
# }

# resource "microsoft365_graph_beta_groups_group" "sales_pilot_b" {
#   display_name     = "Sales Pilot Group B"
#   mail_nickname    = "sales-pilot-b"
#   security_enabled = true
#   
#   group_members = data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_1"]
# }

# resource "microsoft365_graph_beta_groups_group" "sales_pilot_c" {
#   display_name     = "Sales Pilot Group C"
#   mail_nickname    = "sales-pilot-c"
#   security_enabled = true
#   
#   group_members = data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_2"]
# }

output "group_a_members" {
  description = "Members for Pilot Group A (ready for group_members attribute)"
  value       = data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_0"]
}

output "group_a_count" {
  description = "Number of members in Group A"
  value       = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_0"])
}

output "group_b_count" {
  description = "Number of members in Group B"
  value       = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_1"])
}

output "group_c_count" {
  description = "Number of members in Group C"
  value       = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_2"])
}

output "distribution_summary" {
  description = "Member distribution verification"
  value = {
    total_members = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_0"]) + length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_1"]) + length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_2"])
    group_a       = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_0"])
    group_b       = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_1"])
    group_c       = length(data.microsoft365_utility_guid_list_sharder.split_department.shards["shard_2"])
  }
}
