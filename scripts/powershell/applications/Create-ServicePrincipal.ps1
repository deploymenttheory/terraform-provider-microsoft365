<#
.SYNOPSIS
    Creates a service principal for an existing application with configurable properties.

.DESCRIPTION
    Creates a service principal from an existing application's App ID (Client ID),
    optionally setting account enabled status, app role assignment requirements, tags,
    and other properties. Useful for testing service principal creation and updates
    with various configurations.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ApplicationId
    The App ID (Client ID) of the application for which to create the service principal.

.PARAMETER AccountEnabled
    Whether the service principal account is enabled. Defaults to true.

.PARAMETER AppRoleAssignmentRequired
    Whether users must be assigned to the service principal. Defaults to false.

.PARAMETER Tags
    Array of tags to apply to the service principal. Common values include:
    - "HideApp" - Hide from My Apps
    - "WindowsAzureActiveDirectoryIntegratedApp" - Marks as integrated app

.PARAMETER Description
    Description for the service principal.

.PARAMETER Notes
    Internal notes about the service principal.

.PARAMETER PreferredSingleSignOnMode
    Preferred single sign-on mode. Possible values: password, saml, notSupported, oidc.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Create a basic service principal
    .\Create-ServicePrincipal.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationId "12345678-1234-1234-1234-123456789012"

.EXAMPLE
    # Create service principal with tags and app role assignment required
    .\Create-ServicePrincipal.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationId "12345678-1234-1234-1234-123456789012" `
        -AppRoleAssignmentRequired $true `
        -Tags @("HideApp", "WindowsAzureActiveDirectoryIntegratedApp") `
        -ExportToJson $true

.EXAMPLE
    # Create service principal with full configuration
    .\Create-ServicePrincipal.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationId "12345678-1234-1234-1234-123456789012" `
        -AccountEnabled $true `
        -AppRoleAssignmentRequired $true `
        -Description "Test service principal" `
        -Notes "Created for acceptance testing" `
        -Tags @("HideApp")

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
    [string]$ApplicationId,
    
    [Parameter(Mandatory=$false)]
    [bool]$AccountEnabled = $true,
    
    [Parameter(Mandatory=$false)]
    [bool]$AppRoleAssignmentRequired = $false,
    
    [Parameter(Mandatory=$false)]
    [string[]]$Tags,
    
    [Parameter(Mandatory=$false)]
    [string]$Description,
    
    [Parameter(Mandatory=$false)]
    [string]$Notes,
    
    [Parameter(Mandatory=$false)]
    [ValidateSet("password", "saml", "notSupported", "oidc", "")]
    [string]$PreferredSingleSignOnMode,
    
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
        [string]$ServicePrincipalId
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "ServicePrincipal_${ServicePrincipalId}_${timestamp}.json"
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
    Write-Host "   ➕ Create Service Principal" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Verify the application exists
    Write-Host "🔍 Verifying application exists..." -ForegroundColor Cyan
    
    $appSearchUri = "https://graph.microsoft.com/beta/applications?`$filter=appId eq '$ApplicationId'"
    $appSearchResult = Invoke-MgGraphRequest -Method GET -Uri $appSearchUri
    
    if (-not $appSearchResult.value -or $appSearchResult.value.Count -eq 0) {
        Write-Host "❌ Application with App ID '$ApplicationId' not found" -ForegroundColor Red
        Write-Host ""
        exit 1
    }
    
    $application = $appSearchResult.value[0]
    Write-Host "✅ Found application: $($application.displayName)" -ForegroundColor Green
    Write-Host ""
    
    # Build request body
    Write-Host "📝 Building service principal configuration..." -ForegroundColor Cyan
    
    $requestBody = @{
        appId = $ApplicationId
    }
    
    # Add optional properties
    if ($PSBoundParameters.ContainsKey('AccountEnabled')) {
        $requestBody['accountEnabled'] = $AccountEnabled
        Write-Host "   • Account Enabled: $AccountEnabled" -ForegroundColor Gray
    }
    
    if ($PSBoundParameters.ContainsKey('AppRoleAssignmentRequired')) {
        $requestBody['appRoleAssignmentRequired'] = $AppRoleAssignmentRequired
        Write-Host "   • App Role Assignment Required: $AppRoleAssignmentRequired" -ForegroundColor Gray
    }
    
    if ($Tags -and $Tags.Count -gt 0) {
        $requestBody['tags'] = $Tags
        Write-Host "   • Tags: $($Tags -join ', ')" -ForegroundColor Gray
    }
    
    if (-not [string]::IsNullOrEmpty($Description)) {
        $requestBody['description'] = $Description
        Write-Host "   • Description: $Description" -ForegroundColor Gray
    }
    
    if (-not [string]::IsNullOrEmpty($Notes)) {
        $requestBody['notes'] = $Notes
        Write-Host "   • Notes: $Notes" -ForegroundColor Gray
    }
    
    if (-not [string]::IsNullOrEmpty($PreferredSingleSignOnMode)) {
        $requestBody['preferredSingleSignOnMode'] = $PreferredSingleSignOnMode
        Write-Host "   • Preferred SSO Mode: $PreferredSingleSignOnMode" -ForegroundColor Gray
    }
    
    Write-Host ""
    
    # Create service principal
    Write-Host "➕ Creating service principal..." -ForegroundColor Cyan
    
    $createUri = "https://graph.microsoft.com/beta/servicePrincipals"
    $servicePrincipal = Invoke-MgGraphRequest -Method POST -Uri $createUri -Body ($requestBody | ConvertTo-Json -Depth 10)
    
    Write-Host "✅ Service principal created successfully!" -ForegroundColor Green
    Write-Host "   Object ID: $($servicePrincipal.id)" -ForegroundColor Gray
    Write-Host "   App ID: $($servicePrincipal.appId)" -ForegroundColor Gray
    Write-Host ""
    
    # Allow time for eventual consistency
    Write-Host "⏱️  Waiting 15 seconds for eventual consistency..." -ForegroundColor Gray
    Start-Sleep -Seconds 15
    
    # Retrieve full service principal details
    Write-Host "🔍 Retrieving full service principal details..." -ForegroundColor Cyan
    
    $detailUri = "https://graph.microsoft.com/beta/servicePrincipals/$($servicePrincipal.id)?`$select=id,appId,displayName,accountEnabled,appRoleAssignmentRequired,description,homepage,loginUrl,logoutUrl,notes,notificationEmailAddresses,preferredSingleSignOnMode,servicePrincipalNames,servicePrincipalType,signInAudience,tags"
    $fullServicePrincipal = Invoke-MgGraphRequest -Method GET -Uri $detailUri
    
    Write-Host "✅ Retrieved full details" -ForegroundColor Green
    Write-Host ""
    
    # Display service principal details
    Write-Host "📋 Service Principal Details" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    $keyAttributes = @(
        @{ Name = "id"; Label = "Object ID" },
        @{ Name = "appId"; Label = "App ID (Client ID)" },
        @{ Name = "displayName"; Label = "Display Name" },
        @{ Name = "servicePrincipalType"; Label = "Service Principal Type" },
        @{ Name = "accountEnabled"; Label = "Account Enabled" },
        @{ Name = "appRoleAssignmentRequired"; Label = "App Role Assignment Required" },
        @{ Name = "description"; Label = "Description" },
        @{ Name = "homepage"; Label = "Homepage" },
        @{ Name = "loginUrl"; Label = "Login URL" },
        @{ Name = "logoutUrl"; Label = "Logout URL" },
        @{ Name = "notes"; Label = "Notes" },
        @{ Name = "preferredSingleSignOnMode"; Label = "Preferred SSO Mode" },
        @{ Name = "signInAudience"; Label = "Sign-In Audience" }
    )
    
    foreach ($attr in $keyAttributes) {
        $value = $fullServicePrincipal[$attr.Name]
        $formattedValue = Format-AttributeValue -Value $value
        Write-Host "   $($attr.Label): " -NoNewline -ForegroundColor White
        Write-Host "$formattedValue" -ForegroundColor Green
    }
    
    Write-Host ""
    
    # Display collections
    if ($fullServicePrincipal.tags -and $fullServicePrincipal.tags.Count -gt 0) {
        Write-Host "   Tags:" -ForegroundColor White
        foreach ($tag in $fullServicePrincipal.tags) {
            Write-Host "     • $tag" -ForegroundColor Green
        }
    } else {
        Write-Host "   Tags: " -NoNewline -ForegroundColor White
        Write-Host "(none)" -ForegroundColor Gray
    }
    
    if ($fullServicePrincipal.servicePrincipalNames -and $fullServicePrincipal.servicePrincipalNames.Count -gt 0) {
        Write-Host "   Service Principal Names:" -ForegroundColor White
        foreach ($name in $fullServicePrincipal.servicePrincipalNames) {
            Write-Host "     • $name" -ForegroundColor Green
        }
    }
    
    if ($fullServicePrincipal.notificationEmailAddresses -and $fullServicePrincipal.notificationEmailAddresses.Count -gt 0) {
        Write-Host "   Notification Email Addresses:" -ForegroundColor White
        foreach ($email in $fullServicePrincipal.notificationEmailAddresses) {
            Write-Host "     • $email" -ForegroundColor Green
        }
    }
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display summary
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Display Name: $($fullServicePrincipal.displayName)" -ForegroundColor White
    Write-Host "   Object ID: $($fullServicePrincipal.id)" -ForegroundColor White
    Write-Host "   App ID: $($fullServicePrincipal.appId)" -ForegroundColor White
    Write-Host "   Type: $($fullServicePrincipal.servicePrincipalType)" -ForegroundColor White
    Write-Host "   Account Enabled: $(Format-AttributeValue -Value $fullServicePrincipal.accountEnabled)" -ForegroundColor White
    Write-Host "   App Role Assignment Required: $(Format-AttributeValue -Value $fullServicePrincipal.appRoleAssignmentRequired)" -ForegroundColor White
    
    $tagsCount = if ($fullServicePrincipal.tags) { $fullServicePrincipal.tags.Count } else { 0 }
    Write-Host "   Tags Count: $tagsCount" -ForegroundColor White
    
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            Operation = "Create"
            RequestBody = $requestBody
            ServicePrincipal = $fullServicePrincipal
            Summary = @{
                DisplayName = $fullServicePrincipal.displayName
                ObjectId = $fullServicePrincipal.id
                AppId = $fullServicePrincipal.appId
                ServicePrincipalType = $fullServicePrincipal.servicePrincipalType
                AccountEnabled = $fullServicePrincipal.accountEnabled
                AppRoleAssignmentRequired = $fullServicePrincipal.appRoleAssignmentRequired
                TagsCount = $tagsCount
            }
        }
        
        Export-ResultsToJson -Results $exportObject -ServicePrincipalId $fullServicePrincipal.id
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
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
