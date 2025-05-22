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
    HelpMessage="Specify the ID of the role definition to retrieve")]
    [ValidateNotNullOrEmpty()]
    [string]$RoleDefinitionId
)

# Helper function to retrieve all pages of results
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
        }
        
        # Get the next page URL if it exists
        $currentUri = $response.'@odata.nextLink'
    } while ($currentUri)

    return $allItems
}

# Helper function to retrieve a role definition by ID and all its assignments
function Get-RoleDefinitionWithAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$RoleDefinitionId
    )

    try {
        # Get the role definition
        $roleDefUri = "https://graph.microsoft.com/beta/deviceManagement/roleDefinitions/$RoleDefinitionId"
        $roleDef = Invoke-MgGraphRequest -Method GET -Uri $roleDefUri
        
        Write-Host "Retrieved role definition: $($roleDef.displayName)"
        
        # Save the base role definition to a file
        $roleDef | ConvertTo-Json -Depth 100 | Out-File "roleDefinition.json"
        Write-Host "Role definition saved to 'roleDefinition.json'"

        # Get all assignments for this role definition
        $assignmentsUri = "https://graph.microsoft.com/beta/deviceManagement/roleDefinitions/$RoleDefinitionId/roleAssignments"
        $assignmentsList = Get-Paginated -InitialUri $assignmentsUri
        
        Write-Host "Found $($assignmentsList.Count) role assignments"
        
        # Save the assignments list to a file
        $assignmentsList | ConvertTo-Json -Depth 100 | Out-File "roleAssignmentsList.json"
        Write-Host "Role assignments list saved to 'roleAssignmentsList.json'"
        
        # Get detailed info for each assignment
        $detailedAssignments = @()
        
        foreach ($assignment in $assignmentsList) {
            Write-Host "Retrieving details for assignment: $($assignment.id)"
            $assignmentUri = "https://graph.microsoft.com/beta/deviceManagement/roleAssignments/$($assignment.id)"
            $detailedAssignment = Invoke-MgGraphRequest -Method GET -Uri $assignmentUri
            $detailedAssignments += $detailedAssignment
        }
        
        # Save the detailed assignments to a file
        $detailedAssignments | ConvertTo-Json -Depth 100 | Out-File "roleAssignmentsDetailed.json"
        Write-Host "Detailed role assignments saved to 'roleAssignmentsDetailed.json'"
        
        # Return a combined object with all the information
        $result = @{
            roleDefinition = $roleDef
            assignmentsList = $assignmentsList
            detailedAssignments = $detailedAssignments
        }
        
        return $result
    }
    catch {
        Write-Error "Error retrieving role definition and assignments: $_"
        return $null
    }
}

# Script Setup
Import-Module Microsoft.Graph.Authentication

$secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
$clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret

Write-Host "Connecting to Microsoft Graph..."
Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId

Write-Host "Retrieving role definition with ID: $RoleDefinitionId"
$roleData = Get-RoleDefinitionWithAssignments -RoleDefinitionId $RoleDefinitionId

if ($null -ne $roleData) {
    Write-Host "`nRole Definition and Assignments retrieved successfully."
    
    # Save the complete data to a single file
    $roleData | ConvertTo-Json -Depth 100 | Out-File "roleDefinitionComplete.json"
    Write-Host "Complete data saved to 'roleDefinitionComplete.json'"
} else {
    Write-Host "No data found for the specified role definition ID."
}

Disconnect-MgGraph
Write-Host "Disconnected from Microsoft Graph."