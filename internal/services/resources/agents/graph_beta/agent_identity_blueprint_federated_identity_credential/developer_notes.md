# Agent Identity Blueprint Federated Identity Credential - Developer Notes

**Date:** 2025-12-04  
**Status:** ❌ Not Implementable (API Issues)  
**Resource:** `microsoft365_graph_beta_agents_agent_identity_blueprint_federated_identity_credential`

## Summary

This resource cannot be implemented at this time due to Microsoft Graph API behavior where federated identity credentials created via the `agentIdentityBlueprint` cast endpoint are not retrievable via any GET endpoint.

---

## Test Sequence

### Test 1: Standard SDK Call for Create

**Approach:** Use the standard Kiota SDK method to create the federated identity credential.

```go
r.client.
Applications().
ByApplicationId(blueprintID).
FederatedIdentityCredentials().
Post(ctx, requestBody, nil)
```

**Result:** ❌ Failed  
**Error:** `Bad Request - 400: A resource with type 'Microsoft.DirectoryServices.AgentIdentityBlueprint' was found, but it is not assignable to the expected type 'Microsoft.DirectoryServices.FederatedIdentityCredential'`

**Analysis:** The standard endpoint doesn't work for `agentIdentityBlueprint` applications. Microsoft documentation indicates that a cast endpoint is required.

---

### Test 2: Create with Cast Endpoint (Custom Request)

**Approach:** Use a custom POST request with the cast endpoint as per Microsoft documentation.

```
POST /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials
```

**Result:** ✅ Create Succeeded  
**Response:** Credential created with ID, Name, Issuer, Subject, Audiences returned correctly.

```
Created credential - ID: c364b785-b56b-41d3-953d-f35705ae9644, Name: acc-test-fic-minimal-82ty16i8
```

---

### Test 3: Read with Standard Endpoint

**Approach:** Read the created credential using the standard SDK endpoint.

```
GET /applications/{blueprintId}/federatedIdentityCredentials/{credentialId}
```

**Result:** ❌ Failed  
**Error:** `Request_ResourceNotFound - 404`

---

### Test 4: Read with Cast Endpoint

**Approach:** Read using the cast endpoint as suggested by Microsoft documentation.

```
GET /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
```

**Result:** ❌ Failed  
**Error:** `Request_ResourceNotFound - 404: Resource '{credentialId}' does not exist`

---

### Test 5: Add 10 Second Delay for Eventual Consistency

**Approach:** Add a 10-second wait after Create before attempting Read, to allow for eventual consistency.

```go
tflog.Debug(ctx, "Waiting 10 seconds for eventual consistency after create")
time.Sleep(10 * time.Second)
```

**Result:** ❌ Still Failed  
**Error:** Same 404 errors on all Read attempts.

---

### Test 6: List Credentials via Standard Endpoint

**Approach:** List all federated identity credentials to see if the created credential appears.

```
GET /applications/{blueprintId}/federatedIdentityCredentials
```

**Result:** ❌ Empty List  
**Response:**
```json
{"@odata.context":"...","value":[]}
```

**Analysis:** Zero credentials returned even though Create reported success.

---

### Test 7: List Credentials via Cast Endpoint

**Approach:** List using the cast endpoint.

```
GET /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials
```

**Result:** ❌ Empty List  
**Response:**
```json
{"@odata.context":"...","value":[]}
```

**Analysis:** Also returns zero credentials.

---

### Test 8: GET by ID via Cast Endpoint (After List)

**Approach:** After confirming list is empty, attempt GET by ID.

```
GET /applications/{blueprintId}/microsoft.graph.agentIdentityBlueprint/federatedIdentityCredentials/{credentialId}
```

**Result:** ❌ Failed  
**Error:**
```json
{"error":{"code":"Request_ResourceNotFound","message":"Resource '{credentialId}' does not exist..."}}
```

---

### Test 9: Map State Directly from Create Response

**Approach:** Skip the Read after Create entirely and map Terraform state directly from the Create response data.

```go
object.ID = types.StringValue(*createdCredential.GetId())
object.Name = types.StringValue(*createdCredential.GetName())
object.Issuer = types.StringValue(*createdCredential.GetIssuer())
// ... etc
```

**Result:** ⚠️ Partially Works  
**Issue:** Create works, but subsequent `terraform plan` or `terraform refresh` operations will fail because Read cannot retrieve the resource state from the API. This breaks the Terraform lifecycle.

---

## API Behavior Summary

| Operation | Endpoint | Result |
|-----------|----------|--------|
| Create | Cast endpoint | ✅ Works - Returns credential with ID |
| Read by ID (standard) | Standard endpoint | ❌ 404 Not Found |
| Read by ID (cast) | Cast endpoint | ❌ 404 Not Found |
| List (standard) | Standard endpoint | ❌ Empty array |
| List (cast) | Cast endpoint | ❌ Empty array |

---

## Conclusion

**The Microsoft Graph Beta API for `agentIdentityBlueprint` federated identity credentials has a fundamental issue:**

1. **Create works** via the cast endpoint and returns valid data including the credential ID
2. **All Read operations fail** - both List and GET by ID return empty results or 404 errors
3. The credential reportedly appears in the Azure Portal GUI (per user observation), suggesting the data is persisted somewhere, but it is not accessible via the documented API endpoints
4. Even with 10+ seconds of eventual consistency delay, the credential is not retrievable

**This prevents implementation of a functioning Terraform resource because:**
- Terraform requires the ability to Read resource state to detect drift
- Terraform requires Read after Create to verify the resource was created correctly
- Without Read functionality, `terraform plan` and `terraform refresh` cannot work

---

## Recommendations

1. **File a Microsoft Support Ticket** - Report this API behavior as a bug
2. **Monitor Graph API Updates** - Check future beta API releases for fixes
3. **Re-test Periodically** - The beta API may be updated without notice

---

## Code Status

The resource code has been written and tested but is **not registered** in the provider due to the API issues described above. The code remains in the codebase for future use when/if the API is fixed.

**Files:**
- `crud.go` - CRUD operations implemented
- `construct.go` - Request body construction
- `model.go` - Terraform resource model
- `resource.go` - Resource schema definition
- `state.go` - State mapping functions
- `mocks/responders.go` - Test mock responders

---

## References

- [Microsoft Graph API - Create federatedIdentityCredential for agentIdentityBlueprint](https://learn.microsoft.com/en-us/graph/api/agentidentityblueprint-post-federatedidentitycredentials?view=graph-rest-beta&tabs=http)
- [Microsoft Graph API - Get federatedIdentityCredential](https://learn.microsoft.com/en-us/graph/api/federatedidentitycredential-get?view=graph-rest-beta&tabs=http)

