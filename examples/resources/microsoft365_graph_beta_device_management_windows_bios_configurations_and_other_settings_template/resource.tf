// Example: Intune "BIOS configuration and other settings" template
//
// The configuration file is produced by the OEM's own tooling (for Dell, the
// Dell Command | Configure "multiplatform package" export) and uploaded to Intune
// base64 encoded. Use Terraform's filebase64() function to read it from disk.

resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "example" {
  display_name = "Example BIOS configuration"
  description  = "Secure Boot and TPM baseline for Dell workstations"

  file_name                  = "dell-bios-baseline.cctk"
  configuration_file_content = filebase64("${path.module}/dell-bios-baseline.cctk")

  # Optional: the service infers this from file_name when omitted.
  hardware_configuration_format = "dell"

  # Leave false so Intune generates and manages a unique BIOS password per device.
  per_device_password_disabled = false

  role_scope_tag_ids = ["0"]

  # Optional: assignments block. Only group inclusion and exclusion targets are
  # supported; included groups may additionally be narrowed by an assignment filter.
  assignments = [
    # Included group, narrowed by an assignment filter
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000000"
      filter_id   = "11111111-1111-1111-1111-111111111111"
      filter_type = "include"
    },
    # Included group, no filter
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    # Excluded group
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
  ]

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
