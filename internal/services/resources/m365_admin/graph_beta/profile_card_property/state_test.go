package graphBetaM365AdminProfileCardProperty

import (
	"context"
	"testing"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func stringPtr(value string) *string {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}

// The ID is the directoryPropertyName exactly as Graph returns it.
func TestMapRemoteResourceStateKeepsKeyCasing(t *testing.T) {
	remote := graphmodels.NewProfileCardProperty()
	remote.SetDirectoryPropertyName(stringPtr("customAttribute9"))
	remote.SetIsVisible(boolPtr(false))

	state := MapRemoteResourceStateToTerraform(context.Background(), ProfileCardPropertyResourceModel{}, remote)

	if state.ID.ValueString() != "customAttribute9" || state.DirectoryPropertyName.ValueString() != "customAttribute9" {
		t.Errorf("expected ID and directory_property_name 'customAttribute9', got %q and %q", state.ID.ValueString(), state.DirectoryPropertyName.ValueString())
	}
	if state.IsVisible.ValueBool() {
		t.Errorf("expected is_visible false")
	}
}

// Graph returns empty collections for unset annotations and localizations; both map to null to match omitted configuration.
func TestMapRemoteResourceStateEmptyCollectionsAreNull(t *testing.T) {
	remote := graphmodels.NewProfileCardProperty()
	remote.SetDirectoryPropertyName(stringPtr("CustomAttribute1"))
	remote.SetAnnotations([]graphmodels.ProfileCardAnnotationable{})

	if state := MapRemoteResourceStateToTerraform(context.Background(), ProfileCardPropertyResourceModel{}, remote); state.Annotations != nil {
		t.Errorf("expected null annotations, got %+v", state.Annotations)
	}

	annotation := graphmodels.NewProfileCardAnnotation()
	annotation.SetDisplayName(stringPtr("Employee Number"))
	annotation.SetLocalizations([]graphmodels.DisplayNameLocalizationable{})
	remote.SetAnnotations([]graphmodels.ProfileCardAnnotationable{annotation})

	state := MapRemoteResourceStateToTerraform(context.Background(), ProfileCardPropertyResourceModel{}, remote)
	if len(state.Annotations) != 1 || state.Annotations[0].DisplayName.ValueString() != "Employee Number" {
		t.Fatalf("expected one annotation, got %+v", state.Annotations)
	}
	if state.Annotations[0].Localizations != nil {
		t.Errorf("expected null localizations, got %+v", state.Annotations[0].Localizations)
	}
}

func TestMapRemoteResourceStateLocalizations(t *testing.T) {
	localization := graphmodels.NewDisplayNameLocalization()
	localization.SetLanguageTag(stringPtr("nb-NO"))
	localization.SetDisplayName(stringPtr("Kostnadssenter"))
	annotation := graphmodels.NewProfileCardAnnotation()
	annotation.SetDisplayName(stringPtr("Cost Center"))
	annotation.SetLocalizations([]graphmodels.DisplayNameLocalizationable{localization})
	remote := graphmodels.NewProfileCardProperty()
	remote.SetDirectoryPropertyName(stringPtr("CustomAttribute2"))
	remote.SetAnnotations([]graphmodels.ProfileCardAnnotationable{annotation})

	state := MapRemoteResourceStateToTerraform(context.Background(), ProfileCardPropertyResourceModel{}, remote)
	if len(state.Annotations) != 1 || len(state.Annotations[0].Localizations) != 1 {
		t.Fatalf("expected one annotation with one localization, got %+v", state.Annotations)
	}
	if got := state.Annotations[0].Localizations[0]; got.LanguageTag.ValueString() != "nb-NO" || got.DisplayName.ValueString() != "Kostnadssenter" {
		t.Errorf("unexpected localization %+v", got)
	}
}
