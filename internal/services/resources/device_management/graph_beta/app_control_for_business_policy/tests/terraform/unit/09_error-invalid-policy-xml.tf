# Test 09: Error Test - Invalid Policy XML
# Purpose: Test error handling with invalid XML policy content

# App Control Policy with Invalid XML - Should Fail
resource "microsoft365_graph_beta_device_management_app_control_for_business_policy" "error_invalid_xml" {
  name        = "unit-test-app-control-policy-error-invalid-xml"
  description = "Error test - invalid XML policy"

  # Invalid XML - missing closing tags
  policy_xml = <<-EOT
<?xml version="1.0" encoding="utf-8"?>
<SiPolicy xmlns="urn:schemas-microsoft-com:sipolicy" PolicyType="Base Policy">
  <VersionEx>1.0.3.0</VersionEx>
  <PolicyID>{264C0644-19BE-418F-BAED-29E5E36250AD}
  <BasePolicyID>{264C0644-19BE-418F-BAED-29E5E36250AD}</BasePolicyID>
  <!-- Missing closing tags intentionally -->
  EOT

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
