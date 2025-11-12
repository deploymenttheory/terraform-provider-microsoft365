package utilityWindowsRemediationScriptRegistryKeyGenerator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GenerateScripts creates PowerShell detection and remediation scripts based on the provided parameters.
func GenerateScripts(context, registryKeyPath, valueName, valueType, valueData string) (types.String, types.String, error) {
	// Validate and prepare PowerShell property type
	psPropertyType, err := getPowerShellPropertyType(valueType)
	if err != nil {
		return types.StringNull(), types.StringNull(), err
	}

	// Validate value data based on type
	psValueData, err := formatValueData(valueType, valueData)
	if err != nil {
		return types.StringNull(), types.StringNull(), err
	}

	var detectionScript, remediationScript string

	switch context {
	case "current_user":
		detectionScript = generateDetectionScriptCurrentUser(registryKeyPath, valueName, psValueData)
		remediationScript = generateRemediationScriptCurrentUser(registryKeyPath, valueName, psPropertyType, psValueData)
	case "all_users":
		detectionScript = generateDetectionScriptAllUsers(registryKeyPath, valueName, psValueData)
		remediationScript = generateRemediationScriptAllUsers(registryKeyPath, valueName, psPropertyType, psValueData)
	default:
		return types.StringNull(), types.StringNull(), fmt.Errorf("invalid context: %s", context)
	}

	return types.StringValue(detectionScript), types.StringValue(remediationScript), nil
}

// getPowerShellPropertyType maps registry value types to PowerShell property types.
func getPowerShellPropertyType(valueType string) (string, error) {
	switch valueType {
	case "REG_SZ":
		return "String", nil
	case "REG_DWORD":
		return "DWord", nil
	case "REG_QWORD":
		return "QWord", nil
	case "REG_MULTI_SZ":
		return "MultiString", nil
	case "REG_EXPAND_SZ":
		return "ExpandString", nil
	case "REG_BINARY":
		return "Binary", nil
	default:
		return "", fmt.Errorf("unsupported value type: %s", valueType)
	}
}

// formatValueData validates and formats the value data for PowerShell based on type.
func formatValueData(valueType, valueData string) (string, error) {
	switch valueType {
	case "REG_DWORD", "REG_QWORD":
		// Validate it's a valid integer
		_, err := strconv.ParseInt(valueData, 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid %s value '%s': must be a valid integer", valueType, valueData)
		}
		return valueData, nil
	case "REG_SZ", "REG_EXPAND_SZ":
		// String values need to be quoted
		return fmt.Sprintf("'%s'", strings.ReplaceAll(valueData, "'", "''")), nil
	case "REG_MULTI_SZ":
		// Multi-string values need to be split and quoted
		lines := strings.Split(valueData, "\n")
		quotedLines := make([]string, len(lines))
		for i, line := range lines {
			quotedLines[i] = fmt.Sprintf("'%s'", strings.ReplaceAll(strings.TrimSpace(line), "'", "''"))
		}
		return "@(" + strings.Join(quotedLines, ", ") + ")", nil
	case "REG_BINARY":
		// Binary data should be hex string, convert to byte array
		cleaned := strings.ReplaceAll(strings.ReplaceAll(valueData, " ", ""), "0x", "")
		if len(cleaned)%2 != 0 {
			return "", fmt.Errorf("invalid binary data: hex string must have even length")
		}
		bytes := []string{}
		for i := 0; i < len(cleaned); i += 2 {
			bytes = append(bytes, "0x"+cleaned[i:i+2])
		}
		return "@(" + strings.Join(bytes, ", ") + ")", nil
	default:
		return "", fmt.Errorf("unsupported value type: %s", valueType)
	}
}

func generateDetectionScriptCurrentUser(regPath, name, value string) string {
	return fmt.Sprintf(`#Registry Detection Script: Current logged on User
#Get SID of current interactive user
$CurrentLoggedOnUser = (Get-CimInstance win32_computersystem).UserName
if (-not ([string]::IsNullOrEmpty($CurrentLoggedOnUser))) {
    $AdObj = New-Object System.Security.Principal.NTAccount($CurrentLoggedOnUser)
    $strSID = $AdObj.Translate([System.Security.Principal.SecurityIdentifier])
    $UserSid = $strSID.Value
} else {
    $UserSid = $null
}

New-PSDrive -PSProvider Registry -Name "HKU" -Root HKEY_USERS -ErrorAction SilentlyContinue | Out-Null

$regkey = "HKU:\$UserSid\%s"

If (!(Test-Path $regkey)) {
    Write-Output 'RegKey not available - remediate'
    Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
    Exit 1
}

$check = (Get-ItemProperty -Path $regkey -Name %s -ErrorAction SilentlyContinue).%s

if ($check -eq %s) {
    Write-Output 'Setting OK - no remediation required'
    Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
    Exit 0
} else {
    Write-Output 'Value not OK, no value or could not read - go and remediate'
    Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
    Exit 1
}`, regPath, formatPropertyName(name), formatPropertyName(name), value)
}

func generateRemediationScriptCurrentUser(regPath, name, propType, value string) string {
	return fmt.Sprintf(`#Registry Remediation Script: Current logged on User
#Get SID of current interactive user
$CurrentLoggedOnUser = (Get-CimInstance win32_computersystem).UserName
if (-not ([string]::IsNullOrEmpty($CurrentLoggedOnUser))) {
    $AdObj = New-Object System.Security.Principal.NTAccount($CurrentLoggedOnUser)
    $strSID = $AdObj.Translate([System.Security.Principal.SecurityIdentifier])
    $UserSid = $strSID.Value
} else {
    $UserSid = $null
}

New-PSDrive -PSProvider Registry -Name "HKU" -Root HKEY_USERS -ErrorAction SilentlyContinue | Out-Null

$regkey = "HKU:\$UserSid\%s"

If (!(Test-Path $regkey)) {
    New-Item -Path $regkey -ErrorAction Stop
}

if (!(Get-ItemProperty -Path $regkey -Name %s -ErrorAction SilentlyContinue)) {
    New-ItemProperty -Path $regkey -Name %s -Value %s -PropertyType %s -ErrorAction Stop
    Write-Output "Remediation complete"
    Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
    Exit 0
}

Set-ItemProperty -Path $regkey -Name %s -Value %s -ErrorAction Stop
Write-Output "Remediation complete"
Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
Exit 0`, regPath, formatPropertyName(name), formatPropertyName(name), value, propType, formatPropertyName(name), value)
}

func generateDetectionScriptAllUsers(regPath, name, value string) string {
	return fmt.Sprintf(`#Registry Detection Script: All Users
New-PSDrive -PSProvider Registry -Name "HKU" -Root HKEY_USERS -ErrorAction SilentlyContinue | Out-Null
$Users = (Get-ChildItem "HKU:").Name

foreach ($User in $Users) {
    If ($User -ne "HKEY_USERS\S-1-5-18" -and $User -ne "HKEY_USERS\S-1-5-19" -and $User -ne "HKEY_USERS\S-1-5-20" -and $User -ne "HKEY_USERS\.DEFAULT" -and $User -notlike "*_Classes") {
        $regkey = "HKU:\$User\%s"
        
        If (!(Test-Path $regkey)) {
            Write-Output 'RegKey not available - remediate'
            Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
            Exit 1
        }
        
        $check = (Get-ItemProperty -Path $regkey -Name %s -ErrorAction SilentlyContinue).%s
        
        if ($check -eq %s) {
            Write-Output "Setting OK for $User - no remediation required"
        } else {
            Write-Output "Value not OK for $User - go and remediate"
            Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
            Exit 1
        }
    }
}

Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
Exit 0`, regPath, formatPropertyName(name), formatPropertyName(name), value)
}

func generateRemediationScriptAllUsers(regPath, name, propType, value string) string {
	return fmt.Sprintf(`#Registry Remediation Script: All Users
New-PSDrive -PSProvider Registry -Name "HKU" -Root HKEY_USERS -ErrorAction SilentlyContinue | Out-Null
$Users = (Get-ChildItem "HKU:").Name

foreach ($User in $Users) {
    If ($User -ne "HKEY_USERS\S-1-5-18" -and $User -ne "HKEY_USERS\S-1-5-19" -and $User -ne "HKEY_USERS\S-1-5-20" -and $User -ne "HKEY_USERS\.DEFAULT" -and $User -notlike "*_Classes") {
        $regkey = "HKU:\$User\%s"
        
        If (!(Test-Path $regkey)) {
            New-Item -Path $regkey -ErrorAction Stop
        }
        
        if (!(Get-ItemProperty -Path $regkey -Name %s -ErrorAction SilentlyContinue)) {
            New-ItemProperty -Path $regkey -Name %s -Value %s -PropertyType %s -ErrorAction Stop
        } else {
            Set-ItemProperty -Path $regkey -Name %s -Value %s -ErrorAction Stop
        }
    }
}

Write-Output "Remediation complete"
Remove-PSDrive -Name "HKU" -ErrorAction SilentlyContinue | Out-Null
Exit 0`, regPath, formatPropertyName(name), formatPropertyName(name), value, propType, formatPropertyName(name), value)
}

// formatPropertyName handles special case for default value.
func formatPropertyName(name string) string {
	if name == "(Default)" || name == "" {
		return "''"
	}
	return fmt.Sprintf("'%s'", strings.ReplaceAll(name, "'", "''"))
}
