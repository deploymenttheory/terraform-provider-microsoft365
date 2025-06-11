# Microsoft Graph Data Exporter
# PowerShell version of the Python data extraction script
[CmdletBinding()]
param (
    [Parameter(Mandatory=$true,
    HelpMessage="Specify the Entra ID tenant ID (Directory ID)")]
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
    
    [Parameter(Mandatory=$false,
    HelpMessage="Specify the full path where you want to save all extracted data (e.g., 'C:\GraphExports' or './MyExports')")]
    [ValidateNotNullOrEmpty()]
    [string]$ExportPath = ".\GraphData"
)

# =============================================================================
# AUTHENTICATION MODULE
# =============================================================================

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
        Write-Host "🔐 Authenticating with Microsoft Graph..." -ForegroundColor Yellow
        
        $tokenUrl = "https://login.microsoftonline.com/$TenantId/oauth2/v2.0/token"
        
        $body = @{
            client_id     = $ClientId
            client_secret = $ClientSecret
            scope         = "https://graph.microsoft.com/.default"
            grant_type    = "client_credentials"
        }
        
        $response = Invoke-RestMethod -Uri $tokenUrl -Method POST -Body $body -ContentType "application/x-www-form-urlencoded"
        Write-Host "✅ Authentication completed successfully" -ForegroundColor Green
        return $response.access_token
    }
    catch {
        Write-Host "❌ Authentication failed: $_" -ForegroundColor Red
        throw
    }
}

# =============================================================================
# API CALL HELPER FUNCTIONS
# =============================================================================

function Invoke-GraphApiCall {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$Url,
        
        [Parameter(Mandatory=$false)]
        [string]$Method = "GET",
        
        [Parameter(Mandatory=$false)]
        [string]$Body = $null,
        
        [Parameter(Mandatory=$false)]
        [bool]$HandlePaging = $true
    )
    
    $authHeader = @{
        'Authorization' = "Bearer $GraphToken"
        'Content-Type' = 'application/json'
    }
    
    $retryCount = 0
    $maxRetries = 3
    $allResults = @()
    $currentUrl = $Url
    
    do {
        while ($retryCount -le $maxRetries) {
            try {
                if ($Method -eq "POST" -and $Body) {
                    $response = Invoke-RestMethod -Uri $currentUrl -Method $Method -Headers $authHeader -Body $Body
                } else {
                    $response = Invoke-RestMethod -Uri $currentUrl -Method $Method -Headers $authHeader
                }
                
                # Handle different response types
                if ($response.value) {
                    # This is a paginated response
                    $allResults += $response.value
                    
                    if ($HandlePaging -and $response.'@odata.nextLink') {
                        $currentUrl = $response.'@odata.nextLink'
                        $retryCount = 0  # Reset retry count on success
                        break  # Exit retry loop, continue with next page
                    } else {
                        # No more pages or paging disabled
                        return $allResults
                    }
                } else {
                    # This is a single object response
                    return $response
                }
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
                    Write-Host "⚠️ Rate limited. Retrying in $retryAfter seconds... (Attempt $retryCount/$maxRetries)" -ForegroundColor Yellow
                    Start-Sleep -Seconds $retryAfter
                } else {
                    Write-Host "❌ API call failed: $_" -ForegroundColor Red
                    Write-Host "❌ URL: $currentUrl" -ForegroundColor Red
                    throw
                }
            }
        }
        
        if ($retryCount -gt $maxRetries) {
            throw "Max retry attempts reached for $Method request to $currentUrl"
        }
        
    } while ($HandlePaging -and $response.'@odata.nextLink')
    
    return $allResults
}

function Get-IntunePortalData {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Url,
        
        [Parameter(Mandatory=$false)]
        [hashtable]$Headers = @{}
    )
    
    $retryCount = 0
    $maxRetries = 3
    
    while ($retryCount -le $maxRetries) {
        try {
            $response = Invoke-RestMethod -Uri $Url -Method GET -Headers $Headers
            return $response
        }
        catch {
            $retryCount++
            Write-Host "⚠️ Request failed, retrying... (Attempt $retryCount/$maxRetries)" -ForegroundColor Yellow
            Start-Sleep -Seconds ($retryCount * 2)
            
            if ($retryCount -gt $maxRetries) {
                Write-Host "❌ Failed to retrieve data from: $Url" -ForegroundColor Red
                throw
            }
        }
    }
}

function Save-JsonData {
    param(
        [Parameter(Mandatory=$true)]
        [AllowNull()]
        [object]$Data,
        
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    try {
        if ($null -eq $Data) {
            Write-Host "  ⚠️ Skipping $(Split-Path -Leaf $FilePath) - No data to save" -ForegroundColor Yellow
            return
        }
        
        $directory = Split-Path -Path $FilePath -Parent
        if (-not (Test-Path -Path $directory)) {
            New-Item -ItemType Directory -Path $directory -Force | Out-Null
        }
        
        $Data | ConvertTo-Json -Depth 20 | Out-File -FilePath $FilePath -Encoding UTF8
        Write-Host "  ✅ Saved: $(Split-Path -Leaf $FilePath)" -ForegroundColor Green
    }
    catch {
        Write-Host "  ❌ Failed to save: $(Split-Path -Leaf $FilePath) - $_" -ForegroundColor Red
        throw
    }
}

# =============================================================================
# DATA EXTRACTION MODULES
# =============================================================================

function Get-SettingStatusErrors {
    param(
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Setting Status Errors..." -ForegroundColor Cyan
        
        # Get version information from Intune portal
        $versionResponse = Get-IntunePortalData -Url "https://intune.microsoft.com/signin/idpRedirect.js"
        
        if ($versionResponse -match '"extensionsPageVersion":(\{[^}]+\})') {
            $versions = $matches[1] | ConvertFrom-Json
            $deviceSettingsVersion = $versions.'Microsoft_Intune_DeviceSettings'[0]
            
            $root = "https://afd-v2.hosting.portal.azure.net"
            $settingStatusUrl = "$root/intunedevicesettings/Content/$deviceSettingsVersion/Scripts/DeviceConfiguration/Blades/DevicePoliciesStatus/SettingStatus.js"
            
            $settingStatusData = Get-IntunePortalData -Url $settingStatusUrl
            
            if ($settingStatusData -match 'SettingStatusErrorMap = (\{[^}]+\})') {
                $errorMap = $matches[1] | ConvertFrom-Json
                Save-JsonData -Data $errorMap -FilePath "$ExportPath\SettingStatusErrors.json"
            }
        }
        
        Write-Host "✅ Setting Status Errors extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract Setting Status Errors: $_" -ForegroundColor Red
        throw
    }
}

function Get-DCv1Policies {
    param(
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Device Configuration v1 Policies..." -ForegroundColor Cyan
        
        # Get version information
        $versionResponse = Get-IntunePortalData -Url "https://intune.microsoft.com/signin/idpRedirect.js"
        
        if ($versionResponse -match '"extensionsPageVersion":(\{[^}]+\})') {
            $versions = $matches[1] | ConvertFrom-Json
            $deviceSettingsVersion = $versions.'Microsoft_Intune_DeviceSettings'[0]
            
            $root = "https://afd-v2.hosting.portal.azure.net"
            $rootDeviceSettings = "$root/intunedevicesettings/Content/$deviceSettingsVersion/Scripts/DeviceConfiguration"
            
            # Clean existing DCv1 directory
            $dcv1Path = "$ExportPath\DCv1"
            if (Test-Path $dcv1Path) {
                Remove-Item -Path $dcv1Path -Recurse -Force
            }
            
            foreach ($source in @('Configuration', 'Compliance')) {
                Write-Host "  📋 Processing $source metadata..." -ForegroundColor Yellow
                
                $sourcePath = "$dcv1Path\$source"
                New-Item -ItemType Directory -Path $sourcePath -Force | Out-Null
                
                $metadataUrl = "$rootDeviceSettings/Metadata/${source}Metadata.js"
                $metadataResponse = Get-IntunePortalData -Url $metadataUrl
                
                if ($metadataResponse -match '(?s)metadata = (\{.+\});') {
                    $metadata = $matches[1] | ConvertFrom-Json
                    
                    foreach ($family in $metadata.PSObject.Properties) {
                        foreach ($setting in $family.Value) {
                            $cleanId = ($setting.id -split '_')[0..-2] -join '_'
                            $setting.id = $cleanId
                            
                            # Clean nested IDs
                            $setting = Remove-DCv1VersionSuffixes -Setting $setting
                            
                            $filePath = "$sourcePath\$cleanId.json"
                            Save-JsonData -Data $setting -FilePath $filePath
                        }
                    }
                }
            }
        }
        
        Write-Host "✅ Device Configuration v1 extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract DCv1 policies: $_" -ForegroundColor Red
        throw
    }
}

function Remove-DCv1VersionSuffixes {
    param(
        [Parameter(Mandatory=$true)]
        [object]$Setting
    )
    
    # Clean child settings
    if ($Setting.childSettings) {
        foreach ($child in $Setting.childSettings) {
            $child.id = ($child.id -split '_')[0..-2] -join '_'
            $child = Remove-DCv1VersionSuffixes -Setting $child
        }
    }
    
    # Clean options
    if ($Setting.options) {
        foreach ($option in $Setting.options) {
            if ($option.children) {
                foreach ($child in $option.children) {
                    $child.id = ($child.id -split '_')[0..-2] -join '_'
                    $child = Remove-DCv1VersionSuffixes -Setting $child
                }
            }
        }
    }
    
    # Clean complex options
    if ($Setting.complexOptions) {
        foreach ($complexOption in $Setting.complexOptions) {
            $complexOption.id = ($complexOption.id -split '_')[0..-2] -join '_'
            $complexOption = Remove-DCv1VersionSuffixes -Setting $complexOption
        }
    }
    
    # Clean columns
    if ($Setting.columns) {
        foreach ($column in $Setting.columns) {
            if ($column.metadata) {
                $column.metadata.id = ($column.metadata.id -split '_')[0..-2] -join '_'
                $column.metadata = Remove-DCv1VersionSuffixes -Setting $column.metadata
            }
        }
    }
    
    return $Setting
}

function Get-ServicePrincipalsAndEndpoints {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Service Principals and Endpoints..." -ForegroundColor Cyan
        
        # Get endpoints
        Write-Host "  📋 Processing endpoints..." -ForegroundColor Yellow
        $endpointsUrl = "https://graph.microsoft.com/beta/servicePrincipals/appId=0000000a-0000-0000-c000-000000000000/endpoints"
        $endpoints = Invoke-GraphApiCall -GraphToken $GraphToken -Url $endpointsUrl -HandlePaging $false
        
        if ($endpoints -and $endpoints.Count -gt 0) {
            $sortedEndpoints = $endpoints | Sort-Object capability
            Save-JsonData -Data $sortedEndpoints -FilePath "$ExportPath\Endpoints.json"
        } else {
            Write-Host "  ⚠️ No endpoints data returned" -ForegroundColor Yellow
        }
        
        # Get service principals
        Write-Host "  📋 Processing service principals..." -ForegroundColor Yellow
        $servicePrincipalsPath = "$ExportPath\ServicePrincipals"
        if (Test-Path $servicePrincipalsPath) {
            Remove-Item -Path $servicePrincipalsPath -Recurse -Force
        }
        New-Item -ItemType Directory -Path $servicePrincipalsPath -Force | Out-Null
        
        $servicePrincipals = Invoke-GraphApiCall -GraphToken $GraphToken -Url "https://graph.microsoft.com/beta/servicePrincipals"
        
        if ($servicePrincipals -and $servicePrincipals.Count -gt 0) {
            foreach ($sp in $servicePrincipals) {
                if ($sp.appId) {
                    $appId = $sp.appId
                    $filePath = "$servicePrincipalsPath\$appId.json"
                    Save-JsonData -Data $sp -FilePath $filePath
                }
            }
        } else {
            Write-Host "  ⚠️ No service principals data returned" -ForegroundColor Yellow
        }
        
        Write-Host "✅ Service Principals and Endpoints extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract Service Principals and Endpoints: $_" -ForegroundColor Red
        throw
    }
}

function Get-RoleDefinitions {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Role Definitions..." -ForegroundColor Cyan
        
        $roleDefinitionsPath = "$ExportPath\RoleDefinitions"
        if (Test-Path $roleDefinitionsPath) {
            Remove-Item -Path $roleDefinitionsPath -Recurse -Force
        }
        New-Item -ItemType Directory -Path $roleDefinitionsPath -Force | Out-Null
        
        $providers = @('cloudPC', 'deviceManagement', 'directory', 'entitlementManagement', 'exchange')
        $successfulProviders = @()
        $failedProviders = @()
        
        foreach ($provider in $providers) {
            Write-Host "  📋 Processing $provider role definitions..." -ForegroundColor Yellow
            
            try {
                $providerPath = "$roleDefinitionsPath\$provider"
                New-Item -ItemType Directory -Path $providerPath -Force | Out-Null
                
                $roleDefUrl = "https://graph.microsoft.com/beta/roleManagement/$provider/roleDefinitions"
                $roleDefinitions = Invoke-GraphApiCall -GraphToken $GraphToken -Url $roleDefUrl
                
                if ($roleDefinitions -and $roleDefinitions.Count -gt 0) {
                    foreach ($roleDef in $roleDefinitions) {
                        $roleId = $roleDef.id
                        $filePath = "$providerPath\$roleId.json"
                        Save-JsonData -Data $roleDef -FilePath $filePath
                    }
                    $successfulProviders += $provider
                } else {
                    Write-Host "    ⚠️ No role definitions found for $provider" -ForegroundColor Yellow
                    $successfulProviders += $provider
                }
            }
            catch {
                $errorMessage = $_.Exception.Message
                if ($errorMessage -like "*Authorization_RequestDenied*" -or $errorMessage -like "*Insufficient privileges*") {
                    Write-Host "    ⚠️ Insufficient privileges for $provider role definitions - skipping" -ForegroundColor Yellow
                    $failedProviders += $provider
                } else {
                    Write-Host "    ❌ Failed to extract $provider role definitions: $errorMessage" -ForegroundColor Red
                    $failedProviders += $provider
                }
            }
        }
        
        if ($successfulProviders.Count -gt 0) {
            Write-Host "✅ Role Definitions extraction completed for: $($successfulProviders -join ', ')" -ForegroundColor Green
        }
        
        if ($failedProviders.Count -gt 0) {
            Write-Host "⚠️ Role Definitions extraction failed for: $($failedProviders -join ', ') (insufficient permissions)" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "❌ Failed to extract Role Definitions: $_" -ForegroundColor Red
        throw
    }
}

function Get-ResourceOperations {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Resource Operations..." -ForegroundColor Cyan
        
        $resourceOperationsPath = "$ExportPath\ResourceOperations"
        if (Test-Path $resourceOperationsPath) {
            Remove-Item -Path $resourceOperationsPath -Recurse -Force
        }
        New-Item -ItemType Directory -Path $resourceOperationsPath -Force | Out-Null
        
        $resourceOpsUrl = "https://graph.microsoft.com/beta/deviceManagement/resourceOperations"
        $resourceOperations = Invoke-GraphApiCall -GraphToken $GraphToken -Url $resourceOpsUrl
        
        foreach ($operation in $resourceOperations) {
            $operationId = $operation.id
            $filePath = "$resourceOperationsPath\$operationId.json"
            Save-JsonData -Data $operation -FilePath $filePath
        }
        
        Write-Host "✅ Resource Operations extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract Resource Operations: $_" -ForegroundColor Red
        throw
    }
}

function Get-DefenderHuntingTables {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputPath
    )
    
    try {
        Write-Host "🔍 Extracting Defender Hunting Table Schemas..." -ForegroundColor Cyan
        
        $defenderPath = "$OutputPath\Defender"
        if (Test-Path $defenderPath) {
            Remove-Item -Path $defenderPath -Recurse -Force
        }
        New-Item -ItemType Directory -Path $defenderPath -Force | Out-Null
        
        $huntingTables = @(
            'AlertEvidence', 'AlertInfo', 'BehaviorEntities', 'BehaviorInfo',
            'AADSignInEventsBeta', 'AADSpnSignInEventsBeta', 'CloudAppEvents', 
            'IdentityInfo', 'IdentityLogonEvents', 'EmailAttachmentInfo', 
            'EmailEvents', 'EmailPostDeliveryEvents', 'EmailUrlInfo', 
            'UrlClickEvents', 'ExposureGraphEdges', 'ExposureGraphNodes'
        )
        
        $successfulTables = @()
        $failedTables = @()
        
        foreach ($table in $huntingTables) {
            Write-Host "  📋 Processing $table schema..." -ForegroundColor Yellow
            
            $query = "$table | getschema | project Description=`"`", Type=split(DataType, `".`")[1], Entity=`"`", Name=ColumnName"
            
            $requestBody = @{
                Query = $query
            } | ConvertTo-Json
            
            try {
                $response = Invoke-GraphApiCall -GraphToken $GraphToken -Url "https://graph.microsoft.com/beta/security/runHuntingQuery" -Method "POST" -Body $requestBody -HandlePaging $false
                
                if ($response -and $response.results) {
                    Save-JsonData -Data $response.results -FilePath "$defenderPath\$table.json"
                    $successfulTables += $table
                } elseif ($response -and $response.Results) {
                    # Handle case sensitivity
                    Save-JsonData -Data $response.Results -FilePath "$defenderPath\$table.json"
                    $successfulTables += $table
                } else {
                    Write-Host "    ⚠️ No results returned for $table" -ForegroundColor Yellow
                    $failedTables += $table
                }
            }
            catch {
                $errorMessage = $_.Exception.Message
                if ($errorMessage -like "*Authorization_RequestDenied*" -or $errorMessage -like "*Insufficient privileges*") {
                    Write-Host "    ⚠️ Insufficient privileges for $table - skipping" -ForegroundColor Yellow
                } else {
                    Write-Host "    ⚠️ Failed to get schema for $table : $errorMessage" -ForegroundColor Yellow
                }
                $failedTables += $table
            }
        }
        
        if ($successfulTables.Count -gt 0) {
            Write-Host "✅ Defender Hunting Tables extraction completed for: $($successfulTables.Count)/$($huntingTables.Count) tables" -ForegroundColor Green
        }
        
        if ($failedTables.Count -gt 0) {
            Write-Host "⚠️ Some Defender tables could not be accessed (permissions or availability): $($failedTables.Count) tables" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "❌ Failed to extract Defender Hunting Tables: $_" -ForegroundColor Red
        # Don't throw here - this is optional data
        Write-Host "⚠️ Defender hunting tables extraction will be skipped" -ForegroundColor Yellow
    }
}

function Get-DCv2ConfigurationSettings {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Device Configuration v2 Settings..." -ForegroundColor Cyan
        
        $dcv2Path = "$ExportPath\DCv2"
        if (Test-Path $dcv2Path) {
            Remove-Item -Path $dcv2Path -Recurse -Force
        }
        
        # Configuration Settings (Settings Catalog)
        Write-Host "  📋 Processing configuration settings..." -ForegroundColor Yellow
        $settingsPath = "$dcv2Path\Settings"
        New-Item -ItemType Directory -Path $settingsPath -Force | Out-Null
        
        $configSettingsUrl = "https://graph.microsoft.com/beta/deviceManagement/configurationSettings"
        $configSettings = Invoke-GraphApiCall -GraphToken $GraphToken -Url $configSettingsUrl
        
        foreach ($setting in $configSettings) {
            if ($setting.version) {
                $setting.PSObject.Properties.Remove('version')
            }
            $settingId = $setting.id
            $filePath = "$settingsPath\$settingId.json"
            Save-JsonData -Data $setting -FilePath $filePath
        }
        
        # Create backwards compatibility folder
        $backwardsCompatPath = "$ExportPath\settings"
        if (Test-Path $backwardsCompatPath) {
            Remove-Item -Path $backwardsCompatPath -Recurse -Force
        }
        Copy-Item -Path $settingsPath -Destination $backwardsCompatPath -Recurse
        
        Write-Host "✅ DCv2 Configuration Settings extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract DCv2 Configuration Settings: $_" -ForegroundColor Red
        throw
    }
}

function Get-DCv2ComplianceSettings {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Device Configuration v2 Compliance Settings..." -ForegroundColor Cyan
        
        $compliancePath = "$ExportPath\DCv2\Compliance"
        New-Item -ItemType Directory -Path $compliancePath -Force | Out-Null
        
        $complianceSettingsUrl = "https://graph.microsoft.com/beta/deviceManagement/complianceSettings"
        $complianceSettings = Invoke-GraphApiCall -GraphToken $GraphToken -Url $complianceSettingsUrl
        
        foreach ($setting in $complianceSettings) {
            if ($setting.version) {
                $setting.PSObject.Properties.Remove('version')
            }
            $settingId = $setting.id
            $filePath = "$compliancePath\$settingId.json"
            Save-JsonData -Data $setting -FilePath $filePath
        }
        
        Write-Host "✅ DCv2 Compliance Settings extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract DCv2 Compliance Settings: $_" -ForegroundColor Red
        throw
    }
}

function Get-DCv2PolicyTemplates {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Device Configuration v2 Policy Templates..." -ForegroundColor Cyan
        
        $templatesPath = "$ExportPath\DCv2\Templates"
        New-Item -ItemType Directory -Path $templatesPath -Force | Out-Null
        
        $templatesUrl = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicyTemplates"
        $templates = Invoke-GraphApiCall -GraphToken $GraphToken -Url $templatesUrl
        
        foreach ($template in $templates) {
            $templateId = $template.id
            $filePath = "$templatesPath\$templateId.json"
            Save-JsonData -Data $template -FilePath $filePath
        }
        
        Write-Host "✅ DCv2 Policy Templates extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract DCv2 Policy Templates: $_" -ForegroundColor Red
        throw
    }
}

function Get-DCv2InventorySettings {
    param(
        [Parameter(Mandatory=$true)]
        [string]$GraphToken,
        
        [Parameter(Mandatory=$true)]
        [string]$ExportPath
    )
    
    try {
        Write-Host "🔍 Extracting Device Configuration v2 Inventory Settings..." -ForegroundColor Cyan
        
        $inventoryPath = "$ExportPath\DCv2\Inventory"
        New-Item -ItemType Directory -Path $inventoryPath -Force | Out-Null
        
        $inventoryUrl = "https://graph.microsoft.com/beta/deviceManagement/inventorySettings"
        $inventorySettings = Invoke-GraphApiCall -GraphToken $GraphToken -Url $inventoryUrl
        
        foreach ($setting in $inventorySettings) {
            if ($setting.version) {
                $setting.PSObject.Properties.Remove('version')
            }
            $settingId = $setting.id
            $filePath = "$inventoryPath\$settingId.json"
            Save-JsonData -Data $setting -FilePath $filePath
        }
        
        Write-Host "✅ DCv2 Inventory Settings extraction completed" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to extract DCv2 Inventory Settings: $_" -ForegroundColor Red
        throw
    }
}

# =============================================================================
# MAIN EXECUTION MODULE
# =============================================================================

function Invoke-GraphDataExtraction {
    try {
        Write-Host "`n📊 Microsoft Graph Data Extraction Tool" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        
        # Create output directory
        if (Test-Path $ExportPath) {
            Write-Host "⚠️ Export directory exists. Some data may be overwritten." -ForegroundColor Yellow
        } else {
            New-Item -ItemType Directory -Path $ExportPath -Force | Out-Null
            Write-Host "✅ Created export directory: $ExportPath" -ForegroundColor Green
        }
        
        # Get authentication token
        $graphToken = Get-GraphAPIAccessToken -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
        
        # Execute data extraction modules with individual error handling
        Write-Host "`n🚀 Starting data extraction..." -ForegroundColor Cyan
        
        $successfulModules = @()
        $failedModules = @()
        
        # Module 1: Setting Status Errors
        try {
            Get-SettingStatusErrors -ExportPath $ExportPath
            $successfulModules += "Setting Status Errors"
        } catch {
            Write-Host "⚠️ Setting Status Errors module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "Setting Status Errors"
        }
        
        # Module 2: DCv1 Policies
        try {
            Get-DCv1Policies -ExportPath $ExportPath
            $successfulModules += "DCv1 Policies"
        } catch {
            Write-Host "⚠️ DCv1 Policies module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "DCv1 Policies"
        }
        
        # Module 3: Service Principals and Endpoints
        try {
            Get-ServicePrincipalsAndEndpoints -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "Service Principals and Endpoints"
        } catch {
            Write-Host "⚠️ Service Principals and Endpoints module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "Service Principals and Endpoints"
        }
        
        # Module 4: Role Definitions
        try {
            Get-RoleDefinitions -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "Role Definitions"
        } catch {
            Write-Host "⚠️ Role Definitions module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "Role Definitions"
        }
        
        # Module 5: Resource Operations
        try {
            Get-ResourceOperations -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "Resource Operations"
        } catch {
            Write-Host "⚠️ Resource Operations module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "Resource Operations"
        }
        
        # Module 6: Defender Hunting Tables
        try {
            Get-DefenderHuntingTables -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "Defender Hunting Tables"
        } catch {
            Write-Host "⚠️ Defender Hunting Tables module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "Defender Hunting Tables"
        }
        
        # Module 7: DCv2 Configuration Settings
        try {
            Get-DCv2ConfigurationSettings -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "DCv2 Configuration Settings"
        } catch {
            Write-Host "⚠️ DCv2 Configuration Settings module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "DCv2 Configuration Settings"
        }
        
        # Module 8: DCv2 Compliance Settings
        try {
            Get-DCv2ComplianceSettings -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "DCv2 Compliance Settings"
        } catch {
            Write-Host "⚠️ DCv2 Compliance Settings module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "DCv2 Compliance Settings"
        }
        
        # Module 9: DCv2 Policy Templates
        try {
            Get-DCv2PolicyTemplates -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "DCv2 Policy Templates"
        } catch {
            Write-Host "⚠️ DCv2 Policy Templates module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "DCv2 Policy Templates"
        }
        
        # Module 10: DCv2 Inventory Settings
        try {
            Get-DCv2InventorySettings -GraphToken $graphToken -ExportPath $ExportPath
            $successfulModules += "DCv2 Inventory Settings"
        } catch {
            Write-Host "⚠️ DCv2 Inventory Settings module failed - continuing with next module" -ForegroundColor Yellow
            $failedModules += "DCv2 Inventory Settings"
        }
        
        Write-Host "`n✨ Data extraction process completed!" -ForegroundColor Green
        Write-Host "📁 All data saved to: $ExportPath" -ForegroundColor Cyan
        
        # Display detailed summary
        if ($successfulModules.Count -gt 0) {
            Write-Host "`n✅ Successful modules ($($successfulModules.Count)):" -ForegroundColor Green
            foreach ($module in $successfulModules) {
                Write-Host "   • $module" -ForegroundColor Green
            }
        }
        
        if ($failedModules.Count -gt 0) {
            Write-Host "`n⚠️ Failed modules ($($failedModules.Count)):" -ForegroundColor Yellow
            foreach ($module in $failedModules) {
                Write-Host "   • $module" -ForegroundColor Yellow
            }
            Write-Host "`nNote: Some failures may be due to insufficient permissions or feature availability." -ForegroundColor Yellow
        }
        
        # Display file summary
        try {
            $subfolders = Get-ChildItem -Path $ExportPath -Directory | Measure-Object
            $files = Get-ChildItem -Path $ExportPath -File -Recurse | Measure-Object
            
            Write-Host "`n📊 Extraction Summary:" -ForegroundColor Cyan
            Write-Host "   • Total directories: $($subfolders.Count)" -ForegroundColor White
            Write-Host "   • Total files: $($files.Count)" -ForegroundColor White
            Write-Host "   • Export path: $ExportPath" -ForegroundColor White
        } catch {
            Write-Host "`n📊 Data extraction completed (summary unavailable)" -ForegroundColor Cyan
        }
        
    }
    catch {
        Write-Host "`n❌ Data extraction process failed: $_" -ForegroundColor Red
        throw
    }
}

# =============================================================================
# SCRIPT EXECUTION
# =============================================================================

try {
    Write-Host "🚀 Starting Microsoft Graph Data Extraction..." -ForegroundColor Cyan
    Write-Host "📁 Export location: $ExportPath" -ForegroundColor Cyan
    
    Invoke-GraphDataExtraction
    
    Write-Host "`n🎉 Script completed!" -ForegroundColor Green
    Write-Host "📋 Microsoft Graph data extraction process finished" -ForegroundColor Green
    Write-Host "💡 Check the summary above for any modules that may have failed due to permissions" -ForegroundColor Cyan
}
catch {
    Write-Host "`n💥 Script execution encountered errors!" -ForegroundColor Red
    Write-Host "Error: $_" -ForegroundColor Red
    Write-Host "`n💡 Some data may have been successfully extracted. Check your export directory: $ExportPath" -ForegroundColor Cyan
    exit 1
}