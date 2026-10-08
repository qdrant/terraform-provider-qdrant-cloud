package qdrant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	spaceauthv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/auth/v1"
	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
	spacev1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/v1"
)

type mockSpaceServiceClient struct {
	spacev1.SpaceServiceClient
	space *spacev1.Space
	err   error
}

func (m *mockSpaceServiceClient) GetSpace(_ context.Context, _ *spacev1.GetSpaceRequest, _ ...grpc.CallOption) (*spacev1.GetSpaceResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &spacev1.GetSpaceResponse{Space: m.space}, nil
}

func spaceInPhase(phase spacev1.SpaceStatePhase, url string) *spacev1.Space {
	state := &spacev1.SpaceState{Phase: phase, Reason: "some reason"}
	if url != "" {
		state.Endpoint = &spacev1.SpaceEndpoint{Url: url}
	}
	return &spacev1.Space{Id: "space-id", State: state}
}

func TestServerlessSpaceReadyRefreshFunc(t *testing.T) {
	refresh := func(m *mockSpaceServiceClient) (interface{}, string, error) {
		return serverlessSpaceReadyRefreshFunc(m, context.Background(), "account-id", "space-id")()
	}

	t.Run("when the space is ready with an endpoint, it returns ready", func(t *testing.T) {
		_, state, err := refresh(&mockSpaceServiceClient{space: spaceInPhase(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_READY, "https://x")})
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitReady, state)
	})
	t.Run("when the space is ready without an endpoint, it keeps waiting", func(t *testing.T) {
		_, state, err := refresh(&mockSpaceServiceClient{space: spaceInPhase(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_READY, "")})
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitPending, state)
	})
	t.Run("when the space is processing, it keeps waiting", func(t *testing.T) {
		_, state, err := refresh(&mockSpaceServiceClient{space: spaceInPhase(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_PROCESSING, "https://x")})
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitPending, state)
	})
	t.Run("when the space is disabled, it fails with the reason", func(t *testing.T) {
		_, _, err := refresh(&mockSpaceServiceClient{space: spaceInPhase(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_DISABLED, "")})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "SPACE_STATE_PHASE_DISABLED")
		assert.Contains(t, err.Error(), "some reason")
	})
	t.Run("when the API fails, it returns the error", func(t *testing.T) {
		_, _, err := refresh(&mockSpaceServiceClient{err: status.Error(codes.Internal, "boom")})
		assert.Error(t, err)
	})
}

func TestServerlessSpaceDeletedRefreshFunc(t *testing.T) {
	refresh := func(m *mockSpaceServiceClient) (interface{}, string, error) {
		return serverlessSpaceDeletedRefreshFunc(m, context.Background(), "account-id", "space-id")()
	}

	t.Run("when the space is not found, it returns deleted", func(t *testing.T) {
		result, state, err := refresh(&mockSpaceServiceClient{err: status.Error(codes.NotFound, "gone")})
		require.NoError(t, err)
		assert.NotNil(t, result, "a nil result would be treated as not found by the state change conf")
		assert.Equal(t, serverlessWaitDeleted, state)
	})
	t.Run("when the space has a deletion time, it returns deleted", func(t *testing.T) {
		space := spaceInPhase(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_DELETING, "")
		space.DeletedAt = timestamppb.Now()
		_, state, err := refresh(&mockSpaceServiceClient{space: space})
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitDeleted, state)
	})
	t.Run("when the space is still being deleted, it keeps waiting", func(t *testing.T) {
		_, state, err := refresh(&mockSpaceServiceClient{space: spaceInPhase(spacev1.SpaceStatePhase_SPACE_STATE_PHASE_DELETING, "")})
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitPending, state)
	})
	t.Run("when the API fails, it returns the error", func(t *testing.T) {
		_, _, err := refresh(&mockSpaceServiceClient{err: status.Error(codes.Internal, "boom")})
		assert.Error(t, err)
	})
}

type mockSpaceAPIKeyServiceClient struct {
	spaceauthv1.SpaceApiKeyServiceClient
	keys []*spaceauthv1.SpaceApiKey
}

func (m *mockSpaceAPIKeyServiceClient) ListSpaceApiKeys(_ context.Context, _ *spaceauthv1.ListSpaceApiKeysRequest, _ ...grpc.CallOption) (*spaceauthv1.ListSpaceApiKeysResponse, error) {
	return &spaceauthv1.ListSpaceApiKeysResponse{Items: m.keys}, nil
}

func TestServerlessSpaceAPIKeyReadyRefreshFunc(t *testing.T) {
	key := func(id string, phase spaceauthv1.SpaceApiKeyStatePhase) *spaceauthv1.SpaceApiKey {
		return &spaceauthv1.SpaceApiKey{Id: id, State: &spaceauthv1.SpaceApiKeyState{Phase: phase}}
	}
	refresh := func(keys ...*spaceauthv1.SpaceApiKey) (interface{}, string, error) {
		m := &mockSpaceAPIKeyServiceClient{keys: keys}
		return serverlessSpaceAPIKeyReadyRefreshFunc(m, context.Background(), "account-id", "space-id", "key-id")()
	}

	t.Run("when the key is ready, it returns ready", func(t *testing.T) {
		result, state, err := refresh(
			key("other", spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_PROCESSING),
			key("key-id", spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_READY),
		)
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitReady, state)
		assert.Equal(t, "key-id", result.(*spaceauthv1.SpaceApiKey).GetId())
	})
	t.Run("when the key is processing, it keeps waiting", func(t *testing.T) {
		_, state, err := refresh(key("key-id", spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_PROCESSING))
		require.NoError(t, err)
		assert.Equal(t, serverlessWaitPending, state)
	})
	t.Run("when the key is disabled, it fails", func(t *testing.T) {
		_, _, err := refresh(key("key-id", spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_DISABLED))
		assert.ErrorContains(t, err, "SPACE_API_KEY_STATE_PHASE_DISABLED")
	})
	t.Run("when the key is missing, it fails", func(t *testing.T) {
		_, _, err := refresh(key("other", spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_READY))
		assert.ErrorContains(t, err, "not found")
	})
}

type mockServerlessBackupServiceClient struct {
	spacebackupv1.BackupServiceClient
	backup *spacebackupv1.Backup
}

func (m *mockServerlessBackupServiceClient) GetBackup(_ context.Context, _ *spacebackupv1.GetBackupRequest, _ ...grpc.CallOption) (*spacebackupv1.GetBackupResponse, error) {
	return &spacebackupv1.GetBackupResponse{Backup: m.backup}, nil
}

func TestServerlessBackupRefreshFunc(t *testing.T) {
	refresh := func(s spacebackupv1.BackupStatus) (interface{}, string, error) {
		m := &mockServerlessBackupServiceClient{backup: &spacebackupv1.Backup{Id: "backup-id", Status: s}}
		return serverlessBackupRefreshFunc(m, context.Background(), "account-id", "backup-id")()
	}

	_, state, err := refresh(spacebackupv1.BackupStatus_BACKUP_STATUS_SUCCEEDED)
	require.NoError(t, err)
	assert.Equal(t, serverlessWaitReady, state)

	for _, s := range []spacebackupv1.BackupStatus{
		spacebackupv1.BackupStatus_BACKUP_STATUS_UNSPECIFIED,
		spacebackupv1.BackupStatus_BACKUP_STATUS_RUNNING,
	} {
		_, state, err := refresh(s)
		require.NoError(t, err, s.String())
		assert.Equal(t, serverlessWaitPending, state, s.String())
	}
	for _, s := range []spacebackupv1.BackupStatus{
		spacebackupv1.BackupStatus_BACKUP_STATUS_FAILED,
		spacebackupv1.BackupStatus_BACKUP_STATUS_FAILED_TO_SYNC,
		spacebackupv1.BackupStatus_BACKUP_STATUS_SKIPPED,
		spacebackupv1.BackupStatus_BACKUP_STATUS_NOT_FOUND,
	} {
		_, _, err := refresh(s)
		assert.ErrorContains(t, err, s.String())
	}
}
