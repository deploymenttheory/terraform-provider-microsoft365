# Example: Retrieve all Cloud PC Frontline Service Plans

data "microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan" "all" {
  filter_type = "all"
}

# Output: List all frontline service plan IDs
output "all_frontline_service_plan_ids" {
  value = [for plan in data.microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan.all.items : plan.id]
}

# Output: Show all details for the first frontline service plan (if present)
output "first_frontline_service_plan_details" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan.all.items[0]
}

# Example: Retrieve a specific frontline service plan by ID
data "microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan" "by_id" {
  filter_type  = "id"
  filter_value = "12345678-1234-1234-1234-123456789012" # Replace with an actual ID
}

output "frontline_service_plan_by_id" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan.by_id.items[0]
}

# Example: Retrieve frontline service plans by display name substring
data "microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "2vCPU" # This will match plans containing "2vCPU" in their name
}

output "frontline_service_plans_by_display_name" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan.by_display_name.items
}

# Example: Show usage statistics for all frontline service plans
output "frontline_service_plan_usage" {
  value = [for plan in data.microsoft365_graph_beta_cloud_pc_cloud_pc_frontline_service_plan.all.items : {
    display_name  = plan.display_name
    used_count    = plan.used_count
    total_count   = plan.total_count
    usage_percent = plan.total_count > 0 ? (plan.used_count * 100 / plan.total_count) : 0
  }]
} 