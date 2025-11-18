package license

import (
	"testing"
)

func TestGetRequiredLicensesForFeature(t *testing.T) {
	tests := []struct {
		name         string
		featureName  string
		wantMinCount int
	}{
		{
			name:         "NetworkFilteringPolicy requires licenses",
			featureName:  "NetworkFilteringPolicy",
			wantMinCount: 1,
		},
		{
			name:         "ConditionalAccessPolicy requires licenses",
			featureName:  "ConditionalAccessPolicy",
			wantMinCount: 1,
		},
		{
			name:         "IntuneEndpointPrivilege requires licenses",
			featureName:  "IntuneEndpointPrivilege",
			wantMinCount: 1,
		},
		{
			name:         "IntuneCloudPKI requires licenses",
			featureName:  "IntuneCloudPKI",
			wantMinCount: 1,
		},
		{
			name:         "Windows365CloudPC requires licenses",
			featureName:  "Windows365CloudPC",
			wantMinCount: 1,
		},
		{
			name:         "Unknown feature returns empty",
			featureName:  "NonExistentFeature",
			wantMinCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			licenses := GetRequiredLicensesForFeature(tt.featureName)
			if len(licenses) < tt.wantMinCount {
				t.Errorf("GetRequiredLicensesForFeature() returned %d licenses, want at least %d",
					len(licenses), tt.wantMinCount)
			}
		})
	}
}

func TestFormatRequiredLicensesMessage(t *testing.T) {
	tests := []struct {
		name        string
		featureName string
		wantEmpty   bool
	}{
		{
			name:        "NetworkFilteringPolicy returns message",
			featureName: "NetworkFilteringPolicy",
			wantEmpty:   false,
		},
		{
			name:        "IntuneEndpointPrivilege returns message",
			featureName: "IntuneEndpointPrivilege",
			wantEmpty:   false,
		},
		{
			name:        "Windows365CloudPC returns message",
			featureName: "Windows365CloudPC",
			wantEmpty:   false,
		},
		{
			name:        "Unknown feature returns empty",
			featureName: "NonExistentFeature",
			wantEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := FormatRequiredLicensesMessage(tt.featureName)
			isEmpty := message == ""
			if isEmpty != tt.wantEmpty {
				t.Errorf("FormatRequiredLicensesMessage() empty=%v, want %v", isEmpty, tt.wantEmpty)
			}
		})
	}
}

func TestFormatLicensesForError(t *testing.T) {
	tests := []struct {
		name     string
		licenses []LicenseInfo
		wantText string
	}{
		{
			name:     "Empty licenses",
			licenses: []LicenseInfo{},
			wantText: "No licenses found",
		},
		{
			name: "Entra Suite SKU with service plan",
			licenses: []LicenseInfo{
				{
					SkuPartNumber:   "Microsoft_Entra_Suite",
					ServicePlanName: "",
					ConsumedUnits:   10,
					PrepaidUnits:    25,
				},
				{
					SkuPartNumber:   "Microsoft_Entra_Suite",
					ServicePlanName: "Entra_Premium_Internet_Access",
				},
			},
			wantText: "Licenses found",
		},
		{
			name: "Intune Suite SKU with service plans",
			licenses: []LicenseInfo{
				{
					SkuPartNumber:   "Microsoft_Intune_Suite",
					ServicePlanName: "",
					ConsumedUnits:   5,
					PrepaidUnits:    10,
				},
				{
					SkuPartNumber:   "Microsoft_Intune_Suite",
					ServicePlanName: "INTUNE_P2",
				},
				{
					SkuPartNumber:   "Microsoft_Intune_Suite",
					ServicePlanName: "CLOUD_PKI",
				},
			},
			wantText: "Licenses found",
		},
		{
			name: "Windows 365 Business SKU with service plans",
			licenses: []LicenseInfo{
				{
					SkuPartNumber:   "CPC_B_2C_8RAM_128GB",
					ServicePlanName: "",
					ConsumedUnits:   1,
					PrepaidUnits:    0,
				},
				{
					SkuPartNumber:   "CPC_B_2C_8RAM_128GB",
					ServicePlanName: "CPC_SS_2",
				},
				{
					SkuPartNumber:   "CPC_B_2C_8RAM_128GB",
					ServicePlanName: "Windows_10_ESU_Commercial",
				},
			},
			wantText: "Licenses found",
		},
		{
			name: "Multiple different SKUs",
			licenses: []LicenseInfo{
				{
					SkuPartNumber:   "ENTERPRISEPACK",
					ServicePlanName: "",
					ConsumedUnits:   50,
					PrepaidUnits:    100,
				},
				{
					SkuPartNumber:   "Microsoft_Intune_Suite",
					ServicePlanName: "",
					ConsumedUnits:   5,
					PrepaidUnits:    10,
				},
				{
					SkuPartNumber:   "Microsoft_Entra_Suite",
					ServicePlanName: "",
					ConsumedUnits:   10,
					PrepaidUnits:    25,
				},
			},
			wantText: "Licenses found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatLicensesForError(tt.licenses)
			if result == "" {
				t.Errorf("FormatLicensesForError() returned empty string")
			}
		})
	}
}
