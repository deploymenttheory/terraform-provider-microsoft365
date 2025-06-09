# Windows Autopilot Device Preparation Policy Creator
# Based on CloudFlow blog: https://cloudflow.be/windows-autopilot-device-perpetration-with-graph-api/
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
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Display name for the WADP policy")]
    [ValidateNotNullOrEmpty()]
    [string]$PolicyDisplayName,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Description for the WADP policy")]
    [string]$PolicyDescription = "",
    
    [Parameter(Mandatory=$true,
    HelpMessage="Device security group ID for just-in-time assignment (must have Intune Provisioning Client as owner)")]
    [ValidateNotNullOrEmpty()]
    [string]$DeviceSecurityGroupId,
    
    [Parameter(Mandatory=$true,
    HelpMessage="User security group ID for policy assignment")]
    [ValidateNotNullOrEmpty()]
    [string]$UserSecurityGroupId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Deployment mode (Standard=0, Enhanced=1)")]
    [ValidateSet("0", "1")]
    [string]$DeploymentMode = "0",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Deployment type (User-driven=0, Self-deploying=1)")]
    [ValidateSet("0", "1")]
    [string]$DeploymentType = "0",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Join type (Entra ID joined=0, Entra ID hybrid joined=1)")]
    [ValidateSet("0", "1")]
    [string]$JoinType = "0",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Account type (Standard User=0, Administrator=1)")]
    [ValidateSet("0", "1")]
    [string]$AccountType = "0",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Timeout in minutes (15-720)")]
    [ValidateRange(15, 720)]
    [int]$TimeoutInMinutes = 60,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Custom error message")]
    [string]$CustomErrorMessage = "Contact your organization's support person for help.",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Allow skip after failed attempts")]
    [bool]$AllowSkip = $false,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Allow diagnostics access")]
    [bool]$AllowDiagnostics = $false,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Array of allowed app objects with id and type properties")]
    [array]$AllowedApps = @(),
    
    [Parameter(Mandatory=$false,
    HelpMessage="Array of allowed script IDs (GUIDs)")]
    [array]$AllowedScriptIds = @()
)

# Helper function to validate GUID format
function Test-IsValidGuid {
    param (
        [Parameter(Mandatory=$true)]
        [string]$InputGuid
    )
    
    $guidRegex = '^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$'
    return $InputGuid -match $guidRegex
}

# Function to get data from Microsoft Graph API
function Get-GraphData {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$Url
    )
    
    $authHeader = @{
        'Authorization' = "Bearer $GraphToken"
        'Content-Type' = 'application/json'
    }
    
    $retryCount = 0
    $maxRetries = 3
    $results = @()
    
    while ($retryCount -le $maxRetries) {
        try {
            do {
                $response = Invoke-RestMethod -Uri $Url -Method Get -Headers $authHeader
                
                if ($response.'@odata.nextLink') {
                    $Url = $response.'@odata.nextLink'
                    $results += $response
                } else {
                    $results += $response
                    return $results
                }
            } while ($response.'@odata.nextLink')
        }
        catch {
            $statusCode = $_.Exception.Response.StatusCode.value__
            if ($statusCode -eq 429 -or $statusCode -eq 503) {
                $retryCount++
                $retryAfter = if ($_.Exception.Response.Headers.'Retry-After') { 
                    [int]($_.Exception.Response.Headers.'Retry-After') 
                } else { 
                    $retryCount * 10 
                }
                Write-Host "⚠️ Rate limited. Retrying in $retryAfter seconds..." -ForegroundColor Yellow
                Start-Sleep -Seconds $retryAfter
            } else {
                Write-Error "API call failed: $_"
                throw
            }
        }
    }
    
    throw "Max retry attempts reached for GET request"
}

# Function to post data to Microsoft Graph API
function Post-GraphData {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$Url,
        
        [Parameter(Mandatory=$true)]
        [string]$Body
    )
    
    $authHeader = @{
        'Authorization' = "Bearer $GraphToken"
        'Content-Type' = 'application/json'
    }
    
    $retryCount = 0
    $maxRetries = 3
    
    while ($retryCount -le $maxRetries) {
        try {
            $response = Invoke-RestMethod -Uri $Url -Method POST -Headers $authHeader -Body $Body
            return $response
        }
        catch {
            $statusCode = $_.Exception.Response.StatusCode.value__
            if ($statusCode -eq 429 -or $statusCode -eq 503) {
                $retryCount++
                $retryAfter = if ($_.Exception.Response.Headers.'Retry-After') { 
                    [int]($_.Exception.Response.Headers.'Retry-After') 
                } else { 
                    $retryCount * 10 
                }
                Write-Host "⚠️ Rate limited. Retrying in $retryAfter seconds..." -ForegroundColor Yellow
                Start-Sleep -Seconds $retryAfter
            } else {
                Write-Error "API call failed: $_"
                Write-Error "Response: $($_.Exception.Response | Out-String)"
                throw
            }
        }
    }
    
    throw "Max retry attempts reached for POST request"
}

# Function to get Graph API access token
function Get-GraphAPIAccessToken {
    param (
        [Parameter(Mandatory=$true)]
        [string]$TenantId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientSecret
    )
    
    try {
        $tokenUrl = "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token"
        
        $body = @{
            client_id     = $ClientId
            client_secret = $ClientSecret
            scope         = "https://graph.microsoft.com/.default"
            grant_type    = "client_credentials"
        }
        
        $response = Invoke-RestMethod -Uri $tokenUrl -Method POST -Body $body -ContentType "application/x-www-form-urlencoded"
        return $response.access_token
    }
    catch {
        Write-Error "Failed to get access token: $_"
        throw
    }
}

# Function to validate security group ownership
function Test-SecurityGroupOwnership {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$GroupId
    )
    
    try {
        Write-Host "🔍 Validating security group ownership..." -ForegroundColor Yellow
        
        $ownersUrl = "https://graph.microsoft.com/beta/groups/$GroupId/owners"
        $owners = Get-GraphData -GraphToken $GraphToken -Url $ownersUrl
        
        $intuneProvisioningClientAppId = "f1346770-5b25-470b-88bd-d5744ab7952c"
        $hasRequiredOwner = $false
        
        if ($owners.value) {
            foreach ($owner in $owners.value) {
                if ($owner.'@odata.type' -eq "#microsoft.graph.servicePrincipal" -and $owner.appId -eq $intuneProvisioningClientAppId) {
                    $hasRequiredOwner = $true
                    break
                }
            }
        }
        
        if (-not $hasRequiredOwner) {
            throw "Security group $GroupId must have the 'Intune Provisioning Client' service principal (AppId: $intuneProvisioningClientAppId) as an owner"
        }
        
        Write-Host "✅ Security group ownership validated" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Error "Security group validation failed: $_"
        throw
    }
}

# Function to create WADP configuration policy
function New-WADPConfigurationPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [hashtable]$PolicySettings
    )
    
    try {
        Write-Host "🔄 Creating Windows Autopilot Device Preparation Policy..." -ForegroundColor Yellow
        
        # Build settings array
        $settings = @()
        
        # Deployment Mode
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_deploymentmode"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "5180aeab-886e-4589-97d4-40855c646315"
                }
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "enrollment_autopilot_dpp_deploymentmode_$($PolicySettings.DeploymentMode)"
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "5874c2f6-bcf1-463b-a9eb-bee64e2f2d82"
                    }
                    children = @()
                }
            }
        }
        
        # Deployment Type
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_deploymenttype"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "f4184296-fa9f-4b67-8b12-1723b3f8456b"
                }
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "enrollment_autopilot_dpp_deploymenttype_$($PolicySettings.DeploymentType)"
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "e0af022f-37f3-4a40-916d-1ab7281c88d9"
                    }
                    children = @()
                }
            }
        }
        
        # Join Type
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_jointype"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "6310e95d-6cfa-4d2f-aae0-1e7af12e2182"
                }
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "enrollment_autopilot_dpp_jointype_$($PolicySettings.JoinType)"
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "1fa84eb3-fcfa-4ed6-9687-0f3d486402c4"
                    }
                    children = @()
                }
            }
        }
        
        # Account Type
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_accountype"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "d4f2a840-86d5-4162-9a08-fa8cc608b94e"
                }
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "enrollment_autopilot_dpp_accountype_$($PolicySettings.AccountType)"
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "bf13bb47-69ef-4e06-97c1-50c2859a49c2"
                    }
                    children = @()
                }
            }
        }
        
        # Timeout
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_timeout"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "6dec0657-dfb8-4906-a7ee-3ac6ee1edecb"
                }
                simpleSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationIntegerSettingValue"
                    value = $PolicySettings.TimeoutInMinutes
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "0bbcce5b-a55a-4e05-821a-94bf576d6cc8"
                    }
                }
            }
        }
        
        # Custom Error Message
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_customerrormessage"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "2ddf0619-2b7a-46de-b29b-c6191e9dda6e"
                }
                simpleSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    value = $PolicySettings.CustomErrorMessage
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "fe5002d5-fbe9-4920-9e2d-26bfc4b4cc97"
                    }
                }
            }
        }
        
        # Allow Skip
        $allowSkipValue = if ($PolicySettings.AllowSkip) { "1" } else { "0" }
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_allowskip"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "2a71dc89-0f17-4ba9-bb27-af2521d34710"
                }
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "enrollment_autopilot_dpp_allowskip_$allowSkipValue"
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "a2323e5e-ac56-4517-8847-b0a6fdb467e7"
                    }
                    children = @()
                }
            }
        }
        
        # Allow Diagnostics
        $allowDiagnosticsValue = if ($PolicySettings.AllowDiagnostics) { "1" } else { "0" }
        $settings += @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
            settingInstance = @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
                settingDefinitionId = "enrollment_autopilot_dpp_allowdiagnostics"
                settingInstanceTemplateReference = @{
                    settingInstanceTemplateId = "e2b7a81b-f243-4abd-bce3-c1856345f405"
                }
                choiceSettingValue = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingValue"
                    value = "enrollment_autopilot_dpp_allowdiagnostics_$allowDiagnosticsValue"
                    settingValueTemplateReference = @{
                        settingValueTemplateId = "c59d26fd-3460-4b26-b47a-f7e202e7d5a3"
                    }
                    children = @()
                }
            }
        }
        
        # Allowed Apps
        if ($PolicySettings.AllowedApps -and $PolicySettings.AllowedApps.Count -gt 0) {
            $appValues = @()
            foreach ($app in $PolicySettings.AllowedApps) {
                $appJson = "{`"id`":`"$($app.id)`",`"type`":`"#microsoft.graph.$($app.type)`"}"
                $appValues += @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    value = $appJson
                }
            }
            
            $settings += @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
                settingInstance = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
                    settingDefinitionId = "enrollment_autopilot_dpp_allowedappids"
                    settingInstanceTemplateReference = @{
                        settingInstanceTemplateId = "70d22a8a-a03c-4f62-b8df-dded3e327639"
                    }
                    simpleSettingCollectionValue = $appValues
                }
            }
        }
        
        # Allowed Scripts
        if ($PolicySettings.AllowedScriptIds -and $PolicySettings.AllowedScriptIds.Count -gt 0) {
            $scriptValues = @()
            foreach ($scriptId in $PolicySettings.AllowedScriptIds) {
                $scriptValues += @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
                    value = $scriptId
                }
            }
            
            $settings += @{
                "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSetting"
                settingInstance = @{
                    "@odata.type" = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
                    settingDefinitionId = "enrollment_autopilot_dpp_allowedscriptids"
                    settingInstanceTemplateReference = @{
                        settingInstanceTemplateId = "1bc67702-800c-4271-8fd9-609351cc19cf"
                    }
                    simpleSettingCollectionValue = $scriptValues
                }
            }
        }
        
        # Create the policy
        $policyBody = @{
            "@odata.type" = "#microsoft.graph.deviceManagementConfigurationPolicy"
            name = $PolicySettings.DisplayName
            description = $PolicySettings.Description
            settings = $settings
            roleScopeTagIds = @("0")
            platforms = "windows10"
            technologies = "enrollment"
            templateReference = @{
                templateId = "80d33118-b7b4-40d8-b15f-81be745e053f_1"
                templateFamily = "enrollmentConfiguration"
            }
        } | ConvertTo-Json -Depth 20
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
        $response = Post-GraphData -GraphToken $GraphToken -Url $url -Body $policyBody
        
        Write-Host "✅ Policy created successfully (ID: $($response.id))" -ForegroundColor Green
        return $response
    }
    catch {
        Write-Error "Failed to create WADP policy: $_"
        throw
    }
}

# Function to assign device group using alternate method
function Set-WADPDeviceGroupAssignment {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyId,
        
        [Parameter(Mandatory=$true)]
        [string]$DeviceSecurityGroupId
    )
    
    try {
        Write-Host "🔄 Attempting device group assignment with alternate method..." -ForegroundColor Yellow
        
        # This is a more generic assignment approach that might work when the enrollment time
        # membership target method fails due to internal server errors
        $assignmentBody = @{
            assignments = @(
                @{
                    id = ""
                    source = "direct"
                    target = @{
                        groupId = $DeviceSecurityGroupId
                        "@odata.type" = "#microsoft.graph.groupAssignmentTarget"
                        deviceAndAppManagementAssignmentFilterType = "none"
                    }
                }
            )
        } | ConvertTo-Json -Depth 5
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$PolicyId')/assign"
        Post-GraphData -GraphToken $GraphToken -Url $url -Body $assignmentBody | Out-Null
        
        Write-Host "✅ Device group assigned using alternate method" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Error "Failed to assign device group using alternate method: $_"
        return $false
    }
}

# Function to assign just-in-time configuration
function Set-WADPJustInTimeConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyId,
        
        [Parameter(Mandatory=$true)]
        [string]$DeviceSecurityGroupId
    )
    
    try {
        Write-Host "🔄 Assigning enrollment time device membership target..." -ForegroundColor Yellow
        
        # Create a simple string-based body that matches the browser request exactly
        $bodyJson = '{
  "enrollmentTimeDeviceMembershipTargets": [
    {
      "targetType": "staticSecurityGroup",
      "targetId": "' + $DeviceSecurityGroupId + '"
    }
  ]
}'
        
        # Show the exact request body for debugging
        Write-Host "Request body: $bodyJson" -ForegroundColor Cyan
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$PolicyId')/setEnrollmentTimeDeviceMembershipTarget"
        Write-Host "Request URL: $url" -ForegroundColor Cyan
        
        # Add sleep to ensure any background processing has completed
        Start-Sleep -Seconds 2
        
        # Add headers explicitly for more control
        $authHeader = @{
            'Authorization' = "Bearer $GraphToken"
            'Content-Type' = 'application/json'
            'Accept' = 'application/json'
        }
        
        # Use Invoke-RestMethod directly for this call for better error handling
        $response = Invoke-RestMethod -Uri $url -Method POST -Headers $authHeader -Body $bodyJson -ErrorAction Stop
        Write-Host "✅ Enrollment time device membership target assigned successfully" -ForegroundColor Green
        return $response
    }
    catch {
        Write-Error "Failed to assign enrollment time device membership target: $_"
        
        # Enhanced error information
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode.value__
            Write-Host "Status code: $statusCode" -ForegroundColor Red
            
            try {
                $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
                $reader.BaseStream.Position = 0
                $reader.DiscardBufferedData()
                $responseBody = $reader.ReadToEnd()
                Write-Host "Error details: $responseBody" -ForegroundColor Red
            }
            catch {
                Write-Host "Could not read error response: $_" -ForegroundColor Red
            }
        }
        
        # Try alternate approach if this continues to fail
        Write-Host "⚠️ Trying fallback method for device group assignment..." -ForegroundColor Yellow
        $success = Set-WADPDeviceGroupAssignment -GraphToken $GraphToken -PolicyId $PolicyId -DeviceSecurityGroupId $DeviceSecurityGroupId
        
        if ($success) {
            Write-Host "✅ Device group assigned using alternate method" -ForegroundColor Green
            return $null
        } else {
            Write-Host "⚠️ Both assignment methods failed. Policy created but device group assignment will need to be done manually." -ForegroundColor Yellow
            return $null
        }
    }
}

# Function to assign policy to user group
function Set-WADPPolicyAssignment {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyId,
        
        [Parameter(Mandatory=$true)]
        [string]$UserSecurityGroupId
    )
    
    try {
        Write-Host "🔄 Assigning policy to user group..." -ForegroundColor Yellow
        
        $assignmentBody = @{
            assignments = @(
                @{
                    id = ""
                    source = "direct"
                    target = @{
                        groupId = $UserSecurityGroupId
                        "@odata.type" = "#microsoft.graph.groupAssignmentTarget"
                        deviceAndAppManagementAssignmentFilterType = "none"
                    }
                }
            )
        } | ConvertTo-Json -Depth 5
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$PolicyId')/assign"
        Post-GraphData -GraphToken $GraphToken -Url $url -Body $assignmentBody | Out-Null
        
        Write-Host "✅ Policy assignment completed successfully" -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to assign policy to user group: $_"
        throw
    }
}

# Function to verify policy creation
function Get-WADPPolicyDetails {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        Write-Host "🔍 Retrieving policy details..." -ForegroundColor Yellow
        
        $url = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$PolicyId')"
        $policy = Get-GraphData -GraphToken $GraphToken -Url $url
        
        Write-Host "✅ Policy verification completed" -ForegroundColor Green
        return $policy
    }
    catch {
        Write-Error "Failed to retrieve policy details: $_"
        throw
    }
}

# Main execution function
function Invoke-WADPPolicyCreation {
    try {
        Write-Host "`n📋 Windows Autopilot Device Preparation Policy Creation" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        # Input validation
        Write-Host "`n🔍 Validating input parameters..." -ForegroundColor Yellow
        
        if (-not (Test-IsValidGuid -InputGuid $DeviceSecurityGroupId)) {
            throw "Invalid Device Security Group ID format. Must be a valid GUID."
        }
        
        if (-not (Test-IsValidGuid -InputGuid $UserSecurityGroupId)) {
            throw "Invalid User Security Group ID format. Must be a valid GUID."
        }
        
        foreach ($app in $AllowedApps) {
            if (-not $app.id -or -not $app.type) {
                throw "Each allowed app must have 'id' and 'type' properties"
            }
            if (-not (Test-IsValidGuid -InputGuid $app.id)) {
                throw "Invalid app ID format: $($app.id). Must be a valid GUID."
            }
        }
        
        foreach ($scriptId in $AllowedScriptIds) {
            if (-not (Test-IsValidGuid -InputGuid $scriptId)) {
                throw "Invalid script ID format: $scriptId. Must be a valid GUID."
            }
        }
        
        Write-Host "✅ Input validation completed" -ForegroundColor Green
        
        # Display configuration summary
        Write-Host "`n📝 Policy Configuration:" -ForegroundColor Cyan
        Write-Host "   • Policy Name: $PolicyDisplayName" -ForegroundColor Cyan
        Write-Host "   • Description: $PolicyDescription" -ForegroundColor Cyan
        Write-Host "   • Device Security Group: $DeviceSecurityGroupId" -ForegroundColor Cyan
        Write-Host "   • User Security Group: $UserSecurityGroupId" -ForegroundColor Cyan
        Write-Host "   • Deployment Mode: $(if($DeploymentMode -eq '0'){'Standard'}else{'Enhanced'})" -ForegroundColor Cyan
        Write-Host "   • Deployment Type: $(if($DeploymentType -eq '0'){'User-driven'}else{'Self-deploying'})" -ForegroundColor Cyan
        Write-Host "   • Join Type: $(if($JoinType -eq '0'){'Entra ID joined'}else{'Entra ID hybrid joined'})" -ForegroundColor Cyan
        Write-Host "   • Account Type: $(if($AccountType -eq '0'){'Standard User'}else{'Administrator'})" -ForegroundColor Cyan
        Write-Host "   • Timeout: $TimeoutInMinutes minutes" -ForegroundColor Cyan
        Write-Host "   • Allow Skip: $AllowSkip" -ForegroundColor Cyan
        Write-Host "   • Allow Diagnostics: $AllowDiagnostics" -ForegroundColor Cyan
        if ($AllowedApps.Count -gt 0) {
            Write-Host "   • Allowed Apps: $($AllowedApps.Count) app(s)" -ForegroundColor Cyan
        }
        if ($AllowedScriptIds.Count -gt 0) {
            Write-Host "   • Allowed Scripts: $($AllowedScriptIds.Count) script(s)" -ForegroundColor Cyan
        }
        
        # Get access token
        Write-Host "`n🔐 Authenticating with Microsoft Graph..." -ForegroundColor Yellow
        $graphToken = Get-GraphAPIAccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
        Write-Host "✅ Authentication completed" -ForegroundColor Green
        
        # Validate security group ownership
        Test-SecurityGroupOwnership -GraphToken $graphToken -GroupId $DeviceSecurityGroupId
        
        # Prepare policy settings
        $policySettings = @{
            DisplayName = $PolicyDisplayName
            Description = $PolicyDescription
            DeploymentMode = $DeploymentMode
            DeploymentType = $DeploymentType
            JoinType = $JoinType
            AccountType = $AccountType
            TimeoutInMinutes = $TimeoutInMinutes
            CustomErrorMessage = $CustomErrorMessage
            AllowSkip = $AllowSkip
            AllowDiagnostics = $AllowDiagnostics
            AllowedApps = $AllowedApps
            AllowedScriptIds = $AllowedScriptIds
        }
        
        # Create the policy
        $policy = New-WADPConfigurationPolicy -GraphToken $graphToken -PolicySettings $policySettings
        
        # Variable to track if we should proceed despite enrollment time assignment failures
        $proceedWithAssignment = $true
        
        # Assign just-in-time configuration
        try {
            $result = Set-WADPJustInTimeConfiguration -GraphToken $graphToken -PolicyId $policy.id -DeviceSecurityGroupId $DeviceSecurityGroupId
            # If we get here, assignment worked
        } 
        catch {
            Write-Host "⚠️ Enrollment time device membership target assignment failed, but continuing with policy assignment" -ForegroundColor Yellow
            Write-Host "ℹ️ You may need to manually assign the device group in the Intune Portal" -ForegroundColor Yellow
            $proceedWithAssignment = $true  # Still proceed with policy assignment
        }
        
        # Assign policy to user group
        if ($proceedWithAssignment) {
            Set-WADPPolicyAssignment -GraphToken $graphToken -PolicyId $policy.id -UserSecurityGroupId $UserSecurityGroupId
            
            # Verify creation
            $policyDetails = Get-WADPPolicyDetails -GraphToken $graphToken -PolicyId $policy.id
            
            # Final success message
            Write-Host "`n✨ Windows Autopilot Device Preparation Policy created successfully!" -ForegroundColor Green
            Write-Host "🔗 Policy ID: $($policy.id)" -ForegroundColor Cyan
            Write-Host "🔗 Intune Portal: https://intune.microsoft.com/#view/Microsoft_Intune_DeviceSettings/DevicesEnrollmentMenu/~/windowsEnrollment" -ForegroundColor Cyan
            
            return $policy
        } else {
            throw "Policy creation process could not be completed."
        }
    }
    catch {
        Write-Host "`n❌ Policy creation failed: $_" -ForegroundColor Red
        throw
    }
}

# Script execution
try {
    Write-Host "🚀 Starting Windows Autopilot Device Preparation Policy creation..." -ForegroundColor Cyan
    
    $result = Invoke-WADPPolicyCreation
    
    Write-Host "`n🎉 Script completed successfully!" -ForegroundColor Green
    Write-Host "📋 Policy '$PolicyDisplayName' is now ready for use" -ForegroundColor Green
}
catch {
    Write-Host "`n💥 Script execution failed!" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}