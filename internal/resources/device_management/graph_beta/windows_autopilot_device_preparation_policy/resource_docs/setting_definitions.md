# Windows Autopilot Device Preparation Policy - Setting Definitions

This document outlines the setting definitions for Windows Autopilot Device Preparation Policy.

## Template ID

The correct template ID for Windows Autopilot Device Preparation Policy is:

```
80d33118-b7b4-40d8-b15f-81be745e053f_1
```

## Setting Definitions

| Setting Name | Setting Definition ID | Type | Valid Values |
|-------------|----------------------|------|--------------|
| Deployment Mode | `enrollment_autopilot_dpp_deploymentmode` | Choice | `enrollment_autopilot_dpp_deploymentmode_0` (Standard mode)<br>`enrollment_autopilot_dpp_deploymentmode_1` (Enhanced mode) |
| Deployment Type | `enrollment_autopilot_dpp_deploymenttype` | Choice | `enrollment_autopilot_dpp_deploymenttype_0` (User-driven)<br>`enrollment_autopilot_dpp_deploymenttype_1` (Self-deploying) |
| Join Type | `enrollment_autopilot_dpp_jointype` | Choice | `enrollment_autopilot_dpp_jointype_0` (Azure AD joined)<br>`enrollment_autopilot_dpp_jointype_1` (Azure AD hybrid joined) |
| Account Type | `enrollment_autopilot_dpp_accountype` | Choice | `enrollment_autopilot_dpp_accountype_0` (Standard User)<br>`enrollment_autopilot_dpp_accountype_1` (Administrator) |
| Device Security Group | `enrollment_autopilot_dpp_devicegroup` | String | ID of the device security group |
| Timeout (Minutes) | `enrollment_autopilot_dpp_timeout` | Integer | Valid range: 15-720 |
| Custom Error Message | `enrollment_autopilot_dpp_custonerror` | String | Any text, max 1000 chars |
| Allow Skip | `enrollment_autopilot_dpp_allowskip` | Boolean | `true` or `false` |
| Allow Diagnostics | `enrollment_autopilot_dpp_allowdiagnostics` | Boolean | `true` or `false` |
| Allowed Apps | `enrollment_autopilot_dpp_allowedapps` | Collection | List of app IDs |
| Allowed Scripts | `enrollment_autopilot_dpp_allowedscripts` | Collection | List of script IDs |

## Mapping to Terraform Resource Schema

The following table shows how the API setting definition IDs map to the Terraform resource schema:

| API Setting Definition ID | Terraform Schema Path |
|--------------------------|----------------------|
| `enrollment_autopilot_dpp_deploymentmode` | `deployment_settings.deployment_mode` |
| `enrollment_autopilot_dpp_deploymenttype` | `deployment_settings.deployment_type` |
| `enrollment_autopilot_dpp_jointype` | `deployment_settings.join_type` |
| `enrollment_autopilot_dpp_accountype` | `deployment_settings.account_type` |
| `enrollment_autopilot_dpp_devicegroup` | `device_security_group` |
| `enrollment_autopilot_dpp_timeout` | `oobe_settings.timeout_in_minutes` |
| `enrollment_autopilot_dpp_custonerror` | `oobe_settings.custom_error_message` |
| `enrollment_autopilot_dpp_allowskip` | `oobe_settings.allow_skip` |
| `enrollment_autopilot_dpp_allowdiagnostics` | `oobe_settings.allow_diagnostics` |
| `enrollment_autopilot_dpp_allowedapps` | `allowed_apps` |
| `enrollment_autopilot_dpp_allowedscripts` | `allowed_scripts` |

## Notes on Implementation

The mapping between our internal `deviceConfiguration--windows10AutopilotDevicePreparation_*` identifiers and the actual Graph API `enrollment_autopilot_dpp_*` identifiers needs to be fixed in the code.

According to the CloudFlow blog post, the API request structure for settings differs from our current implementation. The API expects settings in a format like:

```json
{
  "@odata.type": "#microsoft.graph.deviceManagementConfigurationSetting",
  "settingInstance": {
    "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
    "settingDefinitionId": "enrollment_autopilot_dpp_deploymentmode",
    "choiceSettingValue": {
      "@odata.type": "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue",
      "value": "enrollment_autopilot_dpp_deploymentmode_0",
      "children": []
    }
  }
}
```
