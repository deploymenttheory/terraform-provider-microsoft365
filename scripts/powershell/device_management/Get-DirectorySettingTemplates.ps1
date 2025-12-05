# Get-DirectorySettingTemplates.ps1
# Script to get directory setting templates from Microsoft Graph API and save to JSON file
# Based on: https://learn.microsoft.com/en-us/graph/api/directorysettingtemplate-list

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
    HelpMessage="Optional ID of a specific directory setting template to retrieve")]
    [string]$TemplateId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional select query to specify which properties to retrieve")]
    [string]$Select,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Path to output JSON file")]
    [string]$OutputFile
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

# Function to get directory setting templates and save as JSON file
function Get-DirectorySettingTemplatesToFile {
    param (
        [Parameter(Mandatory=$false)]
        [string]$TemplateId,
        
        [Parameter(Mandatory=$false)]
        [string]$Select,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFile
    )
    
    try {
        $baseUrl = "https://graph.microsoft.com/beta/directorySettingTemplates"
        $queryParams = @()
        
        # Build the URL based on parameters
        if ($TemplateId) {
            $url = "$baseUrl/$TemplateId"
            
            if ($Select) {
                $queryParams += "`$select=$([System.Web.HttpUtility]::UrlEncode($Select))"
            }
        } else {
            $url = $baseUrl
            
            if ($Select) {
                $queryParams += "`$select=$([System.Web.HttpUtility]::UrlEncode($Select))"
            }
        }
        
        if ($queryParams.Count -gt 0) {
            $url += "?" + ($queryParams -join "&")
        }
        
        Write-Host "Retrieving directory setting templates..." -ForegroundColor Cyan
        Write-Host "URL: $url" -ForegroundColor Gray
        
        # Get response and save directly to file
        $response = Invoke-MgGraphRequest -Method GET -Uri $url
        $prettyJson = ConvertTo-Json -InputObject $response -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFile -Encoding utf8
        
        # Get template count for summary
        if ($TemplateId) {
            $templateCount = 1
            $templateInfo = $response
        } else {
            $templateCount = if ($response.value) { $response.value.Count } else { 0 }
            $templateInfo = $response.value
        }
        
        return @{
            Count = $templateCount
            Templates = $templateInfo
        }
    }
    catch {
        Write-Host "❌ Error retrieving directory setting templates: $_" -ForegroundColor Red
        throw
    }
}

# Function to display template summary
function Show-TemplatesSummary {
    param (
        [Parameter(Mandatory=$true)]
        [object]$Templates
    )
    
    Write-Host "`n📋 Directory Setting Templates Summary:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Templates -is [array]) {
        foreach ($template in $Templates) {
            Write-Host "`n📄 Template: " -NoNewline -ForegroundColor Yellow
            Write-Host "$($template.displayName)" -ForegroundColor White
            Write-Host "   ID: $($template.id)" -ForegroundColor Gray
            Write-Host "   Description: $($template.description)" -ForegroundColor Gray
            
            if ($template.values) {
                Write-Host "   Settings Count: $($template.values.Count)" -ForegroundColor Gray
                Write-Host "   Available Settings:" -ForegroundColor Gray
                foreach ($value in $template.values) {
                    Write-Host "     • $($value.name) ($($value.type)): $($value.description)" -ForegroundColor DarkGray
                }
            }
        }
    } else {
        Write-Host "`n📄 Template: " -NoNewline -ForegroundColor Yellow
        Write-Host "$($Templates.displayName)" -ForegroundColor White
        Write-Host "   ID: $($Templates.id)" -ForegroundColor Gray
        Write-Host "   Description: $($Templates.description)" -ForegroundColor Gray
        
        if ($Templates.values) {
            Write-Host "   Settings Count: $($Templates.values.Count)" -ForegroundColor Gray
            Write-Host "   Available Settings:" -ForegroundColor Gray
            foreach ($value in $Templates.values) {
                Write-Host "     • $($value.name) ($($value.type)): $($value.description)" -ForegroundColor DarkGray
            }
        }
    }
}

# Main script execution
try {
    # Connect to Microsoft Graph
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Get directory setting templates and save to file
    Write-Host "`n🔧 Retrieving directory setting templates..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $result = Get-DirectorySettingTemplatesToFile -TemplateId $TemplateId -Select $Select -OutputFile $OutputFile
    
    # Display summary
    Show-TemplatesSummary -Templates $result.Templates
    
    # Summary message
    Write-Host "`n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
    if ($TemplateId) {
        Write-Host "✨ Successfully saved directory setting template details to: $OutputFile" -ForegroundColor Green
    } else {
        Write-Host "✨ Successfully saved $($result.Count) directory setting templates to: $OutputFile" -ForegroundColor Green
    }
    Write-Host "`n💡 Common Template IDs:" -ForegroundColor Cyan
    Write-Host "   • Group.Unified: 62375ab9-6b52-47ed-826b-58e47e0e304b" -ForegroundColor Gray
    Write-Host "   • Group.Unified.Guest: 08d542b9-071f-4e16-94b0-74abb372e3d9" -ForegroundColor Gray
    Write-Host "   • Prohibited Names Settings: 80661d51-be2f-4d46-9713-98a2fcaec5bc" -ForegroundColor Gray
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

