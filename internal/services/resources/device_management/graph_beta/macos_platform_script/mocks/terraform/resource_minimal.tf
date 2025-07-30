resource "microsoft365_graph_beta_device_management_macos_platform_script" "minimal" {
  display_name   = "Minimal macOS Script"
  script_content = "#!/bin/bash\necho 'Minimal Script'"
  run_as_account = "system"
  file_name      = "minimal-script.sh"
} 