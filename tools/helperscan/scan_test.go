package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func scanSample(t *testing.T) map[string]Func {
	t.Helper()
	funcs, err := scanDir(fixturePackage(t, "sample"))
	if err != nil {
		t.Fatalf("scanDir: %v", err)
	}
	byName := map[string]Func{}
	for _, f := range funcs {
		byName[displayName(f)] = f
	}
	return byName
}

func TestScanDirReportsOnlyNonStandardFunctions(t *testing.T) {
	got := scanSample(t)
	want := []string{
		"constructWidgetRequest", "settings", "SampleResource.refreshWidget", "lockWidget", "mapWidgetsToState", "waitForWidgetDeletion", "findWidgetByName",
		"createWidgetResponseFromDiscriminatorValue", "widgetResponse.GetAdditionalData", "widgetResponse.GetFieldDeserializers", "widgetResponse.Serialize",
		"widgetNameValidator.Description", "widgetNameValidator.MarkdownDescription", "widgetNameValidator.ValidateString",
		"MapRemoteStateToTerraform",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d functions, want %d: %v", len(got), len(want), got)
	}
	for _, name := range want {
		if _, ok := got[name]; !ok {
			t.Errorf("missing %s", name)
		}
	}
}

func TestScanDirClassifiesKinds(t *testing.T) {
	got := scanSample(t)
	cases := map[string]string{
		"constructWidgetRequest":                     kindHelper,
		"SampleResource.refreshWidget":               kindHelper,
		"createWidgetResponseFromDiscriminatorValue": kindKiota,
		"widgetResponse.Serialize":                   kindKiota,
		"widgetResponse.GetFieldDeserializers":       kindKiota,
		"widgetResponse.GetAdditionalData":           kindKiota,
		"widgetNameValidator.Description":            kindFramework,
		"widgetNameValidator.ValidateString":         kindFramework,
		"MapRemoteStateToTerraform":                  kindStateMapper,
	}
	for name, kind := range cases {
		if got[name].Kind != kind {
			t.Errorf("%s: kind %q, want %q", name, got[name].Kind, kind)
		}
	}
}

func TestScanDirCountsCRUDCallers(t *testing.T) {
	got := scanSample(t)
	cases := map[string][]string{
		// Three distinct CRUD methods across two resources, one of them inside a closure.
		"constructWidgetRequest":       {"OtherResource.Read", "SampleResource.Create", "SampleResource.Update"},
		"SampleResource.refreshWidget": {"SampleResource.Create", "SampleResource.Update"},
		"waitForWidgetDeletion":        {"SampleResource.Delete", "SampleResource.Update"},
		"findWidgetByName":             {"SampleResource.ImportState", "SampleResource.Read"},
		"lockWidget":                   {"SampleResource.Create"},
		// Called only from constructResource and MapRemoteStateToTerraform, which are not CRUD methods.
		"settings":          nil,
		"mapWidgetsToState": nil,
	}
	for name, want := range cases {
		if !reflect.DeepEqual(got[name].CRUDCallers, want) {
			t.Errorf("%s: CRUD callers %v, want %v", name, got[name].CRUDCallers, want)
		}
	}
}

func TestScanDirCountsReferences(t *testing.T) {
	got := scanSample(t)
	cases := map[string]int{
		"constructWidgetRequest": 3,
		// The local variable named settings in Read is not a reference to the function.
		"settings": 1,
		// The recursive call is not counted; the test file is not scanned.
		"findWidgetByName": 2,
		"lockWidget":       1,
		// Passing the factory as a value counts as a reference.
		"createWidgetResponseFromDiscriminatorValue": 1,
		"widgetNameValidator.Description":            1,
	}
	for name, want := range cases {
		if got[name].Refs != want {
			t.Errorf("%s: refs %d, want %d", name, got[name].Refs, want)
		}
	}
}

func TestScanDirMeasuresDeclarationSpan(t *testing.T) {
	got := scanSample(t)
	cases := map[string]int{"constructWidgetRequest": 11, "waitForWidgetDeletion": 12, "findWidgetByName": 9, "lockWidget": 4, "SampleResource.refreshWidget": 3}
	for name, want := range cases {
		if got[name].Lines != want {
			t.Errorf("%s: lines %d, want %d", name, got[name].Lines, want)
		}
	}
}

func TestScanDirAppliesRule(t *testing.T) {
	got := scanSample(t)
	cases := map[string]string{
		"constructWidgetRequest":       "pass",
		"waitForWidgetDeletion":        "pass",
		"findWidgetByName":             "fails: lines",
		"SampleResource.refreshWidget": "fails: lines,name",
		"lockWidget":                   "fails: callers,lines,name",
		"settings":                     "fails: callers,lines,name",
		"widgetResponse.Serialize":     "excluded",
		// The state mapper has no CRUD callers and 3 lines, but the caller and line rules don't apply to it.
		"MapRemoteStateToTerraform": "exempt",
	}
	for name, want := range cases {
		if rule := got[name].Rule(); rule != want {
			t.Errorf("%s: rule %q, want %q", name, rule, want)
		}
	}
	if got["mapWidgetsToState"].NameCheck != "map… belongs in state.go" {
		t.Errorf("mapWidgetsToState: name check %q", got["mapWidgetsToState"].NameCheck)
	}
}

func TestScanDirErrors(t *testing.T) {
	if _, err := scanDir(fixturePackage(t, "broken")); err == nil {
		t.Error("expected a parse error for the broken fixture")
	}
	if _, err := scanDir(filepath.Join("testdata", "missing")); err == nil {
		t.Error("expected a read error for a missing directory")
	}
}

func TestScanDirEmptyPackage(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("no Go here\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	funcs, err := scanDir(dir)
	if err != nil || len(funcs) != 0 {
		t.Fatalf("got %v, %v; want no functions and no error", funcs, err)
	}
}

func TestReceiverName(t *testing.T) {
	dir := t.TempDir()
	src := "package generic\n\ntype box[T any] struct{ v T }\n\ntype pair[K comparable, V any] struct{}\n\n" +
		"func (b *box[T]) unwrapBoxValue() T { return b.v }\n\nfunc (p pair[K, V]) describePairKeys() {}\n\nfunc (box[T]) Create() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "generic.go"), []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	funcs, err := scanDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, f := range funcs {
		got[f.Name] = f.Receiver
	}
	want := map[string]string{"unwrapBoxValue": "box", "describePairKeys": "pair"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("receivers %v, want %v", got, want)
	}
}

func TestScanDirChecksStateMapperNames(t *testing.T) {
	funcs, err := scanDir(fixturePackage(t, "legacy"))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]Func{}
	for _, f := range funcs {
		got[f.Name] = f
	}
	rename := "non-standard state mapper: rename to MapRemoteStateToTerraform"
	cases := map[string]struct{ kind, name, rule string }{
		"MapRemoteStateToTerraform":         {kindStateMapper, "MapRemoteStateToTerraform belongs in state.go", "fails: name"},
		"MapRemoteResourceStateToTerraform": {kindHelper, rename, "fails: callers,lines,name"},
		"mapRemoteState":                    {kindHelper, rename, "fails: callers,lines,name"},
	}
	for name, want := range cases {
		f := got[name]
		if f.Kind != want.kind || f.NameCheck != want.name || f.Rule() != want.rule {
			t.Errorf("%s: kind %q, name %q, rule %q; want %q, %q, %q", name, f.Kind, f.NameCheck, f.Rule(), want.kind, want.name, want.rule)
		}
	}
}

func scanByName(t *testing.T, dir string) map[string]Func {
	t.Helper()
	funcs, err := scanDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]Func{}
	for _, f := range funcs {
		got[f.Name] = f
	}
	return got
}

func TestScanDirChecksDataSourceStateMapperNames(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "internal", "services", "datasources", "example", "graph_beta")
	got := scanByName(t, fixturePackageAt(t, "datasource", parent))
	rename := "non-standard state mapper: rename to MapRemoteStateToDataSource"
	cases := map[string]struct{ kind, name, rule string }{
		// One reference and 3 lines, but exempt from the caller and line rules.
		"MapRemoteStateToDataSource":        {kindStateMapper, nameOK, "exempt"},
		"mapRemoteStateToDataSource":        {kindHelper, rename, "fails: callers,lines,name"},
		"MapRemoteResourceStateToTerraform": {kindHelper, rename, "fails: callers,lines,name"},
		"MapRemoteStateToTerraform":         {kindHelper, rename, "fails: callers,lines,name"},
	}
	for name, want := range cases {
		f := got[name]
		if f.Kind != want.kind || f.NameCheck != want.name || f.Rule() != want.rule {
			t.Errorf("%s: kind %q, name %q, rule %q; want %q, %q, %q", name, f.Kind, f.NameCheck, f.Rule(), want.kind, want.name, want.rule)
		}
	}
}

func TestScanDirKeepsResourceRuleOutsideDataSources(t *testing.T) {
	parent := filepath.Join(t.TempDir(), "internal", "services", "resources", "example", "graph_beta")
	got := scanByName(t, fixturePackageAt(t, "datasource", parent))
	if f := got["MapRemoteStateToTerraform"]; f.Kind != kindStateMapper || f.NameCheck != "MapRemoteStateToTerraform belongs in state.go" {
		t.Errorf("MapRemoteStateToTerraform in a resource package: kind %q, name %q", f.Kind, f.NameCheck)
	}
	if f := got["MapRemoteStateToDataSource"]; f.Kind != kindHelper || f.NameCheck == "non-standard state mapper: rename to MapRemoteStateToTerraform" {
		t.Errorf("MapRemoteStateToDataSource in a resource package: kind %q, name %q", f.Kind, f.NameCheck)
	}
}

func TestStateMapperFor(t *testing.T) {
	cases := map[string]string{
		"internal/services/datasources/identity_and_access/graph_beta/x": dataSourceStateMapperName,
		"internal/services/datasources":                                  dataSourceStateMapperName,
		"internal/services/resources/identity_and_access/graph_beta/x":   stateMapperName,
		"internal/services/datasources_extra/x":                          stateMapperName,
	}
	for dir, want := range cases {
		if got := stateMapperFor(filepath.FromSlash(dir)); got != want {
			t.Errorf("stateMapperFor(%q) = %q, want %q", dir, got, want)
		}
	}
}
