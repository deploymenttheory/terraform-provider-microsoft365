# Example with minimal configuration
resource "microsoft365_graph_beta_device_management_macos_device_compliance_policy" "minimal" {
  display_name = "macOS Minimal Compliance Policy"
  description  = "Minimal macOS device compliance policy with basic security requirements"

  # Basic security requirements
  password_required          = true
  storage_require_encryption = true
  firewall_enabled           = true

  # Scheduled actions for rules (required)
  scheduled_actions_for_rule = [
    {
      # rule_name is optional; the provider uses PasswordRequired.
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 0
        }
      ]
    }
  ]
}
