#!/bin/bash

# Import an existing macOS Device Compliance Policy
# Replace the ID with the actual ID of your policy from Microsoft Graph API

terraform import microsoft365_graph_beta_device_management_macos_device_compliance_policy.basic 00000000-0000-0000-0000-000000000001
