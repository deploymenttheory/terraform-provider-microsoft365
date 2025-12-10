// Package directory provides common helpers for msgraph directory beta api operations to
// soft delete and hard delete targetted Microsoft Graph resources.
//
// Supported resources:
// - application
// - agentIdentityBlueprint
// - agentIdentity
// - agentIdentityBlueprintPrincipal
// - agentUser
// - certificateBasedAuthPki
// - certificateAuthorityDetail
// - externalUserProfile
// - group
// - pendingExternalUserProfile
// - servicePrincipal
// - user
//
// This package handles the two-step deletion process required by Microsoft Graph:
// 1. Soft delete - moves the resource to the deleted items collection
// 2. Hard delete - permanently removes the resource from deleted items
//
// Both operations include verification with retry logic and jitter to handle
// eventual consistency and prevent thundering herd problems and they take into account
// the hcl supplied boolean value hard_delete to determine if the hard delete operation should be performed.
//
// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/directory-deleteditems-delete?view=graph-rest-beta
package directory

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// -----------------------------------------------------------------------------
// Constants
// -----------------------------------------------------------------------------

const (
	// defaultMaxRetries is the default number of retry attempts for verification
	defaultMaxRetries = 10

	// defaultRetryInterval is the default time to wait between retries
	defaultRetryInterval = 3 * time.Second

	// defaultHardDeleteRetryInterval is the default time to wait between hard delete verification retries
	defaultHardDeleteRetryInterval = 2 * time.Second

	// jitterFactor determines the maximum additional random delay (0.5 = up to 50% extra)
	jitterFactor = 0.5

	// HTTP status code for not found
	httpStatusNotFound = 404
)

// Known OData error codes that indicate a resource was not found
var notFoundErrorCodes = map[string]bool{
	"ResourceNotFound":         true,
	"Request_ResourceNotFound": true,
	"ItemNotFound":             true,
	"ErrorItemNotFound":        true,
	"NotFound":                 true,
}

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// ResourceType represents the type of resource being deleted.
// This is used for logging and error messages.
type ResourceType string

const (
	ResourceTypeApplication      ResourceType = "application"
	ResourceTypeServicePrincipal ResourceType = "servicePrincipal"
	ResourceTypeUser             ResourceType = "user"
	ResourceTypeGroup            ResourceType = "group"
)

// DeleteOptions configures the delete operation behavior.
type DeleteOptions struct {
	// MaxRetries is the maximum number of retry attempts for verification (default: 10)
	MaxRetries int

	// RetryInterval is the time to wait between retries (default: 3 seconds)
	RetryInterval time.Duration

	// ResourceType is the type of resource being deleted (for logging)
	ResourceType ResourceType

	// ResourceID is the ID of the resource being deleted
	ResourceID string

	// ResourceName is the display name of the resource (optional, for logging)
	ResourceName string
}

// SoftDeleteFunc is a function that performs the initial soft delete operation.
// The function should delete the resource, which moves it to the deleted items collection.
type SoftDeleteFunc func(ctx context.Context) error

// -----------------------------------------------------------------------------
// Constructor Functions
// -----------------------------------------------------------------------------

// DefaultDeleteOptions returns sensible default options for delete operations.
func DefaultDeleteOptions() DeleteOptions {
	return DeleteOptions{
		MaxRetries:    defaultMaxRetries,
		RetryInterval: defaultRetryInterval,
	}
}

// -----------------------------------------------------------------------------
// Public API - Main Entry Point
// -----------------------------------------------------------------------------

// ExecuteDeleteWithVerification orchestrates the full delete workflow:
//  1. Soft delete (moves resource to deleted items) + verification
//  2. Hard delete (permanently removes from deleted items) + verification (if hardDelete is true)
//
// If hardDelete is false, only the soft delete is performed and verified.
// If hardDelete is true, both soft delete and permanent delete are performed with verification.
//
// The function handles eventual consistency through polling with jittered retry intervals.
func ExecuteDeleteWithVerification(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	softDeleteFunc SoftDeleteFunc,
	hardDelete bool,
	opts DeleteOptions,
) error {
	opts = applyDefaults(opts)

	// Step 1: Soft delete with verification
	if err := ExecuteSoftDelete(ctx, client, softDeleteFunc, opts); err != nil {
		return err
	}

	// If hard delete is not requested, we're done
	if !hardDelete {
		tflog.Info(ctx, fmt.Sprintf("Soft delete only - %s %s moved to deleted items (can be restored within 30 days)", opts.ResourceType, opts.ResourceID))
		return nil
	}

	tflog.Info(ctx, fmt.Sprintf("Hard delete enabled - proceeding with permanent deletion of %s %s", opts.ResourceType, opts.ResourceID))

	// Step 2: Hard delete with verification
	if err := ExecuteHardDelete(ctx, client, opts); err != nil {
		return err
	}

	tflog.Info(ctx, fmt.Sprintf("Complete deletion successful for %s %s", opts.ResourceType, opts.ResourceID))
	return nil
}

// -----------------------------------------------------------------------------
// Public API - Individual Operations
// -----------------------------------------------------------------------------

// ExecuteSoftDelete performs a soft delete operation and verifies the resource
// appears in deleted items. This handles the eventual consistency delay.
//
// The function:
//  1. Calls the provided deleteFunc to perform the soft delete
//  2. Polls the deleted items collection until the resource appears
//
// Returns nil when the resource has been successfully soft deleted and verified
// to be in the deleted items collection.
func ExecuteSoftDelete(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deleteFunc SoftDeleteFunc, opts DeleteOptions) error {
	opts = applyDefaults(opts)

	// Step 1: Perform the soft delete
	tflog.Info(ctx, fmt.Sprintf("Performing soft delete for %s %s", opts.ResourceType, opts.ResourceID))

	if err := deleteFunc(ctx); err != nil {
		errorInfo := errors.GraphError(ctx, err)
		return fmt.Errorf("soft delete API call failed for %s %s [HTTP %d, Code: %s]: %s",
			opts.ResourceType, opts.ResourceID, errorInfo.StatusCode, errorInfo.ErrorCode, errorInfo.ErrorMessage)
	}

	tflog.Debug(ctx, fmt.Sprintf("Soft delete API call completed for %s %s, waiting for resource to appear in deleted items", opts.ResourceType, opts.ResourceID))

	// Step 2: Verify the resource appears in deleted items
	return verifySoftDelete(ctx, client, opts)
}

// ExecuteHardDelete permanently deletes a resource from the deleted items collection
// and verifies the deletion was successful by confirming the resource is gone.
//
// The function:
//  1. Calls DELETE on /directory/deletedItems/{id}
//  2. Polls to verify the resource is no longer found (404)
//
// If the resource is already not found during the DELETE call, this is treated as success
// (the resource is already permanently deleted).
//
// Returns nil when the resource has been successfully permanently deleted.
func ExecuteHardDelete(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, opts DeleteOptions) error {
	// Use shorter defaults for hard delete verification
	if opts.MaxRetries <= 0 {
		opts.MaxRetries = 5
	}
	if opts.RetryInterval <= 0 {
		opts.RetryInterval = defaultHardDeleteRetryInterval
	}

	// Step 1: Perform the permanent delete
	tflog.Info(ctx, fmt.Sprintf("Permanently deleting %s %s from deleted items", opts.ResourceType, opts.ResourceID))

	err := client.
		Directory().
		DeletedItems().
		ByDirectoryObjectId(opts.ResourceID).
		Delete(ctx, nil)

	if err != nil {
		// If the resource is already not found, treat as success - it's already permanently deleted
		if isNotFoundError(ctx, err) {
			tflog.Info(ctx, fmt.Sprintf("Resource %s %s already permanently deleted (not found in deleted items)", opts.ResourceType, opts.ResourceID))
			return nil
		}
		errorInfo := errors.GraphError(ctx, err)
		return fmt.Errorf("hard delete API call failed for %s %s [HTTP %d, Code: %s]: %s",
			opts.ResourceType, opts.ResourceID, errorInfo.StatusCode, errorInfo.ErrorCode, errorInfo.ErrorMessage)
	}

	tflog.Debug(ctx, fmt.Sprintf("Hard delete API call completed for %s %s, verifying resource is gone", opts.ResourceType, opts.ResourceID))

	// Step 2: Verify the resource is no longer in deleted items
	return verifyHardDelete(ctx, client, opts)
}

// -----------------------------------------------------------------------------
// Private Helper Functions
// -----------------------------------------------------------------------------

// applyDefaults ensures DeleteOptions has reasonable default values.
func applyDefaults(opts DeleteOptions) DeleteOptions {
	if opts.MaxRetries <= 0 {
		opts.MaxRetries = defaultMaxRetries
	}
	if opts.RetryInterval <= 0 {
		opts.RetryInterval = defaultRetryInterval
	}
	return opts
}

// verifySoftDelete polls the deleted items collection until the resource appears.
func verifySoftDelete(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, opts DeleteOptions) error {
	var lastError error
	var lastErrorInfo errors.GraphErrorInfo

	for attempt := 1; attempt <= opts.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("soft delete verification cancelled after %d/%d attempts for %s %s: context error: %w",
				attempt, opts.MaxRetries, opts.ResourceType, opts.ResourceID, err)
		}

		// Try to get the resource from deleted items
		_, err := client.
			Directory().
			DeletedItems().
			ByDirectoryObjectId(opts.ResourceID).
			Get(ctx, nil)

		if err == nil {
			// Resource found in deleted items - soft delete verified
			tflog.Info(ctx, fmt.Sprintf("Soft delete verified: %s %s found in deleted items (attempt %d/%d)", opts.ResourceType, opts.ResourceID, attempt, opts.MaxRetries))
			return nil
		}

		// Store the last error for final error message
		lastError = err
		lastErrorInfo = errors.GraphError(ctx, err)

		if attempt == opts.MaxRetries {
			return fmt.Errorf("soft delete verification failed for %s %s after %d attempts [last HTTP %d, Code: %s]: resource did not appear in deleted items. Last error: %s",
				opts.ResourceType, opts.ResourceID, opts.MaxRetries, lastErrorInfo.StatusCode, lastErrorInfo.ErrorCode, lastErrorInfo.ErrorMessage)
		}

		sleepWithJitter(ctx, opts.RetryInterval, attempt, opts.MaxRetries, "Resource not yet in deleted items")
	}

	return fmt.Errorf("soft delete verification failed for %s %s after %d attempts [last HTTP %d, Code: %s]: %s",
		opts.ResourceType, opts.ResourceID, opts.MaxRetries, lastErrorInfo.StatusCode, lastErrorInfo.ErrorCode, lastError)
}

// verifyHardDelete polls the deleted items collection until the resource is not found.
func verifyHardDelete(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, opts DeleteOptions) error {
	for attempt := 1; attempt <= opts.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("hard delete verification cancelled after %d/%d attempts for %s %s: context error: %w",
				attempt, opts.MaxRetries, opts.ResourceType, opts.ResourceID, err)
		}

		// Try to get the resource from deleted items
		_, err := client.
			Directory().
			DeletedItems().
			ByDirectoryObjectId(opts.ResourceID).
			Get(ctx, nil)

		if err != nil {
			if isNotFoundError(ctx, err) {
				tflog.Info(ctx, fmt.Sprintf("Hard delete verified: %s %s has been permanently deleted (attempt %d/%d)", opts.ResourceType, opts.ResourceID, attempt, opts.MaxRetries))
				return nil
			}
			// Unexpected error during verification - log details and return
			errorInfo := errors.GraphError(ctx, err)
			return fmt.Errorf("hard delete verification encountered unexpected error for %s %s on attempt %d/%d [HTTP %d, Code: %s]: %s",
				opts.ResourceType, opts.ResourceID, attempt, opts.MaxRetries, errorInfo.StatusCode, errorInfo.ErrorCode, errorInfo.ErrorMessage)
		}

		// Resource still exists
		if attempt == opts.MaxRetries {
			return fmt.Errorf("hard delete verification failed for %s %s: resource still exists in deleted items after %d attempts (resource may need manual cleanup)",
				opts.ResourceType, opts.ResourceID, opts.MaxRetries)
		}

		sleepWithJitter(ctx, opts.RetryInterval, attempt, opts.MaxRetries, "Resource still in deleted items")
	}

	return fmt.Errorf("hard delete verification failed for %s %s after %d attempts: resource may still exist in deleted items",
		opts.ResourceType, opts.ResourceID, opts.MaxRetries)
}

// isNotFoundError checks if the error indicates the resource was not found.
// Uses the kiota errors package for consistent error code handling.
func isNotFoundError(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}

	// Extract error information using the standard error handling package
	errorInfo := errors.GraphError(ctx, err)

	// Check HTTP status code
	if errorInfo.StatusCode == httpStatusNotFound {
		tflog.Debug(ctx, "Not found error detected via HTTP status code 404", map[string]any{
			"status_code": errorInfo.StatusCode,
			"error_code":  errorInfo.ErrorCode,
		})
		return true
	}

	// Check known OData error codes for not found
	if notFoundErrorCodes[errorInfo.ErrorCode] {
		tflog.Debug(ctx, "Not found error detected via OData error code", map[string]any{
			"status_code": errorInfo.StatusCode,
			"error_code":  errorInfo.ErrorCode,
		})
		return true
	}

	// Fallback: check error message for common "not found" patterns
	// This handles cases where the error structure doesn't match expected OData format
	if containsNotFoundPattern(errorInfo.ErrorMessage) {
		tflog.Debug(ctx, "Not found error detected via error message pattern", map[string]any{
			"status_code":   errorInfo.StatusCode,
			"error_code":    errorInfo.ErrorCode,
			"error_message": errorInfo.ErrorMessage,
		})
		return true
	}

	return false
}

// containsNotFoundPattern checks if an error message contains common "not found" patterns.
func containsNotFoundPattern(message string) bool {
	if message == "" {
		return false
	}

	notFoundPatterns := []string{
		"does not exist",
		"not found",
		"could not be found",
		"cannot be found",
		"no longer exists",
	}

	lowerMessage := strings.ToLower(message)
	for _, pattern := range notFoundPatterns {
		if strings.Contains(lowerMessage, pattern) {
			return true
		}
	}
	return false
}

// calculateJitteredInterval returns a duration with random jitter applied.
// The jitter is calculated as: baseInterval + random(0, baseInterval * jitterFactor)
// This helps prevent thundering herd problems when multiple resources are deleted simultaneously.
func calculateJitteredInterval(baseInterval time.Duration) time.Duration {
	maxJitter := float64(baseInterval) * jitterFactor
	jitter := time.Duration(rand.Float64() * maxJitter)
	return baseInterval + jitter
}

// sleepWithJitter sleeps for a jittered interval and logs the retry attempt.
func sleepWithJitter(ctx context.Context, baseInterval time.Duration, attempt, maxRetries int, message string) {
	jitteredDelay := calculateJitteredInterval(baseInterval)
	tflog.Debug(ctx, fmt.Sprintf("%s, retrying in %v (attempt %d/%d)", message, jitteredDelay, attempt, maxRetries))
	time.Sleep(jitteredDelay)
}
