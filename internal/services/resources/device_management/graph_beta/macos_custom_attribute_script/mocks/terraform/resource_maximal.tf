resource "microsoft365_graph_beta_device_management_macos_custom_attribute_script" "maximal" {
  display_name           = "Test Maximal macOS Custom Attribute Script - Unique"
  description            = "Maximal custom attribute script for testing with all features"
  custom_attribute_type  = "string"
  file_name              = "test_maximal.sh"
  script_content         = <<-EOT
    #!/bin/bash
    
    # Comprehensive custom attribute script with all features
    set -e
    
    # Check system requirements
    if [[ $$(uname) != 'Darwin' ]]; then
        echo 'Error: This script requires macOS'
        exit 1
    fi
    
    # Gather system information for custom attribute
    OS_VERSION=$$(sw_vers -productVersion)
    HARDWARE_MODEL=$$(system_profiler SPHardwareDataType | grep 'Model Name' | awk -F': ' '{print $$2}')
    MEMORY_GB=$$(system_profiler SPHardwareDataType | grep 'Memory' | awk -F': ' '{print $$2}' | awk '{print $$1}')
    UPTIME=$$(uptime | awk '{print $$3,$$4}' | sed 's/,//')
    
    # Create comprehensive system info string
    SYSTEM_INFO="macOS $${OS_VERSION} | $${HARDWARE_MODEL} | $${MEMORY_GB} RAM | Uptime: $${UPTIME}"
    
    # Log the information
    echo "Collected system information: $$SYSTEM_INFO" >&2
    
    # Return the custom attribute value
    echo "$$SYSTEM_INFO"
    exit 0
  EOT
  run_as_account         = "user"
  role_scope_tag_ids     = ["0", "1"]

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