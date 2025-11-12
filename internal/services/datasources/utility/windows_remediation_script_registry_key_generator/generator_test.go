package utilityWindowsRemediationScriptRegistryKeyGenerator

import (
	"strings"
	"testing"
)

func TestGenerateScripts_CurrentUser(t *testing.T) {
	t.Parallel()

	detection, remediation, err := GenerateScripts(
		"current_user",
		"Software\\Policies\\Microsoft\\WindowsStore\\",
		"RequirePrivateStoreOnly",
		"REG_DWORD",
		"1",
	)

	if err != nil {
		t.Fatalf("GenerateScripts returned error: %v", err)
	}

	detectionStr := detection.ValueString()
	remediationStr := remediation.ValueString()

	// Verify detection script contains key elements
	if !strings.Contains(detectionStr, "Get-CimInstance win32_computersystem") {
		t.Error("Detection script missing user SID retrieval")
	}
	if !strings.Contains(detectionStr, "Software\\Policies\\Microsoft\\WindowsStore\\") {
		t.Error("Detection script missing registry path")
	}
	if !strings.Contains(detectionStr, "'RequirePrivateStoreOnly'") {
		t.Error("Detection script missing value name")
	}
	if !strings.Contains(detectionStr, "Exit 1") {
		t.Error("Detection script missing exit codes")
	}

	// Verify remediation script contains key elements
	if !strings.Contains(remediationStr, "New-Item -Path $regkey") {
		t.Error("Remediation script missing key creation")
	}
	if !strings.Contains(remediationStr, "New-ItemProperty") {
		t.Error("Remediation script missing property creation")
	}
	if !strings.Contains(remediationStr, "Set-ItemProperty") {
		t.Error("Remediation script missing property update")
	}
	if !strings.Contains(remediationStr, "PropertyType DWord") {
		t.Error("Remediation script missing correct property type")
	}
}

func TestGenerateScripts_AllUsers(t *testing.T) {
	t.Parallel()

	detection, remediation, err := GenerateScripts(
		"all_users",
		"Software\\MyApp\\Settings\\",
		"EnableFeature",
		"REG_SZ",
		"Enabled",
	)

	if err != nil {
		t.Fatalf("GenerateScripts returned error: %v", err)
	}

	detectionStr := detection.ValueString()
	remediationStr := remediation.ValueString()

	// Verify detection script loops through all users
	if !strings.Contains(detectionStr, "Get-ChildItem \"HKU:\"") {
		t.Error("Detection script missing user enumeration")
	}
	if !strings.Contains(detectionStr, "foreach ($User in $Users)") {
		t.Error("Detection script missing user loop")
	}
	if !strings.Contains(detectionStr, "S-1-5-18") {
		t.Error("Detection script missing system account exclusions")
	}

	// Verify remediation script contains key elements
	if !strings.Contains(remediationStr, "foreach ($User in $Users)") {
		t.Error("Remediation script missing user loop")
	}
	if !strings.Contains(remediationStr, "PropertyType String") {
		t.Error("Remediation script missing correct property type")
	}
	if !strings.Contains(remediationStr, "'Enabled'") {
		t.Error("Remediation script missing quoted string value")
	}
}

func TestGetPowerShellPropertyType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    string
		expected string
		wantErr  bool
	}{
		"REG_SZ":        {input: "REG_SZ", expected: "String", wantErr: false},
		"REG_DWORD":     {input: "REG_DWORD", expected: "DWord", wantErr: false},
		"REG_QWORD":     {input: "REG_QWORD", expected: "QWord", wantErr: false},
		"REG_MULTI_SZ":  {input: "REG_MULTI_SZ", expected: "MultiString", wantErr: false},
		"REG_EXPAND_SZ": {input: "REG_EXPAND_SZ", expected: "ExpandString", wantErr: false},
		"REG_BINARY":    {input: "REG_BINARY", expected: "Binary", wantErr: false},
		"Invalid":       {input: "INVALID", expected: "", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := getPowerShellPropertyType(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestFormatValueData(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		valueType string
		valueData string
		wantErr   bool
		contains  string
	}{
		"DWORD valid":        {valueType: "REG_DWORD", valueData: "123", wantErr: false, contains: "123"},
		"DWORD invalid":      {valueType: "REG_DWORD", valueData: "invalid", wantErr: true},
		"String":             {valueType: "REG_SZ", valueData: "Test", wantErr: false, contains: "'Test'"},
		"String with quote":  {valueType: "REG_SZ", valueData: "Test's Value", wantErr: false, contains: "''"},
		"Multi-string":       {valueType: "REG_MULTI_SZ", valueData: "Line1\nLine2", wantErr: false, contains: "@("},
		"Binary valid":       {valueType: "REG_BINARY", valueData: "01AF", wantErr: false, contains: "0x01"},
		"Binary invalid odd": {valueType: "REG_BINARY", valueData: "01A", wantErr: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result, err := formatValueData(tc.valueType, tc.valueData)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !strings.Contains(result, tc.contains) {
				t.Errorf("Expected result to contain %s, got: %s", tc.contains, result)
			}
		})
	}
}

func TestFormatPropertyName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    string
		expected string
	}{
		"Normal name": {input: "MyValue", expected: "'MyValue'"},
		"Default":     {input: "(Default)", expected: "''"},
		"Empty":       {input: "", expected: "''"},
		"With quote":  {input: "Test's", expected: "'Test''s'"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := formatPropertyName(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
