#!/bin/bash

# Import examples for Microsoft 365 Graph Beta Device Management Notification Message Template

# Basic import using template ID
# Replace "12345678-1234-1234-1234-123456789012" with the actual notification message template ID from Microsoft Graph
terraform import microsoft365_graph_beta_device_management_notification_message_template.basic "12345678-1234-1234-1234-123456789012"

# Import multilingual template
terraform import microsoft365_graph_beta_device_management_notification_message_template.multilingual "87654321-4321-4321-4321-210987654321"

# Import advanced template
terraform import microsoft365_graph_beta_device_management_notification_message_template.advanced "11111111-2222-3333-4444-555555555555"

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