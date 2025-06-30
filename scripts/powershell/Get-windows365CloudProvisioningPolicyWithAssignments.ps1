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
    HelpMessage="Cloud PC Provisioning Policy ID")]
    [ValidateNotNullOrEmpty()]
    [string]$ProvisioningPolicyId
)

Import-Module Microsoft.Graph.Authentication

function Get-Windows365CloudProvisioningPolicyWithAssignmentsAndSelect {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    try {
        $uri = "https://graph.microsoft.com/beta/deviceManagement/virtualEndpoint/provisioningPolicies/$PolicyId`?expand=assignments&select=*"
        Write-Host "🔍 Getting Cloud PC Provisioning Policy with assignments and select..." -ForegroundColor Cyan
        Write-Host "   Policy ID: $PolicyId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        return $response
    }
    catch {
        Write-Host "❌ Error getting provisioning policy: $_" -ForegroundColor Red
        throw
    }
}

function Show-ProvisioningPolicyDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Policy
    )
    Write-Host "📋 Provisioning Policy Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    if ($Policy.id) { Write-Host "   • ID: $($Policy.id)" -ForegroundColor Green }
    if ($Policy.displayName) { Write-Host "   • Display Name: $($Policy.displayName)" -ForegroundColor Green }
    if ($Policy.description) { Write-Host "   • Description: $($Policy.description)" -ForegroundColor Green }
    if ($Policy.imageId) { Write-Host "   • Image ID: $($Policy.imageId)" -ForegroundColor Green }
    if ($Policy.imageDisplayName) { Write-Host "   • Image Display Name: $($Policy.imageDisplayName)" -ForegroundColor Green }
    if ($Policy.provisioningType) { Write-Host "   • Provisioning Type: $($Policy.provisioningType)" -ForegroundColor Green }
    if ($Policy.managedBy) { Write-Host "   • Managed By: $($Policy.managedBy)" -ForegroundColor Green }
    if ($Policy.enableSingleSignOn -ne $null) { Write-Host "   • Enable SSO: $($Policy.enableSingleSignOn)" -ForegroundColor Green }
    if ($Policy.cloudPcNamingTemplate) { Write-Host "   • Naming Template: $($Policy.cloudPcNamingTemplate)" -ForegroundColor Green }
    if ($Policy.scopeIds) { Write-Host "   • Scope IDs: $($Policy.scopeIds -join ", ")" -ForegroundColor Green }
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

function Show-AssignmentDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Assignment
    )
    Write-Host "📋 Assignment Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    if ($Assignment.id) { Write-Host "   • ID: $($Assignment.id)" -ForegroundColor Green }
    if ($Assignment.target) {
        $target = $Assignment.target
        if ($target.'@odata.type') {
            $targetType = $target.'@odata.type' -replace '#microsoft.graph.', ''
            Write-Host "   • Target Type: $targetType" -ForegroundColor Yellow
        }
        if ($target.groupId) { Write-Host "   • Group ID: $($target.groupId)" -ForegroundColor Yellow }
        if ($target.servicePlanId) { Write-Host "   • Service Plan ID: $($target.servicePlanId)" -ForegroundColor Yellow }
        if ($target.allotmentDisplayName) { Write-Host "   • Allotment Display Name: $($target.allotmentDisplayName)" -ForegroundColor Yellow }
        if ($target.allotmentLicensesCount) { Write-Host "   • Allotment License Count: $($target.allotmentLicensesCount)" -ForegroundColor Yellow }
    }
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    $policy = Get-Windows365CloudProvisioningPolicyWithAssignmentsAndSelect -PolicyId $ProvisioningPolicyId
    Show-ProvisioningPolicyDetails -Policy $policy
    if ($policy.assignments -and $policy.assignments.Count -gt 0) {
        Write-Host "📊 Found $($policy.assignments.Count) assignment(s)" -ForegroundColor Green
        for ($i = 0; $i -lt $policy.assignments.Count; $i++) {
            Write-Host "Assignment $($i + 1):" -ForegroundColor Magenta
            Show-AssignmentDetails -Assignment $policy.assignments[$i]
        }
    } elseif ($policy.assignments -and $null -eq $policy.assignments.Count) {
        Write-Host "📊 Found 1 assignment" -ForegroundColor Green
        Show-AssignmentDetails -Assignment $policy.assignments
    } else {
        Write-Host "📊 No assignments found for this provisioning policy" -ForegroundColor Yellow
    }
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    }
    catch {
        # Ignore disconnect errors
    }
} 