Unit Test Strategy

What Should the Unit Test Verify?

Resource CRUD Logic:

Does Create() correctly translate Terraform config → Graph API calls?
Does Read() correctly translate Graph API responses → Terraform state?
Does Update() handle state changes properly?
Does Delete() clean up correctly?

State Management:

Are Terraform attributes correctly mapped to/from the Graph API model?
Are computed fields properly handled?
Is the resource ID correctly set and used?


Error Handling:

Does the resource properly handle API errors?
Are appropriate diagnostics returned?


Schema Validation:

Are required fields enforced?
Are field types and constraints working?

--------------------------------

Unit Test Starts
    ↓
Set Mock Environment Variables  
    ↓
Activate httpmock (BEFORE any HTTP clients created)
    ↓
Register Mock HTTP Responses
    ↓
Provider.Configure() detects unit test → uses http.DefaultClient
    ↓
Graph SDK clients created with http.DefaultClient
    ↓
Resource CRUD operations call Graph SDK
    ↓
Graph SDK makes HTTP calls via http.DefaultClient
    ↓
httpmock intercepts ALL HTTP calls → returns mock responses
    ↓
NO REAL API CALLS MADE ✅