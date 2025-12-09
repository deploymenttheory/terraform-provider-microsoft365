<#
.SYNOPSIS
    Gets a service principal by display name with all attributes and assigned permissions.

.DESCRIPTION
    Searches for a service principal by display name, retrieves all its attributes,
    and lists all assigned app role permissions with their details.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER TargetServicePrincipalName
    The display name of the service principal to search for.

.PARAMETER ExportToJson
    Whether to export the results to a JSON file.

.EXAMPLE
    # Get service principal by name
    .\Get-ServicePrincipalByName.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -TargetServicePrincipalName "My Enterprise App"

.EXAMPLE
    # Get service principal and export to JSON
    .\Get-ServicePrincipalByName.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -TargetServicePrincipalName "My Enterprise App" `
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
    [string]$TargetServicePrincipalName,
    
    [Parameter(Mandatory=$false)]
    [bool]$ExportToJson = $false
)

Import-Module Microsoft.Graph.Authentication

#region Helper Functions

function Get-AppRoleDetails {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ResourceId,
        
        [Parameter(Mandatory=$true)]
        [string]$AppRoleId
    )
    
    try {
        $uri = "https://graph.microsoft.com/beta/servicePrincipals/$ResourceId"
        $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        $appRole = $resourceSp.appRoles | Where-Object { $_.id -eq $AppRoleId }
        
        if ($appRole) {
            return @{
                Value = $appRole.value
                DisplayName = $appRole.displayName
                Description = $appRole.description
            }
        } else {
            return @{
                Value = "Unknown"
                DisplayName = "Unknown Permission"
                Description = "App role ID: $AppRoleId"
            }
        }
    }
    catch {
        return @{
            Value = "Unknown"
            DisplayName = "Error retrieving"
            Description = $_.Exception.Message
        }
    }
}

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
        [string]$ServicePrincipalName
    )
    
    try {
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $safeName = $ServicePrincipalName -replace '[^\w\-]', '_'
        $fileName = "ServicePrincipal_${safeName}_${timestamp}.json"
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
    Write-Host "   🔍 Get Service Principal By Name" -ForegroundColor Magenta
    Write-Host "═══════════════════════════════════════════════════════════════════" -ForegroundColor Magenta
    Write-Host ""
    
    # Connect to Microsoft Graph
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    
    # Search for service principal by display name
    Write-Host "🔍 Searching for service principal: '$TargetServicePrincipalName'..." -ForegroundColor Cyan
    
    $searchUri = "https://graph.microsoft.com/beta/servicePrincipals?`$filter=displayName eq '$TargetServicePrincipalName'"
    $searchResult = Invoke-MgGraphRequest -Method GET -Uri $searchUri
    
    if (-not $searchResult.value -or $searchResult.value.Count -eq 0) {
        Write-Host "⚠️  No service principal found with name: '$TargetServicePrincipalName'" -ForegroundColor Yellow
        Write-Host ""
        Write-Host "💡 Tip: Try a partial search..." -ForegroundColor Gray
        
        $partialUri = "https://graph.microsoft.com/beta/servicePrincipals?`$filter=startswith(displayName,'$TargetServicePrincipalName')&`$top=5"
        $partialResult = Invoke-MgGraphRequest -Method GET -Uri $partialUri
        
        if ($partialResult.value -and $partialResult.value.Count -gt 0) {
            Write-Host "   Found similar names:" -ForegroundColor Gray
            foreach ($sp in $partialResult.value) {
                Write-Host "   • $($sp.displayName)" -ForegroundColor Gray
            }
        }
        Write-Host ""
        exit 1
    }
    
    if ($searchResult.value.Count -gt 1) {
        Write-Host "⚠️  Multiple service principals found with name: '$TargetServicePrincipalName'" -ForegroundColor Yellow
        Write-Host "   Using the first match." -ForegroundColor Gray
        Write-Host ""
    }
    
    $objectId = $searchResult.value[0].id
    Write-Host "✅ Found service principal with Object ID: $objectId" -ForegroundColor Green
    Write-Host ""
    
    # Get full service principal details
    Write-Host "🔍 Getting full service principal details..." -ForegroundColor Cyan
    
    $detailUri = "https://graph.microsoft.com/beta/servicePrincipals/$objectId"
    $servicePrincipal = Invoke-MgGraphRequest -Method GET -Uri $detailUri
    
    Write-Host "✅ Retrieved service principal details" -ForegroundColor Green
    Write-Host ""
    
    # Display service principal attributes
    Write-Host "📋 Service Principal Attributes" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Key attributes first
    $keyAttributes = @(
        @{ Name = "id"; Label = "Object ID" },
        @{ Name = "appId"; Label = "App ID (Client ID)" },
        @{ Name = "displayName"; Label = "Display Name" },
        @{ Name = "servicePrincipalType"; Label = "Service Principal Type" },
        @{ Name = "accountEnabled"; Label = "Account Enabled" },
        @{ Name = "appOwnerOrganizationId"; Label = "App Owner Org ID" },
        @{ Name = "createdDateTime"; Label = "Created Date Time" },
        @{ Name = "description"; Label = "Description" },
        @{ Name = "homepage"; Label = "Homepage" },
        @{ Name = "loginUrl"; Label = "Login URL" },
        @{ Name = "logoutUrl"; Label = "Logout URL" },
        @{ Name = "signInAudience"; Label = "Sign-In Audience" },
        @{ Name = "preferredSingleSignOnMode"; Label = "Preferred SSO Mode" },
        @{ Name = "disabledByMicrosoftStatus"; Label = "Disabled By Microsoft" }
    )
    
    Write-Host "   📌 Key Attributes" -ForegroundColor Yellow
    Write-Host ""
    
    foreach ($attr in $keyAttributes) {
        $value = $servicePrincipal[$attr.Name]
        $formattedValue = Format-AttributeValue -Value $value
        Write-Host "   $($attr.Label): " -NoNewline -ForegroundColor White
        Write-Host "$formattedValue" -ForegroundColor Green
    }
    
    Write-Host ""
    Write-Host "   📌 Collection Attributes" -ForegroundColor Yellow
    Write-Host ""
    
    # Collection attributes
    $collectionAttributes = @(
        @{ Name = "tags"; Label = "Tags" },
        @{ Name = "servicePrincipalNames"; Label = "Service Principal Names" },
        @{ Name = "replyUrls"; Label = "Reply URLs" },
        @{ Name = "alternativeNames"; Label = "Alternative Names" },
        @{ Name = "notificationEmailAddresses"; Label = "Notification Emails" }
    )
    
    foreach ($attr in $collectionAttributes) {
        $value = $servicePrincipal[$attr.Name]
        if ($value -and $value.Count -gt 0) {
            Write-Host "   $($attr.Label):" -ForegroundColor White
            foreach ($item in $value) {
                Write-Host "     • $item" -ForegroundColor Green
            }
        } else {
            Write-Host "   $($attr.Label): " -NoNewline -ForegroundColor White
            Write-Host "(none)" -ForegroundColor Gray
        }
    }
    
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
    
    # Get app role assignments
    Write-Host "🔍 Getting app role assignments..." -ForegroundColor Cyan
    
    $assignmentsUri = "https://graph.microsoft.com/beta/servicePrincipals/$objectId/appRoleAssignments"
    $assignments = Invoke-MgGraphRequest -Method GET -Uri $assignmentsUri
    
    $detailedAssignments = @()
    
    if (-not $assignments.value -or $assignments.value.Count -eq 0) {
        Write-Host "⚠️  No app role assignments found" -ForegroundColor Yellow
        Write-Host ""
    }
    else {
        Write-Host "✅ Found $($assignments.value.Count) app role assignment(s)" -ForegroundColor Green
        Write-Host ""
        
        Write-Host "📦 Assigned Permissions" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        $groupedByResource = $assignments.value | Group-Object -Property resourceId
        
        foreach ($group in $groupedByResource) {
            # Get resource service principal details
            $resourceUri = "https://graph.microsoft.com/beta/servicePrincipals/$($group.Name)"
            $resourceSp = Invoke-MgGraphRequest -Method GET -Uri $resourceUri
            
            Write-Host "   📦 Resource: $($resourceSp.displayName)" -ForegroundColor Yellow
            Write-Host "      App ID: $($resourceSp.appId)" -ForegroundColor Gray
            Write-Host "      Object ID: $($resourceSp.id)" -ForegroundColor Gray
            Write-Host ""
            
            foreach ($assignment in $group.Group) {
                $roleDetails = Get-AppRoleDetails -ResourceId $assignment.resourceId -AppRoleId $assignment.appRoleId
                
                Write-Host "      ✓ $($roleDetails.Value)" -ForegroundColor Green
                Write-Host "        Display Name: $($roleDetails.DisplayName)" -ForegroundColor Gray
                Write-Host "        App Role ID: $($assignment.appRoleId)" -ForegroundColor Gray
                Write-Host ""
                
                $detailedAssignments += @{
                    ResourceDisplayName = $resourceSp.displayName
                    ResourceAppId = $resourceSp.appId
                    ResourceObjectId = $resourceSp.id
                    PermissionName = $roleDetails.Value
                    PermissionDisplayName = $roleDetails.DisplayName
                    PermissionDescription = $roleDetails.Description
                    AppRoleId = $assignment.appRoleId
                    AssignmentId = $assignment.id
                    CreatedDateTime = $assignment.createdDateTime
                }
            }
        }
        
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
    }
    
    # Display summary
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   Display Name: $($servicePrincipal.displayName)" -ForegroundColor White
    Write-Host "   Object ID: $($servicePrincipal.id)" -ForegroundColor White
    Write-Host "   App ID: $($servicePrincipal.appId)" -ForegroundColor White
    Write-Host "   Type: $($servicePrincipal.servicePrincipalType)" -ForegroundColor White
    Write-Host "   Total Permissions: $($assignments.value.Count)" -ForegroundColor White
    if ($groupedByResource) {
        Write-Host "   Resource Providers: $($groupedByResource.Count)" -ForegroundColor White
    }
    Write-Host ""
    
    # Export results if requested
    if ($ExportToJson) {
        $exportObject = @{
            Timestamp = Get-Date -Format "o"
            ServicePrincipal = $servicePrincipal
            Summary = @{
                DisplayName = $servicePrincipal.displayName
                ObjectId = $servicePrincipal.id
                AppId = $servicePrincipal.appId
                ServicePrincipalType = $servicePrincipal.servicePrincipalType
                TotalPermissions = $assignments.value.Count
            }
            Permissions = $detailedAssignments
        }
        
        Export-ResultsToJson -Results $exportObject -ServicePrincipalName $servicePrincipal.displayName
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

