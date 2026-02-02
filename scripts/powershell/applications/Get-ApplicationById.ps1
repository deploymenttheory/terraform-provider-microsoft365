<#
.SYNOPSIS
    Gets an application by Object ID or App ID with all attributes.

.DESCRIPTION
    Retrieves a Microsoft Entra application by its Object ID or App ID (Client ID),
    displaying all attributes including app roles, API permissions, credentials, and owners.
    Useful for comparing API values against what appears in the Azure Portal.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER ApplicationId
    The Object ID or App ID (Client ID) of the application to retrieve.
    Will automatically detect which type of ID is provided.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Get application by Object ID
    .\Get-ApplicationById.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationId "12345678-1234-1234-1234-123456789012"

.EXAMPLE
    # Get application and export to JSON
    .\Get-ApplicationById.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -ApplicationId "12345678-1234-1234-1234-123456789012" `
        -ExportToJson $true

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

function Format-ComplexObject {
    param (
        [Parameter(Mandatory=$true)]
        $Object,
        
        [Parameter(Mandatory=$false)]
        [int]$Indent = 6
    )
    
    $indentStr = " " * $Indent
    
    foreach ($key in $Object.Keys) {
        $value = $Object[$key]
        
        if ($null -eq $value) {
            Write-Host "$indentStr${key}: " -NoNewline -ForegroundColor White
            Write-Host "(null)" -ForegroundColor Gray
        }
        elseif ($value -is [array] -and $value.Count -gt 0) {
            Write-Host "$indentStr${key}:" -ForegroundColor White
            foreach ($item in $value) {
                if ($item -is [string]) {
                    Write-Host "$indentStr  • $item" -ForegroundColor Green
                }
                else {
                    Write-Host "$indentStr  • (complex)" -ForegroundColor Gray
                }
            }
        }
        elseif ($value -is [bool]) {
            Write-Host "$indentStr${key}: " -NoNewline -ForegroundColor White
            Write-Host $value.ToString().ToLower() -ForegroundColor Green
        }
        else {
            Write-Host "$indentStr${key}: " -NoNewline -ForegroundColor White
            Write-Host $value -ForegroundColor Green
        }
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
    Write-Host "   🔍 Get Application By ID" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Determine if the provided ID is an Object ID or App ID
    Write-Host "🔍 Retrieving application..." -ForegroundColor Cyan
    
    $application = $null
    $searchByAppId = $false
    
    # Try to get by Object ID first
    try {
        $uri = "https://graph.microsoft.com/beta/applications/$ApplicationId"
        $application = Invoke-MgGraphRequest -Method GET -Uri $uri
        Write-Host "✅ Found application by Object ID" -ForegroundColor Green
    }
    catch {
        # If not found, try searching by App ID
        Write-Host "   Not found by Object ID, trying App ID..." -ForegroundColor Gray
        $searchByAppId = $true
        
        $searchUri = "https://graph.microsoft.com/beta/applications?`$filter=appId eq '$ApplicationId'"
        $searchResult = Invoke-MgGraphRequest -Method GET -Uri $searchUri
        
        if ($searchResult.value -and $searchResult.value.Count -gt 0) {
            $application = $searchResult.value[0]
            Write-Host "✅ Found application by App ID" -ForegroundColor Green
        }
        else {
            Write-Host "❌ Application not found with ID: $ApplicationId" -ForegroundColor Red
            Write-Host ""
            exit 1
        }
    }
    
    Write-Host ""
    
    # Display basic information
    Write-Host "📋 Application Overview" -ForegroundColor Cyan
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
        @{ Name = "disabledByMicrosoftStatus"; Label = "Disabled By Microsoft" },
        @{ Name = "notes"; Label = "Notes" },
        @{ Name = "isDeviceOnlyAuthSupported"; Label = "Device Only Auth Supported" },
        @{ Name = "isFallbackPublicClient"; Label = "Fallback Public Client" }
    )
    
    foreach ($attr in $keyAttributes) {
        $value = $application[$attr.Name]
        $formattedValue = Format-AttributeValue -Value $value
        Write-Host "   $($attr.Label): " -NoNewline -ForegroundColor White
        Write-Host "$formattedValue" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Display App Roles
    if ($application.appRoles -and $application.appRoles.Count -gt 0) {
        Write-Host "👥 App Roles ($($application.appRoles.Count))" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($role in $application.appRoles) {
            Write-Host "   📌 $($role.displayName)" -ForegroundColor Yellow
            Write-Host "      ID: $($role.id)" -ForegroundColor Gray
            Write-Host "      Value: $($role.value)" -ForegroundColor Gray
            Write-Host "      Description: $($role.description)" -ForegroundColor Gray
            Write-Host "      Is Enabled: $($role.isEnabled)" -ForegroundColor $(if ($role.isEnabled) { "Green" } else { "Red" })
            Write-Host "      Allowed Member Types: $($role.allowedMemberTypes -join ', ')" -ForegroundColor Gray
            Write-Host "      Origin: $($role.origin)" -ForegroundColor Gray
            Write-Host ""
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display Identifier URIs
    if ($application.identifierUris -and $application.identifierUris.Count -gt 0) {
        Write-Host "🔗 Identifier URIs" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($uri in $application.identifierUris) {
            Write-Host "   • $uri" -ForegroundColor Green
        }
        
        Write-Host ""
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display API configuration
    if ($application.api) {
        Write-Host "🔌 API Configuration" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        Format-ComplexObject -Object $application.api
        
        Write-Host ""
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display Web configuration
    if ($application.web) {
        Write-Host "🌐 Web Configuration" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        Format-ComplexObject -Object $application.web
        
        Write-Host ""
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display Required Resource Access (API Permissions)
    if ($application.requiredResourceAccess -and $application.requiredResourceAccess.Count -gt 0) {
        Write-Host "🔐 Required Resource Access (API Permissions)" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($resource in $application.requiredResourceAccess) {
            Write-Host "   📦 Resource ID: $($resource.resourceAppId)" -ForegroundColor Yellow
            
            if ($resource.resourceAccess -and $resource.resourceAccess.Count -gt 0) {
                foreach ($access in $resource.resourceAccess) {
                    Write-Host "      • $($access.id) ($($access.type))" -ForegroundColor Green
                }
            }
            Write-Host ""
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display Key Credentials
    if ($application.keyCredentials -and $application.keyCredentials.Count -gt 0) {
        Write-Host "🔑 Key Credentials ($($application.keyCredentials.Count))" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($cred in $application.keyCredentials) {
            Write-Host "   📌 $($cred.displayName)" -ForegroundColor Yellow
            Write-Host "      Key ID: $($cred.keyId)" -ForegroundColor Gray
            Write-Host "      Type: $($cred.type)" -ForegroundColor Gray
            Write-Host "      Usage: $($cred.usage)" -ForegroundColor Gray
            Write-Host "      Start: $($cred.startDateTime)" -ForegroundColor Gray
            Write-Host "      End: $($cred.endDateTime)" -ForegroundColor Gray
            Write-Host ""
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display Password Credentials
    if ($application.passwordCredentials -and $application.passwordCredentials.Count -gt 0) {
        Write-Host "🔐 Password Credentials ($($application.passwordCredentials.Count))" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($cred in $application.passwordCredentials) {
            Write-Host "   📌 $($cred.displayName)" -ForegroundColor Yellow
            Write-Host "      Key ID: $($cred.keyId)" -ForegroundColor Gray
            Write-Host "      Start: $($cred.startDateTime)" -ForegroundColor Gray
            Write-Host "      End: $($cred.endDateTime)" -ForegroundColor Gray
            Write-Host ""
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Get and display owners
    Write-Host "👤 Retrieving application owners..." -ForegroundColor Cyan
    
    $ownersUri = "https://graph.microsoft.com/beta/applications/$($application.id)/owners"
    $owners = Invoke-MgGraphRequest -Method GET -Uri $ownersUri
    
    if ($owners.value -and $owners.value.Count -gt 0) {
        Write-Host "✅ Found $($owners.value.Count) owner(s)" -ForegroundColor Green
        Write-Host ""
        
        Write-Host "👥 Application Owners" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        foreach ($owner in $owners.value) {
            Write-Host "   • $($owner.displayName)" -ForegroundColor Green
            Write-Host "     User Principal Name: $($owner.userPrincipalName)" -ForegroundColor Gray
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
    Write-Host "   Display Name: $($application.displayName)" -ForegroundColor White
    Write-Host "   Object ID: $($application.id)" -ForegroundColor White
    Write-Host "   App ID: $($application.appId)" -ForegroundColor White
    Write-Host "   Sign-In Audience: $($application.signInAudience)" -ForegroundColor White
    
    $appRolesCount = if ($application.appRoles) { $application.appRoles.Count } else { 0 }
    Write-Host "   App Roles: $appRolesCount" -ForegroundColor White
    
    $keyCredsCount = if ($application.keyCredentials) { $application.keyCredentials.Count } else { 0 }
    Write-Host "   Key Credentials: $keyCredsCount" -ForegroundColor White
    
    $passwordCredsCount = if ($application.passwordCredentials) { $application.passwordCredentials.Count } else { 0 }
    Write-Host "   Password Credentials: $passwordCredsCount" -ForegroundColor White
    
    $ownersCount = if ($owners.value) { $owners.value.Count } else { 0 }
    Write-Host "   Owners: $ownersCount" -ForegroundColor White
    
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            Application = $application
            Owners = $owners.value
            Summary = @{
                DisplayName = $application.displayName
                ObjectId = $application.id
                AppId = $application.appId
                SignInAudience = $application.signInAudience
                AppRolesCount = $appRolesCount
                KeyCredentialsCount = $keyCredsCount
                PasswordCredentialsCount = $passwordCredsCount
                OwnersCount = $ownersCount
            }
        }
        
        Export-ResultsToJson -Results $exportObject -ApplicationId $application.id
    }
    
    Write-Host "🎉 Operation completed!" -ForegroundColor Green
    Write-Host ""
}
catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
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
