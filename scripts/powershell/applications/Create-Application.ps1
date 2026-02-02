<#
.SYNOPSIS
    Creates a Microsoft Entra application with configurable properties.

.DESCRIPTION
    Creates a new Microsoft Entra application with specified display name and optional
    properties such as description, sign-in audience, and identifier URIs. Useful for
    testing application creation and for setting up test environments.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER DisplayName
    The display name for the application. Required.

.PARAMETER Description
    Description for the application.

.PARAMETER SignInAudience
    The sign-in audience for the application. Possible values:
    - AzureADMyOrg (Single tenant)
    - AzureADMultipleOrgs (Multi-tenant)
    - AzureADandPersonalMicrosoftAccount (Multi-tenant + personal Microsoft accounts)
    - PersonalMicrosoftAccount (Personal Microsoft accounts only)
    Defaults to AzureADMyOrg.

.PARAMETER IdentifierUris
    Array of identifier URIs for the application.

.PARAMETER Notes
    Internal notes about the application.

.PARAMETER HardDelete
    If true, the application will be permanently deleted when removed (not soft-deleted).

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Create a basic application
    .\Create-Application.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -DisplayName "Test Application"

.EXAMPLE
    # Create application with full configuration
    .\Create-Application.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -DisplayName "My Enterprise App" `
        -Description "Application for testing service principals" `
        -SignInAudience "AzureADMyOrg" `
        -Notes "Created for acceptance testing" `
        -ExportToJson $true

.EXAMPLE
    # Create multi-tenant application
    .\Create-Application.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -DisplayName "Multi-Tenant App" `
        -SignInAudience "AzureADMultipleOrgs" `
        -IdentifierUris @("api://my-app")

.NOTES
    Author: Deployment Theory
    Requires: Microsoft.Graph.Authentication module
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    [string]$DisplayName,
    
    [Parameter(Mandatory=$false)]
    [string]$Description,
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("AzureADMyOrg", "AzureADMultipleOrgs", "AzureADandPersonalMicrosoftAccount", "PersonalMicrosoftAccount")]
    [string]$SignInAudience = "AzureADMyOrg",
    
    [Parameter(Mandatory=$false)]
    [string[]]$IdentifierUris,
    
    [Parameter(Mandatory=$false)]
    [string]$Notes,
    
    [Parameter(Mandatory=$false)]
    [bool]$HardDelete = $false,
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

#region Helper Functions

function Format-AttributeValue {
    param (
        [Parameter(Mandatory=$false)]
        $Value
    )
    
    if ($null -eq $Value) {
        return "(null)"
    }
    elseif ($Value -is [array]) {
        if ($Value.Count -eq 0) {
            return "(empty array)"
        }
        return ($Value -join ", ")
    }
    elseif ($Value -is [hashtable] -or $Value -is [System.Collections.IDictionary]) {
        return "(complex object)"
    }
    elseif ($Value -is [bool]) {
        return $Value.ToString().ToLower()
    }
    elseif ($Value -is [string] -and [string]::IsNullOrEmpty($Value)) {
        return "(empty)"
    }
    else {
        return $Value.ToString()
    }
}

function Export-ResultsToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Results,
        
        [Parameter(Mandatory=$true)]
        [string]$ApplicationId
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "Application_${ApplicationId}_${timestamp}.json"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        $Results | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported results to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting results to JSON: $_" -ForegroundColor Red
        return $null
    }
}

#endregion Helper Functions

#region Main Script Execution

try {
    Write-Host ""
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host "   ➕ Create Application" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Build request body
    Write-Host "📝 Building application configuration..." -ForegroundColor Cyan
    
    $requestBody = @{
        displayName = $DisplayName
        signInAudience = $SignInAudience
    }
    
    Write-Host "   • Display Name: $DisplayName" -ForegroundColor Gray
    Write-Host "   • Sign-In Audience: $SignInAudience" -ForegroundColor Gray
    
    # Add optional properties
    if (-not [string]::IsNullOrEmpty($Description)) {
        $requestBody['description'] = $Description
        Write-Host "   • Description: $Description" -ForegroundColor Gray
    }
    
    if ($IdentifierUris -and $IdentifierUris.Count -gt 0) {
        $requestBody['identifierUris'] = $IdentifierUris
        Write-Host "   • Identifier URIs: $($IdentifierUris -join ', ')" -ForegroundColor Gray
    }
    
    if (-not [string]::IsNullOrEmpty($Notes)) {
        $requestBody['notes'] = $Notes
        Write-Host "   • Notes: $Notes" -ForegroundColor Gray
    }
    
    Write-Host ""
    
    # Create application
    Write-Host "➕ Creating application..." -ForegroundColor Cyan
    
    $createUri = "https://graph.microsoft.com/beta/applications"
    $application = Invoke-MgGraphRequest -Method POST -Uri $createUri -Body ($requestBody | ConvertTo-Json -Depth 10)
    
    Write-Host "✅ Application created successfully!" -ForegroundColor Green
    Write-Host "   Object ID: $($application.id)" -ForegroundColor Gray
    Write-Host "   App ID: $($application.appId)" -ForegroundColor Gray
    Write-Host ""
    
    # Allow time for eventual consistency
    Write-Host "⏱️  Waiting 15 seconds for eventual consistency..." -ForegroundColor Gray
    Start-Sleep -Seconds 15
    
    # Retrieve full application details
    Write-Host "🔍 Retrieving full application details..." -ForegroundColor Cyan
    
    $detailUri = "https://graph.microsoft.com/beta/applications/$($application.id)"
    $fullApplication = Invoke-MgGraphRequest -Method GET -Uri $detailUri
    
    Write-Host "✅ Retrieved full details" -ForegroundColor Green
    Write-Host ""
    
    # Display application details
    Write-Host "📋 Application Details" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $keyAttributes = @(
        @{ Name = "id"; Label = "Object ID" },
        @{ Name = "appId"; Label = "App ID (Client ID)" },
        @{ Name = "displayName"; Label = "Display Name" },
        @{ Name = "description"; Label = "Description" },
        @{ Name = "signInAudience"; Label = "Sign-In Audience" },
        @{ Name = "publisherDomain"; Label = "Publisher Domain" },
        @{ Name = "createdDateTime"; Label = "Created Date Time" },
        @{ Name = "notes"; Label = "Notes" },
        @{ Name = "isDeviceOnlyAuthSupported"; Label = "Device Only Auth Supported" },
        @{ Name = "isFallbackPublicClient"; Label = "Fallback Public Client" },
        @{ Name = "disabledByMicrosoftStatus"; Label = "Disabled By Microsoft" }
    )
    
    foreach ($attr in $keyAttributes) {
        $value = $fullApplication[$attr.Name]
        $formattedValue = Format-AttributeValue -Value $value
        Write-Host "   $($attr.Label): " -NoNewline -ForegroundColor White
        Write-Host "$formattedValue" -ForegroundColor Green
    }
    
    Write-Host ""
    
    # Display collections
    if ($fullApplication.identifierUris -and $fullApplication.identifierUris.Count -gt 0) {
        Write-Host "   Identifier URIs:" -ForegroundColor White
        foreach ($uri in $fullApplication.identifierUris) {
            Write-Host "     • $uri" -ForegroundColor Green
        }
    } else {
        Write-Host "   Identifier URIs: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
    }
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display app roles if any
    if ($fullApplication.appRoles -and $fullApplication.appRoles.Count -gt 0) {
        Write-Host "👥 App Roles ($($fullApplication.appRoles.Count))" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($role in $fullApplication.appRoles) {
            Write-Host "   📌 $($role.displayName)" -ForegroundColor Yellow
            Write-Host "      ID: $($role.id)" -ForegroundColor Gray
            Write-Host "      Value: $($role.value)" -ForegroundColor Gray
            Write-Host "      Is Enabled: $($role.isEnabled)" -ForegroundColor $(if ($role.isEnabled) { "Green" } else { "Red" })
            Write-Host ""
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Get owners
    Write-Host "🔍 Retrieving application owners..." -ForegroundColor Cyan
    
    $ownersUri = "https://graph.microsoft.com/beta/applications/$($fullApplication.id)/owners"
    $owners = Invoke-MgGraphRequest -Method GET -Uri $ownersUri
    
    if ($owners.value -and $owners.value.Count -gt 0) {
        Write-Host "✅ Found $($owners.value.Count) owner(s)" -ForegroundColor Green
        Write-Host ""
        
        Write-Host "👤 Application Owners" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($owner in $owners.value) {
            Write-Host "   • $($owner.displayName)" -ForegroundColor Green
            if ($owner.userPrincipalName) {
                Write-Host "     User Principal Name: $($owner.userPrincipalName)" -ForegroundColor Gray
            }
            Write-Host "     Object ID: $($owner.id)" -ForegroundColor Gray
            Write-Host ""
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    else {
        Write-Host "⚠️  No owners found" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Display summary
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Display Name: $($fullApplication.displayName)" -ForegroundColor White
    Write-Host "   Object ID: $($fullApplication.id)" -ForegroundColor White
    Write-Host "   App ID (Client ID): $($fullApplication.appId)" -ForegroundColor White
    Write-Host "   Sign-In Audience: $($fullApplication.signInAudience)" -ForegroundColor White
    
    $appRolesCount = if ($fullApplication.appRoles) { $fullApplication.appRoles.Count } else { 0 }
    Write-Host "   App Roles: $appRolesCount" -ForegroundColor White
    
    $ownersCount = if ($owners.value) { $owners.value.Count } else { 0 }
    Write-Host "   Owners: $ownersCount" -ForegroundColor White
    
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            Operation = "Create"
            RequestBody = $requestBody
            Application = $fullApplication
            Owners = $owners.value
            Summary = @{
                DisplayName = $fullApplication.displayName
                ObjectId = $fullApplication.id
                AppId = $fullApplication.appId
                SignInAudience = $fullApplication.signInAudience
                AppRolesCount = $appRolesCount
                OwnersCount = $ownersCount
            }
        }
        
        Export-ResultsToJson -Results $exportObject -ApplicationId $fullApplication.id
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "💡 Use these values for creating a service principal:" -ForegroundColor Cyan
    Write-Host "   App ID: $($fullApplication.appId)" -ForegroundColor Yellow
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   Status Code: $statusCode" -ForegroundColor Gray
        
        if ($_.ErrorDetails) {
            $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
            Write-Host "   Error Code: $($errorDetails.error.code)" -ForegroundColor Gray
            Write-Host "   Error Message: $($errorDetails.error.message)" -ForegroundColor Gray
        }
    }
    
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
