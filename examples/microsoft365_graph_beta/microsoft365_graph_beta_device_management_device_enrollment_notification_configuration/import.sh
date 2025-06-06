#!/bin/bash

# Import examples for device enrollment notification configuration resource

# Example 1: Import an existing device enrollment notification configuration by ID
terraform import microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.basic "12345678-1234-1234-1234-123456789abc"

# Example 2: Import with state check
terraform import microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.advanced "87654321-4321-4321-4321-cba987654321"
terraform show microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.advanced

# Example 3: Import and plan to see differences
terraform import microsoft365_graph_beta_device_management_device_enrollment_notification_configuration.ios_push "11111111-2222-3333-4444-555555555555"
terraform plan 