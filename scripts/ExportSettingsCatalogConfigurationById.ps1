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
      $policyUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$SettingsCatalogItemId"
      $policy = Invoke-MgGraphRequest -Method GET -Uri $policyUri

      $settingsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$SettingsCatalogItemId/settings"
      $allSettings = Get-Paginated -InitialUri $settingsUri

      # Arrange the settings into the expected format with sequential IDs for each setting
      $formattedSettings = @()
      for($i = 0; $i -lt $allSettings.Count; $i++) {
          $formattedSettings += @{
              id = $i.ToString()
              settingInstance = $allSettings[$i].settingInstance
          }
      }

      # Add the formatted settings array to the policy object
      $policy | Add-Member -NotePropertyName 'settings' -NotePropertyValue $formattedSettings

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
$catalogItemData = Get-SettingsCatalogPolicyById -SettingsCatalogItemId $SettingsCatalogItemId

if ($null -ne $catalogItemData) {
    Write-Host "`nFull JSON output for settings catalog policy:"
    $jsonString = $catalogItemData | ConvertTo-Json -Depth 100 -Compress
    # Format the JSON for readability
    $jsonFormatted = $jsonString | ConvertFrom-Json | ConvertTo-Json -Depth 100
    
    Write-Output $jsonFormatted
    
    $jsonFormatted | Out-File "settingsCatalogPolicy.json"
    Write-Host "`nJSON output has also been saved to 'settingsCatalogPolicy.json'"
} else {
    Write-Host "No data found for the specified catalog policy ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."