<#
.SYNOPSIS
    Retrieves uploaded ADMX/ADML definition files from Microsoft Intune

.DESCRIPTION
    This script queries the Microsoft Graph API to retrieve custom Group Policy
    definition files (ADMX templates) that have been uploaded to Microsoft Intune.
    
    It supports two modes:
    - LIST mode: Retrieves all uploaded definition files (1 API call)
    - GET BY ID mode: Retrieves specific file details and definitions (2 API calls)
    
    LIST MODE API Call:
    1. GET /deviceManagement/groupPolicyUploadedDefinitionFiles
    
    GET BY ID MODE API Call:
    1. GET /deviceManagement/groupPolicyUploadedDefinitionFiles('{id}')?$expand=GroupPolicyOperations

.PARAMETER TenantId
    Entra ID Tenant ID

.PARAMETER ClientId
    Application (Client) ID of the Entra ID app registration

.PARAMETER ClientSecret
    Client secret for authentication

.PARAMETER OutputDirectory
    Directory path where JSON responses will be saved

.PARAMETER FileId
    Optional. Specific file ID to retrieve details for (enables GET BY ID mode)

.EXAMPLE
    # LIST MODE - Get all uploaded ADMX templates
    pwsh Get-GroupPolicyUploadedDefinitionFiles.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -OutputDirectory "C:\temp\uploaded_admx"
    
    Returns: step1_list_all_files.json

.EXAMPLE
    # GET BY ID MODE - Get specific file details with operations
    pwsh Get-GroupPolicyUploadedDefinitionFiles.ps1 `
        -TenantId "your-tenant-id" `
        -ClientId "your-client-id" `
        -ClientSecret "your-secret" `
        -OutputDirectory "C:\temp\uploaded_admx" `
        -FileId "7dd588c0-c074-483f-bcf6-711f6b4116ba"
    
    Returns: get_file_by_id_with_operations.json

.EXAMPLE
    # Real-world example
    pwsh Get-GroupPolicyUploadedDefinitionFiles.ps1 `
        -TenantId "00000000-0000-0000-0000-000000000000" `
        -ClientId "00000000-0000-0000-0000-000000000000" `
        -ClientSecret "your-secret" `
        -OutputDirectory "/Users/username/localtesting/uploaded_admx" `
        -FileId "2bafdbc3-3c8a-4aa3-8ab7-bfa37136e1a6"

.NOTES
    File Name      : Get-GroupPolicyUploadedDefinitionFiles.ps1
    Prerequisite   : Microsoft.Graph.Authentication PowerShell module
    Copyright      : Deployment Theory
    
    API Permissions Required:
    - DeviceManagementConfiguration.Read.All (minimum)
    - DeviceManagementConfiguration.ReadWrite.All (for full access)

.LINK
    https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicyuploadeddefinitionfile-list
    https://learn.microsoft.com/en-us/graph/api/intune-grouppolicy-grouppolicyuploadeddefinitionfile-get
#>

[CmdletBinding()]
param (
    [Parameter(Mandatory=$true)]
    [string]$TenantId,

    [Parameter(Mandatory=$true)]
    [string]$ClientId,
    
    [Parameter(Mandatory=$true)]
    [string]$ClientSecret,
    
    [Parameter(Mandatory=$true)]
    [string]$OutputDirectory,
    
    [Parameter(Mandatory=$false,
    HelpMessage="Specific file ID to retrieve (enables GET BY ID mode)")]
    [string]$FileId = ""
)

# Import module
Import-Module Microsoft.Graph.Authentication

# Function to authenticate
function Connect-MicrosoftGraph {
    param (
        [string]$TenantId,
        [string]$ClientId,
        [string]$ClientSecret
    )
    
    try {
        Write-Host "🔐 Authenticating to Microsoft Graph..." -ForegroundColor Cyan
        
        $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
        $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
        
        Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId -NoWelcome
        
        Write-Host "✅ Connected successfully" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Authentication failed: $_" -ForegroundColor Red
        throw
    }
}

# Function to list all uploaded definition files
function Get-AllUploadedDefinitionFiles {
    param (
        [string]$OutputFile
    )
    
    try {
        Write-Host "`n📍 API CALL 1: LIST All Uploaded Definition Files" -ForegroundColor Yellow
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        Write-Host "Mode: LIST (retrieving all files)" -ForegroundColor Cyan
        
        # Build URL with expanded properties
        $select = "id,fileName,status,defaultLanguageCode,languageCodes,targetPrefix,targetNamespace,policyType,revision,uploadDateTime,lastModifiedDateTime"
        $url = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles?`$select=$select"
        
        Write-Host "Endpoint: GET /deviceManagement/groupPolicyUploadedDefinitionFiles" -ForegroundColor Gray
        Write-Host "URL: $url" -ForegroundColor DarkGray
        
        $response = Invoke-MgGraphRequest -Uri $url -Method GET
        
        if ($response.value) {
            $files = $response.value
            $fileCount = $files.Count
            
            Write-Host "`n✅ Found $fileCount uploaded definition file(s)" -ForegroundColor Green
            
            # Display summary
            Write-Host "`n📋 Summary:" -ForegroundColor Cyan
            foreach ($file in $files) {
                Write-Host "  • $($file.fileName)" -ForegroundColor White
                Write-Host "    ID: $($file.id)" -ForegroundColor Gray
                Write-Host "    Status: $($file.status)" -ForegroundColor Gray
                Write-Host "    Namespace: $($file.targetNamespace)" -ForegroundColor Gray
                Write-Host "    Language: $($file.defaultLanguageCode)" -ForegroundColor Gray
                Write-Host ""
            }
            
            # Save to file
            $response | ConvertTo-Json -Depth 10 | Out-File -FilePath $OutputFile -Encoding UTF8
            Write-Host "💾 Saved to: $OutputFile" -ForegroundColor Green
            
            return $files
        } else {
            Write-Host "⚠️  No uploaded definition files found" -ForegroundColor Yellow
            return @()
        }
    }
    catch {
        Write-Host "❌ Failed to retrieve uploaded definition files: $_" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        throw
    }
}

# Function to get specific file by ID with operations
function Get-DefinitionFileById {
    param (
        [string]$FileId,
        [string]$OutputFile
    )
    
    try {
        Write-Host "`n📍 API CALL: GET File Details by ID with Operations" -ForegroundColor Yellow
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
        Write-Host "Mode: GET BY ID (specific file with expanded operations)" -ForegroundColor Cyan
        Write-Host "File ID: $FileId" -ForegroundColor White
        
        # Use parentheses syntax with expand
        $url = "https://graph.microsoft.com/beta/deviceManagement/groupPolicyUploadedDefinitionFiles('$FileId')?`$expand=GroupPolicyOperations"
        
        Write-Host "Endpoint: GET /deviceManagement/groupPolicyUploadedDefinitionFiles('{id}')" -ForegroundColor Gray
        Write-Host "Expand: GroupPolicyOperations" -ForegroundColor Gray
        Write-Host "URL: $url" -ForegroundColor DarkGray
        
        $response = Invoke-MgGraphRequest -Uri $url -Method GET
        
        if ($response) {
            Write-Host "`n✅ File found: $($response.fileName)" -ForegroundColor Green
            
            # Display details
            Write-Host "`n📋 File Details:" -ForegroundColor Cyan
            Write-Host "  File Name: $($response.fileName)" -ForegroundColor White
            Write-Host "  Status: $($response.status)" -ForegroundColor Gray
            Write-Host "  Namespace: $($response.targetNamespace)" -ForegroundColor Gray
            Write-Host "  Prefix: $($response.targetPrefix)" -ForegroundColor Gray
            Write-Host "  Language: $($response.defaultLanguageCode)" -ForegroundColor Gray
            Write-Host "  Policy Type: $($response.policyType)" -ForegroundColor Gray
            Write-Host "  Upload Date: $($response.uploadDateTime)" -ForegroundColor Gray
            Write-Host "  Last Modified: $($response.lastModifiedDateTime)" -ForegroundColor Gray
            
            # Display operations if present
            if ($response.groupPolicyOperations) {
                $opsCount = $response.groupPolicyOperations.Count
                Write-Host "`n📋 Group Policy Operations: $opsCount" -ForegroundColor Cyan
                foreach ($op in $response.groupPolicyOperations) {
                    Write-Host "  • Operation Type: $($op.operationType)" -ForegroundColor White
                    Write-Host "    Status: $($op.operationStatus)" -ForegroundColor Gray
                    Write-Host "    Last Modified: $($op.lastModifiedDateTime)" -ForegroundColor Gray
                }
            }
            
            # Save to file
            $response | ConvertTo-Json -Depth 10 | Out-File -FilePath $OutputFile -Encoding UTF8
            Write-Host "`n💾 Saved to: $OutputFile" -ForegroundColor Green
            
            return $response
        }
    }
    catch {
        Write-Host "❌ Failed to retrieve file details: $_" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        throw
    }
}

# Main execution
try {
    Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
    Write-Host "║  Get Uploaded Group Policy ADMX Definition Files              ║" -ForegroundColor Cyan
    Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
    
    # Create output directory if it doesn't exist
    if (-not (Test-Path -Path $OutputDirectory)) {
        New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
        Write-Host "📁 Created output directory: $OutputDirectory" -ForegroundColor Green
    }
    
    # Authenticate
    Connect-MicrosoftGraph -TenantId $TenantId -ClientId $ClientId -ClientSecret $ClientSecret
    
    # Determine mode
    $isListMode = [string]::IsNullOrWhiteSpace($FileId)
    
    if ($isListMode) {
        # LIST MODE - Get all files
        $outputFile = Join-Path -Path $OutputDirectory -ChildPath "step1_list_all_files.json"
        $files = Get-AllUploadedDefinitionFiles -OutputFile $outputFile
        
        # Final summary
        Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
        Write-Host "║                    ✅ EXECUTION COMPLETE                       ║" -ForegroundColor Green
        Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Green
        Write-Host "📊 Mode: LIST" -ForegroundColor White
        Write-Host "📊 Total files found: $($files.Count)" -ForegroundColor White
        Write-Host "📁 Output location: $OutputDirectory" -ForegroundColor White
        Write-Host "📄 Files created: step1_list_all_files.json" -ForegroundColor White
    }
    else {
        # GET BY ID MODE - Get specific file with operations
        $outputFile = Join-Path -Path $OutputDirectory -ChildPath "get_file_by_id_with_operations.json"
        $fileDetails = Get-DefinitionFileById -FileId $FileId -OutputFile $outputFile
        
        if ($fileDetails) {
            # Count operations if present
            $opsCount = 0
            if ($fileDetails.groupPolicyOperations) {
                $opsCount = $fileDetails.groupPolicyOperations.Count
            }
            
            # Final summary
            Write-Host "`n╔═══════════════════════════════════════════════════════════════╗" -ForegroundColor Green
            Write-Host "║                    ✅ EXECUTION COMPLETE                       ║" -ForegroundColor Green
            Write-Host "╚═══════════════════════════════════════════════════════════════╝" -ForegroundColor Green
            Write-Host "📊 Mode: GET BY ID" -ForegroundColor White
            Write-Host "📊 File: $($fileDetails.fileName)" -ForegroundColor White
            Write-Host "📊 Status: $($fileDetails.status)" -ForegroundColor White
            Write-Host "📊 Operations found: $opsCount" -ForegroundColor White
            Write-Host "📁 Output location: $OutputDirectory" -ForegroundColor White
            Write-Host "📄 File created: get_file_by_id_with_operations.json" -ForegroundColor White
        }
    }
    
    Write-Host "`n✨ Done!" -ForegroundColor Cyan
}
catch {
    Write-Host "`n❌ Script execution failed: $_" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}
finally {
    # Disconnect from Graph
    try {
        Disconnect-MgGraph | Out-Null
        Write-Host "`n🔓 Disconnected from Microsoft Graph" -ForegroundColor Gray
    }
    catch {
        # Ignore disconnect errors
    }
}
