# Windows Autopatch Group Creator
# Creates Autopatch Groups via Microsoft Graph API
# Interactive script - no parameters needed, will prompt for required values

# Embedded JSON configuration for testing
$script:EmbeddedConfig = @"
{
  "name": "auto-patch-group",
  "description": "",
  "globalUserManagedAadGroups": [],
  "deploymentGroups": [
    {
      "aadId": "00000000-0000-0000-0000-000000000000",
      "name": "auto-patch-group - Test",
      "userManagedAadGroups": [
        {
          "name": "[Azure]-[ConditonalAccess]-[Prod]-[CAD003-PolicyExclude]-[UG]",
          "id": "410a28bd-9c9f-403f-b1b2-4a0bd04e98d9",
          "type": 0
        }
      ],
      "failedPreRequisiteCheckCount": 0,
      "deploymentGroupPolicySettings": {
        "aadGroupName": "auto-patch-group - Test",
        "isUpdateSettingsModified": false,
        "deviceConfigurationSetting": {
          "policyId": "000",
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deadline": 1,
            "deferral": 0,
            "gracePeriod": 0
          },
          "featureDeploymentSettings": {
            "deadline": 5,
            "deferral": 0
          },
          "updateFrequencyUI": null,
          "installDays": null,
          "installTime": null,
          "activeHourEndTime": null,
          "activeHourStartTime": null
        },
        "dnfUpdateCloudSetting": {
          "policyId": "000",
          "approvalType": "Automatic",
          "deploymentDeferralInDays": 0
        },
        "officeDCv2Setting": {
          "policyId": "000",
          "deadline": 1,
          "deferral": 0,
          "hideUpdateNotifications": false,
          "targetChannel": "MonthlyEnterprise"
        },
        "edgeDCv2Setting": {
          "policyId": "000",
          "targetChannel": "Beta"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 24H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true
        }
      }
    },
    {
      "name": "auto-patch-group - Ring1",
      "userManagedAadGroups": [
        {
          "name": "[Azure]-[ConditonalAccess]-[Prod]-[CAD002-PolicyExclude]-[UG]",
          "id": "35d09841-af73-43e6-a59f-024fef1b6b95",
          "type": 0
        }
      ],
      "aadId": "00000000-0000-0000-0000-000000000000",
      "deploymentGroupPolicySettings": {
        "aadGroupName": "auto-patch-group - Ring1",
        "isUpdateSettingsModified": false,
        "deviceConfigurationSetting": {
          "policyId": "000",
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deadline": 2,
            "deferral": 1,
            "gracePeriod": 2
          },
          "featureDeploymentSettings": {
            "deadline": 5,
            "deferral": 0
          },
          "updateFrequencyUI": null,
          "installDays": null,
          "installTime": null,
          "activeHourEndTime": null,
          "activeHourStartTime": null
        },
        "dnfUpdateCloudSetting": {
          "policyId": "000",
          "approvalType": "Automatic",
          "deploymentDeferralInDays": 1
        },
        "officeDCv2Setting": {
          "policyId": "000",
          "deadline": 2,
          "deferral": 1,
          "hideUpdateNotifications": false,
          "targetChannel": "MonthlyEnterprise"
        },
        "edgeDCv2Setting": {
          "policyId": "000",
          "targetChannel": "Stable"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 24H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true
        }
      }
    },
    {
      "aadId": "00000000-0000-0000-0000-000000000000",
      "name": "auto-patch-group - Last",
      "userManagedAadGroups": [
        {
          "name": "[Azure]-[ConditonalAccess]-[Prod]-[CAD005-PolicyExclude]-[UG]",
          "id": "48fe6d79-f045-448a-bd74-716db27f0783",
          "type": 0
        }
      ],
      "failedPreRequisiteCheckCount": 0,
      "deploymentGroupPolicySettings": {
        "aadGroupName": "auto-patch-group - Last",
        "isUpdateSettingsModified": false,
        "deviceConfigurationSetting": {
          "policyId": "000",
          "updateBehavior": "AutoInstallAndRestart",
          "notificationSetting": "DefaultNotifications",
          "qualityDeploymentSettings": {
            "deadline": 3,
            "deferral": 5,
            "gracePeriod": 2
          },
          "featureDeploymentSettings": {
            "deadline": 5,
            "deferral": 0
          },
          "updateFrequencyUI": null,
          "installDays": null,
          "installTime": null,
          "activeHourEndTime": null,
          "activeHourStartTime": null
        },
        "dnfUpdateCloudSetting": {
          "policyId": "000",
          "approvalType": "Automatic",
          "deploymentDeferralInDays": 5
        },
        "officeDCv2Setting": {
          "policyId": "000",
          "deadline": 3,
          "deferral": 5,
          "hideUpdateNotifications": false,
          "targetChannel": "MonthlyEnterprise"
        },
        "edgeDCv2Setting": {
          "policyId": "000",
          "targetChannel": "Stable"
        },
        "featureUpdateAnchorCloudSetting": {
          "targetOSVersion": "Windows 11, version 24H2",
          "installLatestWindows10OnWindows11IneligibleDevice": true
        }
      }
    }
  ],
  "windowsUpdateSettings": [],
  "status": "Unknown",
  "type": "Unknown",
  "distributionType": "Unknown",
  "driverUpdateSettings": [],
  "enableDriverUpdate": true,
  "scopeTags": [0],
  "enabledContentTypes": 31
}
"@

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

# Function to get Graph API access token using client credentials
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
        Write-Host "🔐 Authenticating using Client Credentials flow..." -ForegroundColor Yellow
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
        Write-Error "Failed to get access token using client credentials: $_"
        throw
    }
}

# Function to get Graph API access token using Device Code Flow (modern user authentication)
function Get-GraphAPIAccessTokenDeviceCode {
    param (
        [Parameter(Mandatory=$true)]
        [string]$TenantId
    )
    
    try {
        Write-Host "🔐 Starting Device Code Flow authentication..." -ForegroundColor Yellow
        
        # Use PowerShell's built-in client ID for public clients
        $powerShellClientId = "1950a258-227b-4e31-a9cf-717495945fc2"
        $scope = "https://graph.microsoft.com/.default"
        
        # Step 1: Request device code
        $deviceCodeUrl = "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/devicecode"
        $deviceCodeBody = @{
            client_id = $powerShellClientId
            scope     = $scope
        }
        
        Write-Host "📱 Requesting device code..." -ForegroundColor Yellow
        $deviceCodeResponse = Invoke-RestMethod -Uri $deviceCodeUrl -Method POST -Body $deviceCodeBody -ContentType "application/x-www-form-urlencoded"
        
        # Step 2: Display user instructions and open browser automatically
        Write-Host "`n" -ForegroundColor Yellow
        Write-Host "🌐 AUTHENTICATION REQUIRED" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host "1. Opening web browser automatically..." -ForegroundColor White
        Write-Host "2. Enter the following code: $($deviceCodeResponse.user_code)" -ForegroundColor Yellow
        Write-Host "3. Sign in with your Azure AD credentials" -ForegroundColor White
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        # Automatically open browser
        try {
            if ($IsWindows -or $env:OS -eq "Windows_NT") {
                Start-Process $deviceCodeResponse.verification_uri
            }
            elseif ($IsMacOS -or (Test-Path "/Applications")) {
                & open $deviceCodeResponse.verification_uri
            }
            elseif ($IsLinux -or (Get-Command xdg-open -ErrorAction SilentlyContinue)) {
                & xdg-open $deviceCodeResponse.verification_uri
            }
            else {
                Write-Host "⚠️ Could not automatically open browser. Please manually navigate to: $($deviceCodeResponse.verification_uri)" -ForegroundColor Yellow
            }
        }
        catch {
            Write-Host "⚠️ Could not automatically open browser. Please manually navigate to: $($deviceCodeResponse.verification_uri)" -ForegroundColor Yellow
        }
        
        Write-Host "`nWaiting for authentication to complete..." -ForegroundColor Yellow
        
        # Step 3: Poll for token
        $tokenUrl = "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token"
        $tokenBody = @{
            grant_type  = "urn:ietf:params:oauth:grant-type:device_code"
            client_id   = $powerShellClientId
            device_code = $deviceCodeResponse.device_code
        }
        
        $timeout = [DateTime]::Now.AddSeconds(60)
        $interval = $deviceCodeResponse.interval
        
        do {
            Start-Sleep -Seconds $interval
            
            try {
                $tokenResponse = Invoke-RestMethod -Uri $tokenUrl -Method POST -Body $tokenBody -ContentType "application/x-www-form-urlencoded"
                
                Write-Host "✅ Authentication successful!" -ForegroundColor Green
                return $tokenResponse.access_token
            }
            catch {
                $statusCode = $_.Exception.Response.StatusCode.value__
                if ($statusCode -eq 400) {
                    # Parse error response from the exception message
                    $errorMessage = $_.Exception.Message
                    if ($errorMessage -like "*authorization_pending*") {
                        Write-Host "⏳ Still waiting for user authentication..." -ForegroundColor Yellow
                        continue
                    }
                    elseif ($errorMessage -like "*slow_down*") {
                        $interval += 5
                        Write-Host "⏳ Slowing down polling interval..." -ForegroundColor Yellow
                        continue
                    }
                    elseif ($errorMessage -like "*expired_token*") {
                        throw "Device code expired. Please restart the authentication process."
                    }
                    elseif ($errorMessage -like "*access_denied*") {
                        throw "User cancelled the authentication process."
                    }
                    else {
                        # Try to extract JSON error from the message
                        try {
                            $jsonMatch = [regex]::Match($errorMessage, '\{.*\}')
                            if ($jsonMatch.Success) {
                                $errorBody = $jsonMatch.Value | ConvertFrom-Json
                                if ($errorBody.error -eq "authorization_pending") {
                                    Write-Host "⏳ Still waiting for user authentication..." -ForegroundColor Yellow
                                    continue
                                }
                                elseif ($errorBody.error -eq "slow_down") {
                                    $interval += 5
                                    Write-Host "⏳ Slowing down polling interval..." -ForegroundColor Yellow
                                    continue
                                }
                                else {
                                    throw "Authentication failed: $($errorBody.error_description)"
                                }
                            }
                        }
                        catch {
                            # If JSON parsing fails, just continue polling
                            Write-Host "⏳ Still waiting for user authentication..." -ForegroundColor Yellow
                            continue
                        }
                    }
                }
                else {
                    throw "Unexpected error during token request: $_"
                }
            }
        } while ([DateTime]::Now -lt $timeout)
        
        throw "Authentication timeout. Please try again."
    }
    catch {
        Write-Error "Failed to get access token using Device Code Flow: $_"
        throw
    }
}

# Function to create Autopatch Group
function New-AutopatchGroup {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ConfigJson
    )
    
    try {
        Write-Host "🔄 Creating Windows Autopatch Group..." -ForegroundColor Yellow
        
        # Parse JSON config to validate structure
        $config = $ConfigJson | ConvertFrom-Json
        Write-Host "✅ Configuration validated" -ForegroundColor Green
        
        # Show configuration summary
        Write-Host "📝 Autopatch Group Configuration:" -ForegroundColor Cyan
        Write-Host "   • Group Name: $($config.name)" -ForegroundColor Cyan
        Write-Host "   • Description: $($config.description)" -ForegroundColor Cyan
        Write-Host "   • Deployment Groups: $($config.deploymentGroups.Count)" -ForegroundColor Cyan
        Write-Host "   • Driver Updates: $($config.enableDriverUpdate)" -ForegroundColor Cyan
        Write-Host "   • Content Types: $($config.enabledContentTypes)" -ForegroundColor Cyan
        
        # Make the API call to create the autopatch group
        $url = "https://services.autopatch.microsoft.com/device/v2/autopatchGroups"
        $response = Post-GraphData -GraphToken $GraphToken -Url $url -Body $ConfigJson
        
        Write-Host "✅ Autopatch Group created successfully" -ForegroundColor Green
        return $response
    }
    catch {
        Write-Error "Failed to create Autopatch Group: $_"
        throw
    }
}

# Function to get Autopatch Groups
function Get-AutopatchGroups {
    param (
        [Parameter(Mandatory=$true)]
        [string]$GraphToken
    )
    
    try {
        Write-Host "🔍 Retrieving existing Autopatch Groups..." -ForegroundColor Yellow
        
        $url = "https://services.autopatch.microsoft.com/device/v2/autopatchGroups"
        $groups = Get-GraphData -GraphToken $GraphToken -Url $url
        
        Write-Host "✅ Retrieved $($groups.Count) Autopatch Groups" -ForegroundColor Green
        return $groups
    }
    catch {
        Write-Error "Failed to retrieve Autopatch Groups: $_"
        throw
    }
}

# Function to display Autopatch Group details
function Show-AutopatchGroupDetails {
    param (
        [Parameter(Mandatory=$true)]
        [object]$Group
    )
    
    Write-Host "`n📋 Autopatch Group Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "🔗 ID: $($Group.id)" -ForegroundColor White
    Write-Host "📝 Name: $($Group.name)" -ForegroundColor White
    Write-Host "📄 Description: $($Group.description)" -ForegroundColor White
    Write-Host "🏷️ Type: $($Group.type)" -ForegroundColor White
    Write-Host "⚡ Status: $($Group.status)" -ForegroundColor White
    Write-Host "🔄 Distribution Type: $($Group.distributionType)" -ForegroundColor White
    
    if ($Group.flowId) {
        Write-Host "🔄 Flow ID: $($Group.flowId)" -ForegroundColor White
        Write-Host "🔄 Flow Type: $($Group.flowType)" -ForegroundColor White
        Write-Host "⚡ Flow Status: $($Group.flowStatus)" -ForegroundColor White
    }
    
    Write-Host "🎯 Scope Tags: $($Group.scopeTags -join ', ')" -ForegroundColor White
    Write-Host "📱 Registered Devices: $($Group.numberOfRegisteredDevices)" -ForegroundColor White
    
    if ($Group.deploymentGroups -and $Group.deploymentGroups.Count -gt 0) {
        Write-Host "`n📦 Deployment Groups:" -ForegroundColor Cyan
        foreach ($deployGroup in $Group.deploymentGroups) {
            Write-Host "   • $($deployGroup.name)" -ForegroundColor White
            if ($deployGroup.userManagedAadGroups -and $deployGroup.userManagedAadGroups.Count -gt 0) {
                Write-Host "     └── Managed Groups: $($deployGroup.userManagedAadGroups.Count)" -ForegroundColor Gray
            }
        }
    }
    
    if ($Group.globalUserManagedAadGroups -and $Group.globalUserManagedAadGroups.Count -gt 0) {
        Write-Host "`n🌐 Global Managed Groups: $($Group.globalUserManagedAadGroups.Count)" -ForegroundColor Cyan
    }
}

# Main execution function
function Invoke-AutopatchGroupCreation {
    try {
        Write-Host "`n📋 Windows Autopatch Group Creation" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        # Prompt for Tenant ID first
        Write-Host "`n🏢 Enter your Azure AD Tenant Information:" -ForegroundColor Yellow
        do {
            $TenantId = Read-Host "Tenant ID (Directory ID)"
            if (-not $TenantId) {
                Write-Host "❌ Tenant ID is required." -ForegroundColor Red
            }
            elseif (-not (Test-IsValidGuid -InputGuid $TenantId)) {
                Write-Host "❌ Invalid Tenant ID format. Must be a valid GUID." -ForegroundColor Red
                $TenantId = $null
            }
        } while (-not $TenantId)
        
        # Authentication method selection
        Write-Host "`n🔐 Select Authentication Method:" -ForegroundColor Yellow
        Write-Host "1. Client Credentials (Service-to-Service)" -ForegroundColor White
        Write-Host "   • Uses Client ID and Client Secret" -ForegroundColor Gray
        Write-Host "   • For automated scripts and service accounts" -ForegroundColor Gray
        Write-Host "   • No user interaction required" -ForegroundColor Gray
        Write-Host ""
        Write-Host "2. Device Code Flow (Interactive User)" -ForegroundColor White
        Write-Host "   • Opens browser for user authentication" -ForegroundColor Gray
        Write-Host "   • Supports MFA and Conditional Access" -ForegroundColor Gray
        Write-Host "   • Modern authentication method" -ForegroundColor Gray
        
        do {
            $authChoice = Read-Host "`nEnter your choice (1 or 2)"
        } while ($authChoice -notin @("1", "2"))
        
        $authMethod = ""
        $ClientId = $null
        $ClientSecret = $null
        
        if ($authChoice -eq "1") {
            $authMethod = "ClientCredentials"
            Write-Host "`n📋 Client Credentials Authentication Selected" -ForegroundColor Cyan
            
            # Prompt for Client ID
            do {
                $ClientId = Read-Host "Client ID (Application ID)"
                if (-not $ClientId) {
                    Write-Host "❌ Client ID is required." -ForegroundColor Red
                }
                elseif (-not (Test-IsValidGuid -InputGuid $ClientId)) {
                    Write-Host "❌ Invalid Client ID format. Must be a valid GUID." -ForegroundColor Red
                    $ClientId = $null
                }
            } while (-not $ClientId)
            
            # Prompt for Client Secret
            do {
                $ClientSecret = Read-Host "Client Secret" -AsSecureString
                if (-not $ClientSecret -or $ClientSecret.Length -eq 0) {
                    Write-Host "❌ Client Secret is required." -ForegroundColor Red
                    $ClientSecret = $null
                }
            } while (-not $ClientSecret)
            
            # Convert SecureString to plain text for API call
            $ClientSecretPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($ClientSecret))
        }
        else {
            $authMethod = "DeviceCode"
            Write-Host "`n📋 Device Code Flow Authentication Selected" -ForegroundColor Cyan
        }
        
        Write-Host "✅ Input collection completed" -ForegroundColor Green
        
        # Get access token based on authentication method
        Write-Host "`n🔐 Authenticating with Microsoft Graph..." -ForegroundColor Yellow
        
        if ($authMethod -eq "ClientCredentials") {
            $graphToken = Get-GraphAPIAccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecretPlain
            # Clear the plain text secret
            $ClientSecretPlain = $null
        }
        elseif ($authMethod -eq "DeviceCode") {
            $graphToken = Get-GraphAPIAccessTokenDeviceCode -TenantId $TenantId
        }
        
        Write-Host "✅ Authentication completed" -ForegroundColor Green
        
        # Get existing groups (optional - for verification)
        Write-Host "`n🔍 Checking existing Autopatch Groups..." -ForegroundColor Yellow
        $existingGroups = Get-AutopatchGroups -GraphToken $graphToken
        
        if ($existingGroups -and $existingGroups.Count -gt 0) {
            Write-Host "📋 Found $($existingGroups.Count) existing groups:" -ForegroundColor Cyan
            foreach ($group in $existingGroups) {
                $statusColor = switch ($group.status) {
                    "Active" { "Green" }
                    "Creating" { "Yellow" }
                    "InProgress" { "Yellow" }
                    default { "Gray" }
                }
                Write-Host "   • $($group.name) - Status: $($group.status)" -ForegroundColor $statusColor
            }
        }
        
        # Create the autopatch group
        Write-Host "`n🚀 Creating new Autopatch Group..." -ForegroundColor Yellow
        $newGroup = New-AutopatchGroup -GraphToken $graphToken -ConfigJson $script:EmbeddedConfig
        
        # Display results
        if ($newGroup) {
            Write-Host "`n✨ Windows Autopatch Group created successfully!" -ForegroundColor Green
            Show-AutopatchGroupDetails -Group $newGroup
            
            Write-Host "`n🔗 Next Steps:" -ForegroundColor Cyan
            Write-Host "   • Monitor group creation status in the Autopatch portal" -ForegroundColor White
            Write-Host "   • Verify deployment group configurations" -ForegroundColor White
            Write-Host "   • Assign devices to the appropriate Azure AD groups" -ForegroundColor White
        }
        
        return $newGroup
    }
    catch {
        Write-Host "`n❌ Autopatch Group creation failed: $_" -ForegroundColor Red
        throw
    }
}

# Script execution
try {
    Write-Host "🚀 Starting Windows Autopatch Group creation..." -ForegroundColor Cyan
    
    $result = Invoke-AutopatchGroupCreation
    
    Write-Host "`n🎉 Script completed successfully!" -ForegroundColor Green
    Write-Host "📋 Autopatch Group '$($result.name)' is now being created" -ForegroundColor Green
}
catch {
    Write-Host "`n💥 Script execution failed!" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}