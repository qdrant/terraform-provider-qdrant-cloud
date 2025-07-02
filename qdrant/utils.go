package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	requestIDTrailerField = "qc-trace-id"
)

// getRequestID fetches the humanized Request ID from the provided metadata (or an empty string if not available).
func getRequestID(metadata metadata.MD) string {
	reqIDs := metadata.Get(requestIDTrailerField)
	if len(reqIDs) == 0 {
		return ""
	}
	return fmt.Sprintf(" [%s]", strings.Join(reqIDs, "|"))
}

// getClientConnection creates a client connection from the provided interface.
// This client need to be invoked with the enriched context, which aleady contains the Authorization needed to invoke the API.
// Returns: The connection from the backend API, the enriched context to use, TF Diagnostics.
func getClientConnection(ctx context.Context, m interface{}) (*grpc.ClientConn, context.Context, diag.Diagnostics) {
	clientConfig, ok := m.(*ProviderConfig)
	if !ok {
		return nil, nil, diag.FromErr(fmt.Errorf("error initializing client: provided interface cannot be casted to ClientConfig"))
	}
	if clientConfig.BaseURL == "" {
		return nil, nil, diag.FromErr(fmt.Errorf("error initializing client: provided ClientConfig.BaseURL not set"))
	}
	// Set up a connection to the server.
	tc := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: clientConfig.Insecure,
	})
	conn, err := grpc.NewClient(clientConfig.BaseURL, grpc.WithTransportCredentials(tc))
	if err != nil {
		return nil, nil, diag.FromErr(fmt.Errorf("error initializing client: cannot create gRPC client: %w", err))
	}
	// Add Access Token
	ctxWithToken := metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("apikey %s", clientConfig.ApiKey))
	// Return result
	return conn, ctxWithToken, nil
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

// getDefaultAccountID fetches the default account ID from the provided interface (containing the ClientConfig).
func getDefaultAccountID(m interface{}) string {
	clientConfig, ok := m.(*ProviderConfig)
	if !ok {
		return ""
	}
	return clientConfig.AccountID
}

// parseTime parses the provided value and returns it as  (or nil if it cannot be parsed).
// The provided string should be in RCF3339 format.
func parseTime(v string) *timestamppb.Timestamp {
	result, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil
	}
	return timestamppb.New(result)
}

// formatTime formats the provided proto timestamp into a string
// The resulted string will be in RCF3339 format, so it can be parsed with parseTime again.
func formatTime(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format(time.RFC3339)
}

// parseDuration parses the provided value and returns it as  (or nil if it cannot be parsed).
// The provided string should be in Go duration format.
func parseDuration(v string) *durationpb.Duration {
	result, err := time.ParseDuration(v)
	if err != nil {
		return nil
	}
	return durationpb.New(result)
}

// formatDuration formats the provided proto duration into a string
// The resulted string will be in Go duration format, so it can be parsed with parseDuration again.
func formatDuration(d *durationpb.Duration) string {
	if d == nil {
		return ""
	}
	return d.AsDuration().String()
}

// suppressDurationDiff is a DiffSuppressFunc that suppresses diffs for duration strings
// if they are semantically equivalent (e.g., "1h" and "60m").
func suppressDurationDiff(k, old, new string, d *schema.ResourceData) bool {
	oldDuration, err := time.ParseDuration(old)
	if err != nil {
		return false
	}
	newDuration, err := time.ParseDuration(new)
	if err != nil {
		return false
	}
	return oldDuration == newDuration
}
