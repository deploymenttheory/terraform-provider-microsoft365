# Test 01: Get compliance changes for an update policy
# This test creates a full dependency chain to test compliance changes retrieval

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create deployment audience
resource "microsoft365_graph_beta_windows_updates_autopatch_deployment_audience" "test" {}

# Wait for audience to propagate
resource "time_sleep" "wait_for_audience" {
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test]
  create_duration = "10s"
}

# Create update policy
resource "microsoft365_graph_beta_windows_updates_update_policy" "test" {
  audience_id = microsoft365_graph_beta_windows_updates_autopatch_deployment_audience.test.id

  compliance_changes = true

  compliance_change_rules = [
    {
      content_filter = {
        filter_type = "driverUpdateFilter"
      }
      duration_before_deployment_start = "P7D"
    }
  ]

  depends_on = [time_sleep.wait_for_audience]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Wait for policy to propagate
resource "time_sleep" "wait_for_policy" {
  depends_on      = [microsoft365_graph_beta_windows_updates_update_policy.test]
  create_duration = "10s"
}

# Get a quality update from catalog
data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"

  timeouts = {
    read = "30s"
  }
}

# Create a content approval
resource "microsoft365_graph_beta_windows_updates_autopatch_content_approval" "test" {
  update_policy_id   = microsoft365_graph_beta_windows_updates_update_policy.test.id
  catalog_entry_id   = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality.entries[0].id
  catalog_entry_type = "qualityUpdate"
  is_revoked         = false

  depends_on = [time_sleep.wait_for_policy]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Wait for content approval to propagate
resource "time_sleep" "wait_for_approval" {
  depends_on      = [microsoft365_graph_beta_windows_updates_autopatch_content_approval.test]
  create_duration = "10s"
}

# Query compliance changes
data "microsoft365_graph_beta_windows_updates_compliance_changes" "test" {
  update_policy_id = microsoft365_graph_beta_windows_updates_update_policy.test.id

  depends_on = [time_sleep.wait_for_approval]

  timeouts = {
    read = "30s"
  }
}

output "compliance_changes_count" {
  value = length(data.microsoft365_graph_beta_windows_updates_compliance_changes.test.compliance_changes)
}
