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
    HelpMessage="Specify the ID of the reusable policy settings to retrieve")]
    [ValidateNotNullOrEmpty()]
    [string]$ReusablePolicySettingId
)

# Helper function to retrieve all pages of data
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
        } else {
            $allItems += $response
        }
        
        $currentUri = $response.'@odata.nextLink'
    } while ($currentUri)

    return $allItems
}

# Helper function to retrieve a specific reusable policy setting by ID
function Get-ReusablePolicySettingById {
  param (
      [Parameter(Mandatory=$true)]
      [string]$ReusablePolicySettingId
  )

  try {
      # Get base policy information with select parameters
      $policyUri = "https://graph.microsoft.com/beta/deviceManagement/reusablePolicySettings/$ReusablePolicySettingId`?`$select=settinginstance,displayname,description"
      $policy = Invoke-MgGraphRequest -Method GET -Uri $policyUri

      return $policy
  }
  catch {
      Write-Error "Error retrieving reusable policy setting by ID: $_"
      return $null
  }
}
# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

Write-Host "Retrieving reusable policy setting with ID: $ReusablePolicySettingId"
$policyData = Get-ReusablePolicySettingById -ReusablePolicySettingId $ReusablePolicySettingId

if ($null -ne $policyData) {
    Write-Host "`nFull policy JSON (including references):"
    $jsonOutput = $policyData | ConvertTo-Json -Depth 100
    Write-Output $jsonOutput
    
    $jsonOutput | Out-File "reusable_policy_settings_export.json"
    Write-Host "`nComplete data has been saved to 'reusable_policy_settings_export.json'"
} else {
    Write-Host "No data found for the specified reusable policy setting ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."