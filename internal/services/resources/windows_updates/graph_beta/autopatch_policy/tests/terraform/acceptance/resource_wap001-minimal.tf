resource "random_string" "wap001" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_windows_updates_autopatch_policy" "wap001_minimal" {
  display_name = "acc-test-wap001-minimal-${random_string.wap001.result}"
  description  = "Acceptance test - minimal policy"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
