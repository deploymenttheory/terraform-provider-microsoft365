# Get-IntuneResourceOperations.ps1
# Script to get all resource operations from Intune via Microsoft Graph API and save to JSON file

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
    HelpMessage="Optional ID of a specific resource operation to retrieve")]
    [string]$ResourceOperationId,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional filter query for resource operations")]
    [string]$Filter,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional select query to specify which properties to retrieve")]
    [string]$Select,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Optional expand query to include related entities")]
    [string]$Expand,
    
    [Parameter(Mandatory=$true,
    HelpMessage="Directory path where the ListResourceOperations.json file will be created")]
    [string]$OutputDirectory
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

# Function to get all resource operations from Intune and save as JSON file
function Get-ResourceOperationsToFile {
    param (
        [Parameter(Mandatory=$false)]
        [string]$ResourceOperationId,
        
        [Parameter(Mandatory=$false)]
        [string]$Filter,
        
        [Parameter(Mandatory=$false)]
        [string]$Select,
        
        [Parameter(Mandatory=$false)]
        [string]$Expand,
        
        [Parameter(Mandatory=$true)]
        [string]$OutputFilePath
    )
    
    try {
        $baseUrl = "https://graph.microsoft.com/beta/deviceManagement/resourceOperations"
        $queryParams = @()
        
        # Build the URL based on parameters
        if ($ResourceOperationId) {
            # If a specific ID is requested, get that single resource operation
            $url = "$baseUrl/$ResourceOperationId"
            
            Write-Host "Retrieving specific resource operation from Intune..." -ForegroundColor Cyan
            Write-Host "URL: $url" -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method GET -Uri $url
            $operationCount = 1
        } else {
            # Get all resource operations without specifying a limit
            $url = $baseUrl
            
            if ($Filter) {
                $queryParams += "`$filter=$([System.Web.HttpUtility]::UrlEncode($Filter))"
            }
            
            if ($Select) {
                $queryParams += "`$select=$([System.Web.HttpUtility]::UrlEncode($Select))"
            }
            
            if ($Expand) {
                $queryParams += "`$expand=$([System.Web.HttpUtility]::UrlEncode($Expand))"
            }
            
            if ($queryParams.Count -gt 0) {
                $url += "?" + ($queryParams -join "&")
            }
            
            Write-Host "Retrieving all resource operations from Intune..." -ForegroundColor Cyan
            Write-Host "URL: $url" -ForegroundColor Gray
            
            $response = Invoke-MgGraphRequest -Method GET -Uri $url
            $operationCount = if ($response.value) { $response.value.Count } else { 0 }
        }
        
        # Save the response to file
        $prettyJson = ConvertTo-Json -InputObject $response -Depth 10 -Compress:$false
        $prettyJson | Out-File -FilePath $OutputFilePath -Encoding utf8
        
        return $operationCount
    }
    catch {
        Write-Host "❌ Error retrieving resource operations: $_" -ForegroundColor Red
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
    $outputFilePath = Join-Path -Path $OutputDirectory -ChildPath "ListResourceOperations.json"
    
    # Get resource operations and save to file
    Write-Host "`n🔑 Retrieving Intune resource operations..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $operationCount = Get-ResourceOperationsToFile -ResourceOperationId $ResourceOperationId -Filter $Filter -Select $Select -Expand $Expand -OutputFilePath $outputFilePath
    
    # Summary message
    if ($ResourceOperationId) {
        Write-Host "`n✨ Successfully saved resource operation details to: $outputFilePath" -ForegroundColor Green
    } else {
        Write-Host "`n✨ Successfully saved $operationCount resource operations to: $outputFilePath" -ForegroundColor Green
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