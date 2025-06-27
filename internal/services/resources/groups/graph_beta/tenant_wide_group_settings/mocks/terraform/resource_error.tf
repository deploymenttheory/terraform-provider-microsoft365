# Error configuration for tenant-wide group settings
resource "microsoft365_graph_beta_groups_tenant_wide_group_settings" "error" {
  # Group.Unified template ID
  template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b"

  # Setting that will trigger an error
  values = [
    {
      name  = "EnableGroupCreation"
      value = "error" # This special value will trigger an error in the mock
    }
  ]
} 