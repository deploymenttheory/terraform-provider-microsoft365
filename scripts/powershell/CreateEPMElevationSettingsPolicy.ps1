[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID) where the application is registered")]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true,
    HelpMessage="Specify the application (client) ID of the Entra ID app registration")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the client secret of the Entra ID app registration")]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret
)

# Script Setup
Import-Module Microsoft.Graph.Authentication
Import-Module Microsoft.Graph.Beta.DeviceManagement

# Connect to Microsoft Graph
$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Create policy parameters
$params = @{
    name = "test - epm policy"
    description = ""
    platforms = "windows10"
    technologies = "mdm,endpointPrivilegeManagement"
    roleScopeTagIds = @(
        "0"
    )
    settings = @(
        @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "device_vendor_msft_policy_elevationclientsettings_enableepm"
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "device_vendor_msft_policy_elevationclientsettings_enableepm_1"
                    children = @(
                        @{
                            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                            settingDefinitionId = "device_vendor_msft_policy_elevationclientsettings_defaultelevationresponse"
                            choiceSettingValue = @{
                                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                                value = "device_vendor_msft_policy_elevationclientsettings_defaultelevationresponse_1"
                                children = @(
                                    @{
                                        "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingCollectionInstance"
                                        settingDefinitionId = "device_vendor_msft_policy_privilegemanagement_elevationclientsettings_defaultelevationresponse_validation"
                                        choiceSettingCollectionValue = @(
                                            @{
                                                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                                                value = "device_vendor_msft_policy_privilegemanagement_elevationclientsettings_defaultelevationresponse_validation_0"
                                                children = @()
                                            }
                                        )
                                    }
                                )
                            }
                        }
                        @{
                            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                            settingDefinitionId = "device_vendor_msft_policy_elevationclientsettings_allowelevationdetection"
                            choiceSettingValue = @{
                                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                                value = "device_vendor_msft_policy_elevationclientsettings_allowelevationdetection_1"
                                children = @()
                            }
                        }
                        @{
                            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                            settingDefinitionId = "device_vendor_msft_policy_elevationclientsettings_senddata"
                            choiceSettingValue = @{
                                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                                value = "device_vendor_msft_policy_elevationclientsettings_senddata_1"
                                children = @(
                                    @{
                                        "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                                        settingDefinitionId = "device_vendor_msft_policy_elevationclientsettings_reportingscope"
                                        choiceSettingValue = @{
                                            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                                            value = "device_vendor_msft_policy_elevationclientsettings_reportingscope_2"
                                            children = @()
                                        }
                                    }
                                )
                            }
                        }
                    )
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "a13cc55c-307a-4962-aaec-20b832bf75c7"
                    }
                }
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "58a79a4b-ba9b-4923-a7a5-6dc1a9f638a4"
                }
            }
        }
    )
    templateReference = @{
        templateId = "e7dcaba4-959b-46ed-88f0-16ba39b14fd8_1"
    }
}

try {
    Write-Host "Creating EPM configuration policy: $PolicyName"
    $newPolicy = New-MgBetaDeviceManagementConfigurationPolicy -BodyParameter $params
    Write-Host "Successfully created policy with ID: $($newPolicy.Id)"
}
catch {
    Write-Error "Error creating configuration policy: $_"
}

# Disconnect from Microsoft Graph
Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."