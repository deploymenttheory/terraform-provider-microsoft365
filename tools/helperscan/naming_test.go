package main

import (
	"reflect"
	"testing"
)

func TestCheckName(t *testing.T) {
	cases := []struct {
		name, file, want string
	}{
		{"constructTargetAppsRequest", "construct.go", nameOK},
		{"findTargetedAppByPackageID", "state.go", nameOK},
		{"mapAssignmentsToState", "state.go", nameOK},
		{"waitForPolicyDeletion", "crud.go", nameOK},
		{"lockPolicyTargetApps", "crud.go", nameOK},
		{"writableApp", "construct.go", "banned generic name"},
		{"matches", "state.go", "banned generic name"},
		{"readbackOptions", "crud.go", "banned generic name"},
		{"packageID", "state.go", "verb not in convention"},
		{"decodeCreateResponse", "response.go", "verb not in convention"},
		{"lockPolicy", "crud.go", "missing entity or qualifier"},
		{"constructRequest", "construct.go", "missing entity or qualifier"},
		{"waitForDeletion", "crud.go", "missing entity or qualifier"},
		{"mapResponse", "state.go", "map… must name an entity and end ToState"},
		{"mapToState", "state.go", "map… must name an entity and end ToState"},
		{"findByID", "state.go", "find… must name an entity then By…"},
		{"findAppBy", "state.go", "find… must name an entity then By…"},
		{"findApp", "state.go", "find… must name an entity then By…"},
		{"constructTargetAppsRequest", "crud.go", "construct… belongs in construct.go"},
		{"waitForPolicyDeletion", "state.go", "waitFor… belongs in crud.go"},
		{"", "crud.go", "verb not in convention"},
	}
	for _, c := range cases {
		if got := checkName(c.name, c.file); got != c.want {
			t.Errorf("checkName(%q, %q) = %q, want %q", c.name, c.file, got, c.want)
		}
	}
}

func TestSplitWords(t *testing.T) {
	cases := map[string][]string{
		"findTargetedAppByPackageID": {"find", "Targeted", "App", "By", "Package", "ID"},
		"constructURLRequest":        {"construct", "URL", "Request"},
		"mapV2PolicyToState":         {"map", "V2", "Policy", "To", "State"},
		"matches":                    {"matches"},
		"":                           nil,
	}
	for name, want := range cases {
		if got := splitWords(name); !reflect.DeepEqual(got, want) {
			t.Errorf("splitWords(%q) = %q, want %q", name, got, want)
		}
	}
}
