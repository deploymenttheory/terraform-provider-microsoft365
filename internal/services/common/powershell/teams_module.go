package powershell

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

// EnsurePowerShellAndTeamsModule checks if pwsh and the MicrosoftTeams module are installed, and installs the module if missing.
func EnsurePowerShellAndTeamsModule() error {
	// Check if pwsh is available
	if _, err := exec.LookPath("pwsh"); err != nil {
		return errors.New("powerShell (pwsh) is not installed or not in PATH. Please install PowerShell 7.2+")
	}

	// Check if MicrosoftTeams module is available
	checkCmd := exec.Command("pwsh", "-Command", "Get-Module -ListAvailable -Name MicrosoftTeams")
	var out, stderr bytes.Buffer
	checkCmd.Stdout = &out
	checkCmd.Stderr = &stderr
	if err := checkCmd.Run(); err != nil || !bytes.Contains(out.Bytes(), []byte("MicrosoftTeams")) {
		// Try to install the module
		installCmd := exec.Command("pwsh", "-Command", "Install-Module -Name MicrosoftTeams -Force -AllowClobber -Scope CurrentUser")
		installCmd.Stdout = &out
		installCmd.Stderr = &stderr
		if err := installCmd.Run(); err != nil {
			return errors.New("Failed to install MicrosoftTeams PowerShell module: " + stderr.String())
		}
	}
	return nil
}

// ConnectMicrosoftTeams authenticates PowerShell to Microsoft Teams using Connect-MicrosoftTeams.
// Supports only App ID + Certificate (thumbprint or file) and App ID + Client Secret authentication.
// Pass empty strings for unused parameters. All other authentication types will return an error.
func ConnectMicrosoftTeams(tenantId, applicationId, clientSecret, certificateThumbprint, certificatePath string) error {
	var psCmd string
	switch {
	case applicationId != "" && clientSecret != "" && tenantId != "":
		// Client secret authentication
		psCmd = fmt.Sprintf(`Connect-MicrosoftTeams -ApplicationId '%s' -ClientSecret '%s' -TenantId '%s'`, applicationId, clientSecret, tenantId)
	case applicationId != "" && certificateThumbprint != "" && tenantId != "":
		// Certificate thumbprint authentication
		psCmd = fmt.Sprintf(`Connect-MicrosoftTeams -ApplicationId '%s' -CertificateThumbprint '%s' -TenantId '%s'`, applicationId, certificateThumbprint, tenantId)
	case applicationId != "" && certificatePath != "" && tenantId != "":
		// Certificate file authentication
		psCmd = fmt.Sprintf(`$cert = New-Object System.Security.Cryptography.X509Certificates.X509Certificate2('%s'); Connect-MicrosoftTeams -ApplicationId '%s' -Certificate $cert -TenantId '%s'`, certificatePath, applicationId, tenantId)
	default:
		return errors.New("app ID + Client Secret or App ID + Certificate authentication is supported. Please provide the required parameters")
	}

	cmd := exec.Command("pwsh", "-Command", psCmd)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return errors.New("Failed to authenticate to Microsoft Teams: " + stderr.String())
	}
	return nil
}
