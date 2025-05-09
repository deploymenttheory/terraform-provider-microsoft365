package constants

import "sync"

// GraphSDKMutex is a global lock used to serialize Microsoft Graph SDK (Kiota) API calls.
//
// Reason:
// Kiota's middleware (e.g., HeadersInspectionHandler) modifies shared header maps during HTTP request processing.
// In Go, maps are not concurrency-safe — concurrent writes cause immediate fatal runtime panics ("concurrent map writes").
//
// Terraform executes Read operations across multiple resources in parallel, leading to multiple Graph API calls at the same time.
// Without locking, simultaneous mutation of the shared headers map causes the plugin to crash.
//
// By acquiring this mutex around Graph API calls (Get, Post, Patch, Delete),
// we ensure only one request is processed at a time through the Kiota pipeline, preventing concurrency issues and plugin crashes.
//
// You cannot set a mutex as a constant in go.
var GraphSDKMutex sync.Mutex
