package qdrant

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// getClient creates a client from the provided interface
// This client already contains the Authorization needed to invoke the API
// Returns: The client to call the backend API, TF Diagnostics
func getClient(m interface{}) (*qc.ClientWithResponses, diag.Diagnostics) {
	clientConfig, ok := m.(*ProviderConfig)
	if !ok {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: provided interface cannot be casted to ClientConfig"))
	}
	optsCallback := qc.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", clientConfig.ApiKey))
		return nil
	})
	apiClient, err := qc.NewClientWithResponses(clientConfig.BaseURL, optsCallback)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: %v", err))
	}
	return apiClient, nil
}

// getAccountUUID get the Account ID as UUID, if defined at resouce level that is used, otherwise it fallback to the default on, specified on provider level.
// if no account ID can be found an error will be returned.
func getAccountUUID(d *schema.ResourceData, m interface{}) (uuid.UUID, error) {
	// Get The account ID as UUID from the resource data
	if v, ok := d.GetOk("account_id"); ok {
		id := v.(string)
		if id != "" {
			return uuid.Parse(id)
		}
	}
	// Get From default (if any)
	if id := getDefaultAccountID(m); id != "" {
		return uuid.Parse(id)
	}
	return uuid.Nil, fmt.Errorf("cannot find account ID")
}

// getDefaultAccountID fetches the default account ID from the provided interface (containing the ClientConfig)
func getDefaultAccountID(m interface{}) string {
	clientConfig, ok := m.(*ProviderConfig)
	if !ok {
		return ""
	}
	return clientConfig.AccountID
}

// formatTime converts a time value to a standardized string format.
// t: The time value which can be of type time.Time or string.
// Returns a formatted time string in RFC3339 format if the input is of type time.Time,
// returns the input string unchanged if it is of type string, or an empty string for other types.
func formatTime(t interface{}) string {
	switch v := t.(type) {
	case time.Time:
		// Format time.Time to RFC3339 standard string format.
		return v.Format(time.RFC3339)
	case *time.Time:
		// Format time.Time to RFC3339 standard string format.
		return v.Format(time.RFC3339)
	case string:
		// Return string as is.
		return v
	default:
		// Return empty string for other types.
		return ""
	}
}

// getError returns a human readable error composed from the given HTTP validation error
func getError(error *qc.HTTPValidationError) string {
	if error == nil {
		return "no error"
	}
	details := error.Detail
	if details == nil {
		return "no details"
	}
	var result []string
	for _, ve := range *details {
		result = append(result, fmt.Sprintf("%s:%s", ve.Type, ve.Msg))
	}
	return strings.Join(result, ",")
}
