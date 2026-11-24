package customrequests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Raw JSON mode
//
// Every …RequestConfig accepts either a Kiota Parsable (RequestBody, and a ParsableFactory passed to the function)
// or plain Go values (RawBody and RawResponse) for endpoints whose types the SDK lacks or mishandles. Both modes send
// through the adapter's pipeline (base URL, authentication, retry, logging and ODataError mapping):
//   - RawBody is sent as "application/json" with RequestInformation.SetStreamContentAndContentType. json.RawMessage
//     and []byte values are sent as they are; any other value is marshalled with encoding/json.
//   - RawResponse is a pointer that the response body is decoded into with encoding/json, read through
//     adapter.SendPrimitive(…, "[]byte", …). *json.RawMessage and *[]byte receive the body as it is. A response with no
//     body (e.g. 204 No Content) leaves the target untouched and returns no error.
//   - A raw GET follows @odata.nextLink and combines every page's "value" array into one {"value": [...]} document.
//
// Setting both a Parsable body and a raw body, or both a ParsableFactory and a raw response target, is an error.

const jsonContentType = "application/json"

var (
	// ErrMixedBodyModes is returned when a config sets both a Parsable RequestBody and a RawBody.
	ErrMixedBodyModes = errors.New("custom request config sets both RequestBody and RawBody: use one body mode")
	// ErrMixedResponseModes is returned when a request is given both a ParsableFactory and a RawResponse target.
	ErrMixedResponseModes = errors.New("custom request sets both a ParsableFactory and RawResponse: use one response mode")
)

// isParsableSet reports whether body holds a Parsable, treating a typed nil pointer as unset.
func isParsableSet(body s.Parsable) bool {
	if body == nil {
		return false
	}
	value := reflect.ValueOf(body)
	return value.Kind() != reflect.Pointer || !value.IsNil()
}

// addRequestHeaders adds every configured header to requestInfo.
func addRequestHeaders(requestInfo *abstractions.RequestInformation, headers map[string]string) {
	for key, value := range headers {
		requestInfo.Headers.Add(key, value)
	}
}

// setRequestContent sets the request body from exactly one of body (Parsable mode) or rawBody (raw JSON mode).
// When rawBody is nil the Parsable is serialised exactly as before raw mode existed.
func setRequestContent(ctx context.Context, adapter abstractions.RequestAdapter, requestInfo *abstractions.RequestInformation, body s.Parsable, rawBody any) error {
	if rawBody == nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, jsonContentType, body); err != nil {
			return fmt.Errorf("error setting content: %w", err)
		}
		return nil
	}

	if isParsableSet(body) {
		return ErrMixedBodyModes
	}

	var content []byte
	switch value := rawBody.(type) {
	case json.RawMessage:
		content = value
	case []byte:
		content = value
	default:
		marshalled, err := json.Marshal(rawBody)
		if err != nil {
			return fmt.Errorf("marshal Graph JSON request: %w", err)
		}
		content = marshalled
	}

	requestInfo.SetStreamContentAndContentType(content, jsonContentType)
	return nil
}

// decodeRawResponse decodes body into target. An empty body (e.g. 204 No Content) leaves target untouched.
func decodeRawResponse(body []byte, target any) error {
	if len(body) == 0 {
		return nil
	}

	switch value := target.(type) {
	case *json.RawMessage:
		*value = append(json.RawMessage(nil), body...)
		return nil
	case *[]byte:
		*value = append([]byte(nil), body...)
		return nil
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode Graph JSON response: %w", err)
	}
	return nil
}

// sendForRawResponseOrNoContent sends requestInfo and decodes the response body into rawResponse. When rawResponse
// is nil the response body is discarded with SendNoContent, exactly as before raw mode existed. Errors are wrapped as
// "custom msgraph <METHOD> request failed: ...", and any ODataError stays reachable with errors.As.
func sendForRawResponseOrNoContent(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	requestInfo *abstractions.RequestInformation,
	rawResponse any,
	errorMappings abstractions.ErrorMappings,
) error {
	if rawResponse == nil {
		if err := adapter.SendNoContent(ctx, requestInfo, withODataErrorMappings(errorMappings)); err != nil {
			return fmt.Errorf("custom msgraph %s request failed: %w", requestInfo.Method.String(), err)
		}
		return nil
	}

	requestInfo.Headers.TryAdd("Accept", jsonContentType)
	body, err := sendGraphRequestForBytes(ctx, adapter, requestInfo, errorMappings)
	if err != nil {
		return err
	}
	return decodeRawResponse(body, rawResponse)
}

// sendGraphRequestForBytes sends requestInfo through the adapter's pipeline (authentication, retry, compression,
// logging and status checks) and returns the raw response body. Any 4xx or 5xx response is returned as an ODataError
// unless errorMappings are given; the error is wrapped as "custom msgraph <METHOD> request failed: ..." and the
// ODataError stays reachable with errors.As. A response with no body (e.g. 204 No Content) returns nil.
func sendGraphRequestForBytes(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	requestInfo *abstractions.RequestInformation,
	errorMappings abstractions.ErrorMappings,
) ([]byte, error) {
	tflog.Debug(ctx, "Making custom msgraph request", map[string]any{"method": requestInfo.Method.String(), "url_template": requestInfo.UrlTemplate})

	result, err := adapter.SendPrimitive(ctx, requestInfo, "[]byte", withODataErrorMappings(errorMappings))
	if err != nil {
		return nil, fmt.Errorf("custom msgraph %s request failed: %w", requestInfo.Method.String(), err)
	}

	body, _ := result.([]byte)
	return body, nil
}

// sendGraphGetCollectingPages sends a GET and, when the response is an OData collection with an @odata.nextLink,
// follows every next link through the adapter and returns one {"value": [...]} document holding every page's items.
// headers are sent with every page. Any other response body is returned as it is.
func sendGraphGetCollectingPages(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	requestInfo *abstractions.RequestInformation,
	headers map[string]string,
	errorMappings abstractions.ErrorMappings,
) ([]byte, error) {
	body, err := sendGraphRequestForBytes(ctx, adapter, requestInfo, errorMappings)
	if err != nil {
		return nil, err
	}

	var firstResponse ODataResponse
	if err := json.Unmarshal(body, &firstResponse); err != nil {
		return body, nil
	}

	// If no NextLink or no Value array, this isn't a paginated response
	if firstResponse.NextLink == "" || firstResponse.Value == nil {
		return body, nil
	}

	allResults := append([]json.RawMessage(nil), firstResponse.Value...)
	nextLink := firstResponse.NextLink

	tflog.Debug(ctx, "Pagination detected, retrieving additional pages", map[string]any{"itemsRetrieved": len(allResults)})

	for nextLink != "" {
		nextURL, err := url.Parse(nextLink)
		if err != nil {
			return nil, fmt.Errorf("error parsing @odata.nextLink %q: %w", nextLink, err)
		}

		pageRequest := abstractions.NewRequestInformation()
		pageRequest.Method = abstractions.GET
		pageRequest.SetUri(*nextURL)
		addRequestHeaders(pageRequest, headers)
		pageRequest.Headers.TryAdd("Accept", jsonContentType)

		body, err = sendGraphRequestForBytes(ctx, adapter, pageRequest, errorMappings)
		if err != nil {
			return nil, err
		}

		var pageResponse ODataResponse
		if err := json.Unmarshal(body, &pageResponse); err != nil {
			return nil, fmt.Errorf("error parsing paginated response: %w", err)
		}

		allResults = append(allResults, pageResponse.Value...)
		nextLink = pageResponse.NextLink

		tflog.Debug(ctx, "Retrieved additional page", map[string]any{"itemsRetrieved": len(allResults), "hasNextPage": nextLink != ""})
	}

	combined, err := json.Marshal(map[string]any{"value": allResults})
	if err != nil {
		return nil, fmt.Errorf("error combining paginated responses: %w", err)
	}
	return combined, nil
}
