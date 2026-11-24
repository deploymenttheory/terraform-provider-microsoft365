#!/bin/bash

# Import using the notification message template ID.
# The import ID must be a GUID, for example 12345678-1234-1234-1234-123456789012.
terraform import microsoft365_graph_beta_device_management_device_compliance_notification_template.example "12345678-1234-1234-1234-123456789012"

# To find existing notification message template IDs, you can use Microsoft Graph Explorer:
# https://developer.microsoft.com/en-us/graph/graph-explorer
#
# Use this query to list all notification message templates:
# GET https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates
#
# Example response will include template IDs like:
# {
#   "value": [
#     {
#       "id": "12345678-1234-1234-1234-123456789012",
#       "displayName": "Basic Compliance Notification",
#       "description": "Basic notification template for device compliance violations",
#       "defaultLocale": "en-US",
#       "brandingOptions": "includeCompanyLogo"
#     }
#   ]
# }

# Note: Ensure you have appropriate permissions to read notification message templates:
# - DeviceManagementServiceConfig.Read.All (for reading)
# - DeviceManagementServiceConfig.ReadWrite.All (for importing and managing)
