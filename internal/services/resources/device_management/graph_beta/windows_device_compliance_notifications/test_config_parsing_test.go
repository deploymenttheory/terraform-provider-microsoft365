package graphBetaWindowsDeviceComplianceNotifications

import (
	"testing"
)

func TestAcceptanceConfigParsing(t *testing.T) {
	err := TestConfigParsing()
	if err != nil {
		t.Fatalf("Config parsing failed: %v", err)
	}
	t.Log("All acceptance test configurations parsed successfully")
}