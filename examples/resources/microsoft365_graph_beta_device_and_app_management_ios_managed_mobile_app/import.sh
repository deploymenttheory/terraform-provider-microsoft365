#!/bin/bash
# Import using composite ID format: {managed_app_protection_id}/{app_id}
# The app ID is assigned by Intune as {bundle_id}.ios
terraform import microsoft365_graph_beta_device_and_app_management_ios_managed_mobile_app.example T_00000000-0000-0000-0000-000000000001/com.microsoft.Office.Outlook.ios
