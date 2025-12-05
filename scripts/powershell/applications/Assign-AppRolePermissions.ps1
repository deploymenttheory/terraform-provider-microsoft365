<#
.SYNOPSIS
    Assigns app role permissions to a service principal.

.DESCRIPTION
    Simple script to assign app role IDs to an enterprise app (service principal).
    Uses the appRoleAssignedTo endpoint to grant permissions.

.PARAMETER TenantId
    The Entra ID tenant ID.

.PARAMETER ClientId
    The application (client) ID for authentication.

.PARAMETER ClientSecret
    The client secret for authentication.

.PARAMETER PrincipalId
    The Object ID of the service principal to grant permissions TO.

.PARAMETER ResourceId
    The Object ID of the resource service principal (e.g., Microsoft Graph SP in your tenant).

.PARAMETER AppRoleIds
    Array of app role IDs (GUIDs) to assign.

.EXAMPLE
REF: https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities
    # Assign Agent Identity permissions to enterprise app
    .\Assign-AppRolePermissions.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id-for-auth" `
        -ClientSecret "your-client-secret-for-auth" `
        -PrincipalId "your-target-service-principal-id-for-permissions" ` typically the SP-CPSS-GLBL-AGENTS-C-01 ID
        -ResourceId "cb57292c-59b4-4e5a-a677-187cb72cb6c6" ` typically the Microsoft Graph SP ID
        -AppRoleIds @(
            # Agent Identity Blueprint permissions
            "7547a7d1-36fa-4479-9c31-559a600eaa4f",  # AgentIdentityBlueprint.Read.All
            "7fddd33b-d884-4ec0-8696-72cff90ff825",  # AgentIdentityBlueprint.ReadWrite.All
            "ea4b2453-ad2d-4d94-9155-10d5d9493ce9",  # AgentIdentityBlueprint.Create
            "0510736e-bdfb-4b37-9a1f-89b4a074763a",  # AgentIdentityBlueprint.AddRemoveCreds.All
            "76232daa-a1e4-4544-b664-495a006513bf",  # AgentIdentityBlueprint.UpdateBranding.All
            
            # Agent User permissions
            "b782c9ad-6f2b-4894-a21b-72bf22417f0a",  # AgentIdUser.ReadWrite.All
            "4aa6e624-eee0-40ab-bdd8-f9639038a614"   # AgentIdUser.ReadWrite.IdentityParentedBy
        )

.NOTES
    Author: Deployment Theory
    
    Agent Identity App-only IDs from:
    https://learn.microsoft.com/en-us/graph/api/resources/agentid-platform-overview?view=graph-rest-beta#permissions-for-managing-agent-identities
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true)]
    [string]$PrincipalId,
    
    [Parameter(Mandatory=$true)]
    [string]$ResourceId,
    
    [Parameter(Mandatory=$true)]
    [string[]]$AppRoleIds
)

Import-Module Microsoft.Graph.Authentication

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    
    $secureSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureSecret
    
    Connect-MgGraph -ClientSecretCredential $credential -TenantId $TenantId -NoWelcome
    
    Write-Host "✅ Connected" -ForegroundColor Green
    Write-Host ""
    Write-Host "📋 Principal ID: $PrincipalId" -ForegroundColor Cyan
    Write-Host "📦 Resource ID: $ResourceId" -ForegroundColor Cyan
    Write-Host "🔑 App Roles to assign: $($AppRoleIds.Count)" -ForegroundColor Cyan
    Write-Host ""
    
    $successCount = 0
    $skipCount = 0
    $failCount = 0
    
    foreach ($appRoleId in $AppRoleIds) {
        Write-Host "▶️  Assigning: $appRoleId" -ForegroundColor White
        
        try {
            $body = @{
                principalId = $PrincipalId
                resourceId = $ResourceId
                appRoleId = $appRoleId
            } | ConvertTo-Json
            
            $uri = "https://graph.microsoft.com/v1.0/servicePrincipals/$ResourceId/appRoleAssignedTo"
            
            $result = Invoke-MgGraphRequest -Method POST -Uri $uri -Body $body -ContentType "application/json"
            
            Write-Host "   ✅ Assigned successfully" -ForegroundColor Green
            $successCount++
        }
        catch {
            $errorMessage = $_.Exception.Message
            $errorDetails = $_.ErrorDetails.Message

            # Check for various indicators that the permission already exists
            $alreadyExists = $false

            if ($errorMessage -like "*Permission being assigned already exists*" -or
                $errorMessage -like "*already exists*" -or
                $errorDetails -like "*already exists*" -or
                $errorDetails -like "*Permission being assigned already exists*") {
                $alreadyExists = $true
            }

            # Parse error details for the specific error code
            if ($errorDetails) {
                try {
                    $errorObj = $errorDetails | ConvertFrom-Json
                    if ($errorObj.error.code -eq "Request_BadRequest" -and
                        $errorObj.error.message -like "*already exists*") {
                        $alreadyExists = $true
                    }
                }
                catch {
                    # If JSON parsing fails, continue with string matching
                }
            }

            if ($alreadyExists) {
                Write-Host "   ⚠️  Already assigned" -ForegroundColor Yellow
                $skipCount++
            }
            else {
                Write-Host "   ❌ Failed: $errorMessage" -ForegroundColor Red
                if ($errorDetails) {
                    Write-Host "      Details: $errorDetails" -ForegroundColor Red
                }
                $failCount++
            }
        }
        Write-Host ""
    }
    
    Write-Host "📊 Summary" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host "   ✅ Assigned: $successCount" -ForegroundColor Green
    Write-Host "   ⚠️  Already existed: $skipCount" -ForegroundColor Yellow
    Write-Host "   ❌ Failed: $failCount" -ForegroundColor Red
    Write-Host ""
}
catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
    exit 1
}
finally {
    Disconnect-MgGraph 2>$null
}

