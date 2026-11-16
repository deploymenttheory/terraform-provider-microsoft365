package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"slices"
	"testing"
)

func TestParseResourceID(t *testing.T) {
	tests := []struct {
		id         string
		policyType string
		policyID   string
		wantErr    bool
	}{
		{id: "device/0000000a-0000-0000-c000-000000000000", policyType: "device", policyID: "0000000a-0000-0000-c000-000000000000"},
		{id: "app/0000000a-0000-0000-c000-000000000000", policyType: "app", policyID: "0000000a-0000-0000-c000-000000000000"},
		{id: "0000000a-0000-0000-c000-000000000000", wantErr: true},
		{id: "windows/0000000a-0000-0000-c000-000000000000", wantErr: true},
		{id: "device/not-a-guid", wantErr: true},
		{id: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			policyType, policyID, err := parseResourceID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %q", tt.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if policyType != tt.policyType || policyID != tt.policyID {
				t.Errorf("got %q %q, want %q %q", policyType, policyID, tt.policyType, tt.policyID)
			}
			if got := buildResourceID(policyType, policyID); got != tt.id {
				t.Errorf("buildResourceID round trip got %q, want %q", got, tt.id)
			}
		})
	}
}

func TestDiffGroupIDs(t *testing.T) {
	tests := []struct {
		name       string
		current    []string
		desired    []string
		wantAdd    []string
		wantRemove []string
	}{
		{name: "no change", current: []string{"a", "b"}, desired: []string{"b", "a"}, wantAdd: []string{}, wantRemove: []string{}},
		{name: "add only", current: []string{}, desired: []string{"b", "a"}, wantAdd: []string{"a", "b"}, wantRemove: []string{}},
		{name: "remove all", current: []string{"a", "b"}, desired: []string{}, wantAdd: []string{}, wantRemove: []string{"a", "b"}},
		{name: "swap one", current: []string{"a", "b"}, desired: []string{"b", "c"}, wantAdd: []string{"c"}, wantRemove: []string{"a"}},
		{name: "case insensitive", current: []string{"ABC"}, desired: []string{"abc"}, wantAdd: []string{}, wantRemove: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toAdd, toRemove := diffGroupIDs(tt.current, tt.desired)
			if !slices.Equal(toAdd, tt.wantAdd) {
				t.Errorf("toAdd got %v, want %v", toAdd, tt.wantAdd)
			}
			if !slices.Equal(toRemove, tt.wantRemove) {
				t.Errorf("toRemove got %v, want %v", toRemove, tt.wantRemove)
			}
		})
	}
}
