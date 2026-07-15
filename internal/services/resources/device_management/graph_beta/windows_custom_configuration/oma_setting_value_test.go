package graphBetaWindowsCustomConfiguration

import "testing"

func TestParseOmaSettingValueRejectsNonFiniteFloats(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"NaN", "+Inf", "-Inf"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			_, err := parseOmaSettingValue("#microsoft.graph.omaSettingFloatingPoint", value)
			if err == nil {
				t.Fatalf("parseOmaSettingValue() accepted non-finite value %q", value)
			}
		})
	}
}
