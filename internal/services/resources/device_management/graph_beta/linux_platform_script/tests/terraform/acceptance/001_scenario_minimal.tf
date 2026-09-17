resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_linux_platform_script" "test_001" {
  name                = "acc-test-linux-platform-script-001-${random_string.test_suffix.result}"
  script_content      = "#!/bin/bash\nprintf \"Linux platform script test\\n\"\n"
  execution_context   = "user"
  execution_frequency = "15minutes"
  execution_retries   = "1"
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "5m"
  }
}
