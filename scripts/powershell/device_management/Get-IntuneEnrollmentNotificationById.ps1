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
    HelpMessage="Intune Enrollment Notification Configuration ID")]
    [ValidateNotNullOrEmpty()]
    [string]$ConfigurationId
)

#   Usage:
# .\Get-IntuneEnrollmentNotificationById.ps1 -TenantId "your-tenant-id" -ClientId "your-client-id" -ClientSecret "your-client-secret" -ConfigurationId "config-id_EnrollmentNotificationsConfiguration"


# Import required modules
Import-Module Microsoft.Graph.Authentication

# Function to get enrollment notification configuration
function Get-EnrollmentNotificationConfiguration {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigId
    )
    
    try {
        # Step 1: GET enrollment notification configuration
        $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations('$ConfigId')"
        Write-Host "🔍 Getting enrollment notification configuration..." -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "❌ Error getting enrollment notification configuration: $_" -ForegroundColor Red
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

# Function to get enrollment notification configuration assignments
function Get-EnrollmentNotificationAssignments {
    param (
        [Parameter(Mandatory=$true)]
        [string]$ConfigId
    )
    
    try {
        # Step 2: GET assignments separately (enrollment configurations don't support expand for assignments)
        $uri = "https://graph.microsoft.com/beta/deviceManagement/deviceEnrollmentConfigurations('$ConfigId')/assignments"
        Write-Host "🔍 Getting enrollment notification assignments..." -ForegroundColor Cyan
        Write-Host "   Configuration ID: $ConfigId" -ForegroundColor Gray
        Write-Host "   Endpoint: $uri" -ForegroundColor Gray
        Write-Host ""
        
        $response = Invoke-MgGraphRequest -Method GET -Uri $uri
        
        return $response
    }
    catch {
        Write-Host "⚠️ Warning: Could not fetch assignments for configuration: $_" -ForegroundColor Yellow
        Write-Host "   This may be normal if no assignments exist" -ForegroundColor Yellow
        Write-Host ""
        
        # Return empty structure for assignments if not found
        return @{ value = @() }
    }
}

# Function to get notification message templates with localized messages
function Get-NotificationMessageTemplates {
    param (
        [Parameter(Mandatory=$true)]
        [array]$TemplateIds
    )
    
    $templates = @()
    
    foreach ($templateId in $TemplateIds) {
        try {
            # Determine template type from ID
            $templateType = ""
            if ($templateId -like "*email*") {
                $templateType = "email"
            } elseif ($templateId -like "*push*") {
                $templateType = "push"
            } else {
                Write-Host "   ⚠️ Unknown template type for ID: $templateId" -ForegroundColor Yellow
                continue
            }

            # Extract GUID part from template ID (format: "Email_GUID" or "Push_GUID")
            $guidPart = $templateId
            if ($templateId -like "*_*") {
                $parts = $templateId -split "_", 2
                if ($parts.Length -eq 2) {
                    $guidPart = $parts[1]
                    Write-Host "   🔍 Extracted GUID part from template ID: $templateId -> $guidPart" -ForegroundColor Gray
                }
            }

            # Step 5: GET notification template with localized messages using the GUID part
            $uri = "https://graph.microsoft.com/beta/deviceManagement/notificationMessageTemplates('$guidPart')"
            Write-Host "🔍 Getting notification template ($templateType)..." -ForegroundColor Cyan
            Write-Host "   Template GUID: $guidPart" -ForegroundColor Gray
            Write-Host "   Endpoint: $uri" -ForegroundColor Gray
            Write-Host ""
            
            $templateResponse = Invoke-MgGraphRequest -Method GET -Uri $uri
            
            # Add template type for display purposes
            $templateResponse | Add-Member -NotePropertyName "templateType" -NotePropertyValue $templateType
            $templateResponse | Add-Member -NotePropertyName "originalTemplateId" -NotePropertyValue $templateId
            
            $templates += $templateResponse
        }
        catch {
            Write-Host "   ⚠️ Warning: Error fetching template $templateId : $_" -ForegroundColor Yellow
            continue
        }
    }
    
    return $templates
}

# Function to display enrollment notification configuration details
function Show-EnrollmentNotificationDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Configuration
    )
    
    Write-Host "📋 Android Enrollment Notification Configuration Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Configuration.id) {
        Write-Host "   • ID: $($Configuration.id)" -ForegroundColor Green
    }
    
    if ($Configuration.displayName) {
        Write-Host "   • Display Name: $($Configuration.displayName)" -ForegroundColor Green
    }
    
    if ($Configuration.description) {
        Write-Host "   • Description: $($Configuration.description)" -ForegroundColor Green
    }
    
    if ($Configuration.platformType) {
        Write-Host "   • Platform Type: $($Configuration.platformType)" -ForegroundColor Green
    }
    
    if ($Configuration.defaultLocale) {
        Write-Host "   • Default Locale: $($Configuration.defaultLocale)" -ForegroundColor Green
    }
    
    if ($Configuration.brandingOptions) {
        Write-Host "   • Branding Options: $($Configuration.brandingOptions)" -ForegroundColor Green
    }
    
    if ($Configuration.notificationTemplates -and $Configuration.notificationTemplates.Count -gt 0) {
        Write-Host "   • Notification Templates:" -ForegroundColor Green
        foreach ($template in $Configuration.notificationTemplates) {
            # Transform from API format back to user-friendly format
            $templateType = "unknown"
            if ($template -like "*email*") {
                $templateType = "email"
            } elseif ($template -like "*push*") {
                $templateType = "push"
            }
            Write-Host "     - $templateType (ID: $template)" -ForegroundColor Yellow
        }
    }
    
    if ($Configuration.createdDateTime) {
        Write-Host "   • Created: $($Configuration.createdDateTime)" -ForegroundColor Green
    }
    
    if ($Configuration.lastModifiedDateTime) {
        Write-Host "   • Last Modified: $($Configuration.lastModifiedDateTime)" -ForegroundColor Green
    }
    
    if ($Configuration.version) {
        Write-Host "   • Version: $($Configuration.version)" -ForegroundColor Green
    }
    
    if ($Configuration.priority) {
        Write-Host "   • Priority: $($Configuration.priority)" -ForegroundColor Green
    }
    
    if ($Configuration.deviceEnrollmentConfigurationType) {
        Write-Host "   • Configuration Type: $($Configuration.deviceEnrollmentConfigurationType)" -ForegroundColor Green
    }
    
    if ($Configuration.roleScopeTagIds -and $Configuration.roleScopeTagIds.Count -gt 0) {
        Write-Host "   • Role Scope Tag IDs: $($Configuration.roleScopeTagIds -join ', ')" -ForegroundColor Green
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Function to display notification message templates details
function Show-NotificationTemplatesDetails {
    param (
        [Parameter(Mandatory=$true)]
        [array]$Templates
    )
    
    Write-Host "📋 Notification Message Templates Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Templates.Count -gt 0) {
        Write-Host "   Found $($Templates.Count) template(s)" -ForegroundColor Green
        Write-Host ""
        
        for ($i = 0; $i -lt $Templates.Count; $i++) {
            $template = $Templates[$i]
            Write-Host "   • Template $($i + 1) ($($template.templateType)):" -ForegroundColor Green
            
            if ($template.id) {
                Write-Host "     - Template ID: $($template.id)" -ForegroundColor Yellow
            }
            
            if ($template.originalTemplateId) {
                Write-Host "     - Original Template ID: $($template.originalTemplateId)" -ForegroundColor Yellow
            }
            
            if ($template.displayName) {
                Write-Host "     - Display Name: $($template.displayName)" -ForegroundColor Yellow
            }
            
            if ($template.brandingOptions) {
                Write-Host "     - Branding Options: $($template.brandingOptions)" -ForegroundColor Yellow
            }
            
            if ($template.defaultLocale) {
                Write-Host "     - Default Locale: $($template.defaultLocale)" -ForegroundColor Yellow
            }
            
            # Display localized notification messages
            if ($template.localizedNotificationMessages -and $template.localizedNotificationMessages.Count -gt 0) {
                Write-Host "     - Localized Messages:" -ForegroundColor Yellow
                foreach ($message in $template.localizedNotificationMessages) {
                    Write-Host "       · Locale: $($message.locale)" -ForegroundColor Magenta
                    if ($message.subject) {
                        Write-Host "         Subject: $($message.subject)" -ForegroundColor Magenta
                    }
                    if ($message.messageTemplate) {
                        Write-Host "         Message: $($message.messageTemplate)" -ForegroundColor Magenta
                    }
                    if ($message.isDefault -ne $null) {
                        Write-Host "         Is Default: $($message.isDefault)" -ForegroundColor Magenta
                    }
                    Write-Host ""
                }
            } else {
                Write-Host "     - No localized messages found" -ForegroundColor Yellow
            }
            
            Write-Host ""
        }
    } else {
        Write-Host "   No templates found for this configuration" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Function to display enrollment notification assignments details
function Show-EnrollmentNotificationAssignmentsDetails {
    param (
        [Parameter(Mandatory=$true)]
        $Assignments
    )
    
    Write-Host "📋 Android Enrollment Notification Assignments Details:" -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    if ($Assignments.value -and $Assignments.value.Count -gt 0) {
        Write-Host "   Found $($Assignments.value.Count) assignment(s)" -ForegroundColor Green
        Write-Host ""
        
        for ($i = 0; $i -lt $Assignments.value.Count; $i++) {
            $assignment = $Assignments.value[$i]
            Write-Host "   • Assignment $($i + 1):" -ForegroundColor Green
            
            if ($assignment.id) {
                Write-Host "     - ID: $($assignment.id)" -ForegroundColor Yellow
            }
            
            if ($assignment.target) {
                Write-Host "     - Target:" -ForegroundColor Yellow
                $target = $assignment.target
                
                if ($target.'@odata.type') {
                    $targetType = $target.'@odata.type' -replace '#microsoft.graph.', ''
                    Write-Host "       · Type: $targetType" -ForegroundColor Yellow
                    
                    switch ($targetType) {
                        "groupAssignmentTarget" {
                            if ($target.groupId) {
                                Write-Host "       · Group ID: $($target.groupId)" -ForegroundColor Yellow
                            }
                        }
                        "exclusionGroupAssignmentTarget" {
                            if ($target.groupId) {
                                Write-Host "       · Exclusion Group ID: $($target.groupId)" -ForegroundColor Yellow
                            }
                        }
                        "allDevicesAssignmentTarget" {
                            Write-Host "       · Targets: All Devices" -ForegroundColor Yellow
                        }
                        "allLicensedUsersAssignmentTarget" {
                            Write-Host "       · Targets: All Licensed Users" -ForegroundColor Yellow
                        }
                    }
                }
                
                if ($target.deviceAndAppManagementAssignmentFilterId) {
                    Write-Host "       · Filter ID: $($target.deviceAndAppManagementAssignmentFilterId)" -ForegroundColor Yellow
                }
                
                if ($target.deviceAndAppManagementAssignmentFilterType) {
                    Write-Host "       · Filter Type: $($target.deviceAndAppManagementAssignmentFilterType)" -ForegroundColor Yellow
                }
            }
            
            Write-Host ""
        }
    } else {
        Write-Host "   No assignments found for this enrollment notification configuration" -ForegroundColor Yellow
    }
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    Write-Host ""
}

# Script Setup
try {
    Write-Host "🔐 Connecting to Microsoft Graph..." -ForegroundColor Cyan
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
    
    $secureClientSecret = ConvertTo-SecureString -String $ClientSecret -AsPlainText -Force
    $clientSecretCredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $ClientId, $secureClientSecret
    
    Connect-MgGraph -ClientSecretCredential $clientSecretCredential -TenantId $TenantId
    
    Write-Host "✅ Connected to Microsoft Graph successfully" -ForegroundColor Green
    Write-Host ""
    
    # Step 1: Get the enrollment notification configuration
    $configuration = Get-EnrollmentNotificationConfiguration -ConfigId $ConfigurationId
    
    # Step 2: Get the enrollment notification assignments
    $assignments = Get-EnrollmentNotificationAssignments -ConfigId $ConfigurationId
    
    # Step 3 & 4: Process notification templates and get their details
    $templates = @()
    if ($configuration -and $configuration.notificationTemplates -and $configuration.notificationTemplates.Count -gt 0) {
        Write-Host "🔍 Found notification templates, retrieving detailed information..." -ForegroundColor Cyan
        Write-Host ""
        
        $templates = Get-NotificationMessageTemplates -TemplateIds $configuration.notificationTemplates
    }
    
    # Display configuration details
    if ($configuration) {
        Show-EnrollmentNotificationDetails -Configuration $configuration
    } else {
        Write-Host "📊 No configuration found with ID: $ConfigurationId" -ForegroundColor Yellow
    }
    
    # Display template details
    if ($templates.Count -gt 0) {
        Show-NotificationTemplatesDetails -Templates $templates
    } else {
        Write-Host "📊 No notification templates found for this configuration" -ForegroundColor Yellow
        Write-Host ""
    }
    
    # Display assignments
    if ($assignments) {
        Show-EnrollmentNotificationAssignmentsDetails -Assignments $assignments
    } else {
        Write-Host "📊 No assignments found for this enrollment notification configuration" -ForegroundColor Yellow
        Write-Host ""
    }
    
    Write-Host "🎉 Operation completed successfully!" -ForegroundColor Green
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