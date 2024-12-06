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

    [Parameter(Mandatory=$false,
    HelpMessage="Specify the template family to filter by (optional)")]
    [string]$TemplateFamily = "",

    [Parameter(Mandatory=$false,
    HelpMessage="Specify the number of policies to retrieve per page")]
    [int]$Top = 25
)

# Helper function to retrieve all pages of settings
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
        
        # Get the next page URL if it exists
        $currentUri = $response.'@odata.nextLink'
    } while ($currentUri)

    return $allItems
}

# Function to get list of settings catalog policies
function Get-SettingsCatalogPolicies {
    param (
        [string]$TemplateFamily,
        [int]$Top
    )

    try {
        $baseUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
        
        # Add select parameters - reordered to put name first, then id
        $select = "?`$select=name,id,description,platforms,lastModifiedDateTime,technologies,settingCount,roleScopeTagIds,isAssigned,templateReference"
        
        # Add top parameter
        $topParameter = "&`$top=$Top"
        
        # Add filter if template family is specified
        $filter = ""
        if ($TemplateFamily) {
            $encodedTemplate = [System.Web.HttpUtility]::UrlEncode($TemplateFamily)
            $filter = "&`$filter=templateReference/TemplateFamily eq '$encodedTemplate'"
        }

        $uri = $baseUri + $select + $topParameter + $filter

        # Get all policies
        $policies = Get-Paginated -InitialUri $uri

        # For each policy, get its assignments
        foreach ($policy in $policies) {
            $assignmentsUri = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/$($policy.id)/assignments"
            $assignments = Get-Paginated -InitialUri $assignmentsUri
            $policy | Add-Member -NotePropertyName 'assignments' -NotePropertyValue $assignments -Force
        }

        return $policies
    }
    catch {
        Write-Error "Error retrieving settings catalog policies: $_"
        return $null
    }
}

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

Write-Host "Retrieving settings catalog policies..."
$catalogPolicies = Get-SettingsCatalogPolicies -TemplateFamily $TemplateFamily -Top $Top

if ($null -ne $catalogPolicies -and $catalogPolicies.Count -gt 0) {
    Write-Host "`nFound $($catalogPolicies.Count) policies"
    $jsonOutput = $catalogPolicies | ConvertTo-Json -Depth 100
    Write-Output $jsonOutput
    
    $jsonOutput | Out-File "settings_catalog_policies_list_export.json"
    Write-Host "`nComplete data has been saved to 'settings_catalog_policies_list_export.json'"
} else {
    Write-Host "No policies found matching the specified criteria."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."