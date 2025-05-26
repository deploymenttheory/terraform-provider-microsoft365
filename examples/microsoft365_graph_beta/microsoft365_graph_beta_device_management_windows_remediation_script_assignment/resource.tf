resource "microsoft365_graph_beta_device_management_windows_remediation_script_assignment" "daily_example" {
  device_health_remediation_id = "00000000-0000-0000-0000-000000000001"

  target {
    target_type = "groupAssignment"
    group_id    = "00000000-0000-0000-0000-000000000002"
  }

  run_remediation_script = true

  run_schedule {
    daily {
      # Number of days between runs. Default is 1, valid 1–30.
      interval = 2

      # Time of day, in HH:MM:SS.
      time = "14:30:00"

      # Whether to use UTC (true) or local device time (false).
      use_utc = true
    }
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "graph_beta_device_management_windows_remediation_script_assignment" "hourly_example" {
  device_health_remediation_id = "00000000-0000-0000-0000-000000000003"


  target {
    target_type = "allDevices"
  }

  # Only the interval + UTC flag apply here
  run_schedule {
    hourly {
      # Hours between runs. Default is 1, valid 1–23.
      interval = 3

      use_utc = false
    }
  }
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "graph_beta_device_management_windows_remediation_script_assignment" "once_example" {
  device_health_remediation_id = "00000000-0000-0000-0000-000000000004"

  target {
    target_type   = "configurationManagerCollection"
    collection_id = "MEMABCDEF01"
  }

  run_schedule {
    once {
      # ISO-8601 timestamp for the one-time run
      start_date_time = "2025-06-01T10:00:00Z"

      use_utc = true
    }
  }
}
