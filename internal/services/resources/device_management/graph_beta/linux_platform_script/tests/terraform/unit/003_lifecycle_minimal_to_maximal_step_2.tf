resource "microsoft365_graph_beta_device_management_linux_platform_script" "test_003" {
  name                = "unit-test-linux-platform-script-003"
  description         = "Maximal Linux platform script"
  script_content      = "#!/bin/bash\n# Updated script: café\nprintf \"Linux platform script updated\\n\"\n"
  execution_context   = "root"
  execution_frequency = "1day"
  execution_retries   = "3"
  role_scope_tag_ids  = ["0", "1"]
  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
