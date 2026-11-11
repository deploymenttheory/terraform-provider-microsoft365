resource "microsoft365_graph_beta_device_management_macos_platform_script" "run_once" {
  display_name   = "Run once macOS script"
  file_name      = "run_once.sh"
  run_as_account = "system"

  script_content = <<-EOT
    #!/bin/sh
    echo "Run once"
    exit 0
  EOT

  # Set PT0S explicitly when changing a recurring script back to run once.
  execution_frequency           = "PT0S"
  retry_count                   = 0
  block_execution_notifications = false
  role_scope_tag_ids            = ["0"]
}
