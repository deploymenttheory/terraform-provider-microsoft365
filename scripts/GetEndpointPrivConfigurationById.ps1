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
    [string]$SettingsCatalogItemId
)

# Helper function to retrieve all pages of settings
function Get-Paginated {
    param (
        [Parameter(Mandatory=$true)]
        [string]$InitialUri
    )

    $allSettings = @()
    $currentUri = $InitialUri

    do {
        $response = Invoke-MgGraphRequest -Method GET -Uri $currentUri
        
        if ($response.value) {
            $allSettings += $response.value
        } else {
            $allSettings += $response
        }
        
        # Get the next page URL if it exists
        $currentUri = $response.'@odata.nextLink'
    } while ($currentUri)

    return $allSettings
}

# Helper function to retrieve a specific settings catalog policy by ID
function Get-SettingsCatalogPolicyById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$SettingsCatalogItemId
    )

    try {
        # Get base policy information
        $policyUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$SettingsCatalogItemId"
        $policy = Invoke-MgGraphRequest -Method GET -Uri $policyUri

        # Get settings
        $settingsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$SettingsCatalogItemId/settings"
        $settings = Get-Paginated -InitialUri $settingsUri

        # Get assignments if they exist
        $assignmentsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$SettingsCatalogItemId/assignments"
        $assignments = Get-Paginated -InitialUri $assignmentsUri

        # Combine everything into a single structure
        $policy | Add-Member -NotePropertyName 'settings' -NotePropertyValue @($settings) -Force
        $policy | Add-Member -NotePropertyName 'assignments' -NotePropertyValue $assignments -Force

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

Write-Host "Retrieving catalog policy with ID: $SettingsCatalogItemId"
$catalogData = Get-SettingsCatalogPolicyById -SettingsCatalogItemId $SettingsCatalogItemId

if ($null -ne $catalogData) {
    Write-Host "`nFull policy JSON (including settings and assignments):"
    $jsonOutput = $catalogData | ConvertTo-Json -Depth 100
    Write-Output $jsonOutput
    
    $jsonOutput | Out-File "settings_catalog_policy_export.json"
    Write-Host "`nComplete data has been saved to 'settings_catalog_policy_export.json'"
} else {
    Write-Host "No data found for the specified catalog policy ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."