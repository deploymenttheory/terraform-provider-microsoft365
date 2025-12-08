# AI Custom Instructions – Terraform Provider Microsoft365

These instructions, guide ai tools to follow our project's conventions and best practices when suggesting code. They cover how to format code, name resources and attributes, structure implementations, sdk usage, data type conversion and how write tests in this repository. By following these guidelines, ai's suggestions should align with the project's style and help contributors produce high-quality, consistent code. Always consider existing patterns in the repository—when in doubt, review similar resources or tests for reference and keep the new code idiomatic to the project's practices.

## Development Setup & Workflow

- Use the provided **Makefile** commands for all build and test tasks:
  - `make install` to compile the provider code.
  - `make lint` to run linters and ensure code style compliance.
  - `make unittest` to run all unit tests (optionally use `TEST=<prefix>` to run tests matching a name prefix, e.g. `make unittest TEST=Environment` to run tests named with that prefix). This filters tests by regex `^(TestAcc|TestUnit)<prefix>`.
  - `make acctest TEST=<prefix>` to run acceptance tests (integration tests) matching a prefix. Always provide a specific test prefix to limit scope, and run these tests **only with user consent** (they run against real cloud resources). Note that `make acctest` automatically sets `TF_ACC=1` (no need to set it manually).
  - `make userdocs` to regenerate documentation
  - `make precommit` to run all checks once code is ready to commit. As a copilot agent you don't want to run this command as it will timeout for you. Read the makefile content and run needed commands manually.
  - `make coverage` to run all unit tests and output a code coverage report. It also shows the files that have changed on this branch to help target coverage suggestions to files in the current PR.
  - For comprehensive testing details including test types, patterns, and infrastructure components, see the **Testing Infrastructure and Organization** section below.
- Always run the above `make` commands from the repository root (e.g. in the `/workspaces/terraform-provider-microsoft365` directory).
- **Never run** `terraform init` inside the provider repo. Terraform is only used in examples or tests; initializing in the provider directory is not needed and may cause conflicts.
- Do not manually edit files under the `/docs` folder. These files are auto-generated from the schema `MarkdownDescription` attributes. Instead, update schema's `MarkdownDescription` in code and run `make userdocs` to regenerate documentation.
- To try out an example configuration, navigate to its directory under `/examples` and run `terraform apply -auto-approve` (ensure you've built the provider and set it in your Terraform plugins path beforehand).

## File and Folder Structure

### Resource Organization

- Organize all resource implementations within the `internal/services/resources` directory, with each resource in its own subdirectory within their respective resource category directory `internal/services/resources/resource_category_placeholder/graph_api_type`. (e.g. `internal/services/resources/device_management/graph_beta/device_management_scripts`)
- There are the following resource categories:
  - `applications`
  - `backup_storage`
  - `device_and_app_management`
  - `device_management`
  - `education`
  - `extensions`
  - `external_data_connections`
  - `files`
  - `financials`
  - `groups`
  - `identity_and_access`
  - `industry_data_etl`
  - `m365_admin`
  - `people_and_workplace_intelligence`
  - `security`
  - `sites_and_lists`
  - `teamwork_and_communications`
  - `users`
- There are two types of `{graph_api_type}`: `graph_beta` and `graph_v1.0`.
- Name resource directories using lowercase words with underscores (e.g., `cloud_pc_user_setting`, `group_member_assignment`).
- Choose resource names that reflect the Microsoft365 domain they represent. Prefer the commonly used term for the resource e.g `settings_catalog` over the api type e.g `configuration_policy`.

### Resource Files

Each resource directory MUST contain:

- **Models File**: Create a single `model.go` file containing all data models for the resource. The main struct should be named `{ResourceName}ResourceModel` (e.g., `MacOSPlatformScriptResourceModel`).
- **Resource Implementation Files**: Include `resource.go` (main resource struct and schema), `crud.go` (CRUD logic), and any additional files for state, construction, validation, or plan modification as needed (e.g., `state_assignment.go`, `construct_assignment.go`, `validate_assignment.go`, `modify_plan.go`).
- **Test Files**: Place unit tests in `resource_test.go` and acceptance tests in `resource_acceptance_test.go`. Create `resource_test_helper.go` in the resource package (not `_test` package) implementing the `types.TestResource` interface for Graph API existence checks used by acceptance tests.
- **Mock Responders**: Create `mocks/responders.go` implementing the `mocks.MockRegistrar` interface with `RegisterMocks()`, `RegisterErrorMocks()`, and `CleanupMockState()` methods.
- **Mock Data Files**: Organize test data/fixtures in a `tests/` subdirectory within the resource directory. Store JSON response files in `tests/responses/[operation]/` (e.g., `validate_create/`, `validate_get/`, `validate_update/`, `validate_delete/`). Name JSON files as `<method>_<object>_<scenario>.json` (e.g., `post_ip_named_location_success.json`).
- **Terraform Configs**: Store Terraform configurations in `tests/terraform/unit/` for unit tests and `tests/terraform/acceptance/` for acceptance tests.
- **Documentation Files**: If present, place resource-specific documentation or model JSONs in a `resource_docs/` subdirectory.

### Data Source Organization

- Organize all datasource implementations within the `internal/services/datasources` directory, with each datasource in its own subdirectory within their respective datasource category directory. e.g. `internal/services/datasources/datasource_category_placeholder/graph_api_type`.
- There are two types of api types: `graph_beta` and `graph_v1.0`.
- Name datasource directories using lowercase words with underscores (e.g., `cloud_pc_user_setting`, `group_member_assignment`).
- Each datasource SHALL align with the resource directory naming convention.

### Data Source Files

Each resource directory MUST contain:

- **Models File**: Create a single `model.go` file containing all data models for the datasource. The main struct should be named `{DataSourceName}DataSourceModel` (e.g., `MobileAppDataSourceModel`).
- **Data Source Files**: Name the main file as `datasource.go` containing the datasource struct definition, metadata, schema, and configuration.
- **Read File**: Create a `read.go` file containing the Read method implementation for the datasource.
- **State File**: Create a `state.go` file containing functions for mapping API responses to datasource models.
- **Helper Files**: Create additional helper files (e.g., `helpers.go`) as needed for utility functions specific to the datasource.

### Example Files

- Place each example in its own directory under `examples/microsoft365_{graph_api_type}/`, named with the full resource or data source name `microsoft365_{graph_api_type}_{resource_category}_{resource_name}`. (e.g., `/examples/resources/microsoft365_graph_beta_device_management_assignment_filter` or `/examples/microsoft365_graph_v1.0/microsoft365_graph_v1.0_device_and_app_management_mobile_app`)
- **Resource Examples:**
  - Use a single `resource.tf` file per directory.
  - The example should align precisely with the resource schema.
  - Use placeholder values or reference data sources for IDs.
  - If the resource supports import, include an `import.sh` script with:
    - A comment describing the placeholder (e.g., `# {resource_id}`)
    - The exact `terraform import` command with the correct resource type and placeholder.
- **Data Source Examples:**
  - Use a single `datasource.tf` file per directory when a datasource is required.
  - Include multiple data source blocks if needed e.g for odata scenarios, each with a descriptive comment of the scenario.
  - Define outputs in the same file, referencing data source attributes.
- **General:**
  - Keep example code simple and focused on demonstrating one clear use case per block.
  - Use comments to explain the purpose of each example and any placeholders.

## Naming Conventions

- **Resource and Data Source Names:** Follow the existing naming pattern of prefixing with `microsoft365_` followed by the API type (`graph_beta` or `graph_v1.0`) and the resource category. For example, a resource is named `microsoft365_graph_beta_device_management_assignment_filter`. Use lowercase with underscores for Terraform resource/data names.
- **Attribute Naming:** Name resource attributes to match Microsoft365 terminology. Prefer the modern, user-friendly terms used in the current Microsoft365 API/UX/[Official Documentation](https://learn.microsoft.com/en-us/graph/) over deprecated names. Keep names concise but descriptive of their function in the resource.
- **Model Field Naming Pattern**: For optional block attributes, use a pointer to a shared or local model struct (e.g., `Assignments *sharedmodels.DeviceManagementScriptAssignmentResourceModel`).
  - Name the Go field in PascalCase (e.g., `Assignments`).
  - Use a type of `<SubResourceName>ResourceModel` (e.g., `DeviceManagementScriptAssignmentResourceModel`).
  - Set the Terraform schema tag to the lower_snake_case version of the field (e.g., `tfsdk:"last_modified_date_time"`).
  - The sub-resource struct should be named `<SubResourceName>ResourceModel`.
- **Test Function Naming:** Name test functions with a prefix indicating their type. **Acceptance test** functions should start with `TestAcc` and **unit test** functions with `TestUnit` (this allows filtering tests by type). Also, name test files' package with a `_test` suffix (e.g. `package environment_test`) to ensure tests access the provider only via its public interface.
- **Resource/Data Source Factory:** For each resource and data source, create a new function named `New<ResourceName>Resource` or `New<DataSourceName>DataSource` that returns the appropriate type.
- **Client Factory:** When implementing a client factory, name it `New<Service>Client` (e.g., `NewSolutionClient`).

## Data Type Conversion for resource construction and state mapping

- Use the data type conversion utilities in the following packages:
  - **Constructors (`internal/services/common/convert/framework_to_graph.go`)**: Contains functions for converting Terraform types to Microsoft Graph SDK types when constructing API requests:
    - `FrameworkToGraphString`, `FrameworkToGraphBool`, `FrameworkToGraphInt32`, etc.: Convert Terraform primitive types to pointers for Graph API setters
    - `FrameworkToGraphEnum`: Convert string values to enumeration types
    - `FrameworkToGraphTime`, `FrameworkToGraphDateOnly`, `FrameworkToGraphTimeOnly`: Convert string values to time/date types
    - `FrameworkToGraphStringList`, `FrameworkToGraphStringSet`: Convert Terraform collections to string slices
    - `FrameworkToGraphBytes`: Convert string values to byte slices
    - `FrameworkToGraphISODuration`: Parse ISO 8601 duration strings
    - `FrameworkToGraphUUID`: Parse and convert string UUIDs
    - `FrameworkToGraphBitmaskEnumFromSet`: Convert string sets to bitmask enums
  - **State (`internal/services/common/convert/graph_to_framework.go`)**: Contains functions for converting Microsoft Graph SDK types to Terraform types when mapping API responses to state:
    - `GraphToFrameworkString`, `GraphToFrameworkBool`, `GraphToFrameworkInt32`: Convert primitive pointers to Terraform types
    - `GraphToFrameworkTime`, `GraphToFrameworkDateOnly`, `GraphToFrameworkTimeOnly`: Convert time types to strings
    - `GraphToFrameworkEnum`: Convert enumeration values to strings
    - `GraphToFrameworkBytes`: Convert byte arrays to strings
    - `GraphToFrameworkStringList`, `GraphToFrameworkStringSet`: Convert string slices to Terraform collections
    - `GraphToFrameworkUUID`: Convert UUID values to strings
    - `GraphToFrameworkBitmaskEnumAsSet`: Convert bitmask enums to string sets
    
## Comments and Documentation

- Write Go comments only on exported functions, types, and methods to explain their purpose, parameters, and return values when it adds clarity.
- Focus comments on **why** something is done if it's not obvious from the code.
- Avoid redundant comments that just restate the code or don't provide additional insight.
- When defining resource or data source schema, **always use** the `MarkdownDescription` field for documentation. Do **not** use the deprecated `Description` field. Markdown descriptions will be used to auto-generate docs, so make them clear and user-friendly, and include links to topics in the [official Microsoft365 docs](https://learn.microsoft.com/en-us/graph/) when helpful.

## Code Organization and Implementation Guidelines

### Frameworks

- **Terraform Plugin Framework:** Use [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework) for implementing resources and data sources. Avoid using legacy Terraform SDK constructs.
- **Azure Identity Client:** Use [Azure Identity Client Module for Go](https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azidentity) for implementing authentication methods
- **Microsoft Graph SDK:** Use the [Microsoft Graph SDK](https://github.com/microsoftgraph/msgraph-sdk-go) for implementing graph v1.0 API calls.
- **Microsoft Graph Beta SDK:** Use the [Microsoft Graph Beta SDK](https://github.com/microsoftgraph/msgraph-beta-sdk-go) for implementing graph beta API calls.
- **Terraform Plugin Framework Validators:** Use the [Terraform Plugin Framework Validators](https://github.com/hashicorp/terraform-plugin-framework-validators) for implementing validation logic.
- **Custom Terraform Plugin Framework Validators:** Use the [Custom Terraform Plugin Framework Validators](/internal/services/common/validators) for implementing custom validation logic. Prefer using the custom validators over the built-in validators only when hashicorp provided validators fall short.
- **Terraform Plugin Framework Plan Modifiers:** Use the [Terraform Plugin Framework Plan Modifiers](https://developer.hashicorp.com/terraform/plugin/framework/plan-modifiers) for implementing plan modification logic.
- **Terraform Plugin Framework Timeouts:** Use the implementation found in `internal/services/common/crud/timeout.go` for implementing timeouts.
- **Terraform Plugin Framework Helpers:** Use the [Terraform Plugin Framework Helpers](https://developer.hashicorp.com/terraform/plugin/framework/helpers) for implementing helper functions.

### Common Utilities

- **API Layer:** Use the common API client functionality in `internal/api` for making Microsoft365 API calls:
  - Use service-specific clients that build upon the common API layer
  - Handle API errors gracefully and return detailed error diagnostics
  - Use the retry mechanisms provided by the API layer for transient failures

- **Constants:** Reference centralized constants from `internal/constants/constants.go` instead of hardcoding values for:
  - API endpoints and URL paths
  - Common string literals and configuration keys
  - Regex patterns and expressions
  - Status codes and enum values used across the provider

- **Error Handling:** Leverage the error handling utilities defined in `internal/resources/common/errors` for consistent error management across all resources and data sources:
- Always use `errors.HandleKiotaGraphError(ctx, err, resp, <operation>, <permissions>)` to process and report errors from Microsoft Graph API calls in CRUD operations. This ensures:
  - Standardized extraction and categorization of error details (status code, error code, message, etc.).
  - Automatic logging of error context and details using `tflog`.
  - Proper handling of special cases (e.g., 404 removes resource from state, 401/403 provides permission hints, 429/503 handles throttling and service unavailability).
  - Addition of user-friendly error messages to diagnostics for Terraform users.
- Do not manually parse or handle Graph API errors in resource CRUD methods; always delegate to the error handling package.
- For custom error handling logic, extend or use the helper functions provided in `internal/resources/common/errors` (e.g., for retry logic, error categorization, or extracting additional error details).
- When adding new error handling, ensure it integrates with the existing error handling framework for consistency and observability.


- **Custom Types:** Utilize the custom types defined in `internal/customtypes` for specialized data handling:
  - Use custom Terraform schema types where appropriate
  - Leverage provided plan modifiers and validators for custom types
  - Follow the patterns established for marshaling/unmarshaling custom types

- **Validators:** Apply common validators from `github.com/hashicorp/terraform-plugin-framework-validators` package or `internal/validators` to ensure consistent validation logic:
  - Use provided validators for common validation requirements (UUID format, string length, etc.)
  - Chain validators for attributes that need multiple validation rules
  - Add resource-specific validation only when generic validators are insufficient

- **Helper Functions:** Make use of utility functions in `internal/helpers` to reduce duplication:
  - Use the helper functions for common tasks like state management and data conversion
  - Leverage the provided resource base types and embedded functionality
  - Follow established patterns for logging, attribute access, and diagnostics handling


### Best Practices

- **Method Scope:** Methods that are not used outside the namespace scope should be kept private (unexported).

- **API Interaction:**
  - Use the service-specific clients from the provider for all API calls.
  - Handle asynchronous operations with proper polling and timeouts.
  - Validate input values before sending API requests when needed.
  - Ensure that you always pass context `ctx` into long-running or asynchronous operations like API calls

- **Error Handling:**
  - Add context to API errors using the provider's error types from `internal/customerrors`.
  - Return detailed diagnostics with `resp.Diagnostics.AddError()` for user-friendly messages.
  - Distinguish between different error types (authentication, validation, not found, etc.).
  - Log API responses and errors at debug level using `tflog.Debug` for troubleshooting.

- **Request Context:**
  - Resources and Data Sources should call `ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)` and `defer exitContext()` near the beginning of any method from resource or datasource interfaces.

### Guidelines for Resources

#### Resource Structure and Interfaces

- Implement `resource.Resource` interface for all resources.
- Implement `resource.ResourceWithImportState` for all resources.
- Implement `resource.ResourceWithConfigure` to configure the resource with the provider client.
- Implement `resource.ResourceWithModifyPlan` when plan modification is needed.
- Structure resources in a consistent pattern by ordering methods: `Metadata`, `FullTypeName`, `Configure`, `ImportState`, `Schema`, `Create`, `Read`, `Update`, `Delete`, `ModifyPlan`.
- Define constants for resource name and timeouts (e.g., `ResourceName`, `CreateTimeout`, `ReadTimeout`, `UpdateTimeout`, `DeleteTimeout`).
- Add required client fields to your resource struct to access APIs.
- Include `ReadPermissions` and `WritePermissions` fields in your resource struct to specify required permissions.
- Include `ResourcePath` field to specify the API endpoint path.
- Name factory functions as `New<ResourceName>Resource()` (e.g., `NewAssignmentFilterResource`).
- Return a new instance of your resource struct from the factory function with required permissions and resource path set.

#### Resource Schema Definition

- Define complete schemas with proper attribute types (String, Int32, Bool, etc.).
- Mark attributes explicitly as `Required`, `Optional`, or `Computed`.
- Use `Int32` for attributes that are integers and are less than 2^31. Never use `Int64` for these attributes.
- Use `Computed: true` for server-generated fields like IDs.
- Use `Optional: true` with `Computed: true` for fields that can be specified or defaulted by the service.
- Apply `RequiresReplace` plan modifier to immutable attributes that necessitate resource recreation when changed.
- Apply `UseStateForUnknown` modifier to computed fields to prevent unnecessary diffs during planning.
- Include standard timeouts using `github.com/hashicorp/terraform-plugin-framework-timeouts`.
- Write clear `MarkdownDescription` for each attribute (do not use the deprecated `Description` field).

#### Resource State Management

- In `Create`, populate state with all resource attributes after successful creation.
- In `Read`, refresh the full state based on the current resource values from the API.
- Check for deleted resources in `Read` - when API returns 404, call `resp.State.RemoveResource(ctx)`.
- In `Update`, apply only the changed attributes and refresh state afterwards.
- Return early with appropriate diagnostics when operations cannot complete successfully.

#### Resource Validation

- Apply built-in validators from `github.com/hashicorp/terraform-plugin-framework-validators` for attribute constraints.
- Implement resource-level validation in the `ValidateConfig` method when validation involves multiple attributes.
- Add custom validators only when built-in validators are insufficient.
- Provide clear validation error messages that explain the specific constraint and how to fix it.

### Guidelines for Data Sources

#### Data Source Structure and Interfaces

- Implement the `datasource.DataSource` interface for all data sources.
- Implement the `datasource.DataSourceWithConfigure` interface to configure the data source with the provider client.
- Order data source methods consistently: `Metadata`, `Schema`, `Configure`, `Read`.
- Add required client fields to your data source struct to access APIs.
- Include `ReadPermissions` field in your data source struct to specify required permissions.
- Include `ProviderTypeName` and `TypeName` fields in your data source struct for type name management.
- Define constants for the datasource name and timeout values (e.g., `datasourceName`, `ReadTimeout`).
- Name factory functions as `New<DataSourceName>DataSource()` (e.g., `NewMobileAppDataSource`).
- Return a new instance of your datasource struct from the factory function with required permissions set.

#### Data Source Schema Definition

- Mark all attributes as `Computed: true` since data sources are read-only by design.
- For optional filter parameters, use `Required: false` and `Optional: true`.
- Define nested schemas for complex return types using appropriate collection types:
  - Use `schema.ListNestedAttribute` for collections of objects like "environments", "applications", etc.
  - Use `schema.SingleNestedAttribute` for single complex objects.
- Only include Read timeouts in timeouts schema (omit Create, Update, Delete).
- Use `map[string]schema.Attribute` for schema attributes that allow extensible field sets.
- Include output-only fields that will assist users in identifying or using the data in further resources.
- For primary list attributes (e.g., "applications", "environments"), use the plural form as the attribute name.

#### Data Source Query Parameters

- For data sources that filter results, define explicit filter attributes:
  - Common patterns include `filter_type` (e.g., "all", "id", "display_name", "odata") to specify how to filter.
  - Include `filter_value` for the value to filter by.
  - For OData filtering, include parameters like `odata_filter`, `odata_top`, `odata_skip`, `odata_select`, and `odata_orderby`.
  - Add domain-specific filters (e.g., `app_type_filter`) when appropriate.
- Support sensible combinations of filter parameters that match Microsoft365 API capabilities.
- Document filter parameters with clear examples in the `MarkdownDescription`.
- Use validators (e.g., `stringvalidator.OneOf()`) to restrict filter values to valid options.

#### Data Source Read Implementation

- Parse all input filter parameters from state at the beginning of the Read method.
- Include context propagation and handle timeouts: `ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)` with a matching `defer cancel()`.
- Validate any required filter parameters and return appropriate diagnostic errors.
- Implement different query strategies based on filter types:
  - For ID-based filtering, fetch a single item directly using the ID.
  - For OData filtering, pass OData parameters to the API.
  - For other filter types (e.g., "all", "display_name"), fetch all items and filter locally.
- Use the appropriate client method to retrieve data based on filter criteria.
- For empty results, set an empty list rather than returning an error.
- Transform API responses to data models using the appropriate state mapping functions.
- Set all fields in the state model, even those that might be nil or empty.
- Log API calls using `tflog.Debug` statements to assist troubleshooting.
- For list-type data sources, return a consistent response structure even when results vary in size.

#### Testing Data Sources

- Test all supported filter combinations in unit tests.
- Create separate test cases for each filter type (e.g., "all", "id", "display_name", "odata").
- Test domain-specific filters (e.g., `app_type_filter`) to ensure they correctly filter results.
- Verify that filtered results return the expected subset of data.
- Test edge cases like empty results, single results, and large result sets.
- For collection data sources, test accessing nested attributes and verify attribute counts.
- Ensure acceptance tests use non-destructive read-only operations.
- For data sources that return lists, test accessing list items with collection syntax.
- Create mock responses in JSON files for each test scenario, following the same directory structure as resources.

#### Data Source Documentation and Examples

- Include a representative example in the `/examples/microsoft365_{graph_api_type}/microsoft365_{graph_api_type}_{resource_category}_{resource_name}` directory (e.g., `/examples/resources/microsoft365_graph_beta_device_and_app_management_mobile_app`).
- For data sources with filter parameters, include examples showing different filtering options:
  - Filtering by ID
  - Filtering by display name
  - Using OData filters for advanced queries
  - Using domain-specific filters
- Include examples demonstrating how to access nested attributes in the output.
- Showcase practical use-cases with supporting resources when applicable.
- Describe the purpose of the data source clearly in the schema's `MarkdownDescription`.
- Link to relevant Microsoft Graph API documentation that explains the underlying API or service.
- Include references to the API endpoint used (e.g., `/deviceAppManagement/mobileApps`).

## Logging

Use the Terraform plugin logger (`tflog`) for logging within resource implementations.

- **Debug Level:** Add `tflog.Debug` statements to log major operations or key values for troubleshooting. Use debug logs liberally to aid in diagnosing issues during development or when users enable verbose logging.
- **Info Level:** Use `tflog.Info` sparingly, only for important events or high-level information that could be useful in normal logs. Too much info-level logging can clutter output, so prefer Debug for most internal details.
- **No Print/Printf:** Do not use `fmt.Println`/`log.Printf` or similar for logging. The `tflog` system ensures logs are structured and can be filtered by Terraform log level.
- **Sensitive Data:** Never log sensitive information (credentials, PII, etc.). Ensure that debug logs do not expose secrets or user data.
- **Request Context:** Do not trace the entry/exit of interface methods in resources or data source.  Instead use `EnterRequestContext` and `exitContext`

## Testing Best Practices

### Test File Organization

Each resource directory should contain the following test-related files:

- **Unit Tests:** Place in `resource_test.go` with `TestUnit*` or `Test*` naming pattern.
- **Acceptance Tests:** Place in `resource_acceptance_test.go` with `TestAcc*` naming pattern.
- **Test Resource Helper:** Create `resource_test_helper.go` in the resource package (not the `_test` package) implementing the `types.TestResource` interface for Graph API existence checks.
- **Mock Responders:** Store in `mocks/responders.go` file within the resource directory.
- **Terraform Configs:** Organize in `tests/terraform/unit/` for unit tests and `tests/terraform/acceptance/` for acceptance tests.
- **JSON Responses:** Store mock API responses in `tests/responses/[operation]/[method]_[object]_[scenario].json`.
- **Package Naming:** Use `<package>_test` suffix (e.g., `graphBetaNamedLocation_test`) to ensure tests access the provider only via its public interface.

### Fluent Assertion Pattern

Use the fluent assertion helpers from `internal/acceptance/check` for all test checks:

- Define the resource type as a package-level variable referencing the resource's `ResourceName` constant.
- Use `check.That(resourceType+".resource_name").Key("attribute").HasValue("expected_value")` for value assertions.
- Use `check.That(...).Key("attribute").Exists()` to verify computed attributes are populated.
- Use `check.That(...).Key("attribute").DoesNotExist()` to verify attributes are null or not set.
- Use `check.That(...).Key("attribute").MatchesRegex(pattern)` for pattern-based validation.
- Use `check.That(...).Key("set_attr.*").ContainsTypeSetElement("value")` to check set membership.
- Use `check.That(...).Key("list_attr.#").HasValue("count")` to verify collection lengths.

### Unit Tests

Write unit tests in `resource_test.go` covering all operations and edge cases using the `jarcoal/httpmock` library:

- **Mock Setup:** Create `setupMockEnvironment()` and `setupErrorMockEnvironment()` functions that activate httpmock, register authentication mocks, and return the resource-specific mock struct.
- **Environment Setup:** Call `mocks.SetupUnitTestEnvironment(t)` at the start of each test for consistent environment configuration.
- **Provider Factories:** Use `mocks.TestUnitTestProtoV6ProviderFactories` for the test case.
- **Cleanup:** Use `defer httpmock.DeactivateAndReset()` and `defer resourceMock.CleanupMockState()` to ensure proper cleanup.
- **Config Loading:** Use `helpers.ParseHCLFile("tests/terraform/unit/resource_name.tf")` to load Terraform configurations from the unit test directory.
- **Test Steps:** Structure tests with `resource.UnitTest()` using sequential steps for create, optionally update, and import verification.
- **Import Testing:** Include an import step with `ImportState: true` and `ImportStateVerify: true`.

### Mock Responders

Implement resource-specific mock responders in `mocks/responders.go`:

- Create a mock struct implementing the `mocks.MockRegistrar` interface with `RegisterMocks()`, `RegisterErrorMocks()`, and `CleanupMockState()` methods.
- Register the mock with `mocks.GlobalRegistry.Register()` in the package's `init()` function.
- Maintain mock state using a mutex-protected map to track created resources across CRUD operations.
- Use `mocks.LoadJSONResponse()` to load response templates from `tests/responses/` directory.
- Update response objects with request data before returning to simulate Graph API behavior.
- Handle different resource variants (e.g., IP vs country named locations) based on request body content.

### Acceptance Tests

Write acceptance tests in `resource_acceptance_test.go` for end-to-end validation against real Microsoft 365 APIs:

- **Test Resource Helper:** Create a struct in `resource_test_helper.go` implementing `types.TestResource` with an `Exists()` method that checks resource existence via HTTP call to Microsoft Graph.
- **Resource Type Variable:** Define `resourceType` and `testResource` as package-level variables referencing the resource package constants.
- **PreCheck:** Include `PreCheck: func() { mocks.TestAccPreCheck(t) }` for environment validation.
- **Provider Factories:** Use `mocks.TestAccProtoV6ProviderFactories` for the test case.
- **Destroy Verification:** Use `destroy.CheckDestroyedAllFunc(testResource, resourceType, duration)` to verify all resources of the type are destroyed after the test.
- **External Providers:** Use `ExternalProviders` map for dependencies like `random` provider for unique resource naming.
- **Test Logging:** Use `testlog.StepAction(resourceType, "action description")` in `PreConfig` to log test step actions.
- **Eventual Consistency:** Use `testlog.WaitForConsistency("resource", duration)` followed by `time.Sleep()` to handle Microsoft Graph eventual consistency.
- **Graph Existence Check:** Use `check.That(...).ExistsInGraph(testResource)` to verify resource existence in Microsoft Graph.
- **Config Loading:** Use `mocks.LoadTerraformConfigFile("config_name.tf")` to load configs from the acceptance test directory, wrapped with `acceptance.ConfiguredM365ProviderBlock(config)` to add provider configuration.
- **Import Testing:** Include import steps with `ImportStateVerify: true` and `ImportStateVerifyIgnore: []string{"timeouts"}` for computed/transient fields.

### Acceptance Test Resource Dependencies

Acceptance test Terraform configurations must create all dependent resources dynamically as part of the test itself:

- **Never rely on environment variables** for resource IDs or pre-existing resources. All test dependencies should be created within the Terraform configuration.
- **Always use `random` provider** for unique naming in all test scenarios to avoid conflicts between parallel test runs.
- **Prefer this provider's resources and data sources first** when building dependency trees. Use existing microsoft365 resources and data sources to create or reference prerequisite resources wherever possible.
- **Use external providers as fallback** (e.g., `azuread`) only when this provider does not have a suitable resource or data source for the dependency.
- **Use data sources** for stable, well-known resources (e.g., Microsoft Graph service principal via its well-known appId).
- **Order resources correctly** so Terraform handles dependencies automatically through resource references.

### Environment Requirements

- **IMPORTANT: If you don't have access to a test tenant, DO NOT modify existing acceptance tests.** Focus on unit tests only when lacking test environment access.
- Acceptance tests require valid `M365_*` environment variables for authentication only, use resource dependencies in all other circumstances.
- Use resource naming with random suffixes (via `random` provider) to avoid conflicts between parallel test runs.
- **Do not use real customer data** in test fixtures – anonymize all IDs and personal information.

### Test Coverage

Aim for **at least 80%** code coverage for unit tests on new code. `make unittest` will return a coverage score by service and overall. Focus on the service that is currently being worked on when adding tests to improve coverage.

### Examples and Documentation

Whenever a new resource or data source is added, provide an example configuration under the `/examples` directory to demonstrate usage. This helps both in documentation and in manually verifying the resource behavior. After implementing and testing, run `make userdocs` to update the documentation in `/docs` from your schema comments.

## Testing Infrastructure and Organization

### Test Types

- **Unit Tests:**
  - **Naming Pattern:** `TestUnit[ResourceName]_[Operation]_[Scenario]` (e.g., `TestUnitUserResource_Create_Minimal`)
  - **Characteristics:** Mock all API calls using `httpmock`, no real Microsoft 365 API interactions
  - **Timeout:** 10 minutes per test
  - **Parallelism:** Run with `-p 16` for fast execution
  - **Environment:** Run with `TF_ACC=0` or omitted (default)
  - **Purpose:** Test resource CRUD logic, state management, and schema validation

- **Acceptance Tests:**
  - **Naming Pattern:** `TestAcc[ResourceName]_[Operation]_[Scenario]` (e.g., `TestAccUserResource_Create_Minimal`)
  - **Characteristics:** Make real API calls to Microsoft 365 services
  - **Timeout:** 300 minutes (5 hours) to accommodate complex operations
  - **Parallelism:** Run with `-p 10` to avoid API rate limits
  - **Environment:** Require `TF_ACC=1` and valid authentication credentials
  - **Purpose:** Verify actual resource creation, modification, and deletion in Microsoft 365

### Test Commands

- **Unit Tests:**
  ```bash
  make unittest                    # Run all unit tests
  make unittest TEST=MyTest        # Run specific unit test by prefix
  go test -v -run TestUnitUserResource_Create ./path/to/package
  ```

- **Acceptance Tests:**
  ```bash
  make acctest                     # Run all acceptance tests
  make acctest TEST=MyTest         # Run specific acceptance test by prefix
  TF_ACC=1 go test -v -timeout 30m -run TestAccUserResource_Create ./path/to/package
  ```

- **Coverage and Full Suite:**
  ```bash
  make test                        # Run all tests (unit + acceptance)
  make coverage                    # Generate test coverage report with branch diff
  ```

### Test Infrastructure Components

- **Global Mock System (`/internal/mocks/`):**
  - `AuthenticationMocks`: Mock authentication endpoints and token exchanges
  - `MockGraphClients`: Mock Microsoft Graph API clients for unit testing
  - `TestUnitTestProtoV6ProviderFactories`: Provider factories for unit tests
  - `TestAccProtoV6ProviderFactories`: Provider factories for acceptance tests
  - `SetupUnitTestEnvironment(t)`: Standardized unit test environment setup
  - `TestAccPreCheck(t)`: Validate acceptance test environment and credentials
  - `GlobalRegistry`: Central registry for resource-specific mock implementations
  - `LoadJSONResponse(path)`: Load JSON response templates for mock responders
  - `LoadTerraformConfigFile(filename)`: Load Terraform configs from acceptance test directories

- **Fluent Assertion System (`/internal/acceptance/check/`):**
  - `check.That(resourceName)`: Entry point for fluent assertions on Terraform resources
  - `.Key(attributeName)`: Chain to check specific resource attribute
  - `.HasValue(value)`, `.Exists()`, `.DoesNotExist()`: Value and existence assertions
  - `.MatchesRegex(pattern)`: Pattern-based validation for dynamic values
  - `.ContainsTypeSetElement(value)`: Verify set membership
  - `.ExistsInGraph(testResource)`: Verify resource existence in Microsoft Graph (acceptance tests only)

- **Acceptance Test Infrastructure (`/internal/acceptance/`):**
  - **Destroy Verification (`destroy/`):** `CheckDestroyedAllFunc(testResource, resourceType, waitDuration)` for verifying cleanup
  - **Test Logging (`testlog/`):** `StepAction()` and `WaitForConsistency()` for test observability
  - **Provider Configuration:** `ConfiguredM365ProviderBlock(config)` for wrapping test configs with provider block
  - **HTTP Client:** `TestHTTPClient()` for direct Graph API calls in test helpers
  - **Type Definitions (`types/`):** `TestResource` interface contract for resource existence checking

- **Resource-Specific Mock System (per resource directory):**
  - **Mock Responders:** `mocks/responders.go` implementing `mocks.MockRegistrar` interface
  - **Mock State Management:** Mutex-protected state maps tracking resources across CRUD operations
  - **Cleanup Functions:** `CleanupMockState()` for resetting state between tests
  - **Setup Functions:** `setupMockEnvironment()` and `setupErrorMockEnvironment()` in resource test files

- **Test Resource Helper (per resource directory):**
  - **File:** `resource_test_helper.go` in resource package (not `_test` package)
  - **Implementation:** Struct implementing `types.TestResource` interface
  - **Method:** `Exists(ctx, client, state)` for checking resource existence via HTTP call to Graph API
  - **Usage:** Passed to `destroy.CheckDestroyedAllFunc()` and `check.That().ExistsInGraph()`

- **Test Data Organization (per resource directory):**
  - **JSON Response Files:** `tests/responses/[operation]/[method]_[object]_[scenario].json`
  - **Terraform Unit Configs:** `tests/terraform/unit/` for configs loaded via `helpers.ParseHCLFile()`
  - **Terraform Acceptance Configs:** `tests/terraform/acceptance/` for configs loaded via `mocks.LoadTerraformConfigFile()`
  - **Operation Directories:** `validate_create/`, `validate_get/`, `validate_update/`, `validate_delete/`
  - **Scenario Coverage:** Success, error, minimal, maximal, and hybrid configurations

### Required Environment Variables for Acceptance Tests

```bash
M365_TENANT_ID                   # Azure AD tenant ID
M365_AUTH_METHOD                 # Authentication method to use
M365_CLOUD                       # Cloud environment (public, gcc, gcchigh, dod, china)
M365_CLIENT_ID                   # Application client ID
# Additional auth-specific variables based on M365_AUTH_METHOD
```

### Testing Best Practices Specific to This Provider

- **Test File Organization:** Use distinct files for different purposes within each resource directory
  - **Unit Tests:** `resource_test.go` with `TestUnit*` or `Test*` functions using httpmock
  - **Acceptance Tests:** `resource_acceptance_test.go` with `TestAcc*` functions using real API calls
  - **Test Resource Helper:** `resource_test_helper.go` in resource package implementing `types.TestResource`
  - **Package Naming:** Unit and acceptance tests use `<package>_test` suffix; test helper uses resource package

- **Fluent Assertions:** Use `check.That()` pattern from `internal/acceptance/check` for all test assertions
  - Reference resource type via constant from resource package (e.g., `resourceType = graphBetaNamedLocation.ResourceName`)
  - Chain `.Key("attribute")` for attribute-specific assertions
  - Use `.ExistsInGraph(testResource)` for acceptance tests to verify Graph API state

- **Mock Data Organization:** Store mock responses in resource-specific directory structure
  - **Directory Structure:** `tests/responses/[operation]/[method]_[object]_[scenario].json`
  - **Operations:** `validate_create`, `validate_get`, `validate_update`, `validate_delete`
  - **Loading:** Use `mocks.LoadJSONResponse(path)` in mock responders
  - **Scenario Coverage:** Include success, error, minimal, maximal, and hybrid scenarios

- **Terraform Config Loading:** Use appropriate helpers based on test type
  - **Unit Tests:** `helpers.ParseHCLFile("tests/terraform/unit/resource_name.tf")`
  - **Acceptance Tests:** `mocks.LoadTerraformConfigFile("config_name.tf")` wrapped with `acceptance.ConfiguredM365ProviderBlock(config)`

- **Mock Infrastructure:** Implement resource-specific mock systems following the `MockRegistrar` interface
  - **Mock Responders:** Create `mocks/responders.go` with struct implementing `RegisterMocks()`, `RegisterErrorMocks()`, and `CleanupMockState()`
  - **State Tracking:** Use mutex-protected maps for maintaining resource state across CRUD operations
  - **Registration:** Register mock with `mocks.GlobalRegistry.Register()` in package `init()` function

- **Acceptance Test Infrastructure:** Leverage shared helpers from `internal/acceptance/`
  - **Destroy Verification:** Use `destroy.CheckDestroyedAllFunc(testResource, resourceType, duration)` instead of custom destroy functions
  - **Test Logging:** Use `testlog.StepAction()` in `PreConfig` and `testlog.WaitForConsistency()` for observability
  - **External Dependencies:** Always include `random` provider; add `azuread` or others only when this provider lacks suitable resources
  - **Dependency Trees:** Create prerequisite resources dynamically using this provider first, external providers as fallback; never rely on environment variables for resource IDs

- **State Management Testing:** Verify Terraform state correctly reflects API responses across operations
- **Assignment Testing:** Include assignment validation for resources with assignment capabilities  
- **Import Testing:** Verify resource import functionality with `ImportState: true`, `ImportStateVerify: true`, and `ImportStateVerifyIgnore` for computed fields
- **Edge Case Testing:** Test minimal/maximal configurations, validation errors, and API error scenarios
- **Eventual Consistency:** Use `testlog.WaitForConsistency()` with `time.Sleep()` to handle Microsoft Graph eventual consistency
- **Permission Testing:** Verify appropriate error messages and hints for insufficient permissions using error handling package

### Network Debugging

```bash
make netdump  # Start mitmproxy for capturing and debugging API traffic
```

This tool is invaluable for debugging API interactions during test development.

### Shared Logic and Utilities

- **Use shared logic:** Use shared logic from `internal/resources/common/` for error handling, CRUD retry/timeout, validation, plan modifiers, schema helpers, state management, model construction, normalization, and shared models.
- **Location for new logic:** Place new reusable logic in the appropriate `common` subdirectory.
- **Shared models:** Use or extend shared models for any attribute/block reused across resources.
- **Helpers for transformation:** Use or extend helpers in `normalize/` or `constructors/` for data transformation/building.
- **Error handling and retry:** Always use shared error handling, retry, and timeout logic from `crud/` and `errors/`.
- **No duplication:** Do not duplicate logic—reference or extend shared implementations.
