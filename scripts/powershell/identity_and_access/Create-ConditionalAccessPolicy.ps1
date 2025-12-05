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
    HelpMessage="Test policy name suffix (will be prefixed with 'TEST-')")]
    [string]$PolicyNameSuffix = "PowerShell-Creation-Test",
    
    [Parameter(Mandatory=$true,
    HelpMessage="Export results to JSON file")]
    [bool]$ExportToJson
)

# Import required modules
Import-Module Microsoft.Graph.Authentication

# Test payload based on Terraform examples - simplified version for testing
$testPolicyPayload = @'
{
    "displayName": "TEST-PowerShell-Creation-Test",
    "state": "disabled",
    "conditions": {
        "applications": {
            "includeApplications": ["All"]
        },
        "users": {
            "includeUsers": ["All"],
            "excludeGroups": ["11111111-1111-1111-1111-111111111111"]
        },
        "clientAppTypes": ["browser", "mobileAppsAndDesktopClients"],
        "locations": {
            "includeLocations": ["All"]
        }
    },
    "grantControls": {
        "operator": "OR",
        "builtInControls": ["mfa"]
    }
}
'@

# Function to create conditional access policy
function New-ConditionalAccessPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyJson,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyName
    )
    
    try {
        # Parse and modify the JSON payload
        $policyObject = $PolicyJson | ConvertFrom-Json
        $policyObject.displayName = "TEST-$PolicyName"
        
        # Convert back to JSON
        $finalPayload = $policyObject | ConvertTo-Json -Depth 10
        
        Write-Host "🔄 Creating conditional access policy..." -ForegroundColor Cyan
        Write-Host "   Policy Name: TEST-$PolicyName" -ForegroundColor Gray
        Write-Host "   Endpoint: https://graph.microsoft.com/beta/identity/conditionalAccess/policies" -ForegroundColor Gray
        Write-Host ""
        
        Write-Host "📋 Request Payload:" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host $finalPayload -ForegroundColor Gray
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        Write-Host ""
        
        # Make the POST request
        $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies"
        $response = Invoke-MgGraphRequest -Method POST -Uri $uri -Body $finalPayload -ContentType "application/json"
        
        Write-Host "✅ Policy creation request completed!" -ForegroundColor Green
        Write-Host ""
        
        return $response
    }
    catch {
        Write-Host "❌ Error creating conditional access policy: $_" -ForegroundColor Red
        Write-Host ""
        
        # Enhanced error handling
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            $statusDescription = $_.Exception.Response.StatusDescription
            Write-Host "   Status Code: $statusCode" -ForegroundColor Red
            Write-Host "   Status Description: $statusDescription" -ForegroundColor Red
            
            # Try to get the response content
            try {
                $responseContent = $_.Exception.Response.GetResponseStream()
                if ($responseContent) {
                    $reader = [System.IO.StreamReader]::new($responseContent)
                    $errorDetails = $reader.ReadToEnd()
                    $reader.Close()
                    Write-Host "   Error Details: $errorDetails" -ForegroundColor Red
                }
            }
            catch {
                Write-Host "   Could not read error details" -ForegroundColor Red
            }
        }
        
        throw
    }
}

# Function to get policy by name to verify creation
function Get-PolicyByName {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyName
    )
    
    try {
        Write-Host "🔍 Searching for created policy..." -ForegroundColor Cyan
        Write-Host "   Looking for: TEST-$PolicyName" -ForegroundColor Gray
        
        $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies"
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        if ($response.value) {
            $foundPolicy = $response.value | Where-Object { $_.displayName -eq "TEST-$PolicyName" }
            if ($foundPolicy) {
                Write-Host "✅ Found created policy!" -ForegroundColor Green
                Write-Host "   Policy ID: $($foundPolicy.id)" -ForegroundColor Green
                return $foundPolicy
            } else {
                Write-Host "❌ Policy not found in list" -ForegroundColor Red
                return $null
            }
        } else {
            Write-Host "❌ No policies returned from API" -ForegroundColor Red
            return $null
        }
    }
    catch {
        Write-Host "❌ Error searching for policy: $_" -ForegroundColor Red
        return $null
    }
}

# Function to delete test policy (cleanup)
function Remove-TestPolicy {
    param (
        [Parameter(Mandatory=$true)]
        [string]$PolicyId
    )
    
    try {
        Write-Host "🗑️  Cleaning up test policy..." -ForegroundColor Yellow
        Write-Host "   Policy ID: $PolicyId" -ForegroundColor Gray
        
        $uri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies/$PolicyId"
        Invoke-MgGraphRequest -Method DELETE -Uri $uri
        
        Write-Host "✅ Test policy deleted successfully" -ForegroundColor Green
        Write-Host ""
    }
    catch {
        Write-Host "❌ Error deleting test policy: $_" -ForegroundColor Red
        Write-Host "   Manual cleanup may be required for policy ID: $PolicyId" -ForegroundColor Yellow
        Write-Host ""
    }
}

# Function to export response to JSON
function Export-ResponseToJson {
    param (
        [Parameter(Mandatory=$true)]
        $Response,
        
        [Parameter(Mandatory=$true)]
        [string]$PolicyName,
        
        [Parameter(Mandatory=$false)]
        [string]$Suffix = "CreateResponse"
    )
    
    try {
        # Create output directory if it doesn't exist
        $outputDir = Join-Path -Path (Get-Location) -ChildPath "output"
        if (-not (Test-Path -Path $outputDir)) {
            New-Item -Path $outputDir -ItemType Directory | Out-Null
            Write-Host "📁 Created output directory: $outputDir" -ForegroundColor Gray
        }
        
        # Generate timestamp for filename
        $timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
        $fileName = "ConditionalAccessPolicy_${PolicyName}_${Suffix}_${timestamp}.json"
        $filePath = Join-Path -Path $outputDir -ChildPath $fileName
        
        $Response | ConvertTo-Json -Depth 10 | Out-File -FilePath $filePath -Encoding UTF8
        
        Write-Host "💾 Exported response to: $filePath" -ForegroundColor Green
        return $filePath
    }
    catch {
        Write-Host "❌ Error exporting response to JSON: $_" -ForegroundColor Red
        return $null
    }
}

# Function to display response details
function Show-ResponseDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Response,
        
        [Parameter(Mandatory=$true)]
        [string]$ResponseType
    )
    
    Write-Host "📋 $ResponseType Response Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Response) {
        # Check if response is empty
        if ($Response -eq $null -or ($Response | Get-Member | Measure-Object).Count -eq 0) {
            Write-Host "   ⚠️  Response is empty or null" -ForegroundColor Yellow
        } else {
            # Display response properties
            $Response | Get-Member -MemberType Properties | ForEach-Object {
                $propertyName = $_.Name
                $propertyValue = $Response.$propertyName
                
                if ($propertyValue -ne $null) {
                    Write-Host "   • $propertyName : $propertyValue" -ForegroundColor Green
                } else {
                    Write-Host "   • $propertyName : <null>" -ForegroundColor Gray
                }
            }
        }
        
        # Convert to JSON to see full structure
        Write-Host ""
        Write-Host "📄 Full Response JSON:" -ForegroundColor Cyan
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
        $jsonOutput = $Response | ConvertTo-Json -Depth 10
        if ($jsonOutput -eq "{}" -or $jsonOutput -eq "null") {
            Write-Host "   ⚠️  Response converts to empty JSON: $jsonOutput" -ForegroundColor Yellow
        } else {
            Write-Host $jsonOutput -ForegroundColor Gray
        }
    } else {
        Write-Host "   ⚠️  Response is null" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Script Setup
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Create the test policy
    Write-Host "🚀 Starting Conditional Access Policy Creation Test" -ForegroundColor Magenta
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Magenta
    Write-Host ""
    
    $createResponse = New-ConditionalAccessPolicy -PolicyJson $testPolicyPayload -PolicyName $PolicyNameSuffix
    
    # Display create response details
    Show-ResponseDetails -Response $createResponse -ResponseType "CREATE"
    
    # Export create response if requested
    if ($ExportToJson) {
        Export-ResponseToJson -Response $createResponse -PolicyName $PolicyNameSuffix -Suffix "CreateResponse"
    }
    
    # Wait a moment for the policy to be available
    Write-Host "⏳ Waiting 5 seconds for policy to be available..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    
    # Try to find the created policy
    $foundPolicy = Get-PolicyByName -PolicyName $PolicyNameSuffix
    
    if ($foundPolicy) {
        # Display found policy details
        Show-ResponseDetails -Response $foundPolicy -ResponseType "FOUND POLICY"
        
        # Export found policy if requested
        if ($ExportToJson) {
            Export-ResponseToJson -Response $foundPolicy -PolicyName $PolicyNameSuffix -Suffix "FoundPolicy"
        }
        
        # Ask if user wants to delete the test policy
        Write-Host "❓ Do you want to delete the test policy? (y/N): " -ForegroundColor Yellow -NoNewline
        $deleteChoice = Read-Host
        
        if ($deleteChoice -eq "y" -or $deleteChoice -eq "Y") {
            Remove-TestPolicy -PolicyId $foundPolicy.id
        } else {
            Write-Host "⚠️  Test policy left in place. Manual cleanup required." -ForegroundColor Yellow
            Write-Host "   Policy ID: $($foundPolicy.id)" -ForegroundColor Yellow
            Write-Host "   Policy Name: $($foundPolicy.displayName)" -ForegroundColor Yellow
        }
    } else {
        Write-Host "❌ Could not find the created policy for verification" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "🎉 Test completed!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Script execution failed: $_" -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Microsoft Graph
    Write-Host "🔌 Disconnecting from Microsoft Graph..." -ForegroundColor Cyan
    try {
        Disconnect-MgGraph 2>$null
        Write-Host "✅ Disconnected from Microsoft Graph" -ForegroundColor Green
    }
    catch {
        # Ignore disconnect errors
    }
} 