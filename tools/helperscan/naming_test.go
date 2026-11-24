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
		{"MapRemoteResourceStateToTerraform", "state.go", "non-standard state mapper: rename to MapRemoteStateToTerraform"},
		{"mapRemoteStateToTerraform", "state.go", "non-standard state mapper: rename to MapRemoteStateToTerraform"},
		{"mapRemoteState", "state.go", "non-standard state mapper: rename to MapRemoteStateToTerraform"},
		{"mapRemoteResourceStateToTerraform", "state.go", "non-standard state mapper: rename to MapRemoteStateToTerraform"},
		{"MapRemoteStateToDataSource", "state.go", "map… must name an entity and end ToState"},
		{"MapRemoteOwnersToTerraform", "state.go", "map… must name an entity and end ToState"},
		{"mapRemoteOwnersToTerraform", "state.go", "map… must name an entity and end ToState"},
	}
	for _, c := range cases {
		if got := checkName(c.name, c.file, stateMapperName); got != c.want {
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

func TestCheckStateMapperFile(t *testing.T) {
	if got := checkStateMapperFile("MapRemoteStateToTerraform", "state.go"); got != nameOK {
		t.Errorf("state.go: %q", got)
	}
	if got := checkStateMapperFile("MapRemoteStateToTerraform", "crud.go"); got != "MapRemoteStateToTerraform belongs in state.go" {
		t.Errorf("crud.go: %q", got)
	}
}

func TestCheckNameInDataSource(t *testing.T) {
	rename := "non-standard state mapper: rename to MapRemoteStateToDataSource"
	cases := map[string]string{
		"mapRemoteStateToDataSource":        rename,
		"MapRemoteResourceStateToTerraform": rename,
		"MapRemoteStateToTerraform":         rename,
		"mapRemoteState":                    rename,
		"mapRemoteOwnersToTerraform":        "map… must name an entity and end ToState",
	}
	for name, want := range cases {
		if got := checkName(name, "state.go", dataSourceStateMapperName); got != want {
			t.Errorf("checkName(%q) in a data source = %q, want %q", name, got, want)
		}
	}
	if got := checkStateMapperFile("MapRemoteStateToDataSource", "crud.go"); got != "MapRemoteStateToDataSource belongs in state.go" {
		t.Errorf("checkStateMapperFile = %q", got)
	}
}
