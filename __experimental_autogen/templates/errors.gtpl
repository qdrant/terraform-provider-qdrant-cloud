package qdrantcloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/moul/http2curl"
)

// ValidationError describes detailed errors in a request,
// typically associated with specific fields or parameters.
type ValidationError struct {
	Loc  []interface{} `json:"loc"`  // Location of the error, e.g., ["body", "password"]
	Msg  string        `json:"msg"`  // Human-readable error message
	Type string        `json:"type"` // Error type or identifier
}

// HTTPValidationError captures an array of ValidationError,
// which corresponds to the 422 Unprocessable Entity HTTP error.
type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

// HTTPValidationErrorString captures a string error message,
// which corresponds to the 4xx Unprocessable Entity HTTP error.
type HTTPValidationErrorString struct {
	Detail string `json:"detail"`
}

// Helper function to handle non-200 HTTP responses.
func handleNon422Response(resp *http.Response) diag.Diagnostics {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading response body: %s", err))
	}

	var errResp HTTPValidationError
	err = json.Unmarshal(bodyBytes, &errResp)
	if err != nil {
		var errRespString HTTPValidationErrorString
		err = json.Unmarshal(bodyBytes, &errRespString)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error decoding validation error response: %s", err))
		}

		return diag.FromErr(fmt.Errorf("API error: %s", errRespString.Detail))
	}

	return diag.Errorf("API error: %v", errResp)
}

// DebugPrintCurlCommand function to print the curl command for a given HTTP request
// req: The HTTP request for which the curl command is to be generated.
func DebugPrintCurlCommand(req *http.Request) {
	curlCommand, err := http2curl.GetCurlCommand(req)
	if err != nil {
		fmt.Println("Error getting curl command:", err)
	} else {
		fmt.Println(curlCommand)
	}
}

// GetStringOrEmpty returns the string value of an interface or an empty string if the interface is nil.
// v: The interface to extract the string value from.
func GetStringOrEmpty(v interface{}) string {
	if v == nil {
		return ""
	}
	switch value := v.(type) {
	case string:
		return value
	case *string:
		if value != nil {
			return *value
		}
	}
	return ""
}
