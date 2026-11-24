package errors

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	betaodataerrors "github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// odataErrorFixture describes one Graph error that is built with both the v1.0 and beta SDK models.
type odataErrorFixture struct {
	statusCode int
	code       string
	message    string
	target     string
	headers    map[string]string
	details    []ErrorDetailInfo
	inner      *InnerErrorInfo
	innerDate  *time.Time
}

func (f odataErrorFixture) responseHeaders() *abstractions.ResponseHeaders {
	headers := abstractions.NewResponseHeaders()
	for key, value := range f.headers {
		headers.Add(key, value)
	}
	return headers
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (f odataErrorFixture) v1() *odataerrors.ODataError {
	mainError := odataerrors.NewMainError()
	mainError.SetCode(optionalString(f.code))
	mainError.SetMessage(optionalString(f.message))
	mainError.SetTarget(optionalString(f.target))

	details := make([]odataerrors.ErrorDetailsable, 0, len(f.details))
	for _, d := range f.details {
		detail := odataerrors.NewErrorDetails()
		detail.SetCode(optionalString(d.Code))
		detail.SetMessage(optionalString(d.Message))
		detail.SetTarget(optionalString(d.Target))
		details = append(details, detail)
	}
	if len(details) > 0 {
		mainError.SetDetails(details)
	}

	if f.inner != nil {
		innerError := odataerrors.NewInnerError()
		innerError.SetRequestId(optionalString(f.inner.RequestID))
		innerError.SetClientRequestId(optionalString(f.inner.ClientReqID))
		innerError.SetOdataType(optionalString(f.inner.ODataType))
		innerError.SetDate(f.innerDate)
		mainError.SetInnerError(innerError)
	}

	odataErr := odataerrors.NewODataError()
	odataErr.SetErrorEscaped(mainError)
	odataErr.SetStatusCode(f.statusCode)
	odataErr.SetResponseHeaders(f.responseHeaders())
	return odataErr
}

func (f odataErrorFixture) beta() *betaodataerrors.ODataError {
	mainError := betaodataerrors.NewMainError()
	mainError.SetCode(optionalString(f.code))
	mainError.SetMessage(optionalString(f.message))
	mainError.SetTarget(optionalString(f.target))

	details := make([]betaodataerrors.ErrorDetailsable, 0, len(f.details))
	for _, d := range f.details {
		detail := betaodataerrors.NewErrorDetails()
		detail.SetCode(optionalString(d.Code))
		detail.SetMessage(optionalString(d.Message))
		detail.SetTarget(optionalString(d.Target))
		details = append(details, detail)
	}
	if len(details) > 0 {
		mainError.SetDetails(details)
	}

	if f.inner != nil {
		innerError := betaodataerrors.NewInnerError()
		innerError.SetRequestId(optionalString(f.inner.RequestID))
		innerError.SetClientRequestId(optionalString(f.inner.ClientReqID))
		innerError.SetOdataType(optionalString(f.inner.ODataType))
		innerError.SetDate(f.innerDate)
		mainError.SetInnerError(innerError)
	}

	odataErr := betaodataerrors.NewODataError()
	odataErr.SetErrorEscaped(mainError)
	odataErr.SetStatusCode(f.statusCode)
	odataErr.SetResponseHeaders(f.responseHeaders())
	return odataErr
}

func odataErrorParityFixtures() map[string]odataErrorFixture {
	innerDate := time.Date(2026, 11, 24, 10, 30, 0, 0, time.UTC)

	return map[string]odataErrorFixture{
		"404 Request_ResourceNotFound": {
			statusCode: 404,
			code:       "Request_ResourceNotFound",
			message:    "Resource '00000000-0000-0000-0000-000000000001' does not exist.",
			headers:    map[string]string{"request-id": "req-404"},
		},
		"412 PreconditionFailed": {
			statusCode: 412,
			code:       "PreconditionFailed",
			message:    "The resource has been modified since it was last read.",
		},
		"429 with Retry-After": {
			statusCode: 429,
			code:       "TooManyRequests",
			message:    "Too many requests.",
			headers:    map[string]string{"Retry-After": "30", "x-throttled-reason": "ApplicationRequestsExceeded"},
		},
		"400 with details and inner error": {
			statusCode: 400,
			code:       "BadRequest",
			message:    "One or more properties are invalid.",
			target:     "displayName",
			details: []ErrorDetailInfo{
				{Code: "InvalidValue", Message: "displayName is required.", Target: "displayName"},
				{Code: "InvalidLength", Message: "description is too long."},
			},
			inner: &InnerErrorInfo{
				RequestID:   "inner-request-id",
				ClientReqID: "inner-client-request-id",
				ODataType:   "#microsoft.graph.innerError",
			},
			innerDate: &innerDate,
		},
	}
}

// TestUnitGraphError_BetaODataErrorParity asserts that a beta ODataError yields the same
// GraphErrorInfo as the equivalent v1.0 ODataError, directly and when wrapped.
func TestUnitGraphError_BetaODataErrorParity(t *testing.T) {
	ctx := context.Background()

	for name, fixture := range odataErrorParityFixtures() {
		t.Run(name, func(t *testing.T) {
			v1Info := GraphError(ctx, fixture.v1())
			betaInfo := GraphError(ctx, fixture.beta())

			assert.Equal(t, v1Info, betaInfo)
			assert.Equal(t, fixture.statusCode, betaInfo.StatusCode)
			assert.Equal(t, fixture.code, betaInfo.ErrorCode)
			assert.Equal(t, fixture.message, betaInfo.ErrorMessage)
			assert.Equal(t, fixture.target, betaInfo.Target)
			assert.True(t, betaInfo.IsODataError)

			wrappedV1Info := GraphError(ctx, fmt.Errorf("wrapped: %w", fixture.v1()))
			wrappedBetaInfo := GraphError(ctx, fmt.Errorf("wrapped: %w", fixture.beta()))
			assert.Equal(t, wrappedV1Info, wrappedBetaInfo)
			assert.Equal(t, fixture.code, wrappedBetaInfo.ErrorCode)
		})
	}
}

// TestUnitGraphError_BetaODataErrorFields checks the extracted details, inner error and headers explicitly.
func TestUnitGraphError_BetaODataErrorFields(t *testing.T) {
	ctx := context.Background()
	fixtures := odataErrorParityFixtures()

	badRequest := GraphError(ctx, fixtures["400 with details and inner error"].beta())
	require.Len(t, badRequest.ErrorDetails, 2)
	assert.Equal(t, ErrorDetailInfo{Code: "InvalidValue", Message: "displayName is required.", Target: "displayName"}, badRequest.ErrorDetails[0])
	assert.Equal(t, ErrorDetailInfo{Code: "InvalidLength", Message: "description is too long."}, badRequest.ErrorDetails[1])
	require.Len(t, badRequest.InnerErrors, 1)
	assert.Equal(t, "inner-request-id", badRequest.InnerErrors[0].RequestID)
	assert.Equal(t, "inner-client-request-id", badRequest.InnerErrors[0].ClientReqID)
	assert.Equal(t, "#microsoft.graph.innerError", badRequest.InnerErrors[0].ODataType)
	assert.Equal(t, "2026-11-24 10:30:00 +0000 UTC", badRequest.InnerErrors[0].Date)
	assert.Equal(t, "inner-request-id", badRequest.RequestID)
	assert.Equal(t, "inner-client-request-id", badRequest.ClientRequestID)
	assert.Equal(t, "2026-11-24 10:30:00 +0000 UTC", badRequest.ErrorDate)
	assert.Equal(t, CategoryValidation, badRequest.Category)

	throttled := GraphError(ctx, fixtures["429 with Retry-After"].beta())
	assert.Equal(t, "30", throttled.RetryAfter)
	assert.Equal(t, "ApplicationRequestsExceeded", throttled.ThrottledReason)
	assert.Equal(t, CategoryThrottling, throttled.Category)

	notFound := GraphError(ctx, fixtures["404 Request_ResourceNotFound"].beta())
	assert.Equal(t, "req-404", notFound.RequestID)
	assert.Equal(t, CategoryService, notFound.Category)
}

// graphResponseErrorStub implements only the status and main-error accessors, without ApiErrorable headers.
type graphResponseErrorStub[T any] struct {
	statusCode int
	mainError  T
}

func (e *graphResponseErrorStub[T]) Error() string      { return "graph response error" }
func (e *graphResponseErrorStub[T]) GetStatusCode() int { return e.statusCode }
func (e *graphResponseErrorStub[T]) GetErrorEscaped() T { return e.mainError }

// TestUnitGraphError_GraphResponseErrorInterfaceParity covers errors that expose GetErrorEscaped
// but are not the SDK ODataError type, for both SDK model packages.
func TestUnitGraphError_GraphResponseErrorInterfaceParity(t *testing.T) {
	ctx := context.Background()

	v1Main := odataerrors.NewMainError()
	v1Main.SetCode(optionalString("Request_ResourceNotFound"))
	v1Main.SetMessage(optionalString("not found"))
	betaMain := betaodataerrors.NewMainError()
	betaMain.SetCode(optionalString("Request_ResourceNotFound"))
	betaMain.SetMessage(optionalString("not found"))

	v1Info := GraphError(ctx, &graphResponseErrorStub[odataerrors.MainErrorable]{statusCode: 404, mainError: v1Main})
	betaInfo := GraphError(ctx, &graphResponseErrorStub[betaodataerrors.MainErrorable]{statusCode: 404, mainError: betaMain})

	assert.Equal(t, v1Info, betaInfo)
	assert.Equal(t, 404, betaInfo.StatusCode)
	assert.Equal(t, "Request_ResourceNotFound", betaInfo.ErrorCode)
	assert.Equal(t, "not found", betaInfo.ErrorMessage)

	nilMainInfo := GraphError(ctx, &graphResponseErrorStub[betaodataerrors.MainErrorable]{statusCode: 404})
	assert.Equal(t, 404, nilMainInfo.StatusCode)
	assert.Empty(t, nilMainInfo.ErrorCode)
}

// TestUnitHandleKiotaGraphError_BetaODataErrorDiagnosticsParity asserts that diagnostics for a beta error
// carry the same code, details and tracking information as the v1.0 error.
func TestUnitHandleKiotaGraphError_BetaODataErrorDiagnosticsParity(t *testing.T) {
	ctx := context.Background()
	fixtures := odataErrorParityFixtures()

	for _, name := range []string{"400 with details and inner error", "412 PreconditionFailed", "429 with Retry-After"} {
		t.Run(name, func(t *testing.T) {
			fixture := fixtures[name]
			v1Resp := &resource.CreateResponse{}
			betaResp := &resource.CreateResponse{}

			HandleKiotaGraphError(ctx, fixture.v1(), v1Resp, constants.TfOperationCreate, nil)
			HandleKiotaGraphError(ctx, fixture.beta(), betaResp, constants.TfOperationCreate, nil)

			require.True(t, betaResp.Diagnostics.HasError())
			assert.Equal(t, v1Resp.Diagnostics, betaResp.Diagnostics)
			assert.Contains(t, betaResp.Diagnostics.Errors()[0].Detail(), "Code: "+fixture.code)
		})
	}

	readResp := &resource.ReadResponse{}
	HandleKiotaGraphError(ctx, fixtures["429 with Retry-After"].beta(), readResp, constants.TfOperationRead, nil)
	require.True(t, readResp.Diagnostics.HasError())
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "Retry-After: 30 seconds")
	assert.Contains(t, readResp.Diagnostics.Errors()[0].Detail(), "Code: TooManyRequests")

	actionMessage := HandleKiotaGraphErrorForAction(ctx, fixtures["404 Request_ResourceNotFound"].beta())
	assert.Equal(t, HandleKiotaGraphErrorForAction(ctx, fixtures["404 Request_ResourceNotFound"].v1()), actionMessage)
	assert.Contains(t, actionMessage, "Code: Request_ResourceNotFound")
}

// TestUnitRetryable_BetaODataErrorParity asserts that retry classification is identical for v1.0 and beta errors.
func TestUnitRetryable_BetaODataErrorParity(t *testing.T) {
	ctx := context.Background()

	for name, fixture := range odataErrorParityFixtures() {
		t.Run(name, func(t *testing.T) {
			v1Info := GraphError(ctx, fixture.v1())
			betaInfo := GraphError(ctx, fixture.beta())

			assert.Equal(t, IsRetryableReadError(&v1Info), IsRetryableReadError(&betaInfo))
			assert.Equal(t, IsNonRetryableReadError(&v1Info), IsNonRetryableReadError(&betaInfo))
			assert.Equal(t, IsRetryableDeleteError(&v1Info), IsRetryableDeleteError(&betaInfo))
			assert.Equal(t, IsNonRetryableDeleteError(&v1Info), IsNonRetryableDeleteError(&betaInfo))
		})
	}
}
