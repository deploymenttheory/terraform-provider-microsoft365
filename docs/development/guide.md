# Development Guide

This guide describes the recommended workflow and best practices for developing new resources for the Community Terraform Provider for Microsoft 365.

## Typical Development Workflow

1. **Understand the Microsoft Graph API for the resource type**
   - Read the official Microsoft Graph documentation for [Microsoft Graph v1.0](https://learn.microsoft.com/en-us/graph/api/overview?view=graph-rest-1.0) / [Microsoft Graph beta](https://learn.microsoft.com/en-us/graph/api/overview?view=graph-rest-beta) for the resource you want to implement.
   - Use tools like [Graph X-Ray](https://graphxray.merill.net/) (set to Go language) to observe the API calls made by the Microsoft 365 portal.
   - Determine whether to use the Graph v1.0/beta endpoint based on API availability and stability.
   - If the resource is only available in the beta api, you will need to use the beta sdk.
   - Import the Microsoft Graph SDK for the resource you want to implement.
   - Example:

     ```go
     import (
         "github.com/microsoftgraph/msgraph-sdk-go"
         "github.com/microsoftgraph/msgraph-beta-sdk-go"
     )
     ```


2. **Design the Data Model**
   - Define a Go struct representing the Terraform resource model. Include all required and optional fields, and timeouts.
   - Example:

     ```go
     type ResourceTemplateResourceModel struct {
         ID       types.String   `tfsdk:"id"`
         // Add other fields here
         Timeouts timeouts.Value `tfsdk:"timeouts"` // Always include timeouts
     }
     ```

     - Always name the model following the pattern `ResourceNameResourceModel`.

3. **Implement CRUD Operations**
   - Implement the Create, Read, Update, and Delete methods for the resource. Use the resource template as a starting point.
   - **Use comments** in each CRUD function to describe the logic flow and clarify the purpose of each step.
   - **Do not build API request bodies directly in the CRUD functions.** Instead, delegate request construction to separate functions (e.g., `constructResource`). The CRUD functions should focus on the API calls and the overall logic flow, not on the details of request assembly.
   - Use the `crud.HandleTimeout` helper for operation timeouts.
   - Use `ReadWithRetry` for consistent state refresh after create and update.
   - Use the Graph SDK for every call it can make. Calls the SDK cannot make go through the Config API in `internal/services/common/custom_requests`, and only for one of the four SDK deviations, documented with the comment block in the package doc of `crud.go`. Never build `RequestInformation` or call `adapter.Send*` in a resource package. See [ADR-004](./adr/adr_004_custom_requests_transport.md).
   - Example (Create):
     ```go
     // Create handles the Create operation for the resource.
     // - Retrieves the plan
     // - Constructs the request body using a helper
     // - Calls the API
     // - Handles errors and sets state
     func (r *ResourceTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
         var plan ResourceTemplateResourceModel
         resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Parse plan
         if resp.Diagnostics.HasError() { return }
         ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
         if cancel == nil { return }
         defer cancel()
         requestBody, err := constructResource(ctx, &plan) // Request construction is separated
         if err != nil {
             resp.Diagnostics.AddError("Error constructing resource", err.Error())
             return
         }
         resource, err := r.client.
             DeviceManagement().
             ResourceTemplates().
             Post(ctx, requestBody, nil)
         if err != nil {
             errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
             return
         }
         plan.ID = types.StringValue(*resource.GetId())
         resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...) // Save state
         if resp.Diagnostics.HasError() { return }
         // Call Read with retry to get the initial resource state
         readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
         stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
         opts := crud.DefaultReadWithRetryOptions()
         opts.Operation = constants.TfOperationCreate
         opts.ResourceTypeName = ResourceName
         err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
         if err != nil {
             resp.Diagnostics.AddError("Error reading resource state after create", err.Error())
             return
         }
     }
     ```
   - Example (Read):
     ```go
     // Read handles the Read operation for the resource.
     // - Retrieves the state
     // - Calls the API to get the latest data
     // - Maps the API response to Terraform state
     func (r *ResourceTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
         var state ResourceTemplateResourceModel
         resp.Diagnostics.Append(req.State.Get(ctx, &state)...) // Parse state
         if resp.Diagnostics.HasError() { return }
         ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
         if cancel == nil { return }
         defer cancel()
         resource, err := r.client.
             DeviceManagement().
             ResourceTemplates().
             ByDeviceAndAppManagementResourceTemplateId(state.ID.ValueString()).
             Get(ctx, nil)
         if err != nil {
             errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
             return
         }
         mapRemoteStateToTerraform(ctx, &state, resource)
         resp.Diagnostics.Append(resp.State.Set(ctx, &state)...) // Save state
     }
     ```
   - Example (Update):
     ```go
     // Update handles the Update operation for the resource.
     // - Retrieves the plan
     // - Constructs the request body using a helper
     // - Calls the API
     // - Handles errors and refreshes state
     func (r *ResourceTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
         var plan ResourceTemplateResourceModel
         resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Parse plan
         if resp.Diagnostics.HasError() { return }
         ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
         if cancel == nil { return }
         defer cancel()
         requestBody, err := constructResource(ctx, &plan)
         if err != nil {
             resp.Diagnostics.AddError("Error constructing resource for update method", err.Error())
             return
         }
         _, err = r.client.
             DeviceManagement().
             ResourceTemplates().
             ByDeviceAndAppManagementResourceTemplateId(plan.ID.ValueString()).
             Patch(ctx, requestBody, nil)

         if err != nil {
             errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
             return
         }

         // Call Read with retry to get the updated resource state
         readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
         stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}
         opts := crud.DefaultReadWithRetryOptions()
         opts.Operation = constants.TfOperationUpdate
         opts.ResourceTypeName = ResourceName
         err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
         if err != nil {
             resp.Diagnostics.AddError("Error reading resource state after update", err.Error())
             return
         }
     }
     ```
     
   - Example (Delete):
     ```go
     // Delete handles the Delete operation for the resource.
     // - Retrieves the state
     // - Calls the API to delete the resource
     // - Handles errors and removes state
     func (r *ResourceTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
         var data ResourceTemplateResourceModel
         resp.Diagnostics.Append(req.State.Get(ctx, &data)...) // Parse state
         if resp.Diagnostics.HasError() { return }
         ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
         if cancel == nil { return }
         defer cancel()
         err := r.client.
             DeviceManagement().
             ResourceTemplates().
             ByDeviceAndAppManagementResourceTemplateId(data.ID.ValueString()).
             Delete(ctx, nil)
         if err != nil {
             errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
             return
         }
         resp.State.RemoveResource(ctx)
     }
     ```
   - For more details and up-to-date patterns, see the [resource template implementation](../../internal/services/resources/_resource_template/).

4. **Map Remote State to Terraform State**
   - Implement a function to map the API response to the Terraform state model.
   - Example:

     - Always map the remote state to the Terraform state model.
     - Always use the `tflog` package for debug and trace logging.
     - Always use the `types.StringValue` function to set the value of the Terraform state model.
     - Always use the `convert.GraphToFrameworkString
` function to convert the string pointer to a string.
     - Always use the `types.StringValue` function to set the value of the Terraform state model.

     ```go
     func mapRemoteStateToTerraform(ctx context.Context, data *ResourceTemplateResourceModel, remoteResource graphmodels.DeviceAndAppManagementAssignmentFilterable) {
         if remoteResource == nil {
             tflog.Debug(ctx, "Remote resource is nil")
             return
         }
         data.ID = types.StringValue(convert.GraphToFrameworkString
(remoteResource.GetId()))
         // Map other fields as needed
     }
     ```

5. **Write Constructors and Register the Resource**
   - Implement the resource constructor and register it in the provider's resource list.
   - Example:

     ```go
     func NewResourceTemplateResource() resource.Resource {
         return &ResourceTemplateResource{}
     }
     ```

6. **Testing**
   - Write acceptance and unit tests for your resource. Use the test helpers and patterns found in the provider's test files.

## Reference: Resource Template

- CRUD: `internal/resources/_resource_template/crud.go`
- Model: `internal/resources/_resource_template/model.go`
- State: `internal/resources/_resource_template/state.go`

Use these files as a starting point for new resources. Follow the patterns and structure to ensure consistency across the provider.

## Additional Tips

- Always include a link to the relevant Microsoft Graph API documentation in your resource files.
- Align your graph api implementation with the behaviour of the gui for that resource.
You will often find that resources only exist in the beta api with no v1.0 equivilent.
- Prefer v1.0 endpoints for stable features; use beta only when necessary.
- Use the `tflog` package for debug and trace logging.
- Handle errors and diagnostics carefully to provide clear feedback to users.
- Keep resource logic minimal and focused; avoid unnecessary abstraction until needed.

For questions or to discuss development, join the [community Discord](https://discord.gg/Uq8zG6g7WE).

## Typical Package Structure

A typical resource package in this provider is organized as follows:

- `resource.go`: Resource registration, schema definition, and provider integration logic.
- `model.go`: Data model definitions for the resource and any nested objects.
- `crud.go`: Implementation of the Create, Read, Update, and Delete operations. Only API calls and logic flow should be here; request construction and state mapping are delegated to helpers.
- `construct_*.go`: Functions for building API request bodies from the Terraform model. Keeps CRUD logic clean and focused.
- `state_*.go`: Functions for mapping API responses to the Terraform state model.
- `modify_plan.go`: (Optional) Logic for plan modification or diff suppression.
- `*_assignment.go`: (Optional) Assignment-specific logic for the resource.
- `model_*.go`: (Optional) Additional model definitions for complex nested objects.
- `resource_docs/`: (Optional) Directory for resource-specific documentation.

**Best Practice:**

- Keep each file focused on a single responsibility (CRUD, model, construction, state mapping, schema, etc.).
- Use comments only for describing unclear code or api insights. Do not use comments that
describe the code. The code should be self explanatory. Use docstrings to explain the purpose of each function instead.
- For complex resources, break out logic into multiple files as needed to maintain readability and maintainability.

## Helpers

These rules decide when a helper may exist, where it lives and what it is called. The standard functions (the CRUD methods, `constructResource`, `MapRemoteStateToTerraform`, and the schema and metadata methods) are not helpers.

### DRY: shared mechanics go to common packages

- Generic, configurable CRUD mechanics go in `internal/services/common/crud`: retries, read-back after a write, per-key write locks, precondition retries. If two resource packages need the same mechanism, add it there with options for the parts that differ, instead of copying it.
- Every Graph call the SDK cannot make goes through the Config API in `internal/services/common/custom_requests`. A resource package never builds `RequestInformation` or calls `adapter.Send*`. See [ADR-004](./adr/adr_004_custom_requests_transport.md).
- Logic specific to one resource stays in that resource's package.

### When a new common package is justified

Only for a genuinely new class of helper, one that fits neither `crud` nor `custom_requests` nor an existing package under `internal/services/common`. A helper that two resources happen to share is not a new class: generalise it into the package that already owns that kind of mechanism.

### The local helper rule

A package-local helper may exist only if it meets all three conditions:

1. It is called directly from 2 or more CRUD methods (`Create`, `Read`, `Update`, `Delete`, `ImportState`, `ModifyPlan`). A call from another helper doesn't count.
2. It holds about 10 or more lines of logic.
3. Its name follows the naming convention below.

Otherwise the code stays inline in the CRUD method, even if that makes the method longer.

### Naming: verb + Graph entity + qualifier

A helper's name is a verb, then the Graph entity it works on, then a qualifier, for example `findTargetedAppByPackageID` or `constructTargetAppsRequest`. The verb decides the file:

| Verb | Form | File |
| --- | --- | --- |
| `construct` | `construct<Entity><Qualifier>`, e.g. `constructTargetAppsRequest` | `construct.go` |
| `map` | `map<Entity>ToState`, e.g. `mapAssignmentsToState` | `state.go` |
| `find` | `find<Entity>By<Key>`, e.g. `findTargetedAppByPackageID` | `state.go` |
| `waitFor` | `waitFor<Entity><Qualifier>`, e.g. `waitForPolicyDeletion` | `crud.go` |
| `lock` | `lock<Entity><Qualifier>`, e.g. `lockPolicyTargetApps` | `crud.go` |

Generic names are banned, for example `matches`, `settings`, `requestBody`, `tryClaim`, `readbackOptions` and `writableApp`. They say nothing about the Graph entity, so a reader has to open the function to learn what it does.

### The state mapper: `MapRemoteStateToTerraform`

The function that maps the Graph response onto the resource's Terraform state is always named `MapRemoteStateToTerraform`, and it lives in `state.go`. It is a required part of the CRUD pattern, so the caller-count and line-count parts of the helper rule don't apply to it.

Other mappers keep the `map<Entity>ToState` convention, for example `mapAssignmentsToState`.

`MapRemoteResourceStateToTerraform` is a legacy name, and so are lower-case variants such as `mapRemoteStateToTerraform` and `mapRemoteState`. Rename them to `MapRemoteStateToTerraform` when you next touch their package.

### Before and after: Android managed mobile app

Before, the Android managed mobile app resource had two helpers for one request. `writableApp` has a banned generic name, and `constructPolicyApps` is a five-line wrapper:

```go
// construct.go

// writableApp rebuilds an entry without the server-assigned id and version, which
// targetApps derives itself.
func writableApp(app graphmodels.ManagedMobileAppable) graphmodels.ManagedMobileAppable {
	pkg := packageID(app)
	if pkg == "" {
		return app
	}
	identifier := graphmodels.NewAndroidMobileAppIdentifier()
	identifier.SetPackageId(&pkg)
	entry := graphmodels.NewManagedMobileApp()
	entry.SetMobileAppIdentifier(identifier)
	return entry
}

// constructPolicyApps builds the complete targetApps replacement request.
func constructPolicyApps(apps []graphmodels.ManagedMobileAppable) *deviceappmanagement.ManagedAppPoliciesItemTargetAppsPostRequestBody {
	body := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	body.SetApps(apps)
	return body
}

// crud.go (Create; Delete repeats the loop)
desired := make([]graphmodels.ManagedMobileAppable, 0, len(currentApps)+1)
for _, app := range currentApps {
	desired = append(desired, writableApp(app))
}
desired = append(desired, newApp)

if err := r.client.DeviceAppManagement().ManagedAppPolicies().ByManagedAppPolicyId(policyID).TargetApps().Post(ctx, constructPolicyApps(desired), nil); err != nil {
```

After, one helper builds the whole request. It is called from `Create` and `Delete`, holds the logic of both old helpers, and is named for the Graph entity and what it builds:

```go
// construct.go

// constructTargetAppsRequest builds the targetApps replacement request. Graph assigns each
// entry's id and version itself, so Android entries are rebuilt with only their package ID.
func constructTargetAppsRequest(apps []graphmodels.ManagedMobileAppable) *deviceappmanagement.ManagedAppPoliciesItemTargetAppsPostRequestBody {
	entries := make([]graphmodels.ManagedMobileAppable, 0, len(apps))
	for _, app := range apps {
		identifier, ok := app.GetMobileAppIdentifier().(graphmodels.AndroidMobileAppIdentifierable)
		if !ok || identifier.GetPackageId() == nil {
			entries = append(entries, app)
			continue
		}
		packageID := graphmodels.NewAndroidMobileAppIdentifier()
		packageID.SetPackageId(identifier.GetPackageId())
		entry := graphmodels.NewManagedMobileApp()
		entry.SetMobileAppIdentifier(packageID)
		entries = append(entries, entry)
	}
	body := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	body.SetApps(entries)
	return body
}

// crud.go (Create)
desired := append(slices.Clone(currentApps), newApp)

if err := r.client.DeviceAppManagement().ManagedAppPolicies().ByManagedAppPolicyId(policyID).TargetApps().Post(ctx, constructTargetAppsRequest(desired), nil); err != nil {

// crud.go (Delete)
remaining := make([]graphmodels.ManagedMobileAppable, 0, len(currentApps))
for _, app := range currentApps {
	if app.GetId() == nil || *app.GetId() != appID {
		remaining = append(remaining, app)
	}
}

if err := r.client.DeviceAppManagement().ManagedAppPolicies().ByManagedAppPolicyId(policyID).TargetApps().Post(ctx, constructTargetAppsRequest(remaining), nil); err != nil {
```

### The helper report

Every PR that changes a package under `internal/services` runs `tools/helperscan` on the changed packages, except the shared packages in `internal/services/common`. The report appears in the log of the "Plan test partitions" job and in its step summary. It is report-only for now and never fails the job. For each helper it shows:

- `LINES`: the span of the declaration, from `func` to the closing brace. A helper of 10 lines or fewer counts as short.
- `REFS`: how many places in the package's non-test files reference it. Recursive calls don't count.
- `CRUD CALLERS`: the distinct CRUD methods that reference it directly. Calls inside a closure in a CRUD method count for that method.
- `NAME`: `ok`, or why the name fails the convention above.
- `RULE`: `pass`, or which parts of the helper rule it misses (`callers`, `lines`, `name`). The state mapper shows `exempt`, or `fails: name` if it is not in `state.go`. The legacy state mapper names fail `NAME` with a message to rename them to `MapRemoteStateToTerraform`.

Kiota interface methods (`Serialize`, `GetFieldDeserializers`, `GetAdditionalData`, the `create…FromDiscriminatorValue` factories) and framework validator methods (`Description`, `MarkdownDescription`, `Validate…`, `PlanModify…`) are counted in the summary but left out of the table. Run it locally with:

```sh
go run ./tools/helperscan internal/services/resources/<path to package>
```

## Architecture Diagram

See the [diagram](./resource_architecture.mmd) for a visual overview of function relationships and flow in a complex resource package.

## Integrating Your Resource

After building your resource, you must register it in `internal/provider/resources.go`.

- Import your resource package at the top of the file, following the existing naming conventions.
- Add your resource's constructor (e.g., `graphBetaDeviceManagementSettingsCatalog.NewSettingsCatalogResource`) to the returned slice in the `Resources` function.
- The naming pattern is typically: `graphBeta<DeviceManagement|DeviceAndAppManagement|Groups|...><ResourceName>.New<ResourceName>Resource`.

## Terraform Registry Resource Template

A template for the Terraform Registry documentation is required for each resource. tf registry docs are generated using tfplugindocs and will be populated when the Generate Docs pipeline workflow is run on pushes into main.

The template should be added to the `templates/resources` directory. See `templates/resources/graph_beta_device_management_settings_catalog.md.tmpl` for an example. Your template should include:

- Title, subcategory, and description
- Microsoft documentation links
- API permissions
- Version history
- Example usage (referencing an example `.tf` file)
- Import instructions

The template should be named `graph_beta_<resource_name>.md.tmpl` for beta resources and `graph_<resource_name>.md.tmpl` for v1.0 resources.

## Example Usage and Import

Below is an example of a resource usage and import for terraform registry documentation:

```hcl
resource "microsoft365_graph_beta_device_management_settings_catalog" "example" {
  name        = "Example Catalog Policy"
  description = "Example policy for demonstration purposes"
  platforms   = "windows10"
  configuration_policy {
    settings {
      id = "example-setting-id"
      setting_instance {
        odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingInstance"
        setting_definition_id = "example-definition-id"
        simple_setting_value {
          odata_type = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
          value      = "example-value"
        }
      }
    }
  }
}
```

To import an existing resource:

```sh
terraform import microsoft365_graph_beta_your_resource.example <resource_id>
```
