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
      $templateUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$DeviceConfigurationId')"
      $template = Invoke-MgGraphRequest -Method GET -Uri $templateUri

      return $template
  }
  catch {
      Write-Error "Error retrieving settings template: $_"
      return $null
  }
}

function Get-ConfigurationPolicyAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationPolicyId
    )

    try {
        $assignmentsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$ConfigurationPolicyId')/assignments"
        $assignments = Get-Paginated -InitialUri $assignmentsUri

        Write-Host "Assignments retrieved successfully."
        return $assignments
    }
    catch {
        Write-Error "Error retrieving configuration policy assignments: $_"
        return $null
    }
}

function Get-ConfigurationPolicySettings {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigurationPolicyId
    )

    try {
        $settingsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('$ConfigurationPolicyId')/settings?$expand=settingDefinitions&$top=1000"
        $settings = Get-Paginated -InitialUri $settingsUri

        Write-Host "Settings retrieved successfully."
        return $settings
    }
    catch {
        Write-Error "Error retrieving configuration policy settings: $_"
        return $null
    }
}

function Get-SettingTemplates {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyTemplateId
    )

    try {
        $settingTemplatesUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicyTemplates('$PolicyTemplateId')/settingTemplates?$expand=settingDefinitions&$top=1000"
        $settingTemplates = Get-Paginated -InitialUri $settingTemplatesUri

        Write-Host "Setting templates retrieved successfully."
        return $settingTemplates
    }
    catch {
        Write-Error "Error retrieving setting templates: $_"
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
    # Extract IDs dynamically
    $configurationPolicyId = $templateData.id
    $policyTemplateId = $templateData.templateReference.templateId

    Write-Host "Retrieving assignments for configuration policy..."
    $assignments = Get-ConfigurationPolicyAssignments -ConfigurationPolicyId $configurationPolicyId

    Write-Host "Retrieving settings for configuration policy..."
    $settings = Get-ConfigurationPolicySettings -ConfigurationPolicyId $configurationPolicyId

    # Ensure settings are always wrapped in an array
    if ($null -eq $settings) {
        $settings = @() # Empty array if no settings are found
    } elseif ($settings -isnot [Array]) {
        $settings = @($settings) # Wrap single setting in an array
    }

    Write-Host "Retrieving setting templates..."
    $settingTemplates = Get-SettingTemplates -PolicyTemplateId $policyTemplateId

    # Consolidate data into a single object
    $outputData = [PSCustomObject]@{
        baseResource      = $templateData
        assignments       = $assignments
        settings          = $settings
        settingTemplates  = $settingTemplates
    }

    # Export to a single JSON file
    $outputData | ConvertTo-Json -Depth 100 | Out-File "settingCatalogTemplate.json"
    Write-Host "Consolidated JSON output has been saved to 'settingCatalogTemplate.json'"
} else {
    Write-Host "No data found for the specified template ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."
