resource "random_id" "test_002" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_002" {
  display_name = "acc-test-gpd-bool-max-${random_id.test_002.hex}"
  description  = "Acceptance test for boolean maximal"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

resource "microsoft365_graph_beta_device_management_group_policy_definition" "test_002" {
  group_policy_configuration_id = microsoft365_graph_beta_device_management_group_policy_configuration.test_002.id
  policy_name                   = "Remove Default Microsoft Store packages from the system."
  class_type                    = "machine"
  category_path                 = "\\Windows Components\\App Package Deployment"
  enabled                       = true

  values = [
    { label = "Feedback Hub", value = "true" },
    { label = "Microsoft 365 Copilot", value = "false" },
    { label = "Microsoft Clipchamp", value = "true" },
    { label = "Microsoft Copilot", value = "false" },
    { label = "Microsoft News", value = "true" },
    { label = "Microsoft Photos **", value = "false" },
    { label = "Microsoft Solitaire Collection", value = "true" },
    { label = "Microsoft Sticky Notes", value = "false" },
    { label = "Microsoft Teams", value = "true" },
    { label = "Microsoft To Do", value = "false" },
    { label = "MSN Weather", value = "true" },
    { label = "Outlook for Windows", value = "false" },
    { label = "Paint", value = "true" },
    { label = "Quick Assist", value = "false" },
    { label = "Snipping Tool", value = "true" },
    { label = "Windows Calculator", value = "false" },
    { label = "Windows Camera **", value = "true" },
    { label = "Windows Media Player **", value = "false" },
    { label = "Windows Notepad **", value = "true" },
    { label = "Windows Sound Recorder", value = "false" },
    { label = "Windows Terminal", value = "true" },
    { label = "Xbox Gaming App", value = "false" },
    { label = "Xbox Identity Provider *", value = "true" },
    { label = "Xbox Speech To Text Overlay *", value = "false" },
    { label = "Xbox TCUI *", value = "true" }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
