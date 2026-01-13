package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

// mockState tracks the state of resources for consistent responses
var mockState struct {
	sync.Mutex
	configurations     map[string]map[string]any
	definitionValues   map[string]map[string]any
	presentationValues map[string][]any  // key: configID/defValueID, value: array of presentation values
	operationPhase     map[string]string // tracks if resource is in create/read/update/delete phase
}

func init() {
	mockState.configurations = make(map[string]map[string]any)
	mockState.definitionValues = make(map[string]map[string]any)
	mockState.presentationValues = make(map[string][]any)
	mockState.operationPhase = make(map[string]string)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("group_policy_definition", &GroupPolicyDefinitionMock{})
}

// GroupPolicyDefinitionMock provides mock responses for Group Policy Definition operations
type GroupPolicyDefinitionMock struct{}

var _ mocks.MockRegistrar = (*GroupPolicyDefinitionMock)(nil)

// RegisterMocks registers HTTP mock responses for Group Policy Definition operations
func (m *GroupPolicyDefinitionMock) RegisterMocks() {
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]any)
	mockState.definitionValues = make(map[string]map[string]any)
	mockState.presentationValues = make(map[string][]any)
	mockState.operationPhase = make(map[string]string)
	mockState.Unlock()

	// Register GET for group policy configurations (parent resource)
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configID := urlParts[len(urlParts)-1]

			mockState.Lock()
			config, exists := mockState.configurations[configID]
			mockState.Unlock()

			if !exists {
				config = map[string]any{
					"id":          configID,
					"displayName": "Mock Configuration",
				}
			}

			return httpmock.NewJsonResponse(200, config)
		})

	// Register GET for group policy definitions catalog search
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions\?`,
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query()
			filter := query.Get("$filter")

			scenario := determineCatalogScenario(filter)
			if scenario == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine catalog scenario"}}`), nil
			}

			// Determine operation phase from filter content to load correct folder
			phase := determinePhaseFromContext(filter)

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", phase, "get_catalog_"+scenario+".json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load catalog JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse catalog JSON: `+err.Error()+`"}}`), nil
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for single group policy definition by ID
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions/[^/?]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defID := urlParts[len(urlParts)-1]

			scenario := determinePresentationsScenario(defID)
			if scenario == "" {
				return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Definition not found"}}`), nil
			}

			// Load the catalog JSON for this definition
			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "get_catalog_"+scenario+".json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load definition JSON: `+err.Error()+`"}}`), nil
			}

			var catalogResponse map[string]any
			if err := json.Unmarshal([]byte(content), &catalogResponse); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse definition JSON: `+err.Error()+`"}}`), nil
			}

			// Extract the first definition from the catalog response
			if values, hasValues := catalogResponse["value"].([]any); hasValues && len(values) > 0 {
				if definition, ok := values[0].(map[string]any); ok {
					return httpmock.NewJsonResponse(200, definition)
				}
			}

			return httpmock.NewStringResponse(404, `{"error":{"code":"NotFound","message":"Definition not found"}}`), nil
		})

	// Register GET for presentation templates
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyDefinitions/[^/]+/presentations`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defID := urlParts[len(urlParts)-2]

			scenario := determinePresentationsScenario(defID)
			if scenario == "" {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Could not determine presentations scenario"}}`), nil
			}

			// Determine phase - if definition value doesn't exist yet, it's create
			phase := "validate_create"
			mockState.Lock()
			if len(mockState.definitionValues) > 0 {
				phase = "validate_read"
			}
			mockState.Unlock()

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", phase, "get_presentations_"+scenario+".json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load presentations JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse presentations JSON: `+err.Error()+`"}}`), nil
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register POST for updateDefinitionValues (CREATE/UPDATE/DELETE operation)
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/updateDefinitionValues$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configID := urlParts[len(urlParts)-2]

			var requestBody map[string]any
			err := json.NewDecoder(req.Body).Decode(&requestBody)
			if err != nil {
				return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
			}

			// Determine if this is CREATE, UPDATE, or DELETE
			added, hasAdded := requestBody["added"].([]any)
			updated, hasUpdated := requestBody["updated"].([]any)
			deleted, hasDeleted := requestBody["deleted"].([]any)

			if hasDeleted && len(deleted) > 0 {
				// DELETE operation
				for _, deletedItem := range deleted {
					if id, ok := deletedItem.(string); ok {
						stateKey := configID + "/" + id
						mockState.Lock()
						delete(mockState.definitionValues, stateKey)
						delete(mockState.presentationValues, stateKey)
						delete(mockState.operationPhase, stateKey)
						mockState.Unlock()
					}
				}
			}

			// Handle UPDATED items
			if hasUpdated && len(updated) > 0 {
				for _, updatedItem := range updated {
					if defVal, ok := updatedItem.(map[string]any); ok {
						// Get the existing ID from the update request
						var defValueID string
						if id, hasID := defVal["id"].(string); hasID {
							defValueID = id
						} else {
							// If no ID provided, this is an error
							continue
						}

						storedDefVal := make(map[string]any)
						for k, v := range defVal {
							storedDefVal[k] = v
						}
						storedDefVal["@odata.type"] = "#microsoft.graph.groupPolicyDefinitionValue"
						storedDefVal["id"] = defValueID
						storedDefVal["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

						// Extract definition ID from @odata.bind or definition object
						var definitionID string
						if definition, hasDef := defVal["definition"].(map[string]any); hasDef {
							if id, hasID := definition["id"].(string); hasID {
								definitionID = id
							}
						}

						// Check for definition@odata.bind
						if definitionID == "" {
							if bindURL, hasBinding := defVal["definition@odata.bind"].(string); hasBinding {
								parts := strings.Split(bindURL, "'")
								if len(parts) >= 2 {
									definitionID = parts[len(parts)-2]
								}
							}

							if definitionID == "" {
								if additionalData, hasAdditional := defVal["additionalData"].(map[string]any); hasAdditional {
									if bindURL, hasBinding := additionalData["definition@odata.bind"].(string); hasBinding {
										parts := strings.Split(bindURL, "'")
										if len(parts) >= 2 {
											definitionID = parts[len(parts)-2]
										}
									}
								}
							}
						}

						// Store definition reference with ID
						if definitionID != "" {
							storedDefVal["definition"] = map[string]any{
								"@odata.id": fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionID),
								"id":        definitionID,
							}
						}

						stateKey := configID + "/" + defValueID

						// Extract and store presentation values
						var presValues []any
						if presValuesArray, hasPres := defVal["presentationValues"].([]any); hasPres {
							presValues = make([]any, len(presValuesArray))
							for i, pv := range presValuesArray {
								if pvMap, ok := pv.(map[string]any); ok {
									// Create a copy to avoid mutating the original
									enrichedPV := make(map[string]any)
									for k, v := range pvMap {
										enrichedPV[k] = v
									}

									// Generate ID if not present
									if _, hasID := enrichedPV["id"]; !hasID {
										enrichedPV["id"] = uuid.New().String()
									}

									// Extract presentationId from @odata.bind if present
									if bindURL, hasBinding := enrichedPV["presentation@odata.bind"].(string); hasBinding {
										parts := strings.Split(bindURL, "'")
										if len(parts) >= 2 {
											presID := parts[len(parts)-2]
											enrichedPV["presentationId"] = presID
											enrichedPV["presentation"] = map[string]any{"id": presID}
										}
									}

									presValues[i] = enrichedPV
								}
							}
						}

						mockState.Lock()
						mockState.definitionValues[stateKey] = storedDefVal
						if len(presValues) > 0 {
							mockState.presentationValues[stateKey] = presValues
						}
						mockState.operationPhase[stateKey] = constants.TfOperationUpdate
						mockState.Unlock()
					}
				}
			}

			// Handle ADDED items
			if hasAdded && len(added) > 0 {
				// CREATE operation
				for _, addedItem := range added {
					if defVal, ok := addedItem.(map[string]any); ok {
						defValueID := uuid.New().String()

						storedDefVal := make(map[string]any)
						for k, v := range defVal {
							storedDefVal[k] = v
						}
						storedDefVal["@odata.type"] = "#microsoft.graph.groupPolicyDefinitionValue"
						storedDefVal["id"] = defValueID
						storedDefVal["createdDateTime"] = "2024-01-01T00:00:00Z"
						storedDefVal["lastModifiedDateTime"] = "2024-01-01T00:00:00Z"

						// Extract definition ID from @odata.bind or definition object
						var definitionID string
						if definition, hasDef := defVal["definition"].(map[string]any); hasDef {
							if id, hasID := definition["id"].(string); hasID {
								definitionID = id
							}
						}

						// Check for definition@odata.bind in the request (could be at root level or in additionalData)
						if definitionID == "" {
							// Try at root level first
							if bindURL, hasBinding := defVal["definition@odata.bind"].(string); hasBinding {
								// Extract ID from URL like "https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('def-boolean-minimal')"
								parts := strings.Split(bindURL, "'")
								if len(parts) >= 2 {
									definitionID = parts[len(parts)-2]
								}
							}

							// Try in additionalData if not found at root
							if definitionID == "" {
								if additionalData, hasAdditional := defVal["additionalData"].(map[string]any); hasAdditional {
									if bindURL, hasBinding := additionalData["definition@odata.bind"].(string); hasBinding {
										parts := strings.Split(bindURL, "'")
										if len(parts) >= 2 {
											definitionID = parts[len(parts)-2]
										}
									}
								}
							}
						}

						// Store definition reference with ID
						if definitionID != "" {
							storedDefVal["definition"] = map[string]any{
								"@odata.id": fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionID),
								"id":        definitionID,
							}
						}

						stateKey := configID + "/" + defValueID

						// Extract and store presentation values
						var presValues []any
						if presValuesArray, hasPres := defVal["presentationValues"].([]any); hasPres {
							presValues = make([]any, len(presValuesArray))
							for i, pv := range presValuesArray {
								if pvMap, ok := pv.(map[string]any); ok {
									// Create a copy to avoid mutating the original
									enrichedPV := make(map[string]any)
									for k, v := range pvMap {
										enrichedPV[k] = v
									}

									// Generate ID if not present
									if _, hasID := enrichedPV["id"]; !hasID {
										enrichedPV["id"] = uuid.New().String()
									}

									// Extract presentationId from @odata.bind if present
									if bindURL, hasBinding := enrichedPV["presentation@odata.bind"].(string); hasBinding {
										parts := strings.Split(bindURL, "'")
										if len(parts) >= 2 {
											presID := parts[len(parts)-2]
											enrichedPV["presentationId"] = presID
											enrichedPV["presentation"] = map[string]any{"id": presID}
										}
									}

									presValues[i] = enrichedPV
								}
							}
						}

						mockState.Lock()
						mockState.definitionValues[stateKey] = storedDefVal
						if len(presValues) > 0 {
							mockState.presentationValues[stateKey] = presValues
						}
						mockState.operationPhase[stateKey] = constants.TfOperationCreate
						mockState.Unlock()
					}
				}
			}

			return httpmock.NewStringResponse(204, ""), nil
		})

	// Register GET for listing definition values in a configuration
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			configID := urlParts[len(urlParts)-2]

			// Check if $expand=definition is requested
			query := req.URL.Query()
			expand := query.Get("$expand")
			needsExpandedDefinition := strings.Contains(expand, "definition")

			mockState.Lock()
			definitionValues := []map[string]any{}
			for stateKey, defValue := range mockState.definitionValues {
				if strings.HasPrefix(stateKey, configID+"/") {
					// Make a copy to avoid mutation
					defValueCopy := make(map[string]any)
					for k, v := range defValue {
						defValueCopy[k] = v
					}

					// If definition expansion is requested, ensure full definition object is present
					if needsExpandedDefinition {
						var defID string
						if definition, hasDef := defValue["definition"].(map[string]any); hasDef {
							if id, hasID := definition["id"].(string); hasID {
								defID = id
							}
						}

						// If we have a definition ID, load the full definition from catalog
						if defID != "" {
							scenario := determinePresentationsScenario(defID)
							if scenario != "" {
								content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_create", "get_catalog_"+scenario+".json"))
								if err == nil {
									var catalogResponse map[string]any
									if json.Unmarshal([]byte(content), &catalogResponse) == nil {
										if values, hasValues := catalogResponse["value"].([]any); hasValues && len(values) > 0 {
											if fullDef, ok := values[0].(map[string]any); ok {
												defValueCopy["definition"] = fullDef
											}
										}
									}
								}
							}
						}
					}

					definitionValues = append(definitionValues, defValueCopy)
				}
			}
			mockState.Unlock()

			response := map[string]any{
				"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations('" + configID + "')/definitionValues",
				"value":          definitionValues,
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for individual definition value
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defValueID := urlParts[len(urlParts)-1]
			configID := urlParts[len(urlParts)-3]
			stateKey := configID + "/" + defValueID

			mockState.Lock()
			defValueData, exists := mockState.definitionValues[stateKey]
			phase := mockState.operationPhase[stateKey]
			if phase == "" {
				phase = constants.TfOperationRead
			}
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_definition_value_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			scenario := determineScenarioFromDefValue(defValueData)
			phaseFolder := "validate_" + phase
			if phase == constants.TfOperationCreate {
				phaseFolder = "validate_create"
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", phaseFolder, "get_definition_value_"+scenario+".json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load definition value JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse definition value JSON: `+err.Error()+`"}}`), nil
			}

			response["id"] = defValueID
			if def, hasDef := defValueData["definition"]; hasDef {
				response["definition"] = def
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register GET for presentation values collection
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues/[^/]+/presentationValues`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defValueID := urlParts[len(urlParts)-2]
			configID := urlParts[len(urlParts)-4]
			stateKey := configID + "/" + defValueID

			mockState.Lock()
			defValueData, exists := mockState.definitionValues[stateKey]
			storedPresValues, hasPresValues := mockState.presentationValues[stateKey]
			phase := mockState.operationPhase[stateKey]
			if phase == "" {
				phase = constants.TfOperationRead
			}
			mockState.Unlock()

			if !exists {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations/definitionValues/presentationValues",
					"value":          []any{},
				})
			}

			// If we have stored presentation values, return those
			if hasPresValues && len(storedPresValues) > 0 {
				return httpmock.NewJsonResponse(200, map[string]any{
					"@odata.context": "https://graph.microsoft.com/beta/$metadata#deviceManagement/groupPolicyConfigurations/definitionValues/presentationValues",
					"value":          storedPresValues,
				})
			}

			// Fallback to JSON files for initial test data
			scenario := determineScenarioFromDefValue(defValueData)
			phaseFolder := "validate_" + phase
			if phase == constants.TfOperationCreate {
				phaseFolder = "validate_create"
			}

			content, err := helpers.ParseJSONFile(filepath.Join("../tests", "responses", phaseFolder, "get_presentation_values_"+scenario+".json"))
			if err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to load presentation values JSON: `+err.Error()+`"}}`), nil
			}

			var response map[string]any
			if err := json.Unmarshal([]byte(content), &response); err != nil {
				return httpmock.NewStringResponse(500, `{"error":{"code":"InternalServerError","message":"Failed to parse presentation values JSON: `+err.Error()+`"}}`), nil
			}

			return httpmock.NewJsonResponse(200, response)
		})

	// Register DELETE for definition value
	httpmock.RegisterResponder(constants.TfTfOperationDelete, `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			urlParts := strings.Split(req.URL.Path, "/")
			defValueID := urlParts[len(urlParts)-1]
			configID := urlParts[len(urlParts)-3]
			stateKey := configID + "/" + defValueID

			mockState.Lock()
			_, exists := mockState.definitionValues[stateKey]
			if exists {
				delete(mockState.definitionValues, stateKey)
				delete(mockState.presentationValues, stateKey)
				delete(mockState.operationPhase, stateKey)
			}
			mockState.Unlock()

			if !exists {
				content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_definition_value_not_found.json"))
				var errorResponse map[string]any
				json.Unmarshal([]byte(content), &errorResponse)
				return httpmock.NewJsonResponse(404, errorResponse)
			}

			return httpmock.NewStringResponse(204, ""), nil
		})
}

// Helper functions

func determineCatalogScenario(filter string) string {
	if strings.Contains(filter, "Test Policy Boolean Minimal") {
		return "boolean_minimal"
	}
	if strings.Contains(filter, "Test Policy Boolean Maximal") {
		return "boolean_maximal"
	}
	if strings.Contains(filter, "Test Policy TextBox") {
		return "textbox"
	}
	if strings.Contains(filter, "Test Policy Decimal") {
		return "decimal"
	}
	if strings.Contains(filter, "Test Policy MultiText") {
		return "multitext"
	}
	if strings.Contains(filter, "Test Policy Dropdown") {
		return "dropdown"
	}
	if strings.Contains(filter, "Test Policy With Read-Only Text") {
		return "readonly"
	}
	return ""
}

func determinePresentationsScenario(defID string) string {
	switch defID {
	case "def-boolean-minimal":
		return "boolean_minimal"
	case "def-boolean-maximal":
		return "boolean_maximal"
	case "def-textbox":
		return "textbox"
	case "def-decimal":
		return "decimal"
	case "def-multitext":
		return "multitext"
	case "def-dropdown":
		return "dropdown"
	case "def-readonly":
		return "readonly"
	default:
		return ""
	}
}

func determineScenarioFromDefValue(defValueData map[string]any) string {
	if definition, hasDef := defValueData["definition"].(map[string]any); hasDef {
		if id, hasID := definition["id"].(string); hasID {
			return determinePresentationsScenario(id)
		}
	}
	return "boolean_minimal"
}

func determinePhaseFromContext(filter string) string {
	// During first call (validation), use validate_create
	// This is a simplification - in reality we'd track state better
	return "validate_create"
}

// CleanupMockState clears the mock state for clean test runs
func (m *GroupPolicyDefinitionMock) CleanupMockState() {
	mockState.Lock()
	defer mockState.Unlock()

	for id := range mockState.configurations {
		delete(mockState.configurations, id)
	}
	for id := range mockState.definitionValues {
		delete(mockState.definitionValues, id)
	}
	for id := range mockState.presentationValues {
		delete(mockState.presentationValues, id)
	}
	for id := range mockState.operationPhase {
		delete(mockState.operationPhase, id)
	}
}

// RegisterErrorMocks registers HTTP mock responses for error scenarios
func (m *GroupPolicyDefinitionMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.configurations = make(map[string]map[string]any)
	mockState.definitionValues = make(map[string]map[string]any)
	mockState.operationPhase = make(map[string]string)
	mockState.Unlock()

	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/updateDefinitionValues$`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, `{"error":{"code":"BadRequest","message":"Invalid request body"}}`), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/groupPolicyConfigurations/[^/]+/definitionValues/[^/]+$`,
		func(req *http.Request) (*http.Response, error) {
			content, _ := helpers.ParseJSONFile(filepath.Join("../tests", "responses", "validate_delete", "get_definition_value_not_found.json"))
			var errorResponse map[string]any
			json.Unmarshal([]byte(content), &errorResponse)
			return httpmock.NewJsonResponse(404, errorResponse)
		})
}
