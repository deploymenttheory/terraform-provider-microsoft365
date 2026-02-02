<#
.SYNOPSIS
    Updates the tenant default app management policy with comprehensive field support.

.DESCRIPTION
    Updates the tenant-wide default app management policy (ID: defaultAppManagementPolicy).
    Supports all policy fields and can accept configuration via JSON file or individual parameters.
    
    The policy controls:
    - Password and key credential restrictions for applications
    - Password and key credential restrictions for service principals
    - Identifier URI restrictions
    - Federated identity credential restrictions
    - Redirect URI restrictions
    
    Use -JsonConfigPath for complex configurations or individual parameters for simple updates.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER JsonConfigPath
    Path to a JSON file containing the complete or partial policy configuration.
    See examples folder for sample JSON configurations.

.PARAMETER IsEnabled
    Enable or disable the policy. Default: true.

.PARAMETER DisplayName
    Custom display name for the policy.

.PARAMETER Description
    Custom description for the policy.

.PARAMETER ApplicationPasswordRestrictions
    Array of password credential restrictions for applications.
    Each restriction should be a hashtable with:
    - restrictionType (required): passwordAddition, passwordLifetime, symmetricKeyLifetime, customPasswordAddition
    - maxLifetime (required for lifetime types): ISO 8601 duration (e.g., "P90D")
    - restrictForAppsCreatedAfterDateTime (required): ISO 8601 datetime
    - state (optional): enabled, disabled
    - excludeActors (optional): array of actor exemptions

.PARAMETER ApplicationKeyRestrictions
    Array of key credential restrictions for applications.
    Each restriction should be a hashtable with:
    - restrictionType (required): asymmetricKeyLifetime, trustedCertificateAuthority
    - maxLifetime (required for lifetime): ISO 8601 duration
    - restrictForAppsCreatedAfterDateTime (required): ISO 8601 datetime
    - certificateBasedApplicationConfigurationIds (optional): array of config IDs
    - state (optional): enabled, disabled
    - excludeActors (optional): array of actor exemptions

.PARAMETER ServicePrincipalPasswordRestrictions
    Array of password credential restrictions for service principals (same format as application restrictions).

.PARAMETER ServicePrincipalKeyRestrictions
    Array of key credential restrictions for service principals (same format as application restrictions).

.PARAMETER MergeWithExisting
    If true, merges the provided configuration with the existing policy.
    If false (default), replaces the policy with the provided configuration.

.PARAMETER ExportToJson
    Export the updated policy to a JSON file in the output directory.

.PARAMETER ValidateOnly
    Only validate the configuration without applying changes.

.EXAMPLE
    # Update using JSON configuration file
    .\Update-TenantAppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -JsonConfigPath "./config/policy-config.json"

.EXAMPLE
    # Simple update: Set 90-day password lifetime for applications
    .\Update-TenantAppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationPasswordRestrictions @(
            @{
                restrictionType = "passwordLifetime"
                maxLifetime = "P90D"
                restrictForAppsCreatedAfterDateTime = "2024-01-01T00:00:00Z"
            }
        )

.EXAMPLE
    # Comprehensive update with all restriction types
    $appPwdRestrictions = @(
        @{
            restrictionType = "passwordLifetime"
            maxLifetime = "P1095D"
            restrictForAppsCreatedAfterDateTime = "2024-01-01T00:00:00Z"
            state = "enabled"
        }
    )
    
    $appKeyRestrictions = @(
        @{
            restrictionType = "asymmetricKeyLifetime"
            maxLifetime = "P1095D"
            restrictForAppsCreatedAfterDateTime = "2024-01-01T00:00:00Z"
            state = "enabled"
        }
    )
    
    .\Update-TenantAppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationPasswordRestrictions $appPwdRestrictions `
        -ApplicationKeyRestrictions $appKeyRestrictions `
        -ServicePrincipalPasswordRestrictions $appPwdRestrictions `
        -ServicePrincipalKeyRestrictions $appKeyRestrictions `
        -ExportToJson

.EXAMPLE
    # Validate configuration without applying
    .\Update-TenantAppManagementPolicy.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -JsonConfigPath "./config/policy-config.json" `
        -ValidateOnly

.NOTES
    Author: Deployment Theory
    Version: 2.0
    Requires: Microsoft.Graph.Authentication module
    API Reference: https://learn.microsoft.com/en-us/graph/api/tenantappmanagementpolicy-update?view=graph-rest-beta
    
    Valid restriction types:
    Password Credentials:
    - passwordAddition: Block new password creation (maxLifetime must be null)
    - passwordLifetime: Restrict password lifetime (maxLifetime required)
    - symmetricKeyLifetime: Restrict symmetric key lifetime (maxLifetime required)
    - customPasswordAddition: Block custom passwords (maxLifetime must be null)
    
    Key Credentials:
    - asymmetricKeyLifetime: Restrict certificate/key lifetime (maxLifetime required)
    - trustedCertificateAuthority: Restrict to specific CAs
    
    Duration Format (ISO 8601):
    - P90D = 90 days
    - P365D = 365 days (1 year)
    - P730D = 730 days (2 years)
    - P1095D = 1095 days (3 years)
#>

[CmdletBinding(DefaultParameterSetName='Parameters')]
param (
    [Parameter(Mandatory=$true, ParameterSetName='Parameters')]
    [Parameter(Mandatory=$true, ParameterSetName='JsonConfig')]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true, ParameterSetName='Parameters')]
    [Parameter(Mandatory=$true, ParameterSetName='JsonConfig')]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true, ParameterSetName='Parameters')]
    [Parameter(Mandatory=$true, ParameterSetName='JsonConfig')]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true, ParameterSetName='JsonConfig')]
    [ValidateScript({
        if (-not (Test-Path $_)) {
            throw "JSON configuration file not found: $_"
        }
        if (-not ($_ -match '\.json$')) {
            throw "File must have .json extension: $_"
        }
        return $true
    })]
    [string]$JsonConfigPath,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [bool]$IsEnabled = $true,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [string]$DisplayName,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [string]$Description,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [array]$ApplicationPasswordRestrictions,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [array]$ApplicationKeyRestrictions,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [array]$ServicePrincipalPasswordRestrictions,
    
    [Parameter(Mandatory=$false, ParameterSetName='Parameters')]
    [array]$ServicePrincipalKeyRestrictions,
    
    [Parameter(Mandatory=$false)]
    [bool]$MergeWithExisting = $false,
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false,
    
    [Parameter(Mandatory=$false)]
    [switch]$ValidateOnly
)

Import-Module Microsoft.Graph.Authentication

#region Helper Functions

function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Results,
        
        [Parameter(Mandatory=$false)]
        [string]$Suffix = "Update"
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "TenantAppManagementPolicy_${Suffix}_${timestamp}.json"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        $Results | ConvertTo-Json -Depth 20 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported results to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting results to JSON: $_" -ForegroundColor Red
        return $null
    }
}

function Read-JsonConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [string]$FilePath
    )
    
    try {
        Write-Host "📄 Reading JSON configuration from: $FilePath" -ForegroundColor Cyan
        
        $jsonContent = Get-Content -Path $FilePath -Raw -ErrorAction Stop
        $config = $jsonContent | ConvertFrom-Json -ErrorAction Stop
        
        Write-Host "✅ JSON configuration loaded successfully" -ForegroundColor Green
        return $config
    }
    catch {
        Write-Host "❌ Error reading JSON configuration: $_" -ForegroundColor Red
        throw
    }
}

function ConvertTo-PolicyRequestBody {
    param (
        [Parameter(Mandatory=$false)]
        $JsonConfig,
        
        [Parameter(Mandatory=$false)]
        $ParameterConfig
    )
    
    $requestBody = @{}
    
    if ($JsonConfig) {
        # Convert JSON object to hashtable
        Write-Host "📝 Building policy from JSON configuration..." -ForegroundColor Cyan
        
        if ($JsonConfig.isEnabled -ne $null) {
            $requestBody['isEnabled'] = $JsonConfig.isEnabled
        }
        
        if ($JsonConfig.displayName) {
            $requestBody['displayName'] = $JsonConfig.displayName
        }
        
        if ($JsonConfig.description) {
            $requestBody['description'] = $JsonConfig.description
        }
        
        # Application restrictions
        if ($JsonConfig.applicationRestrictions) {
            $requestBody['applicationRestrictions'] = @{}
            
            if ($JsonConfig.applicationRestrictions.passwordCredentials) {
                $requestBody['applicationRestrictions']['passwordCredentials'] = @()
                foreach ($cred in $JsonConfig.applicationRestrictions.passwordCredentials) {
                    $requestBody['applicationRestrictions']['passwordCredentials'] += ConvertTo-Hashtable $cred
                }
            }
            
            if ($JsonConfig.applicationRestrictions.keyCredentials) {
                $requestBody['applicationRestrictions']['keyCredentials'] = @()
                foreach ($cred in $JsonConfig.applicationRestrictions.keyCredentials) {
                    $requestBody['applicationRestrictions']['keyCredentials'] += ConvertTo-Hashtable $cred
                }
            }
            
            if ($JsonConfig.applicationRestrictions.identifierUris) {
                $requestBody['applicationRestrictions']['identifierUris'] = ConvertTo-Hashtable $JsonConfig.applicationRestrictions.identifierUris
            }
            
            if ($JsonConfig.applicationRestrictions.redirectUris) {
                $requestBody['applicationRestrictions']['redirectUris'] = ConvertTo-Hashtable $JsonConfig.applicationRestrictions.redirectUris
            }
            
            if ($JsonConfig.applicationRestrictions.federatedIdentityCredentials) {
                $requestBody['applicationRestrictions']['federatedIdentityCredentials'] = ConvertTo-Hashtable $JsonConfig.applicationRestrictions.federatedIdentityCredentials
            }
            
            if ($JsonConfig.applicationRestrictions.audiences) {
                $requestBody['applicationRestrictions']['audiences'] = @($JsonConfig.applicationRestrictions.audiences)
            }
        }
        
        # Service principal restrictions
        if ($JsonConfig.servicePrincipalRestrictions) {
            $requestBody['servicePrincipalRestrictions'] = @{}
            
            if ($JsonConfig.servicePrincipalRestrictions.passwordCredentials) {
                $requestBody['servicePrincipalRestrictions']['passwordCredentials'] = @()
                foreach ($cred in $JsonConfig.servicePrincipalRestrictions.passwordCredentials) {
                    $requestBody['servicePrincipalRestrictions']['passwordCredentials'] += ConvertTo-Hashtable $cred
                }
            }
            
            if ($JsonConfig.servicePrincipalRestrictions.keyCredentials) {
                $requestBody['servicePrincipalRestrictions']['keyCredentials'] = @()
                foreach ($cred in $JsonConfig.servicePrincipalRestrictions.keyCredentials) {
                    $requestBody['servicePrincipalRestrictions']['keyCredentials'] += ConvertTo-Hashtable $cred
                }
            }
        }
    }
    elseif ($ParameterConfig) {
        # Build from individual parameters
        Write-Host "📝 Building policy from parameters..." -ForegroundColor Cyan
        
        $requestBody['isEnabled'] = $ParameterConfig.IsEnabled
        
        if ($ParameterConfig.DisplayName) {
            $requestBody['displayName'] = $ParameterConfig.DisplayName
        }
        
        if ($ParameterConfig.Description) {
            $requestBody['description'] = $ParameterConfig.Description
        }
        
        # Application restrictions
        if ($ParameterConfig.ApplicationPasswordRestrictions -or $ParameterConfig.ApplicationKeyRestrictions) {
            $requestBody['applicationRestrictions'] = @{}
            
            if ($ParameterConfig.ApplicationPasswordRestrictions) {
                $requestBody['applicationRestrictions']['passwordCredentials'] = $ParameterConfig.ApplicationPasswordRestrictions
                Write-Host "   • Application password restrictions: $($ParameterConfig.ApplicationPasswordRestrictions.Count)" -ForegroundColor Gray
            }
            
            if ($ParameterConfig.ApplicationKeyRestrictions) {
                $requestBody['applicationRestrictions']['keyCredentials'] = $ParameterConfig.ApplicationKeyRestrictions
                Write-Host "   • Application key restrictions: $($ParameterConfig.ApplicationKeyRestrictions.Count)" -ForegroundColor Gray
            }
        }
        
        # Service principal restrictions
        if ($ParameterConfig.ServicePrincipalPasswordRestrictions -or $ParameterConfig.ServicePrincipalKeyRestrictions) {
            $requestBody['servicePrincipalRestrictions'] = @{}
            
            if ($ParameterConfig.ServicePrincipalPasswordRestrictions) {
                $requestBody['servicePrincipalRestrictions']['passwordCredentials'] = $ParameterConfig.ServicePrincipalPasswordRestrictions
                Write-Host "   • Service principal password restrictions: $($ParameterConfig.ServicePrincipalPasswordRestrictions.Count)" -ForegroundColor Gray
            }
            
            if ($ParameterConfig.ServicePrincipalKeyRestrictions) {
                $requestBody['servicePrincipalRestrictions']['keyCredentials'] = $ParameterConfig.ServicePrincipalKeyRestrictions
                Write-Host "   • Service principal key restrictions: $($ParameterConfig.ServicePrincipalKeyRestrictions.Count)" -ForegroundColor Gray
            }
        }
    }
    
    return $requestBody
}

function ConvertTo-Hashtable {
    param (
        [Parameter(Mandatory=$true)]
        $InputObject
    )
    
    if ($InputObject -is [hashtable]) {
        return $InputObject
    }
    
    $hashtable = @{}
    
    $InputObject.PSObject.Properties | ForEach-Object {
        $value = $_.Value
        
        if ($value -is [PSCustomObject]) {
            $value = ConvertTo-Hashtable $value
        }
        elseif ($value -is [System.Collections.IEnumerable] -and $value -isnot [string]) {
            $arrayValues = @()
            foreach ($item in $value) {
                if ($item -is [PSCustomObject]) {
                    $arrayValues += ConvertTo-Hashtable $item
                }
                else {
                    $arrayValues += $item
                }
            }
            $value = $arrayValues
        }
        
        $hashtable[$_.Name] = $value
    }
    
    return $hashtable
}

function Test-PolicyConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [hashtable]$RequestBody
    )
    
    Write-Host "🔍 Validating policy configuration..." -ForegroundColor Cyan
    $isValid = $true
    $validationErrors = @()
    
    # Validate password credential restrictions
    if ($RequestBody.applicationRestrictions.passwordCredentials) {
        foreach ($cred in $RequestBody.applicationRestrictions.passwordCredentials) {
            if (-not $cred.restrictionType) {
                $validationErrors += "Application password restriction missing restrictionType"
                $isValid = $false
            }
            
            if ($cred.restrictionType -in @('passwordLifetime', 'symmetricKeyLifetime')) {
                if (-not $cred.maxLifetime) {
                    $validationErrors += "Restriction type '$($cred.restrictionType)' requires maxLifetime"
                    $isValid = $false
                }
            }
            
            if (-not $cred.restrictForAppsCreatedAfterDateTime) {
                $validationErrors += "Application password restriction missing restrictForAppsCreatedAfterDateTime"
                $isValid = $false
            }
        }
    }
    
    # Validate key credential restrictions
    if ($RequestBody.applicationRestrictions.keyCredentials) {
        foreach ($cred in $RequestBody.applicationRestrictions.keyCredentials) {
            if (-not $cred.restrictionType) {
                $validationErrors += "Application key restriction missing restrictionType"
                $isValid = $false
            }
            
            if ($cred.restrictionType -eq 'asymmetricKeyLifetime') {
                if (-not $cred.maxLifetime) {
                    $validationErrors += "asymmetricKeyLifetime requires maxLifetime"
                    $isValid = $false
                }
            }
        }
    }
    
    # Validate service principal restrictions (same rules)
    if ($RequestBody.servicePrincipalRestrictions.passwordCredentials) {
        foreach ($cred in $RequestBody.servicePrincipalRestrictions.passwordCredentials) {
            if (-not $cred.restrictionType) {
                $validationErrors += "Service principal password restriction missing restrictionType"
                $isValid = $false
            }
        }
    }
    
    if ($validationErrors.Count -gt 0) {
        Write-Host ""
        Write-Host "❌ Validation failed:" -ForegroundColor Red
        foreach ($error in $validationErrors) {
            Write-Host "   • $error" -ForegroundColor Red
        }
        Write-Host ""
    }
    else {
        Write-Host "✅ Configuration is valid" -ForegroundColor Green
    }
    
    return $isValid
}

function Merge-PolicyConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [hashtable]$ExistingPolicy,
        
        [Parameter(Mandatory=$true)]
        [hashtable]$NewConfig
    )
    
    Write-Host "🔄 Merging with existing policy..." -ForegroundColor Cyan
    
    $merged = @{}
    
    # Copy all existing properties
    foreach ($key in $ExistingPolicy.Keys) {
        if ($key -notin @('@odata.context', 'id')) {
            $merged[$key] = $ExistingPolicy[$key]
        }
    }
    
    # Override with new config
    foreach ($key in $NewConfig.Keys) {
        $merged[$key] = $NewConfig[$key]
    }
    
    return $merged
}

function Format-PolicySummary {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    
    Write-Host ""
    Write-Host "📋 Policy Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    Write-Host "   Policy ID: $($Policy.id)" -ForegroundColor White
    Write-Host "   Display Name: $($Policy.displayName)" -ForegroundColor White
    Write-Host "   Description: $($Policy.description)" -ForegroundColor White
    Write-Host "   Is Enabled: " -NoNewline -ForegroundColor White
    Write-Host "$($Policy.isEnabled)" -ForegroundColor $(if ($Policy.isEnabled) { "Green" } else { "Yellow" })
    Write-Host ""
    
    # Application restrictions
    $appPwdCount = if ($Policy.applicationRestrictions.passwordCredentials) { $Policy.applicationRestrictions.passwordCredentials.Count } else { 0 }
    $appKeyCount = if ($Policy.applicationRestrictions.keyCredentials) { $Policy.applicationRestrictions.keyCredentials.Count } else { 0 }
    
    Write-Host "   📱 Application Restrictions:" -ForegroundColor White
    Write-Host "      • Password Credential Restrictions: $appPwdCount" -ForegroundColor Gray
    Write-Host "      • Key Credential Restrictions: $appKeyCount" -ForegroundColor Gray
    
    Write-Host ""
    
    # Service principal restrictions
    $spPwdCount = if ($Policy.servicePrincipalRestrictions.passwordCredentials) { $Policy.servicePrincipalRestrictions.passwordCredentials.Count } else { 0 }
    $spKeyCount = if ($Policy.servicePrincipalRestrictions.keyCredentials) { $Policy.servicePrincipalRestrictions.keyCredentials.Count } else { 0 }
    
    Write-Host "   ⚙️  Service Principal Restrictions:" -ForegroundColor White
    Write-Host "      • Password Credential Restrictions: $spPwdCount" -ForegroundColor Gray
    Write-Host "      • Key Credential Restrictions: $spKeyCount" -ForegroundColor Gray
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

#endregion Helper Functions

#region Main Script Execution

try {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   🔧 Update Tenant App Management Policy" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    $uri = "https://graph.microsoft.com/beta/policies/defaultAppManagementPolicy"
    
    # Get existing policy if merging
    $existingPolicy = $null
    if ($MergeWithExisting) {
        Write-Host "🔍 Retrieving existing policy for merge..." -ForegroundColor Cyan
        $existingPolicy = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "✅ Retrieved existing policy" -ForegroundColor Green
        Write-Host ""
    }
    
    # Build request body based on input method
    $requestBody = $null
    
    if ($PSCmdlet.ParameterSetName -eq 'JsonConfig') {
        $jsonConfig = Read-JsonConfiguration -FilePath $JsonConfigPath
        $requestBody = ConvertTo-PolicyRequestBody -JsonConfig $jsonConfig
    }
    else {
        $paramConfig = @{
            IsEnabled = $IsEnabled
            DisplayName = $DisplayName
            Description = $Description
            ApplicationPasswordRestrictions = $ApplicationPasswordRestrictions
            ApplicationKeyRestrictions = $ApplicationKeyRestrictions
            ServicePrincipalPasswordRestrictions = $ServicePrincipalPasswordRestrictions
            ServicePrincipalKeyRestrictions = $ServicePrincipalKeyRestrictions
        }
        $requestBody = ConvertTo-PolicyRequestBody -ParameterConfig $paramConfig
    }
    
    # Merge with existing if requested
    if ($MergeWithExisting -and $existingPolicy) {
        $requestBody = Merge-PolicyConfiguration -ExistingPolicy $existingPolicy -NewConfig $requestBody
    }
    
    Write-Host ""
    
    # Validate configuration
    $isValid = Test-PolicyConfiguration -RequestBody $requestBody
    
    if (-not $isValid) {
        throw "Policy configuration validation failed. Please fix the errors and try again."
    }
    
    Write-Host ""
    
    # Display request body
    Write-Host "📤 Request Body" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    Write-Host ($requestBody | ConvertTo-Json -Depth 20) -ForegroundColor Gray
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    if ($ValidateOnly) {
        Write-Host "✅ Validation-only mode: Configuration is valid, no changes applied" -ForegroundColor Green
        Write-Host ""
    }
    else {
        # Apply the policy update
        Write-Host "🔧 Updating policy..." -ForegroundColor Cyan
        
        $response = Invoke-MgGraphRequest -Method PATCH -Uri $uri -Body ($requestBody | ConvertTo-Json -Depth 20) -ContentType "application/json"
        
        Write-Host "✅ Policy updated successfully!" -ForegroundColor Green
        Write-Host "   Note: PATCH returns 204 No Content" -ForegroundColor Gray
        Write-Host ""
        
        # Wait for eventual consistency
        Write-Host "⏱️  Waiting 5 seconds for eventual consistency..." -ForegroundColor Gray
        Start-Sleep -Seconds 5
        
        # Retrieve updated policy
        Write-Host "🔍 Retrieving updated policy..." -ForegroundColor Cyan
        
        $updatedPolicy = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        Write-Host "✅ Retrieved updated policy" -ForegroundColor Green
        
        # Display summary
        Format-PolicySummary -Policy $updatedPolicy
        
        # Display full JSON
        Write-Host "📄 Full Policy JSON" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        Write-Host ($updatedPolicy | ConvertTo-Json -Depth 20) -ForegroundColor Gray
        Write-Host ""
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        # Export results if requested
        if ($ExportToJson) {
            $exportObject = @{
                Timestamp = Get-Date -Format "o"
                Operation = "Update"
                RequestBody = $requestBody
                UpdatedPolicy = $updatedPolicy
            }
            
            Export-ResultsToJson -Results $exportObject
        }
        
        Write-Host "🎉 Policy update completed successfully!" -ForegroundColor Green
        Write-Host ""
    }
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   Status Code: $statusCode" -ForegroundColor Gray
        
        if ($_.ErrorDetails) {
            try {
                $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
                Write-Host "   Error Code: $($errorDetails.error.code)" -ForegroundColor Gray
                Write-Host "   Error Message: $($errorDetails.error.message)" -ForegroundColor Gray
                
                if ($errorDetails.error.innerError) {
                    Write-Host "   Request ID: $($errorDetails.error.innerError.'request-id')" -ForegroundColor Gray
                    Write-Host "   Date: $($errorDetails.error.innerError.date)" -ForegroundColor Gray
                }
            }
            catch {
                Write-Host "   Error Details: $($_.ErrorDetails.Message)" -ForegroundColor Gray
            }
        }
    }
    
    Write-Host ""
    Write-Host "Stack Trace:" -ForegroundColor Gray
    Write-Host $_.ScriptStackTrace -ForegroundColor Gray
    Write-Host ""
    
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph 2>$null
    Write-Host "✅ Disconnected" -ForegroundColor Green
    Write-Host ""
}

#endregion Main Script Execution
