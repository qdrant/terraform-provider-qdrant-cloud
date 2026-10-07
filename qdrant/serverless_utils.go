package qdrant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	serverlessPollInterval   = 5 * time.Second
	serverlessDefaultTimeout = 20 * time.Minute
	// serverlessListPageSize is the page size used when listing serverless resources (max allowed by the API).
	serverlessListPageSize = int32(250)

	serverlessWaitPending = "pending"
	serverlessWaitReady   = "ready"
	serverlessWaitDeleted = "deleted"
)

// protoEnum is implemented by all generated protobuf enums.
type protoEnum interface {
	Number() protoreflect.EnumNumber
	String() string
}

// flattenProtoEnum returns the enum name, or an empty string for the zero (UNSPECIFIED) value,
// so unset Computed enums never produce a diff.
func flattenProtoEnum(e protoEnum) string {
	if e.Number() == 0 {
		return ""
	}
	return e.String()
}

// listAllPages calls fetch until no next page token is returned and collects all items.
// fetch receives the page token to use (nil for the first page).
func listAllPages[T any](fetch func(pageToken *string) ([]T, string, error)) ([]T, error) {
	var result []T
	var pageToken *string
	for {
		items, next, err := fetch(pageToken)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
		if next == "" {
			return result, nil
		}
		pageToken = newPointer(next)
	}
}

// grpcErrorDiag converts a gRPC error into diagnostics, using a friendlier message for invalid arguments.
func grpcErrorDiag(errorPrefix, reqID string, err error) diag.Diagnostics {
	if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
		return diag.Errorf("%s: invalid argument%s: %s", errorPrefix, reqID, st.Message())
	}
	return diag.FromErr(fmt.Errorf("%s%s: %w", errorPrefix, reqID, err))
}

// isNotFound returns true if the provided error is a gRPC NotFound error.
func isNotFound(err error) bool {
	st, ok := status.FromError(err)
	return ok && st.Code() == codes.NotFound
}

// setFlattened stores all flattened values into the resource data.
func setFlattened(d *schema.ResourceData, values map[string]interface{}) error {
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("cannot set %s: %w", k, err)
		}
	}
	return nil
}

// configurableMaxItems returns n for resources, and 0 for data sources (MaxItems is not allowed on computed-only fields).
func configurableMaxItems(n int, asDataSource bool) int {
	if asDataSource {
		return 0
	}
	return n
}

// importSpaceChildState returns an importer for resources identified by `<space_id>/<id>`,
// which sets the given space ID field and uses the second part as resource ID.
func importSpaceChildState(spaceIDField, idName string) schema.StateContextFunc {
	return func(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
		parts := strings.Split(d.Id(), "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("unexpected format of ID (%s), expected <space_id>/<%s>", d.Id(), idName)
		}
		if err := d.Set(spaceIDField, parts[0]); err != nil {
			return nil, fmt.Errorf("error setting %s: %w", spaceIDField, err)
		}
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, nil
	}
}
