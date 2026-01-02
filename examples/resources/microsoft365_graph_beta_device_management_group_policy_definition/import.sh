#!/bin/bash

# Import a Group Policy Definition resource using the composite ID format
# The ID format is: configurationID/definitionValueID
#
# To find the composite ID:
# 1. Navigate to the Intune portal
# 2. Go to Devices > Group Policy Configurations
# 3. Select your configuration and view the definition values
# 4. Use the configuration GUID and definition value GUID

terraform import microsoft365_graph_beta_device_management_group_policy_definition.example \
  "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6/x1y2z3a4-b5c6-d7e8-f9g0-h1i2j3k4l5m6"

# Alternative: Use the PowerShell export script to generate import commands automatically
# pwsh Export-GroupPolicyDefinitionToHCLForImport.ps1 -TenantId "<tenant-id>" -ClientId "<client-id>" -ClientSecret "<secret>"

