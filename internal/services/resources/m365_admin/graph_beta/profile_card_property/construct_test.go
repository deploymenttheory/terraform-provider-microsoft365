package graphBetaM365AdminProfileCardProperty

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func serializeProfileCardProperty(t *testing.T, property graphmodels.ProfileCardPropertyable) string {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", property); err != nil {
		t.Fatalf("serialize profile card property: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized profile card property: %v", err)
	}
	return string(content)
}

func testModel() *ProfileCardPropertyResourceModel {
	return &ProfileCardPropertyResourceModel{
		ID:                    types.StringUnknown(),
		DirectoryPropertyName: types.StringValue("CustomAttribute1"),
		IsVisible:             types.BoolValue(true),
	}
}

func testAnnotations() []ProfileCardAnnotationModel {
	return []ProfileCardAnnotationModel{
		{
			DisplayName: types.StringValue("Cost Center"),
			Localizations: []DisplayNameLocalizationModel{
				{LanguageTag: types.StringValue("de-DE"), DisplayName: types.StringValue("Kostenstelle")},
			},
		},
	}
}

// directoryPropertyName is the key, so it is only sent on create.
func TestConstructRequestBodyDirectoryPropertyNameOnlyOnCreate(t *testing.T) {
	model := testModel()

	if body := serializeProfileCardProperty(t, constructRequestBody(model, true)); !strings.Contains(body, `"directoryPropertyName":"CustomAttribute1"`) {
		t.Errorf("expected directoryPropertyName on create, got %s", body)
	}
	if body := serializeProfileCardProperty(t, constructRequestBody(model, false)); strings.Contains(body, `"directoryPropertyName"`) {
		t.Errorf("expected no directoryPropertyName on update, got %s", body)
	}
}

// On create, unset annotations are omitted rather than sent as null or an empty collection.
func TestConstructRequestBodyOmitsUnsetAnnotationsOnCreate(t *testing.T) {
	body := serializeProfileCardProperty(t, constructRequestBody(testModel(), true))
	if strings.Contains(body, `"annotations"`) {
		t.Errorf("expected annotations to be omitted on create, got %s", body)
	}
	if !strings.Contains(body, `"isVisible":true`) {
		t.Errorf("expected isVisible on create, got %s", body)
	}
}

// On update, PATCH keeps omitted properties, so removed annotations are cleared with an empty collection.
func TestConstructRequestBodyClearsRemovedAnnotationsOnUpdate(t *testing.T) {
	body := serializeProfileCardProperty(t, constructRequestBody(testModel(), false))
	if !strings.Contains(body, `"annotations":[]`) {
		t.Errorf("expected an empty annotations collection on update, got %s", body)
	}
}

func TestConstructAnnotations(t *testing.T) {
	annotations := constructAnnotations(testAnnotations())
	if len(annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(annotations))
	}
	if got := *annotations[0].GetDisplayName(); got != "Cost Center" {
		t.Errorf("expected display name 'Cost Center', got %q", got)
	}
	localizations := annotations[0].GetLocalizations()
	if len(localizations) != 1 || *localizations[0].GetLanguageTag() != "de-DE" || *localizations[0].GetDisplayName() != "Kostenstelle" {
		t.Errorf("unexpected localizations: %+v", localizations)
	}
}

func TestConstructResourceIncludesAnnotationsOnCreate(t *testing.T) {
	model := testModel()
	model.Annotations = testAnnotations()

	requestBody, err := constructResource(context.Background(), model)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	body := serializeProfileCardProperty(t, requestBody)
	if !strings.Contains(body, `"languageTag":"de-DE"`) || !strings.Contains(body, `"directoryPropertyName":"CustomAttribute1"`) {
		t.Errorf("expected annotations and directoryPropertyName on create, got %s", body)
	}
}
