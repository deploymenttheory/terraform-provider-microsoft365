[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,

    [Parameter(Mandatory=$true)]
    [string]$DeviceConfigurationId
)

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
        
        $currentUri = $response.'@odata.nextLink'
    } while ($currentUri)

    return $allSettings
}

function Get-SettingsTemplateById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$DeviceConfigurationId
    )

    try {
        # Get template details
        $templateUri = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations/$DeviceConfigurationId"
        $template = Invoke-MgGraphRequest -Method GET -Uri $templateUri

        return $template
    }
    catch {
        Write-Error "Error retrieving settings template: $_"
        return $null
    }
}

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

Write-Host "Retrieving template with ID: $DeviceConfigurationId"
$templateData = Get-SettingsTemplateById -DeviceConfigurationId $DeviceConfigurationId

if ($null -ne $templateData) {
    Write-Host "`nFull JSON output for settings template:"
    $jsonString = $templateData | ConvertTo-Json -Depth 100 -Compress
    # Format the JSON for readability
    $jsonFormatted = $jsonString | ConvertFrom-Json | ConvertTo-Json -Depth 100
    
    Write-Output $jsonFormatted
    $jsonFormatted | Out-File "settingsTemplate.json"
    Write-Host "`nJSON output has been saved to 'settingsTemplate.json'"
} else {
    Write-Host "No data found for the specified template ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."