package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func serializeBody(t *testing.T, body serialization.Parsable) map[string]any {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize body: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized body: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("decode serialized body %s: %v", string(content), err)
	}
	return result
}

func ruleBasedModel() *LifecycleWorkflowResourceModel {
	return &LifecycleWorkflowResourceModel{
		Category:            types.StringValue("joiner"),
		DisplayName:         types.StringValue("joiner"),
		Description:         types.StringValue(""),
		IsEnabled:           types.BoolValue(false),
		IsSchedulingEnabled: types.BoolValue(false),
		ExecutionConditions: &WorkflowExecutionConditions{
			OnDemandExecutionOnly: types.BoolValue(false),
			Scope:                 &WorkflowScope{Rule: types.StringValue("(department eq 'acc-test-none')"), GroupIds: types.SetNull(types.StringType)},
			Trigger:               &WorkflowTrigger{TimeBasedAttribute: &TimeBasedAttributeTrigger{Attribute: types.StringValue("employeeHireDate"), OffsetInDays: types.Int64Value(-3)}},
		},
		Tasks: []WorkflowTask{
			{
				TaskDefinitionId: types.StringValue("6fc52c9d-398b-4305-9763-15f42c1676fc"),
				DisplayName:      types.StringValue("Enable user account"),
				Description:      types.StringValue(""),
				IsEnabled:        types.BoolValue(true),
				ContinueOnError:  types.BoolValue(false),
			},
			{
				TaskDefinitionId: types.StringValue("70b29d51-b59a-4773-9280-8841dfd3f2ea"),
				DisplayName:      types.StringValue("Send welcome email"),
				Description:      types.StringValue(""),
				IsEnabled:        types.BoolValue(true),
				ContinueOnError:  types.BoolValue(true),
				Arguments:        []WorkflowTaskArgument{{Name: types.StringValue("locale"), Value: types.StringValue("en-us")}},
			},
		},
	}
}

// Graph rejects both a null arguments collection and a task with no arguments property
// ("Task arguments missing or invalid"), so unset arguments must be sent as an empty array.
func TestUnitConstructTasksSendsUnsetArgumentsAsEmptyArray(t *testing.T) {
	workflow, err := constructResource(context.Background(), ruleBasedModel())
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	body := serializeBody(t, workflow)

	tasks, ok := body["tasks"].([]any)
	if !ok || len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %v", body["tasks"])
	}
	first := tasks[0].(map[string]any)
	firstArguments, present := first["arguments"]
	if !present || firstArguments == nil {
		t.Fatalf("expected an empty arguments array for a task without arguments, got %v", first)
	}
	if list, ok := firstArguments.([]any); !ok || len(list) != 0 {
		t.Errorf("expected an empty arguments array, got %v", firstArguments)
	}
	second := tasks[1].(map[string]any)
	arguments, ok := second["arguments"].([]any)
	if !ok || len(arguments) != 1 {
		t.Errorf("expected one argument on the second task, got %v", second["arguments"])
	}
	if first["executionSequence"] != float64(1) || second["executionSequence"] != float64(2) {
		t.Errorf("expected execution sequence to follow list order, got %v and %v", first["executionSequence"], second["executionSequence"])
	}
	if body["category"] != "joiner" {
		t.Errorf("expected category to be sent for createNewVersion, got %v", body["category"])
	}
}

// An explicitly empty argument list is also sent as [] rather than null.
func TestUnitConstructTasksSendsExplicitEmptyArguments(t *testing.T) {
	data := ruleBasedModel()
	data.Tasks[0].Arguments = []WorkflowTaskArgument{}

	workflow, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	body := serializeBody(t, workflow)
	first := body["tasks"].([]any)[0].(map[string]any)
	arguments, present := first["arguments"]
	if !present || arguments == nil {
		t.Fatalf("expected an empty arguments array, got %v", first)
	}
	if list, ok := arguments.([]any); !ok || len(list) != 0 {
		t.Errorf("expected an empty arguments array, got %v", arguments)
	}
}

func TestUnitConstructExecutionConditionsOnDemandOmitsScopeAndTrigger(t *testing.T) {
	data := ruleBasedModel()
	data.ExecutionConditions = &WorkflowExecutionConditions{OnDemandExecutionOnly: types.BoolValue(true)}

	workflow, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	conditions := serializeBody(t, workflow)["executionConditions"].(map[string]any)
	if conditions["@odata.type"] != "#microsoft.graph.identityGovernance.onDemandExecutionOnly" {
		t.Errorf("expected onDemandExecutionOnly, got %v", conditions)
	}
	for _, field := range []string{"scope", "trigger"} {
		if _, present := conditions[field]; present {
			t.Errorf("expected %s to be omitted, got %v", field, conditions)
		}
	}
}

func TestUnitConstructExecutionConditionsGroupScopeAndMembershipTrigger(t *testing.T) {
	data := ruleBasedModel()
	data.Category = types.StringValue("mover")
	data.ExecutionConditions = &WorkflowExecutionConditions{
		OnDemandExecutionOnly: types.BoolValue(false),
		Scope:                 &WorkflowScope{Rule: types.StringNull(), GroupIds: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("11111111-1111-1111-1111-111111111111")})},
		Trigger:               &WorkflowTrigger{MembershipChange: &MembershipChangeTrigger{ChangeType: types.StringValue("remove")}},
	}

	workflow, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	conditions := serializeBody(t, workflow)["executionConditions"].(map[string]any)
	scope := conditions["scope"].(map[string]any)
	if scope["@odata.type"] != "#microsoft.graph.identityGovernance.groupBasedSubjectSet" {
		t.Errorf("expected groupBasedSubjectSet, got %v", scope)
	}
	if _, present := scope["rule"]; present {
		t.Errorf("expected rule to be omitted for a group scope, got %v", scope)
	}
	groups := scope["groups"].([]any)
	if len(groups) != 1 || groups[0].(map[string]any)["id"] != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("unexpected groups %v", groups)
	}
	trigger := conditions["trigger"].(map[string]any)
	if trigger["@odata.type"] != "#microsoft.graph.identityGovernance.membershipChangeTrigger" || trigger["changeType"] != "remove" {
		t.Errorf("unexpected trigger %v", trigger)
	}
}

func TestUnitConstructAttributeChangeTriggerSendsSingleAttribute(t *testing.T) {
	data := ruleBasedModel()
	data.ExecutionConditions.Trigger = &WorkflowTrigger{AttributeChange: &AttributeChangeTrigger{TriggerAttribute: types.StringValue("jobTitle")}}

	workflow, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	trigger := serializeBody(t, workflow)["executionConditions"].(map[string]any)["trigger"].(map[string]any)
	attributes := trigger["triggerAttributes"].([]any)
	if len(attributes) != 1 || attributes[0].(map[string]any)["name"] != "jobTitle" {
		t.Errorf("expected a single jobTitle trigger attribute, got %v", trigger)
	}
}

// Graph rejects any property other than the four basic ones on PATCH.
func TestUnitConstructBasicUpdateOnlySendsPatchableProperties(t *testing.T) {
	workflow, err := constructBasicUpdate(context.Background(), ruleBasedModel())
	if err != nil {
		t.Fatalf("constructBasicUpdate returned error: %v", err)
	}
	body := serializeBody(t, workflow)
	allowed := map[string]bool{"displayName": true, "description": true, "isEnabled": true, "isSchedulingEnabled": true}
	for key := range body {
		if !allowed[key] && !strings.HasPrefix(key, "@odata") {
			t.Errorf("unexpected PATCH property %s in %v", key, body)
		}
	}
	for key := range allowed {
		if _, present := body[key]; !present {
			t.Errorf("expected PATCH property %s in %v", key, body)
		}
	}
}

func TestUnitConstructRejectsInvalidExecutionConditions(t *testing.T) {
	cases := map[string]*WorkflowExecutionConditions{
		"on demand with scope": {
			OnDemandExecutionOnly: types.BoolValue(true),
			Scope:                 &WorkflowScope{Rule: types.StringValue("(department eq 'x')"), GroupIds: types.SetNull(types.StringType)},
		},
		"missing trigger": {
			OnDemandExecutionOnly: types.BoolValue(false),
			Scope:                 &WorkflowScope{Rule: types.StringValue("(department eq 'x')"), GroupIds: types.SetNull(types.StringType)},
		},
		"membership change with rule scope": {
			OnDemandExecutionOnly: types.BoolValue(false),
			Scope:                 &WorkflowScope{Rule: types.StringValue("(department eq 'x')"), GroupIds: types.SetNull(types.StringType)},
			Trigger:               &WorkflowTrigger{MembershipChange: &MembershipChangeTrigger{ChangeType: types.StringValue("add")}},
		},
	}

	for name, conditions := range cases {
		t.Run(name, func(t *testing.T) {
			data := ruleBasedModel()
			data.ExecutionConditions = conditions
			if _, err := constructResource(context.Background(), data); err == nil {
				t.Fatalf("expected an error for %s", name)
			}
		})
	}
}
