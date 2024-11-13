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
    HelpMessage="Specify the ID of the settings catalog policy to retrieve")]
    [ValidateNotNullOrEmpty()]
    [string]$CatalogItemId
)

# Helper function to retrieve a specific settings catalog policy by ID
function Get-SettingsCatalogPolicyById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$CatalogItemId
    )

    try {
        $policyUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$CatalogItemId"
        $policy = Invoke-MgGraphRequest -Method GET -Uri $policyUri

        $settingsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$CatalogItemId/settings"
        $settings = Invoke-MgGraphRequest -Method GET -Uri $settingsUri

        $policy | Add-Member -NotePropertyName 'settingsDetails' -NotePropertyValue $settings.value

        return $policy
    }
    catch {
        Write-Error "Error retrieving settings catalog policy by ID: $_"
        return $null
    }
}

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

# Retrieve and output the specified catalog policy
Write-Host "Retrieving catalog policy with ID: $CatalogItemId"
$catalogItemData = Get-SettingsCatalogPolicyById -CatalogItemId $CatalogItemId

if ($null -ne $catalogItemData) {
    Write-Host "`nFull JSON output for settings catalog policy:"
    $jsonString = $catalogItemData | ConvertTo-Json -Depth 1000 -Compress
    
    $jsonFormatted = $jsonString | ConvertFrom-Json | ConvertTo-Json -Depth 1000
    
    Write-Output $jsonFormatted
    
    $jsonFormatted | Out-File "settingsCatalogPolicy.json"
    Write-Host "`nJSON output has also been saved to 'settingsCatalogPolicy.json'"
} else {
    Write-Host "No data found for the specified catalog policy ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."