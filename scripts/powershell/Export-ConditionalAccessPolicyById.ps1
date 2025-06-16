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
    HelpMessage="Specify the ID of the conditional access policy to retrieve (uses beta endpoint for full property support)")]
    [ValidateNotNullOrEmpty()]
    [string]$ConditionalAccessPolicyId
)

# Helper function to retrieve a specific conditional access policy by ID
function Get-ConditionalAccessPolicyById {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConditionalAccessPolicyId
    )

    try {
        $policyUri = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies/$ConditionalAccessPolicyId"
        $policy = Invoke-MgGraphRequest -Method GET -Uri $policyUri

        return $policy
    }
    catch {
        Write-Error "Error retrieving conditional access policy by ID: $_"
        return $null
    }
}

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

Write-Host "Retrieving conditional access policy with ID: $ConditionalAccessPolicyId"
$policyData = Get-ConditionalAccessPolicyById -ConditionalAccessPolicyId $ConditionalAccessPolicyId

if ($null -ne $policyData) {
    Write-Host "`nFull JSON output for conditional access policy:"
    $jsonString = $policyData | ConvertTo-Json -Depth 100 -Compress
    # Format the JSON for readability
    $jsonFormatted = $jsonString | ConvertFrom-Json | ConvertTo-Json -Depth 100
    
    Write-Output $jsonFormatted
    
    $jsonFormatted | Out-File "conditionalAccessPolicy.json"
    Write-Host "`nJSON output has also been saved to 'conditionalAccessPolicy.json'"
} else {
    Write-Host "No data found for the specified conditional access policy ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."