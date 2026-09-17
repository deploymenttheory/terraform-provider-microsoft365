resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_linux_platform_script" "test_003" {
  name                = "acc-test-linux-platform-script-003-${random_string.test_suffix.result}"
  description         = "Maximal Linux platform script"
  script_content      = "#!/bin/bash\n# Updated script: café\nprintf \"Linux platform script updated\\n\"\n"
  execution_context   = "root"
  execution_frequency = "1day"
  execution_retries   = "3"
  role_scope_tag_ids  = ["0"]
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "5m"
  }
}
