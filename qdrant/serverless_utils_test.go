package qdrant

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spacev1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/v1"
)

func TestListAllPages(t *testing.T) {
	t.Run("when multiple pages are returned, it collects all items and passes the tokens", func(t *testing.T) {
		pages := map[string]struct {
			items []int
			next  string
		}{
			"":   {items: []int{1, 2}, next: "p2"},
			"p2": {items: []int{3}, next: "p3"},
			"p3": {items: []int{4, 5}, next: ""},
		}
		var tokens []string
		result, err := listAllPages(func(pageToken *string) ([]int, string, error) {
			token := ""
			if pageToken != nil {
				token = *pageToken
			}
			tokens = append(tokens, token)
			p := pages[token]
			return p.items, p.next, nil
		})
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
		assert.Equal(t, []string{"", "p2", "p3"}, tokens)
	})
	t.Run("when fetching a page fails, it returns the error", func(t *testing.T) {
		calls := 0
		_, err := listAllPages(func(_ *string) ([]int, string, error) {
			calls++
			if calls == 2 {
				return nil, "", errors.New("boom")
			}
			return []int{1}, "next", nil
		})
		assert.EqualError(t, err, "boom")
	})
}

func TestFlattenProtoEnum(t *testing.T) {
	assert.Empty(t, flattenProtoEnum(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_UNSPECIFIED))
	assert.Equal(t, "SPACE_STATE_PHASE_READY", flattenProtoEnum(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_READY))
}

func TestIsNotFound(t *testing.T) {
	assert.True(t, isNotFound(status.Error(codes.NotFound, "gone")))
	assert.False(t, isNotFound(status.Error(codes.Internal, "boom")))
	assert.False(t, isNotFound(errors.New("plain")))
}

func TestGrpcErrorDiag(t *testing.T) {
	d := grpcErrorDiag("error creating", " (request ID: x)", status.Error(codes.InvalidArgument, "bad name"))
	require.Len(t, d, 1)
	assert.Equal(t, "error creating: invalid argument (request ID: x): bad name", d[0].Summary)

	d = grpcErrorDiag("error creating", "", status.Error(codes.Internal, "boom"))
	require.Len(t, d, 1)
	assert.Contains(t, d[0].Summary, "error creating: ")
	assert.Contains(t, d[0].Summary, "boom")
}

func TestImportSpaceChildState(t *testing.T) {
	s := map[string]*schema.Schema{
		"space_id": {Type: schema.TypeString, Required: true},
	}
	importer := importSpaceChildState("space_id", "key_id")

	t.Run("when the ID has the expected format, it sets the space ID and resource ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, s, map[string]interface{}{})
		d.SetId("space-1/key-1")
		result, err := importer(context.Background(), d, nil)
		require.NoError(t, err)
		require.Len(t, result, 1)
		assert.Equal(t, "key-1", result[0].Id())
		assert.Equal(t, "space-1", result[0].Get("space_id"))
	})
	for _, id := range []string{"key-1", "/key-1", "space-1/", "a/b/c"} {
		t.Run("when the ID is "+id+", it returns an error", func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, map[string]interface{}{})
			d.SetId(id)
			_, err := importer(context.Background(), d, nil)
			assert.ErrorContains(t, err, "expected <space_id>/<key_id>")
		})
	}
}
