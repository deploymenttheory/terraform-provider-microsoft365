resource "microsoft365_graph_beta_device_management_macos_platform_script" "maximal" {
  display_name                    = "Test Maximal macOS Platform Script - Unique"
  description                     = "Maximal platform script for testing with all features"
  file_name                       = "test_maximal.sh"
  script_content                  = "#!/bin/bash\n\n# Comprehensive test script with all features\nset -e\n\n# Check system requirements\nif [[ $(uname) != 'Darwin' ]]; then\n    echo 'This script requires macOS'\n    exit 1\nfi\n\n# Log start\necho 'Starting maximal test script execution'\ndate\n\n# System information\nsystem_profiler SPSoftwareDataType\n\n# Check disk space\ndf -h\n\n# Network connectivity test\nping -c 3 apple.com\n\n# Install example package using Homebrew (if available)\nif command -v brew &> /dev/null; then\n    echo 'Homebrew is available'\n    brew --version\nelse\n    echo 'Homebrew not found'\nfi\n\n# Create test directory\nmkdir -p /tmp/macos_script_test\necho 'test content' > /tmp/macos_script_test/test.txt\n\n# Cleanup\nrm -f /tmp/macos_script_test/test.txt\nrmdir /tmp/macos_script_test\n\necho 'Maximal test script completed successfully'\nexit 0"
  run_as_account                  = "user"
  role_scope_tag_ids             = ["0", "1"]
  block_execution_notifications  = true
  execution_frequency            = "P1D"
  retry_count                    = 3

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}