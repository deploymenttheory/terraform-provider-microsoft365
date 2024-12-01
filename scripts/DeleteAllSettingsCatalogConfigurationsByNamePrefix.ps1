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
    HelpMessage="Specify the settings catalog name prefix to match for deletion (e.g., 'test_collection-')")]
    [ValidateNotNullOrEmpty()]
    [string]$SettingsCatalogNamePrefix
)

# Helper function to retrieve all pages of items
function Get-Paginated {
    param (
        [Parameter(Mandatory=$true)]
        [string]$InitialUri
    )

    $allItems = @()
    $currentUri = $InitialUri

    do {
        $response = Invoke-MgGraphRequest -Method GET -Uri $currentUri
        
        if ($response.value) {
            $allItems += $response.value
        }
        
        # Get the next page URL if it exists
        $currentUri = $response.'@odata.nextLink'
    } while ($currentUri)

    return $allItems
}

# Helper function to get all settings catalog policies
function Get-AllSettingsCatalogPolicies {
    try {
        $policiesUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
        return Get-Paginated -InitialUri $policiesUri
    }
    catch {
        Write-Error "Error retrieving settings catalog policies: $_"
        return $null
    }
}

# Helper function to delete a settings catalog policy
function Remove-SettingsCatalogPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId,
        [string]$PolicyName
    )

    try {
        $policyUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$PolicyId"
        Invoke-MgGraphRequest -Method DELETE -Uri $policyUri
        Write-Host "Successfully deleted policy: $PolicyName (ID: $PolicyId)"
        return $true
    }
    catch {
        Write-Error "Error deleting policy $PolicyName (ID: $PolicyId): $_"
        return $false
    }
}

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Get all settings catalog policies
Write-Host "Retrieving all settings catalog policies..."
$allPolicies = Get-AllSettingsCatalogPolicies

if ($null -ne $allPolicies) {
    # Filter policies by name prefix
    $matchingPolicies = $allPolicies | Where-Object { $_.name -like "$SettingsCatalogNamePrefix*" }
    
    if ($matchingPolicies.Count -gt 0) {
        Write-Host "`nFound $($matchingPolicies.Count) policies matching prefix '$SettingsCatalogNamePrefix'"
        
        # Confirm before deletion
        $confirmation = Read-Host "Do you want to proceed with deletion? (Y/N)"
        if ($confirmation -eq 'Y') {
            $deletedCount = 0
            $failedCount = 0
            
            foreach ($policy in $matchingPolicies) {
                Write-Host "`nDeleting policy: $($policy.name)..."
                $result = Remove-SettingsCatalogPolicy -PolicyId $policy.id -PolicyName $policy.name
                if ($result) {
                    $deletedCount++
                } else {
                    $failedCount++
                }
            }
            
            Write-Host "`nDeletion complete:"
            Write-Host "Successfully deleted: $deletedCount"
            Write-Host "Failed to delete: $failedCount"
        } else {
            Write-Host "Operation cancelled by user."
        }
    } else {
        Write-Host "No policies found matching prefix '$SettingsCatalogNamePrefix'"
    }
} else {
    Write-Host "No settings catalog policies found or error occurred."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."