resource "microsoft365_graph_beta_device_management_linux_platform_script" "test_004" {
  name                = "unit-test-linux-platform-script-004"
  description         = ""
  script_content      = "#!/bin/bash\nprintf \"Linux platform script test\\n\"\n"
  execution_context   = "user"
  execution_frequency = "15minutes"
  execution_retries   = "1"
  role_scope_tag_ids  = ["0"]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
