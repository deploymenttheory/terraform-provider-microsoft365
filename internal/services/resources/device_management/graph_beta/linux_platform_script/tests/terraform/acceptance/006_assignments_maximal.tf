resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# All assignment targets are newly created empty groups.
resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_1" {
  display_name     = "acc-test-linux-platform-script-006-include-${random_string.test_suffix.result}"
  mail_nickname    = "acc-linux-006-include-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_2" {
  display_name     = "acc-test-linux-platform-script-006-second-${random_string.test_suffix.result}"
  mail_nickname    = "acc-linux-006-second-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_006_3" {
  display_name     = "acc-test-linux-platform-script-006-exclude-${random_string.test_suffix.result}"
  mail_nickname    = "acc-linux-006-exclude-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_linux_platform_script" "test_006" {
  name                = "acc-test-linux-platform-script-006-${random_string.test_suffix.result}"
  script_content      = "#!/bin/bash\nprintf \"Linux platform script test\\n\"\n"
  execution_context   = "user"
  execution_frequency = "15minutes"
  execution_retries   = "1"
  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_006_3.id
    },
  ]
  timeouts = {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "5m"
  }
}
