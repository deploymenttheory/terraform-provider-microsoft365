# Get-IntuneSettingsCatalogPolicies.ps1
# Script to get settings catalog policies from Intune via Microsoft Graph API and save to JSON file

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
    HelpMessage="Optional ID of a specific settings catalog policy to retrieve")]
    [string]$PolicyId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional template family to filter policies by")]
    [string]$TemplateFamily,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional filter query for policies")]
    [string]$Filter,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional select query to specify which properties to retrieve")]
    [string]$Select = "name,id,description,platforms,lastModifiedDateTime,technologies,settingCount,roleScopeTagIds,isAssigned,templateReference",
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional expand query to include related entities")]
    [string]$Expand,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Maximum number of policies to return per page (for pagination)")]
    [int]$Top = 25,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Directory path where the SettingsCatalogPolicies.json file will be created")]
    [string]$OutputDirectory,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Whether to include policy assignments")]
    [bool]$IncludeAssignments = $true
)

# Function to authenticate and get access token
function Connect-MicrosoftGraph {
    param (
        [Parameter(Mandatory=$true)]
        [string]$TenantId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientId,
        
        [Parameter(Mandatory=$true)]
        [string]$ClientSecret
    )
    
    try {
        Write-Host "Connecting to Microsoft Graph..." -ForegroundColor Cyan
        
        # Create secure credential
        $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
        $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
        
        # Connect to Microsoft Graph
        Import-Module Microsoft.Graph.Authentication
        Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
        
        Write-Host "✅ Connected to Microsoft Graph" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Error connecting to Microsoft Graph: $_" -ForegroundColor Red
        throw
    }
}

# Function to retrieve all pages of data from a paginated API
function Get-AllPaginatedResults {
    param (
        [Parameter(Mandatory=$true)]
        [string]$InitialUri
    )
    
    try {
        $allItems = @()
        $currentUri = $InitialUri
        $pageCount = 1
        
        do {
            Write-Host "  Retrieving page $pageCount..." -ForegroundColor Gray
            $response = Invoke-MgGraphRequest -Method GET -Uri $currentUri
            
            if ($response.value) {
                $allItems += $response.value
                Write-Host "    Got $($response.value.Count) items" -ForegroundColor Gray
            } else {
                $allItems += $response
                Write-Host "    Got single item" -ForegroundColor Gray
            }
            
            # Get the next page URL if it exists
            $currentUri = $response.'@odata.nextLink'
            $pageCount++
        } while ($currentUri)
        
        return $allItems
    }
    catch {
        Write-Host "❌ Error retrieving paginated results: $_" -ForegroundColor Red
        throw
    }
}

# Function to get settings catalog policies from Intune and save as JSON file
function Get-SettingsCatalogPoliciesToFile {
    param (
        [Parameter(Mandatory=$false)]
        [string]$PolicyId,
        
        [Parameter(Mandatory=$false)]
        [string]$TemplateFamily,
        
        [Parameter(Mandatory=$false)]
        [string]$Filter,
        
        [Parameter(Mandatory=$false)]
        [string]$Select,
        
        [Parameter(Mandatory=$false)]
        [string]$Expand,
        
        [Parameter(Mandatory=$false)]
        [int]$Top = 25,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFilePath,
        
        [Parameter(Mandatory=$false)]
        [bool]$IncludeAssignments = $true
    )
    
    try {
        $baseUrl = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
        $queryParams = @()
        
        # Build the URL based on parameters
        if ($PolicyId) {
            $url = "$baseUrl/$PolicyId"
        } else {
            $url = $baseUrl
            
            if ($Select) {
                $queryParams += "`$select=$([System.Web.HttpUtility]::UrlEncode($Select))"
            }
            
            if ($Expand) {
                $queryParams += "`$expand=$([System.Web.HttpUtility]::UrlEncode($Expand))"
            }
            
            if ($Top -gt 0) {
                $queryParams += "`$top=$Top"
            }
            
            # Add filter based on input parameters
            $filterParts = @()
            
            if ($Filter) {
                $filterParts += $Filter
            }
            
            if ($TemplateFamily) {
                $encodedTemplate = [System.Web.HttpUtility]::UrlEncode($TemplateFamily)
                $filterParts += "templateReference/TemplateFamily eq '$encodedTemplate'"
            }
            
            if ($filterParts.Count -gt 0) {
                $queryParams += "`$filter=" + ($filterParts -join " and ")
            }
            
            if ($queryParams.Count -gt 0) {
                $url += "?" + ($queryParams -join "&")
            }
        }
        
        Write-Host "Retrieving settings catalog policies from Intune..." -ForegroundColor Cyan
        Write-Host "URL: $url" -ForegroundColor Gray
        
        # Get all policies with pagination support
        $policies = Get-AllPaginatedResults -InitialUri $url
        
        # Get assignments if requested and if we have policies
        if ($IncludeAssignments -and $policies) {
            Write-Host "Retrieving assignments for each policy..." -ForegroundColor Cyan
            $policyCount = $policies.Count
            
            # Handle single policy case that doesn't come as an array
            if ($null -eq $policyCount) {
                $policyCount = 1
                $policies = @($policies)
            }
            
            for ($i = 0; $i -lt $policyCount; $i++) {
                $policy = $policies[$i]
                $policyName = $policy.name
                $policyId = $policy.id
                
                Write-Host "  Getting assignments for policy '$policyName' ($($i+1)/$policyCount)..." -ForegroundColor Gray
                $assignmentsUri = "$baseUrl/$policyId/assignments"
                $assignments = Get-AllPaginatedResults -InitialUri $assignmentsUri
                
                # Add assignments to the policy object
                $policy | Add-Member -NotePropertyName 'assignments' -NotePropertyValue $assignments -Force
            }
        }
        
        # Save policies to file
        $prettyJson = ConvertTo-Json -InputObject $policies -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFilePath -Encoding utf8
        
        # Get policy count for summary
        if ($PolicyId) {
            $policyCount = 1
        } else {
            $policyCount = if ($policies) { $policies.Count } else { 0 }
        }
        
        return $policyCount
    }
    catch {
        Write-Host "❌ Error retrieving settings catalog policies: $_" -ForegroundColor Red
        throw
    }
}

# Main script execution
try {
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Ensure the output directory exists
    if (-not (Test-Path -Path $OutputDirectory)) {
        New-Item -Path $OutputDirectory -ItemType Directory -Force | Out-Null
        Write-Host "Created output directory: $OutputDirectory" -ForegroundColor Yellow
    }
    
    # Create the fixed output filename
    $outputFilePath = Join-Path -Path $OutputDirectory -ChildPath "SettingsCatalogPolicies.json"
    
    # Get settings catalog policies and save to file
    Write-Host "`n📝 Retrieving Intune settings catalog policies..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $policyCount = Get-SettingsCatalogPoliciesToFile -PolicyId $PolicyId -TemplateFamily $TemplateFamily -Filter $Filter -Select $Select -Expand $Expand -Top $Top -OutputFilePath $outputFilePath -IncludeAssignments $IncludeAssignments
    
    # Summary message
    if ($PolicyId) {
        Write-Host "`n✨ Successfully saved settings catalog policy details to: $outputFilePath" -ForegroundColor Green
    } else {
        Write-Host "`n✨ Successfully saved $policyCount settings catalog policies to: $outputFilePath" -ForegroundColor Green
    }
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "`nDisconnecting from Microsoft Graph..." -ForegroundColor Cyan
    Disconnect-MgGraph > $null 2>&1
    Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
}